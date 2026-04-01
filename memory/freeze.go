package memory

import (
	"time"

	"gomeow/process"
)

// Freeze repeatedly writes value to addr at the given interval, holding the
// memory location pinned to that value. Returns a cancel function that stops
// the goroutine. The returned cancel is safe to call multiple times.
//
// Typical use with a hotkey toggle:
//
//	var cancel func()
//	hk.RegisterToggle(overlay.VK_F3, func(on bool) {
//	    if on {
//	        buf := make([]byte, 4)
//	        binary.LittleEndian.PutUint32(buf, 100)
//	        cancel = memory.Freeze(p, healthAddr, buf, 50*time.Millisecond)
//	    } else if cancel != nil {
//	        cancel()
//	    }
//	})
//
// Write errors are silently swallowed — a freezer that stops on the first
// transient EFAULT during a level transition would be useless. If you need
// failure telemetry, wrap Write yourself.
func Freeze(p *process.Process, addr uintptr, value []byte, interval time.Duration) (cancel func()) {
	done := make(chan struct{})
	buf := make([]byte, len(value))
	copy(buf, value)

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				_ = Write(p, addr, buf)
			}
		}
	}()

	var once bool
	return func() {
		if once {
			return
		}
		once = true
		close(done)
	}
}
