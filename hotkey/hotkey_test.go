package hotkey

import "testing"

// newWithPoll constructs a Manager with a fake poll source. Internal helper
// kept here (not exposed in the public API) so production callers can't
// accidentally bypass the input package.
func newWithPoll(p func(int) bool) *Manager {
	m := New()
	m.poll = p
	return m
}

func TestPress(t *testing.T) {
	state := map[int]bool{}
	m := newWithPoll(func(k int) bool { return state[k] })

	fires := 0
	m.Register(1, func() { fires++ })

	m.Poll() // initial sample, key up — no fire
	if fires != 0 {
		t.Fatalf("expected 0 fires, got %d", fires)
	}

	state[1] = true
	m.Poll() // rising edge — fires once
	if fires != 1 {
		t.Fatalf("expected 1 fire on rising edge, got %d", fires)
	}

	m.Poll() // still down — no additional fire
	if fires != 1 {
		t.Fatalf("expected no fire while held, got %d", fires)
	}

	state[1] = false
	m.Poll() // falling edge — no fire for press
	state[1] = true
	m.Poll() // second rising edge
	if fires != 2 {
		t.Fatalf("expected 2 fires after re-press, got %d", fires)
	}
}

func TestToggle(t *testing.T) {
	state := map[int]bool{}
	m := newWithPoll(func(k int) bool { return state[k] })

	var got []bool
	m.RegisterToggle(2, func(on bool) { got = append(got, on) })

	state[2] = true
	m.Poll() // rising → on
	state[2] = false
	m.Poll() // release, no fire
	state[2] = true
	m.Poll() // rising → off

	if len(got) != 2 || got[0] != true || got[1] != false {
		t.Fatalf("expected [true false], got %v", got)
	}
}

func TestHold(t *testing.T) {
	state := map[int]bool{}
	m := newWithPoll(func(k int) bool { return state[k] })

	var got []bool
	m.RegisterHold(3, func(down bool) { got = append(got, down) })

	state[3] = true
	m.Poll() // down=true
	m.Poll() // still down, no fire
	state[3] = false
	m.Poll() // down=false

	if len(got) != 2 || got[0] != true || got[1] != false {
		t.Fatalf("expected [true false], got %v", got)
	}
}

func TestUnregisterAndClear(t *testing.T) {
	state := map[int]bool{}
	m := newWithPoll(func(k int) bool { return state[k] })

	fires := 0
	m.Register(4, func() { fires++ })
	m.Unregister(4)

	state[4] = true
	m.Poll()
	if fires != 0 {
		t.Fatalf("unregistered key fired: %d", fires)
	}

	m.Register(4, func() { fires++ })
	m.Clear()
	m.Poll()
	if fires != 0 {
		t.Fatalf("cleared key fired: %d", fires)
	}
}
