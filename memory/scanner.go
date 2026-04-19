package memory

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"os"

	"gomeow/process"
)

// ValueKind identifies the encoding of a scanned value.
type ValueKind int

const (
	KindBytes ValueKind = iota
	KindInt32
	KindUint32
	KindInt64
	KindUint64
	KindFloat32
	KindFloat64
)

// ScanOptions configures a scan run. Zero values are sane defaults.
type ScanOptions struct {
	// Module restricts the scan to a single named module. Empty string =
	// the process's own base module (Process.Name).
	Module string

	// Step is the byte stride between scan positions. 0 → use the size of
	// the value type (type-aligned). Set to 1 to find misaligned matches.
	Step int

	// ChunkSize caps the per-region read buffer in bytes. 0 → 16 MiB.
	// Regions larger than ChunkSize are read in successive chunks.
	ChunkSize int
}

// Scanner runs Cheat Engine-style first/next value scans against a process.
//
// Typical usage:
//
//	s := memory.NewScanner(p)
//	s.Value = int32(100)
//	if err := s.FirstScan(); err != nil { ... }     // initial filter
//	// player takes damage in-game
//	s.Value = int32(87)
//	if err := s.NextScan(); err != nil { ... }      // narrow
//	fmt.Println(s.Results())                         // surviving addresses
//
// Hotkey-friendly bindings are also provided:
//
//	hk.Register(VK_F1, s.BindFirst())
//	hk.Register(VK_F2, s.BindNext())
//	hk.Register(VK_F3, s.BindReset())
type Scanner struct {
	p       *process.Process
	results []uintptr
	kind    ValueKind
	size    int
	target  []byte // last encoded value
	opts    ScanOptions

	// Value is the target value for the next FirstScan / NextScan call.
	// Set this before invoking the no-arg variants. Supported types:
	// int8/16/32/64, uint8/16/32/64, float32, float64, []byte.
	Value any

	// Options is applied on the next First call. Changing Options between
	// First and Next has no effect — the kind/size are locked after First.
	Options ScanOptions

	// OnResults, if non-nil, is called after each successful scan with the
	// new result count and a copy of the addresses.
	OnResults func(count int, addrs []uintptr)

	// OnError, if non-nil, is called when a scan errors. Used by the Bind*
	// helpers so hotkey-driven scans can report errors without panicking.
	OnError func(error)

	// Test seams. nil → use defaults backed by memory.Read and
	// p.EnumMemoryRegions.
	reader    func(addr uintptr, buf []byte) error
	regionsFn func(opts ScanOptions) ([]process.Page, error)
}

// NewScanner constructs a Scanner attached to p.
func NewScanner(p *process.Process) *Scanner {
	return &Scanner{p: p}
}

// First runs an initial scan for value, replacing any prior results.
// opts overrides s.Options for this call.
func (s *Scanner) First(value any, opts ScanOptions) error {
	encoded, kind, size, err := encodeValue(value)
	if err != nil {
		return err
	}
	s.kind = kind
	s.size = size
	s.target = encoded
	s.opts = opts
	s.results = nil

	pages, err := s.regions(opts)
	if err != nil {
		return err
	}

	step := opts.Step
	if step <= 0 {
		step = size
	}
	chunkSize := opts.ChunkSize
	if chunkSize <= 0 {
		chunkSize = 16 * 1024 * 1024
	}

	for _, page := range pages {
		s.scanRegion(page.Start, page.End, chunkSize, step, encoded, true)
	}

	if s.OnResults != nil {
		out := make([]uintptr, len(s.results))
		copy(out, s.results)
		s.OnResults(len(s.results), out)
	}
	return nil
}

// Next narrows results to addresses currently equal to value. The encoded
// size of value must match the size used in First.
func (s *Scanner) Next(value any) error {
	encoded, kind, size, err := encodeValue(value)
	if err != nil {
		return err
	}
	if size != s.size {
		return fmt.Errorf("Next value size %d does not match initial scan size %d", size, s.size)
	}
	s.kind = kind
	s.target = encoded

	buf := make([]byte, size)
	kept := s.results[:0]
	for _, addr := range s.results {
		if err := s.read(addr, buf); err != nil {
			continue
		}
		if bytes.Equal(buf, encoded) {
			kept = append(kept, addr)
		}
	}
	s.results = kept

	if s.OnResults != nil {
		out := make([]uintptr, len(s.results))
		copy(out, s.results)
		s.OnResults(len(s.results), out)
	}
	return nil
}

// FirstScan runs First using s.Value and s.Options. For hotkey use.
func (s *Scanner) FirstScan() error {
	if s.Value == nil {
		return fmt.Errorf("Scanner.Value is nil — set it before scanning")
	}
	return s.First(s.Value, s.Options)
}

// NextScan runs Next using s.Value. For hotkey use.
func (s *Scanner) NextScan() error {
	if s.Value == nil {
		return fmt.Errorf("Scanner.Value is nil — set it before scanning")
	}
	return s.Next(s.Value)
}

// Results returns a copy of the current address list.
func (s *Scanner) Results() []uintptr {
	out := make([]uintptr, len(s.results))
	copy(out, s.results)
	return out
}

// Count returns the number of remaining addresses.
func (s *Scanner) Count() int { return len(s.results) }

// Reset clears results so a fresh First can run.
func (s *Scanner) Reset() {
	s.results = nil
	s.target = nil
	s.size = 0
}

// BindFirst returns a func() suitable for hotkey.Register. Errors and
// results route through OnError / OnResults.
func (s *Scanner) BindFirst() func() {
	return func() {
		if err := s.FirstScan(); err != nil {
			s.reportErr(err)
		}
	}
}

// BindNext returns a func() suitable for hotkey.Register.
func (s *Scanner) BindNext() func() {
	return func() {
		if err := s.NextScan(); err != nil {
			s.reportErr(err)
		}
	}
}

// BindReset returns a func() suitable for hotkey.Register.
func (s *Scanner) BindReset() func() {
	return func() { s.Reset() }
}

func (s *Scanner) reportErr(err error) {
	if s.OnError != nil {
		s.OnError(err)
		return
	}
	fmt.Fprintln(os.Stderr, "scanner:", err)
}

func (s *Scanner) regions(opts ScanOptions) ([]process.Page, error) {
	if s.regionsFn != nil {
		return s.regionsFn(opts)
	}
	if s.p == nil {
		return nil, fmt.Errorf("scanner has no process and no regionsFn")
	}
	name := opts.Module
	if name == "" {
		name = s.p.Name
	}
	mod, err := s.p.GetModule(name)
	if err != nil {
		return nil, err
	}
	return s.p.EnumMemoryRegions(mod)
}

func (s *Scanner) read(addr uintptr, buf []byte) error {
	if s.reader != nil {
		return s.reader(addr, buf)
	}
	return Read(s.p, addr, buf)
}

func (s *Scanner) scanRegion(start, end uintptr, chunkSize, step int, target []byte, accumulate bool) {
	regionLen := int(end - start)
	if regionLen <= 0 {
		return
	}

	buf := make([]byte, min(regionLen, chunkSize))
	pos := uintptr(0)
	tlen := len(target)

	for pos < uintptr(regionLen) {
		want := min(int(uintptr(regionLen)-pos), len(buf))
		slice := buf[:want]
		if err := s.read(start+pos, slice); err != nil {
			// Region unreadable (guard pages, etc.) — skip it
			return
		}

		// Scan within slice with the configured step. End limit ensures we
		// don't read past the slice when comparing tlen bytes.
		end := want - tlen
		for i := 0; i <= end; i += step {
			if bytes.Equal(slice[i:i+tlen], target) {
				s.results = append(s.results, start+pos+uintptr(i))
			}
		}

		// Advance, leaving a (tlen-1)-byte overlap so a match crossing the
		// chunk boundary is caught next iteration.
		if want < len(buf) {
			break
		}
		pos += uintptr(want - (tlen - 1))
		if want-(tlen-1) <= 0 {
			break
		}
	}
}

func encodeValue(v any) ([]byte, ValueKind, int, error) {
	switch x := v.(type) {
	case []byte:
		if len(x) == 0 {
			return nil, 0, 0, fmt.Errorf("empty []byte value")
		}
		out := make([]byte, len(x))
		copy(out, x)
		return out, KindBytes, len(out), nil
	case int8:
		return []byte{byte(x)}, KindBytes, 1, nil
	case uint8:
		return []byte{x}, KindBytes, 1, nil
	case int16:
		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b, uint16(x))
		return b, KindBytes, 2, nil
	case uint16:
		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b, x)
		return b, KindBytes, 2, nil
	case int32:
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(x))
		return b, KindInt32, 4, nil
	case uint32:
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, x)
		return b, KindUint32, 4, nil
	case int64:
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(x))
		return b, KindInt64, 8, nil
	case uint64:
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, x)
		return b, KindUint64, 8, nil
	case int:
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(x))
		return b, KindInt64, 8, nil
	case float32:
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, math.Float32bits(x))
		return b, KindFloat32, 4, nil
	case float64:
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, math.Float64bits(x))
		return b, KindFloat64, 8, nil
	default:
		return nil, 0, 0, fmt.Errorf("unsupported value type %T", v)
	}
}

