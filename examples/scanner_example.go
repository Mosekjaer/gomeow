//go:build ignore

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"

	"gomeow/hotkey"
	"gomeow/memory"
	"gomeow/process"
)

// Run with: go run examples/scanner_example.go <process-name>
//
// Workflow once attached:
//   - Type a value at the prompt and press Enter to set the target.
//   - Press F1 to run the first scan.
//   - Change the value in the target process (take damage, spend ammo, etc.),
//     type the new value, press F2 to narrow the result list.
//   - Repeat F2 until one address remains.
//   - F3 resets so you can scan for a different value.
//   - End or Ctrl+C exits.
//
// On Linux you'll need root (or the relevant ptrace caps) for memory.Read to
// work against another process.
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: scanner_example <process-name>")
		os.Exit(2)
	}

	p, err := process.OpenProcessByName(os.Args[1])
	if err != nil {
		log.Fatalf("open process: %v", err)
	}
	defer p.Close()
	fmt.Printf("attached to %s (PID %d)\n", p.Name, p.PID)

	s := memory.NewScanner(p)
	s.OnResults = func(count int, addrs []uintptr) {
		switch {
		case count == 0:
			fmt.Println("  → no matches")
		case count <= 8:
			fmt.Printf("  → %d match(es): %s\n", count, fmtAddrs(addrs))
		default:
			fmt.Printf("  → %d matches (showing first 8): %s\n", count, fmtAddrs(addrs[:8]))
		}
	}
	s.OnError = func(err error) {
		fmt.Fprintln(os.Stderr, "scan error:", err)
	}

	const (
		keyF1  = 0x70 // VK_F1 / first scan
		keyF2  = 0x71 // VK_F2 / next scan
		keyF3  = 0x72 // VK_F3 / reset
		keyEnd = 0x23 // VK_END / quit
	)

	hk := hotkey.New()
	hk.Register(keyF1, s.BindFirst())
	hk.Register(keyF2, s.BindNext())
	hk.Register(keyF3, s.BindReset())

	done := make(chan struct{})
	hk.Register(keyEnd, func() { close(done) })

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	// Stdin loop: read int32 values, store on Scanner.Value.
	var mu sync.Mutex
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("value> ")
			line, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			v, err := strconv.ParseInt(line, 10, 32)
			if err != nil {
				fmt.Fprintln(os.Stderr, "  bad number:", err)
				continue
			}
			mu.Lock()
			s.Value = int32(v)
			mu.Unlock()
			fmt.Printf("  set Value = %d. F1=first scan, F2=next, F3=reset.\n", v)
		}
	}()

	fmt.Println("F1=first scan, F2=next scan, F3=reset, End=quit")

	tick := time.NewTicker(16 * time.Millisecond)
	defer tick.Stop()

	for {
		select {
		case <-sig:
			return
		case <-done:
			fmt.Println("bye")
			return
		case <-tick.C:
			hk.Poll()
		}
	}
}

func fmtAddrs(a []uintptr) string {
	parts := make([]string, len(a))
	for i, addr := range a {
		parts[i] = fmt.Sprintf("0x%X", addr)
	}
	return "[" + strings.Join(parts, " ") + "]"
}
