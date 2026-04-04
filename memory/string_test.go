package memory

import "testing"

func TestUTF16Roundtrip(t *testing.T) {
	cases := []string{
		"",
		"hello",
		"blåbær",       // Danish multi-byte
		"日本語",          // multi-byte CJK
		"emoji 🦊 fox", // BMP-outside (surrogate pair)
	}
	for _, in := range cases {
		out := decodeUTF16LE(encodeUTF16LE(in))
		if out != in {
			t.Errorf("roundtrip mismatch:\n  in:  %q\n  out: %q", in, out)
		}
	}
}

func TestUTF16NullTerminator(t *testing.T) {
	encoded := encodeUTF16LE("hi")
	// Append garbage past the terminator — decode should ignore it.
	garbage := append(encoded, 0x41, 0x00, 0x42, 0x00)
	if got := decodeUTF16LE(garbage); got != "hi" {
		t.Fatalf("expected %q, got %q", "hi", got)
	}
}

func TestUTF16OddTrailingByte(t *testing.T) {
	encoded := encodeUTF16LE("ok")
	// Drop the last byte — odd length, decoder should not panic.
	truncated := encoded[:len(encoded)-1]
	_ = decodeUTF16LE(truncated)
}

func TestEncodeUTF16LayoutLE(t *testing.T) {
	// 'A' = 0x41 → bytes 0x41 0x00 in LE, then NUL.
	got := encodeUTF16LE("A")
	want := []byte{0x41, 0x00, 0x00, 0x00}
	if len(got) != len(want) {
		t.Fatalf("len: got %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("byte %d: got 0x%02X, want 0x%02X", i, got[i], want[i])
		}
	}
}
