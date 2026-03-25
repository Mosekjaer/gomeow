//go:build ignore

package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"gomeow/hotkey"
)

// Run with: go run examples/hotkey_example.go
//
// VK codes match the constants in gomeow/overlay (VK_F1=0x70, VK_F2=0x71,
// VK_END=0x23). On Windows they map directly to virtual keys via
// GetAsyncKeyState; on Linux input.KeyPressed expects X11 keysyms instead, so
// for Linux either use overlay.IsKeyPressedGlobal as the poll source or pass
// XK_F1=0xFFBE etc. directly. This example uses Windows VK codes — it's the
// shortest readable form.
func main() {
	const (
		keyF1  = 0x70
		keyF2  = 0x71
		keyF3  = 0x72
		keyEnd = 0x23
	)

	m := hotkey.New()

	m.RegisterToggle(keyF1, func(on bool) {
		fmt.Println("ESP:", onOff(on))
	})

	m.Register(keyF2, func() {
		fmt.Println("shoot")
	})

	m.RegisterHold(keyF3, func(down bool) {
		if down {
			fmt.Println("aim: hold")
		} else {
			fmt.Println("aim: release")
		}
	})

	done := make(chan struct{})
	m.Register(keyEnd, func() {
		fmt.Println("bye")
		close(done)
	})

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	fmt.Println("F1=ESP toggle, F2=shoot, F3=hold-to-aim, End=quit (Ctrl+C also works)")

	tick := time.NewTicker(16 * time.Millisecond) // ~60 Hz
	defer tick.Stop()

	for {
		select {
		case <-sig:
			return
		case <-done:
			return
		case <-tick.C:
			m.Poll()
		}
	}
}

func onOff(b bool) string {
	if b {
		return "on"
	}
	return "off"
}
