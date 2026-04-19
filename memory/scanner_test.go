package memory

import (
	"encoding/binary"
	"testing"

	"gomeow/process"
)

// fakeScanner builds a Scanner backed by an in-memory byte slice. The "process"
// address space is rebased so that addr 0 → buf[0]. regionsFn returns one
// page covering the whole buffer.
func fakeScanner(buf []byte) *Scanner {
	s := &Scanner{}
	s.reader = func(addr uintptr, out []byte) error {
		copy(out, buf[addr:addr+uintptr(len(out))])
		return nil
	}
	s.regionsFn = func(_ ScanOptions) ([]process.Page, error) {
		return []process.Page{{Start: 0, End: uintptr(len(buf)), Size: uintptr(len(buf))}}, nil
	}
	return s
}

func putI32(buf []byte, off int, v int32) {
	binary.LittleEndian.PutUint32(buf[off:], uint32(v))
}

func TestScannerFirstAndNext(t *testing.T) {
	buf := make([]byte, 64)
	// Plant the value 100 at offsets 8 and 32, and unrelated values elsewhere.
	putI32(buf, 0, 1)
	putI32(buf, 4, 2)
	putI32(buf, 8, 100)
	putI32(buf, 12, 3)
	putI32(buf, 16, 4)
	putI32(buf, 20, 5)
	putI32(buf, 24, 6)
	putI32(buf, 28, 7)
	putI32(buf, 32, 100)
	putI32(buf, 36, 8)

	s := fakeScanner(buf)

	if err := s.First(int32(100), ScanOptions{}); err != nil {
		t.Fatalf("First: %v", err)
	}
	got := s.Results()
	if len(got) != 2 || got[0] != 8 || got[1] != 32 {
		t.Fatalf("First: expected [8 32], got %v", got)
	}

	// Mutate offset 8 to a new value, leave 32 alone, then Next(87) and
	// expect only the address that now holds 87 to remain — but here we want
	// to model the natural cheat-engine flow: the value transitioned from
	// 100 → 87 at one location. So mutate 8 to 87 and then Next(87).
	putI32(buf, 8, 87)
	if err := s.Next(int32(87)); err != nil {
		t.Fatalf("Next: %v", err)
	}
	got = s.Results()
	if len(got) != 1 || got[0] != 8 {
		t.Fatalf("Next: expected [8], got %v", got)
	}
}

func TestScannerStepAlignment(t *testing.T) {
	buf := make([]byte, 32)
	// Place the int32 100 starting at byte offset 1 — misaligned.
	binary.LittleEndian.PutUint32(buf[1:], 100)

	s := fakeScanner(buf)

	// Default Step (0 → 4-byte aligned for int32) should miss it.
	if err := s.First(int32(100), ScanOptions{}); err != nil {
		t.Fatalf("First: %v", err)
	}
	if s.Count() != 0 {
		t.Fatalf("aligned scan: expected 0 hits, got %d", s.Count())
	}

	// Step=1 should find it.
	if err := s.First(int32(100), ScanOptions{Step: 1}); err != nil {
		t.Fatalf("First step=1: %v", err)
	}
	got := s.Results()
	if len(got) != 1 || got[0] != 1 {
		t.Fatalf("step=1: expected [1], got %v", got)
	}
}

func TestScannerBytes(t *testing.T) {
	buf := []byte("hello world hello world")
	s := fakeScanner(buf)

	if err := s.First([]byte("hello"), ScanOptions{Step: 1}); err != nil {
		t.Fatalf("First: %v", err)
	}
	got := s.Results()
	if len(got) != 2 || got[0] != 0 || got[1] != 12 {
		t.Fatalf("expected [0 12], got %v", got)
	}
}

func TestScannerNextSizeMismatch(t *testing.T) {
	buf := make([]byte, 32)
	putI32(buf, 0, 5)
	s := fakeScanner(buf)
	if err := s.First(int32(5), ScanOptions{}); err != nil {
		t.Fatal(err)
	}
	if err := s.Next(int64(5)); err == nil {
		t.Fatal("expected size mismatch error")
	}
}

func TestScannerCallbacks(t *testing.T) {
	buf := make([]byte, 16)
	putI32(buf, 0, 42)
	s := fakeScanner(buf)

	var lastCount int
	s.OnResults = func(c int, _ []uintptr) { lastCount = c }
	s.Value = int32(42)
	if err := s.FirstScan(); err != nil {
		t.Fatal(err)
	}
	if lastCount != 1 {
		t.Fatalf("OnResults: expected 1, got %d", lastCount)
	}
}

func TestScannerHotkeyBindings(t *testing.T) {
	buf := make([]byte, 16)
	putI32(buf, 0, 7)
	s := fakeScanner(buf)

	var errs []error
	s.OnError = func(e error) { errs = append(errs, e) }

	first := s.BindFirst()
	next := s.BindNext()
	reset := s.BindReset()

	// No value set yet → BindFirst should route to OnError.
	first()
	if len(errs) != 1 {
		t.Fatalf("expected 1 error from unset value, got %d", len(errs))
	}

	s.Value = int32(7)
	first()
	if s.Count() != 1 {
		t.Fatalf("after First: expected 1 result, got %d", s.Count())
	}
	next()
	if s.Count() != 1 {
		t.Fatalf("after Next: expected 1 result, got %d", s.Count())
	}
	reset()
	if s.Count() != 0 {
		t.Fatalf("after Reset: expected 0, got %d", s.Count())
	}
}
