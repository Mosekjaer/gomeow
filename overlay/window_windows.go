//go:build windows

package overlay

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32                = windows.NewLazySystemDLL("user32.dll")
	procFindWindowW       = user32.NewProc("FindWindowW")
	procFindWindowExW     = user32.NewProc("FindWindowExW")
	procGetClientRect     = user32.NewProc("GetClientRect")
	procGetWindowRect     = user32.NewProc("GetWindowRect")
	procGetWindowInfo     = user32.NewProc("GetWindowInfo")
	procGetWindowTextW    = user32.NewProc("GetWindowTextW")
	procGetWindowTextLengthW = user32.NewProc("GetWindowTextLengthW")
	procEnumWindows       = user32.NewProc("EnumWindows")
	procIsWindowVisible   = user32.NewProc("IsWindowVisible")
	procClientToScreen    = user32.NewProc("ClientToScreen")
)

type rect struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

type point struct {
	X int32
	Y int32
}

type windowInfo struct {
	CbSize          uint32
	RcWindow        rect
	RcClient        rect
	DwStyle         uint32
	DwExStyle       uint32
	DwWindowStatus  uint32
	CxWindowBorders uint32
	CyWindowBorders uint32
	AtomWindowType  uint16
	WCreatorVersion uint16
}

// WindowInfo holds window position and size
type WindowInfo struct {
	X      int
	Y      int
	Width  int
	Height int
}

// foundWindow is used during window enumeration
var foundWindow uintptr
var searchName string

// enumWindowsCallback is called for each window during enumeration
func enumWindowsCallback(hwnd uintptr, lParam uintptr) uintptr {
	// Check if window is visible
	visible, _, _ := procIsWindowVisible.Call(hwnd)
	if visible == 0 {
		return 1 // Continue enumeration
	}

	// Get window title length
	length, _, _ := procGetWindowTextLengthW.Call(hwnd)
	if length == 0 {
		return 1 // Continue enumeration
	}

	// Get window title
	buf := make([]uint16, length+1)
	procGetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), length+1)
	title := syscall.UTF16ToString(buf)

	// Check if title contains our search string (case-insensitive)
	if strings.Contains(strings.ToLower(title), strings.ToLower(searchName)) {
		foundWindow = hwnd
		return 0 // Stop enumeration
	}

	return 1 // Continue enumeration
}

// findWindowByTitle searches for a window by partial title match
func findWindowByTitle(name string) (uintptr, error) {
	// First try exact match with FindWindowW
	namePtr, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return 0, err
	}

	hwnd, _, _ := procFindWindowW.Call(0, uintptr(unsafe.Pointer(namePtr)))
	if hwnd != 0 {
		return hwnd, nil
	}

	// If exact match fails, enumerate all windows and search by partial title
	foundWindow = 0
	searchName = name

	callback := syscall.NewCallback(enumWindowsCallback)
	procEnumWindows.Call(callback, 0)

	if foundWindow == 0 {
		return 0, fmt.Errorf("window '%s' not found", name)
	}

	return foundWindow, nil
}

var (
	procGetForegroundWindow = user32.NewProc("GetForegroundWindow")
	procGetAsyncKeyState    = user32.NewProc("GetAsyncKeyState")
)

// Key state tracking for detecting key press (not just held)
var keyStates = make(map[int]bool)

// IsKeyPressed checks if a key was just pressed (global, works even when window not focused)
// Uses GetAsyncKeyState to check key state globally
func IsKeyPressedGlobal(vkCode int) bool {
	state, _, _ := procGetAsyncKeyState.Call(uintptr(vkCode))

	// High bit set means key is currently down
	isDown := (state & 0x8000) != 0

	// Check if this is a new press (wasn't down before)
	wasDown := keyStates[vkCode]
	keyStates[vkCode] = isDown

	// Return true only on the transition from not-pressed to pressed
	return isDown && !wasDown
}

// IsKeyDownGlobal checks if a key is currently held down (global)
func IsKeyDownGlobal(vkCode int) bool {
	state, _, _ := procGetAsyncKeyState.Call(uintptr(vkCode))
	return (state & 0x8000) != 0
}

// Common Virtual Key codes
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

// Stores our overlay window handle
var overlayHwnd uintptr

// SetOverlayHandle stores the overlay window handle for focus checking
func SetOverlayHandle(hwnd uintptr) {
	overlayHwnd = hwnd
}

// GetOverlayHandle returns the stored overlay window handle
func GetOverlayHandle() uintptr {
	return overlayHwnd
}

// IsWindowFocused checks if the target window OR our overlay is currently focused
func IsWindowFocused(name string) bool {
	// Get the currently focused window
	foreground, _, _ := procGetForegroundWindow.Call()
	if foreground == 0 {
		return false
	}

	// If our overlay is focused, that's fine too
	if overlayHwnd != 0 && foreground == overlayHwnd {
		return true
	}

	// Find our target window
	targetHwnd, err := findWindowByTitle(name)
	if err != nil {
		return false
	}

	return foreground == targetHwnd
}

// IsTargetWindowFocused checks if ONLY the target window is focused (not overlay)
func IsTargetWindowFocused(name string) bool {
	foreground, _, _ := procGetForegroundWindow.Call()
	if foreground == 0 {
		return false
	}

	targetHwnd, err := findWindowByTitle(name)
	if err != nil {
		return false
	}

	return foreground == targetHwnd
}

// GetWindowInfo returns information about a window by name
func GetWindowInfo(name string) (*WindowInfo, error) {
	hwnd, err := findWindowByTitle(name)
	if err != nil {
		return nil, err
	}

	// Get client rect (relative coordinates, gives us width/height)
	var clientRect rect
	ret, _, _ := procGetClientRect.Call(hwnd, uintptr(unsafe.Pointer(&clientRect)))
	if ret == 0 {
		return nil, fmt.Errorf("failed to get client rect for '%s'", name)
	}

	width := int(clientRect.Right - clientRect.Left)
	height := int(clientRect.Bottom - clientRect.Top)

	// Validate dimensions
	if width <= 0 || height <= 0 {
		// Window might be minimized, try getting window rect instead
		var windowRect rect
		procGetWindowRect.Call(hwnd, uintptr(unsafe.Pointer(&windowRect)))
		width = int(windowRect.Right - windowRect.Left)
		height = int(windowRect.Bottom - windowRect.Top)

		if width <= 0 || height <= 0 {
			return nil, fmt.Errorf("window '%s' has invalid dimensions (minimized?)", name)
		}

		return &WindowInfo{
			X:      int(windowRect.Left),
			Y:      int(windowRect.Top),
			Width:  width,
			Height: height,
		}, nil
	}

	// Convert client top-left corner to screen coordinates
	var topLeft point
	procClientToScreen.Call(hwnd, uintptr(unsafe.Pointer(&topLeft)))

	return &WindowInfo{
		X:      int(topLeft.X),
		Y:      int(topLeft.Y),
		Width:  width,
		Height: height,
	}, nil
}
