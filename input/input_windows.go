//go:build windows

package input

import (
	"time"
	"unsafe"

	"golang.org/x/sys/windows"

	"gomeow/vec"
)

var (
	user32               = windows.NewLazySystemDLL("user32.dll")
	procGetAsyncKeyState = user32.NewProc("GetAsyncKeyState")
	procGetCursorPos     = user32.NewProc("GetCursorPos")
	procSetCursorPos     = user32.NewProc("SetCursorPos")
	procSendInput        = user32.NewProc("SendInput")
)

const (
	INPUT_MOUSE    = 0
	INPUT_KEYBOARD = 1

	MOUSEEVENTF_MOVE       = 0x0001
	MOUSEEVENTF_LEFTDOWN   = 0x0002
	MOUSEEVENTF_LEFTUP     = 0x0004
	MOUSEEVENTF_RIGHTDOWN  = 0x0008
	MOUSEEVENTF_RIGHTUP    = 0x0010
	MOUSEEVENTF_MIDDLEDOWN = 0x0020
	MOUSEEVENTF_MIDDLEUP   = 0x0040
	MOUSEEVENTF_ABSOLUTE   = 0x8000

	KEYEVENTF_KEYUP = 0x0002

	// Virtual key codes
	VK_LBUTTON = 0x01
	VK_RBUTTON = 0x02
	VK_MBUTTON = 0x04
)

type point struct {
	X int32
	Y int32
}

type mouseInput struct {
	Type uint32
	Mi   struct {
		Dx          int32
		Dy          int32
		MouseData   uint32
		DwFlags     uint32
		Time        uint32
		DwExtraInfo uintptr
	}
}

type keyboardInput struct {
	Type uint32
	Ki   struct {
		WVk         uint16
		WScan       uint16
		DwFlags     uint32
		Time        uint32
		DwExtraInfo uintptr
		_           [8]byte // padding
	}
}

// KeyPressed checks if a key is currently pressed
func KeyPressed(vKey int) bool {
	ret, _, _ := procGetAsyncKeyState.Call(uintptr(vKey))
	return int16(ret) < 0
}

// MousePressed checks if a mouse button is pressed
// button: "left", "right", "middle"
func MousePressed(button string) bool {
	var key int
	switch button {
	case "left":
		key = VK_LBUTTON
	case "right":
		key = VK_RBUTTON
	case "middle":
		key = VK_MBUTTON
	default:
		key = VK_LBUTTON
	}
	return KeyPressed(key)
}

// MousePosition returns the current mouse cursor position
func MousePosition() vec.Vec2 {
	var p point
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&p)))
	return vec.Vec2{X: float32(p.X), Y: float32(p.Y)}
}

// MouseMove moves the mouse cursor to the specified position
// If relative is true, moves relative to current position
func MouseMove(x, y int, relative bool) {
	if relative {
		var p point
		procGetCursorPos.Call(uintptr(unsafe.Pointer(&p)))
		x += int(p.X)
		y += int(p.Y)
	}
	procSetCursorPos.Call(uintptr(x), uintptr(y))
}

// PressKey presses a key (keydown + keyup)
func PressKey(vKey int) {
	var input keyboardInput
	input.Type = INPUT_KEYBOARD
	input.Ki.WVk = uint16(vKey)

	// Key down
	procSendInput.Call(1, uintptr(unsafe.Pointer(&input)), unsafe.Sizeof(input))

	// Key up
	input.Ki.DwFlags = KEYEVENTF_KEYUP
	procSendInput.Call(1, uintptr(unsafe.Pointer(&input)), unsafe.Sizeof(input))
}

// KeyDown presses a key down (without release)
func KeyDown(vKey int) {
	var input keyboardInput
	input.Type = INPUT_KEYBOARD
	input.Ki.WVk = uint16(vKey)
	procSendInput.Call(1, uintptr(unsafe.Pointer(&input)), unsafe.Sizeof(input))
}

// KeyUp releases a key
func KeyUp(vKey int) {
	var input keyboardInput
	input.Type = INPUT_KEYBOARD
	input.Ki.WVk = uint16(vKey)
	input.Ki.DwFlags = KEYEVENTF_KEYUP
	procSendInput.Call(1, uintptr(unsafe.Pointer(&input)), unsafe.Sizeof(input))
}

// MouseDown presses a mouse button down
func MouseDown(button string) {
	var input mouseInput
	input.Type = INPUT_MOUSE

	switch button {
	case "left":
		input.Mi.DwFlags = MOUSEEVENTF_LEFTDOWN
	case "right":
		input.Mi.DwFlags = MOUSEEVENTF_RIGHTDOWN
	case "middle":
		input.Mi.DwFlags = MOUSEEVENTF_MIDDLEDOWN
	default:
		input.Mi.DwFlags = MOUSEEVENTF_LEFTDOWN
	}

	procSendInput.Call(1, uintptr(unsafe.Pointer(&input)), unsafe.Sizeof(input))
}

// MouseUp releases a mouse button
func MouseUp(button string) {
	var input mouseInput
	input.Type = INPUT_MOUSE

	switch button {
	case "left":
		input.Mi.DwFlags = MOUSEEVENTF_LEFTUP
	case "right":
		input.Mi.DwFlags = MOUSEEVENTF_RIGHTUP
	case "middle":
		input.Mi.DwFlags = MOUSEEVENTF_MIDDLEUP
	default:
		input.Mi.DwFlags = MOUSEEVENTF_LEFTUP
	}

	procSendInput.Call(1, uintptr(unsafe.Pointer(&input)), unsafe.Sizeof(input))
}

// MouseClick performs a mouse click (down + small delay + up)
func MouseClick(button string) {
	MouseDown(button)
	time.Sleep(3 * time.Millisecond)
	MouseUp(button)
}

// TypeString types a string by sending key events
func TypeString(text string) {
	for _, char := range text {
		// This is a simplified implementation
		// For full support, would need to handle shift key for uppercase, etc.
		vKey := int(char)
		if char >= 'a' && char <= 'z' {
			vKey = int(char) - 32 // Convert to uppercase VK code
		}
		PressKey(vKey)
		time.Sleep(10 * time.Millisecond)
	}
}
