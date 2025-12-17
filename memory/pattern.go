package memory

import (
	"fmt"
	"strconv"
	"strings"

	"gomeow/process"
)

const (
	wildCardByte = 0x100 // Value that can't be a real byte
)

// parsePattern converts a hex pattern string to byte values
// Supports wildcards: "48 8B ?? 90" or "48 8B ? 90"
func parsePattern(pattern string) ([]int, error) {
	pattern = strings.TrimSpace(pattern)
	pattern = strings.ReplaceAll(pattern, " ", "")

	if len(pattern)%2 != 0 {
		return nil, fmt.Errorf("invalid pattern length")
	}

	var result []int
	for i := 0; i < len(pattern); i += 2 {
		hex := pattern[i : i+2]
		if hex == "??" || hex == "**" {
			result = append(result, wildCardByte)
		} else if hex[0] == '?' || hex[1] == '?' {
			// Partial wildcard not fully supported, treat as full wildcard
			result = append(result, wildCardByte)
		} else {
			val, err := strconv.ParseUint(hex, 16, 8)
			if err != nil {
				return nil, fmt.Errorf("invalid hex byte: %s", hex)
			}
			result = append(result, int(val))
		}
	}

	return result, nil
}

// AOBScanBytes scans a byte slice for a pattern
// Returns offsets of all matches (or just first if single=true)
func AOBScanBytes(pattern string, data []byte, single bool) ([]uintptr, error) {
	patternBytes, err := parsePattern(pattern)
	if err != nil {
		return nil, err
	}

	if len(patternBytes) == 0 {
		return nil, fmt.Errorf("empty pattern")
	}

	if len(patternBytes) > len(data) {
		return nil, nil
	}

	var results []uintptr

	for i := 0; i < len(data)-len(patternBytes); i++ {
		match := true
		for j, pb := range patternBytes {
			if pb != wildCardByte && int(data[i+j]) != pb {
				match = false
				break
			}
		}
		if match {
			results = append(results, uintptr(i))
			if single {
				return results, nil
			}
		}
	}

	return results, nil
}

// AOBScanModule scans a module for a pattern
func AOBScanModule(p *process.Process, moduleName, pattern string, relative, single bool) ([]uintptr, error) {
	module, err := p.GetModule(moduleName)
	if err != nil {
		return nil, err
	}

	// Read entire module
	data, err := ReadBytes(p, module.Base, int(module.Size))
	if err != nil {
		return nil, fmt.Errorf("failed to read module memory: %v", err)
	}

	results, err := AOBScanBytes(pattern, data, single)
	if err != nil {
		return nil, err
	}

	// Convert relative offsets to absolute addresses
	if !relative {
		for i := range results {
			results[i] += module.Base
		}
	}

	return results, nil
}

// AOBScanRange scans a memory range for a pattern
func AOBScanRange(p *process.Process, pattern string, startAddr, endAddr uintptr, relative, single bool) ([]uintptr, error) {
	if startAddr >= endAddr {
		return nil, fmt.Errorf("invalid range: start >= end")
	}

	size := int(endAddr - startAddr)
	data, err := ReadBytes(p, startAddr, size)
	if err != nil {
		return nil, fmt.Errorf("failed to read memory range: %v", err)
	}

	results, err := AOBScanBytes(pattern, data, single)
	if err != nil {
		return nil, err
	}

	// Convert relative offsets to absolute addresses
	if !relative {
		for i := range results {
			results[i] += startAddr
		}
	}

	return results, nil
}

// AOBScanFirst is a convenience function that returns only the first match
func AOBScanFirst(p *process.Process, moduleName, pattern string) (uintptr, error) {
	results, err := AOBScanModule(p, moduleName, pattern, false, true)
	if err != nil {
		return 0, err
	}
	if len(results) == 0 {
		return 0, fmt.Errorf("pattern not found")
	}
	return results[0], nil
}

// PatternToMask converts a pattern to IDA-style signature and mask
// e.g., "48 8B ?? 90" -> sig=[]byte{0x48, 0x8B, 0x00, 0x90}, mask="xx?x"
func PatternToMask(pattern string) ([]byte, string, error) {
	patternBytes, err := parsePattern(pattern)
	if err != nil {
		return nil, "", err
	}

	sig := make([]byte, len(patternBytes))
	mask := make([]byte, len(patternBytes))

	for i, pb := range patternBytes {
		if pb == wildCardByte {
			sig[i] = 0x00
			mask[i] = '?'
		} else {
			sig[i] = byte(pb)
			mask[i] = 'x'
		}
	}

	return sig, string(mask), nil
}

// ScanWithMask scans using IDA-style signature and mask
func ScanWithMask(data, signature []byte, mask string) ([]uintptr, error) {
	if len(signature) != len(mask) {
		return nil, fmt.Errorf("signature and mask length mismatch")
	}

	if len(signature) > len(data) {
		return nil, nil
	}

	var results []uintptr

	for i := 0; i <= len(data)-len(signature); i++ {
		match := true
		for j := 0; j < len(signature); j++ {
			if mask[j] == 'x' && data[i+j] != signature[j] {
				match = false
				break
			}
		}
		if match {
			results = append(results, uintptr(i))
		}
	}

	return results, nil
}
