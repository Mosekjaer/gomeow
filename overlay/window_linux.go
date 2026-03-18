//go:build linux

package overlay

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"gomeow/input"
)

// WindowInfo holds window position and size
type WindowInfo struct {
	X      int
	Y      int
	Width  int
	Height int
}

// Virtual key codes — kept identical to window_windows.go so cross-platform
// callers (e.g. examples) compile against the same constants. On Linux these
// are translated to X11 keysyms via vkToKeysym before polling.
const (
	VK_END    = 0x23
	VK_HOME   = 0x24
	VK_INSERT = 0x2D
	VK_DELETE = 0x2E
	VK_F1     = 0x70
	VK_F2     = 0x71
	VK_F3     = 0x72
	VK_F4     = 0x73
	VK_F5     = 0x74
	VK_F6     = 0x75
	VK_F7     = 0x76
	VK_F8     = 0x77
	VK_F9     = 0x78
	VK_F10    = 0x79
	VK_F11    = 0x7A
	VK_F12    = 0x7B
)

// vkToKeysym maps Windows-style virtual key codes to X11 keysyms.
// Values from /usr/include/X11/keysymdef.h.
var vkToKeysym = map[int]int{
	VK_END:    0xFF57, // XK_End
	VK_HOME:   0xFF50, // XK_Home
	VK_INSERT: 0xFF63, // XK_Insert
	VK_DELETE: 0xFFFF, // XK_Delete
	VK_F1:     0xFFBE,
	VK_F2:     0xFFBF,
	VK_F3:     0xFFC0,
	VK_F4:     0xFFC1,
	VK_F5:     0xFFC2,
	VK_F6:     0xFFC3,
	VK_F7:     0xFFC4,
	VK_F8:     0xFFC5,
	VK_F9:     0xFFC6,
	VK_F10:    0xFFC7,
	VK_F11:    0xFFC8,
	VK_F12:    0xFFC9,
}

var (
	overlayHwnd uintptr
	keyStates   = map[int]bool{}
)

// findWindowByTitle is a stub on Linux — the X11 backend uses xwininfo by name
// and has no exposed handle concept. Returning (0, nil) lets overlay.go's
// optional handle bookkeeping no-op cleanly.
func findWindowByTitle(name string) (uintptr, error) {
	return 0, nil
}

// SetOverlayHandle stores the overlay window handle for focus checking.
// No-op on Linux (no handle to track).
func SetOverlayHandle(hwnd uintptr) {
	overlayHwnd = hwnd
}

// GetOverlayHandle returns the stored overlay window handle.
func GetOverlayHandle() uintptr {
	return overlayHwnd
}

// IsWindowFocused returns whether the target window or our overlay is focused.
// On Linux there is no focus check yet — always returns true.
func IsWindowFocused(name string) bool {
	return true
}

// IsTargetWindowFocused returns whether ONLY the target window is focused.
// On Linux there is no focus check yet — always returns true.
func IsTargetWindowFocused(name string) bool {
	return true
}

// IsKeyDownGlobal checks if a Windows-style virtual key is currently held,
// translated to the matching X11 keysym. Unknown VK codes return false.
func IsKeyDownGlobal(vkCode int) bool {
	keysym, ok := vkToKeysym[vkCode]
	if !ok {
		return false
	}
	return input.KeyPressed(keysym)
}

// IsKeyPressedGlobal returns true on the rising edge of a key (down transition).
// Mirrors the Windows implementation in window_windows.go.
func IsKeyPressedGlobal(vkCode int) bool {
	isDown := IsKeyDownGlobal(vkCode)
	wasDown := keyStates[vkCode]
	keyStates[vkCode] = isDown
	return isDown && !wasDown
}

// GetWindowInfo returns information about a window by name using xwininfo
func GetWindowInfo(name string) (*WindowInfo, error) {
	cmd := exec.Command("xwininfo", "-name", name)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("xwininfo failed (is it installed?): %v", err)
	}

	info := &WindowInfo{}
	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "error") {
			return nil, fmt.Errorf("window '%s' not found", name)
		}

		if strings.Contains(line, "te upper-left X:") {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				info.X, _ = strconv.Atoi(parts[len(parts)-1])
			}
		} else if strings.Contains(line, "te upper-left Y:") {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				info.Y, _ = strconv.Atoi(parts[len(parts)-1])
			}
		} else if strings.Contains(line, "Width:") {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				info.Width, _ = strconv.Atoi(parts[len(parts)-1])
			}
		} else if strings.Contains(line, "Height:") {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				info.Height, _ = strconv.Atoi(parts[len(parts)-1])
			}
		}
	}

	return info, nil
}
