package process

import (
	"context"
	"fmt"
	"time"
)

// WaitForProcess polls until a process named name exists, then opens and
// returns it. interval defaults to 250ms when <= 0. Honors ctx.Done() — pass
// a context.WithTimeout to cap the wait.
//
//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//	p, err := process.WaitForProcess(ctx, "target.exe", 0)
func WaitForProcess(ctx context.Context, name string, interval time.Duration) (*Process, error) {
	if interval <= 0 {
		interval = 250 * time.Millisecond
	}

	if p, err := OpenProcessByName(name); err == nil {
		return p, nil
	}

	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("waiting for process %q: %w", name, ctx.Err())
		case <-t.C:
			if p, err := OpenProcessByName(name); err == nil {
				return p, nil
			}
		}
	}
}
