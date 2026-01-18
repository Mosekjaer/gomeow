//go:build linux

package input

/*
#cgo LDFLAGS: -lX11 -lXtst

#include <X11/Xlib.h>
#include <X11/keysym.h>
#include <X11/extensions/XTest.h>
#include <stdlib.h>

static Display* display = NULL;
static Window root;

void init_display() {
    if (display == NULL) {
        display = XOpenDisplay(NULL);
        if (display != NULL) {
            root = XRootWindow(display, 0);
        }
    }
}

int key_pressed(int keysym) {
    init_display();
    if (display == NULL) return 0;

    char keys[32];
    XQueryKeymap(display, keys);
    KeyCode keycode = XKeysymToKeycode(display, keysym);
    return (keys[keycode / 8] & (1 << (keycode % 8))) != 0;
}

int mouse_pressed(int button) {
    init_display();
    if (display == NULL) return 0;

    Window qRoot, qChild;
    int qRootX, qRootY, qChildX, qChildY;
    unsigned int qMask;

    XQueryPointer(display, root, &qRoot, &qChild, &qRootX, &qRootY, &qChildX, &qChildY, &qMask);

    unsigned int buttonMask[] = {Button1Mask, Button2Mask, Button3Mask};
    if (button >= 0 && button < 3) {
        return (qMask & buttonMask[button]) != 0;
    }
    return 0;
}

void get_mouse_position(int* x, int* y) {
    init_display();
    if (display == NULL) {
        *x = 0;
        *y = 0;
        return;
    }

    Window qRoot, qChild;
    int qRootX, qRootY, qChildX, qChildY;
    unsigned int qMask;

    XQueryPointer(display, root, &qRoot, &qChild, &qRootX, &qRootY, &qChildX, &qChildY, &qMask);
    *x = qRootX;
    *y = qRootY;
}

void press_key(int keysym, int hold) {
    init_display();
    if (display == NULL) return;

    KeyCode keycode = XKeysymToKeycode(display, keysym);
    XTestFakeKeyEvent(display, keycode, True, CurrentTime);
    if (!hold) {
        XTestFakeKeyEvent(display, keycode, False, CurrentTime);
    }
    XFlush(display);
}

void release_key(int keysym) {
    init_display();
    if (display == NULL) return;

    KeyCode keycode = XKeysymToKeycode(display, keysym);
    XTestFakeKeyEvent(display, keycode, False, CurrentTime);
    XFlush(display);
}

void mouse_move(int x, int y, int relative) {
    init_display();
    if (display == NULL) return;

    if (relative) {
        XTestFakeRelativeMotionEvent(display, x, y, CurrentTime);
    } else {
        XTestFakeMotionEvent(display, -1, x, y, CurrentTime);
    }
    XFlush(display);
}

void mouse_button(int button, int down) {
    init_display();
    if (display == NULL) return;

    XTestFakeButtonEvent(display, button, down, 0);
    XFlush(display);
}
*/
import "C"
import (
	"time"

	"gomeow/vec"
)

// KeyPressed checks if a key is currently pressed (X11 keysym)
func KeyPressed(keysym int) bool {
	return C.key_pressed(C.int(keysym)) != 0
}

// MousePressed checks if a mouse button is pressed
// button: "left", "right", "middle"
func MousePressed(button string) bool {
	var btn int
	switch button {
	case "left":
		btn = 0
	case "middle":
		btn = 1
	case "right":
		btn = 2
	default:
		btn = 0
	}
	return C.mouse_pressed(C.int(btn)) != 0
}

// MousePosition returns the current mouse cursor position
func MousePosition() vec.Vec2 {
	var x, y C.int
	C.get_mouse_position(&x, &y)
	return vec.Vec2{X: float32(x), Y: float32(y)}
}

// MouseMove moves the mouse cursor
func MouseMove(x, y int, relative bool) {
	rel := 0
	if relative {
		rel = 1
	}
	C.mouse_move(C.int(x), C.int(y), C.int(rel))
}

// PressKey presses and releases a key
func PressKey(keysym int) {
	C.press_key(C.int(keysym), 0)
}

// KeyDown presses a key down (without release)
func KeyDown(keysym int) {
	C.press_key(C.int(keysym), 1)
}

// KeyUp releases a key
func KeyUp(keysym int) {
	C.release_key(C.int(keysym))
}

// MouseDown presses a mouse button down
func MouseDown(button string) {
	var btn int
	switch button {
	case "left":
		btn = 1
	case "middle":
		btn = 2
	case "right":
		btn = 3
	default:
		btn = 1
	}
	C.mouse_button(C.int(btn), 1)
}

// MouseUp releases a mouse button
func MouseUp(button string) {
	var btn int
	switch button {
	case "left":
		btn = 1
	case "middle":
		btn = 2
	case "right":
		btn = 3
	default:
		btn = 1
	}
	C.mouse_button(C.int(btn), 0)
}

// MouseClick performs a mouse click
func MouseClick(button string) {
	MouseDown(button)
	time.Sleep(3 * time.Millisecond)
	MouseUp(button)
}

// TypeString types a string (stub - would need keysym mapping)
func TypeString(text string) {
	// Would need proper character to keysym mapping
	// This is a placeholder
}
