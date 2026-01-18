//go:build ignore

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gomeow/input"
)

// Windows Virtual Key Codes
const (
	VK_LBUTTON = 0x01
	VK_RBUTTON = 0x02
	VK_MBUTTON = 0x04
	VK_ESCAPE  = 0x1B
	VK_SPACE   = 0x20
	VK_LEFT    = 0x25
	VK_UP      = 0x26
	VK_RIGHT   = 0x27
	VK_DOWN    = 0x28
	VK_A       = 0x41
	VK_D       = 0x44
	VK_S       = 0x53
	VK_W       = 0x57
	VK_F1      = 0x70
	VK_F2      = 0x71
	VK_F3      = 0x72
	VK_F4      = 0x73
)

func main() {
	fmt.Println("=== goMeow Input Detection Example ===")
	fmt.Println()
	fmt.Println("Detecting input... Press ESC to exit")
	fmt.Println()

	// Handle Ctrl+C gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nExiting...")
		os.Exit(0)
	}()

	// Track previous states to detect edges
	var prevMouseLeft, prevMouseRight, prevSpace bool
	var prevW, prevA, prevS, prevD bool

	ticker := time.NewTicker(16 * time.Millisecond) // ~60 FPS
	defer ticker.Stop()

	for range ticker.C {
		// Check for exit
		if input.KeyPressed(VK_ESCAPE) {
			fmt.Println("\nESC pressed, exiting...")
			break
		}

		// Get mouse position
		mousePos := input.MousePosition()

		// Check mouse buttons (detect on press)
		mouseLeft := input.MousePressed("left")
		mouseRight := input.MousePressed("right")

		if mouseLeft && !prevMouseLeft {
			fmt.Printf("[CLICK] Left mouse at (%.0f, %.0f)\n", mousePos.X, mousePos.Y)
		}
		if mouseRight && !prevMouseRight {
			fmt.Printf("[CLICK] Right mouse at (%.0f, %.0f)\n", mousePos.X, mousePos.Y)
		}
		prevMouseLeft = mouseLeft
		prevMouseRight = mouseRight

		// Check WASD keys
		w := input.KeyPressed(VK_W)
		a := input.KeyPressed(VK_A)
		s := input.KeyPressed(VK_S)
		d := input.KeyPressed(VK_D)

		if w && !prevW {
			fmt.Println("[KEY] W pressed (forward)")
		}
		if a && !prevA {
			fmt.Println("[KEY] A pressed (left)")
		}
		if s && !prevS {
			fmt.Println("[KEY] S pressed (backward)")
		}
		if d && !prevD {
			fmt.Println("[KEY] D pressed (right)")
		}
		prevW, prevA, prevS, prevD = w, a, s, d

		// Check space
		space := input.KeyPressed(VK_SPACE)
		if space && !prevSpace {
			fmt.Println("[KEY] SPACE pressed (jump)")
		}
		prevSpace = space

		// Check function keys for actions
		if input.KeyPressed(VK_F1) {
			fmt.Println("[ACTION] F1 - Simulating mouse click...")
			input.MouseClick("left")
			time.Sleep(200 * time.Millisecond) // Debounce
		}

		if input.KeyPressed(VK_F2) {
			fmt.Println("[ACTION] F2 - Moving mouse to center...")
			input.MouseMove(960, 540, false) // Move to center (assuming 1920x1080)
			time.Sleep(200 * time.Millisecond)
		}

		if input.KeyPressed(VK_F3) {
			fmt.Println("[ACTION] F3 - Moving mouse relatively...")
			input.MouseMove(50, 0, true) // Move 50px right
			time.Sleep(200 * time.Millisecond)
		}
	}
}
