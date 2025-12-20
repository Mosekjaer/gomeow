package memory

import (
	"reflect"
	"testing"
)

func TestParsePattern(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		expected []int
		wantErr  bool
	}{
		{
			name:     "simple pattern",
			pattern:  "48 8B 05",
			expected: []int{0x48, 0x8B, 0x05},
			wantErr:  false,
		},
		{
			name:     "pattern with wildcard",
			pattern:  "48 ?? 05",
			expected: []int{0x48, wildCardByte, 0x05},
			wantErr:  false,
		},
		{
			name:     "pattern without spaces",
			pattern:  "488B05",
			expected: []int{0x48, 0x8B, 0x05},
			wantErr:  false,
		},
		{
			name:     "pattern with multiple wildcards",
			pattern:  "?? ?? 90",
			expected: []int{wildCardByte, wildCardByte, 0x90},
			wantErr:  false,
		},
		{
			name:     "lowercase hex",
			pattern:  "ab cd ef",
			expected: []int{0xAB, 0xCD, 0xEF},
			wantErr:  false,
		},
		{
			name:     "mixed case",
			pattern:  "Ab cD Ef",
			expected: []int{0xAB, 0xCD, 0xEF},
			wantErr:  false,
		},
		{
			name:    "invalid pattern - odd length",
			pattern: "48 8B 0",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parsePattern(tt.pattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePattern() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("parsePattern() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestAOBScanBytes(t *testing.T) {
	data := []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x48, 0x8B, 0x05, 0x90, 0x90, 0x90, 0x90, 0x90,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x48, 0x8B, 0x05, 0xAB, 0xCD, 0xEF, 0x00, 0x00,
	}

	tests := []struct {
		name     string
		pattern  string
		single   bool
		expected []uintptr
	}{
		{
			name:     "exact match single",
			pattern:  "48 8B 05",
			single:   true,
			expected: []uintptr{8},
		},
		{
			name:     "exact match all",
			pattern:  "48 8B 05",
			single:   false,
			expected: []uintptr{8, 24},
		},
		{
			name:     "wildcard match",
			pattern:  "48 ?? 05",
			single:   false,
			expected: []uintptr{8, 24},
		},
		{
			name:     "no match",
			pattern:  "FF FF FF",
			single:   false,
			expected: nil,
		},
		{
			name:     "single byte match",
			pattern:  "90",
			single:   true,
			expected: []uintptr{11},
		},
		{
			name:     "match at start",
			pattern:  "00 01 02",
			single:   true,
			expected: []uintptr{0},
		},
		{
			name:     "match at end",
			pattern:  "EF 00 00",
			single:   true,
			expected: []uintptr{29},
		},
		{
			name:     "all wildcards",
			pattern:  "?? ?? ??",
			single:   true,
			expected: []uintptr{0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := AOBScanBytes(tt.pattern, data, tt.single)
			if err != nil {
				t.Errorf("AOBScanBytes() error = %v", err)
				return
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("AOBScanBytes() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestAOBScanBytesErrors(t *testing.T) {
	data := []byte{0x00, 0x01, 0x02}

	tests := []struct {
		name    string
		pattern string
		wantErr bool
	}{
		{
			name:    "invalid hex",
			pattern: "GG HH",
			wantErr: true,
		},
		{
			name:    "empty pattern",
			pattern: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := AOBScanBytes(tt.pattern, data, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("AOBScanBytes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPatternToMask(t *testing.T) {
	tests := []struct {
		name         string
		pattern      string
		expectedSig  []byte
		expectedMask string
	}{
		{
			name:         "no wildcards",
			pattern:      "48 8B 05",
			expectedSig:  []byte{0x48, 0x8B, 0x05},
			expectedMask: "xxx",
		},
		{
			name:         "with wildcards",
			pattern:      "48 ?? 05",
			expectedSig:  []byte{0x48, 0x00, 0x05},
			expectedMask: "x?x",
		},
		{
			name:         "all wildcards",
			pattern:      "?? ?? ??",
			expectedSig:  []byte{0x00, 0x00, 0x00},
			expectedMask: "???",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sig, mask, err := PatternToMask(tt.pattern)
			if err != nil {
				t.Errorf("PatternToMask() error = %v", err)
				return
			}
			if !reflect.DeepEqual(sig, tt.expectedSig) {
				t.Errorf("PatternToMask() sig = %v, want %v", sig, tt.expectedSig)
			}
			if mask != tt.expectedMask {
				t.Errorf("PatternToMask() mask = %v, want %v", mask, tt.expectedMask)
			}
		})
	}
}

func TestScanWithMask(t *testing.T) {
	data := []byte{
		0x00, 0x01, 0x02, 0x03,
		0x48, 0x8B, 0x05, 0x90,
		0x48, 0xAB, 0x05, 0x90,
	}

	sig := []byte{0x48, 0x00, 0x05}
	mask := "x?x"

	results, err := ScanWithMask(data, sig, mask)
	if err != nil {
		t.Errorf("ScanWithMask() error = %v", err)
	}

	expected := []uintptr{4, 8}
	if !reflect.DeepEqual(results, expected) {
		t.Errorf("ScanWithMask() = %v, want %v", results, expected)
	}
}

func TestScanWithMaskMismatch(t *testing.T) {
	data := []byte{0x00, 0x01, 0x02}
	sig := []byte{0x00, 0x01}
	mask := "xxx" // Mask longer than signature

	_, err := ScanWithMask(data, sig, mask)
	if err == nil {
		t.Error("ScanWithMask() should error on sig/mask length mismatch")
	}
}

// Benchmarks

func BenchmarkAOBScanBytes(b *testing.B) {
	// Create a large data buffer
	data := make([]byte, 1024*1024) // 1MB
	for i := range data {
		data[i] = byte(i % 256)
	}
	// Place pattern near the end
	copy(data[len(data)-100:], []byte{0x48, 0x8B, 0x05, 0x90})

	pattern := "48 8B 05 90"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = AOBScanBytes(pattern, data, true)
	}
}

func BenchmarkAOBScanBytesWithWildcard(b *testing.B) {
	data := make([]byte, 1024*1024)
	for i := range data {
		data[i] = byte(i % 256)
	}
	copy(data[len(data)-100:], []byte{0x48, 0x8B, 0x05, 0x90})

	pattern := "48 ?? 05 90"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = AOBScanBytes(pattern, data, true)
	}
}

func BenchmarkParsePattern(b *testing.B) {
	pattern := "48 8B ?? 90 ?? ?? 05 AB CD EF"
	for i := 0; i < b.N; i++ {
		_, _ = parsePattern(pattern)
	}
}
