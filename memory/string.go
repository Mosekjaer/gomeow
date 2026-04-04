package memory

import (
	"encoding/binary"
	"unicode/utf16"

	"gomeow/process"
)

// ReadStringW reads a null-terminated UTF-16LE wide string and returns its
// UTF-8 form. maxChars caps the read length in 16-bit characters (not bytes);
// 0 → 256.
//
// Use this for Win32-style LPCWSTR strings stored by most native Windows
// applications.
func ReadStringW(p *process.Process, address uintptr, maxChars int) (string, error) {
	if maxChars <= 0 {
		maxChars = 256
	}
	buf := make([]byte, maxChars*2)
	if err := Read(p, address, buf); err != nil {
		return "", err
	}
	return decodeUTF16LE(buf), nil
}

// WriteStringW writes value as UTF-16LE with a null terminator.
func WriteStringW(p *process.Process, address uintptr, value string) error {
	return Write(p, address, encodeUTF16LE(value))
}

// decodeUTF16LE decodes raw UTF-16LE bytes up to (but not including) the
// first NUL code unit. Trailing junk past NUL is ignored. An odd trailing
// byte is dropped.
func decodeUTF16LE(buf []byte) string {
	n := len(buf) / 2
	u16 := make([]uint16, 0, n)
	for i := 0; i < n; i++ {
		c := binary.LittleEndian.Uint16(buf[i*2:])
		if c == 0 {
			break
		}
		u16 = append(u16, c)
	}
	return string(utf16.Decode(u16))
}

// encodeUTF16LE returns value as UTF-16LE bytes with a trailing NUL.
func encodeUTF16LE(value string) []byte {
	u16 := utf16.Encode([]rune(value))
	out := make([]byte, len(u16)*2+2)
	for i, c := range u16 {
		binary.LittleEndian.PutUint16(out[i*2:], c)
	}
	// last 2 bytes already zero → NUL terminator
	return out
}
