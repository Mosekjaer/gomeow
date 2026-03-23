// Package hotkey provides a polling-based hotkey manager with callback edge
// detection. It wraps gomeow/input.KeyPressed so it works on both Windows and
// Linux without platform-specific code.
//
// Typical usage inside an overlay loop:
//
//	m := hotkey.New()
//	m.RegisterToggle(overlay.VK_F1, func(on bool) {
//	    fmt.Println("ESP:", on)
//	})
//	m.Register(overlay.VK_F2, func() { fmt.Println("shoot") })
//
//	for overlay.Loop() {
//	    m.Poll()
//	    overlay.BeginDrawing()
//	    // ... draw ...
//	    overlay.EndDrawing()
//	}
//
// Callbacks fire synchronously inside Poll. If a callback is slow, wrap its
// body in `go func() { ... }()` to avoid stalling the caller's frame.
//
// On Windows the key codes are virtual key codes (VK_*). On Linux the
// underlying input.KeyPressed expects X11 keysyms — use the VK constants
// from gomeow/overlay, which are translated automatically via
// overlay.IsKeyDownGlobal, or pass keysyms directly with NewWithPoll if you
// want to bypass the translation.
package hotkey

import (
	"sync"

	"gomeow/input"
)

type kind int

const (
	kindPress kind = iota
	kindToggle
	kindHold
)

type entry struct {
	kind     kind
	on       func()
	toggleFn func(state bool)
	holdFn   func(down bool)
	toggle   bool
}

// Manager tracks registered hotkeys and dispatches callbacks on key edges.
// Safe for concurrent Register/Unregister against Poll. Callbacks themselves
// run on the goroutine that calls Poll.
type Manager struct {
	mu      sync.Mutex
	entries map[int]*entry
	prev    map[int]bool
	poll    func(int) bool
}

// New returns a Manager that samples key state via input.KeyPressed.
func New() *Manager {
	return &Manager{
		entries: map[int]*entry{},
		prev:    map[int]bool{},
		poll:    input.KeyPressed,
	}
}

// Register fires fn on every rising edge (key down transition) of key.
func (m *Manager) Register(key int, fn func()) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries[key] = &entry{kind: kindPress, on: fn}
}

// RegisterToggle flips a boolean on every rising edge and calls fn with the
// new state.
func (m *Manager) RegisterToggle(key int, fn func(state bool)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries[key] = &entry{kind: kindToggle, toggleFn: fn}
}

// RegisterHold fires fn on both edges (down=true on press, down=false on
// release). Useful for hold-to-aim style features.
func (m *Manager) RegisterHold(key int, fn func(down bool)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries[key] = &entry{kind: kindHold, holdFn: fn}
}

// Unregister removes the binding for key (no-op if not registered).
func (m *Manager) Unregister(key int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.entries, key)
	delete(m.prev, key)
}

// Clear removes all bindings.
func (m *Manager) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries = map[int]*entry{}
	m.prev = map[int]bool{}
}

// Snapshot returns a copy of the current down state for every registered
// key. State only updates when Poll runs, so call Poll first if you want
// fresh data. Useful for HUDs and debug overlays.
func (m *Manager) Snapshot() map[int]bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make(map[int]bool, len(m.entries))
	for key := range m.entries {
		out[key] = m.prev[key]
	}
	return out
}

// ToggleStates returns a copy of the current toggle state for every
// RegisterToggle binding. Keys not bound as toggles are omitted.
func (m *Manager) ToggleStates() map[int]bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make(map[int]bool, len(m.entries))
	for key, e := range m.entries {
		if e.kind == kindToggle {
			out[key] = e.toggle
		}
	}
	return out
}

// Poll samples every registered key once and dispatches callbacks. Call once
// per frame from the overlay loop (or any tight loop ~60Hz).
func (m *Manager) Poll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for key, e := range m.entries {
		isDown := m.poll(key)
		wasDown := m.prev[key]
		m.prev[key] = isDown

		switch e.kind {
		case kindPress:
			if isDown && !wasDown && e.on != nil {
				e.on()
			}
		case kindToggle:
			if isDown && !wasDown {
				e.toggle = !e.toggle
				if e.toggleFn != nil {
					e.toggleFn(e.toggle)
				}
			}
		case kindHold:
			if isDown != wasDown && e.holdFn != nil {
				e.holdFn(isDown)
			}
		}
	}
}
