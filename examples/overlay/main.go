// Overlay Drawing Example
// This example demonstrates how to create an overlay on top of a specific window
// and draw various shapes on it.
//
// Features:
// - Only shows overlay when target window is focused
// - Press F2 to toggle menu (global hotkey - works even when game is focused)
// - Press END to exit (global hotkey)
//
// Change the target window name below to match your application.
//
// Run with: go run examples/overlay/main.go

package main

import (
	"fmt"
	"math"

	"gomeow/gui"
	"gomeow/overlay"
	"gomeow/utils"
)

// Menu settings
type Settings struct {
	DrawBorder    bool
	DrawCorners   bool
	DrawCrosshair bool
	DrawCircles   bool
	DrawWave      bool
	DrawBoxes     bool
	CircleCount   int
	CrosshairSize int
	BorderColor   utils.Color
	CircleSpeed   float32
}

func main() {
	// Target window name - change this to your target application
	targetWindow := "Counter-Strike 2"

	fmt.Printf("Looking for window: %s\n", targetWindow)

	// Initialize overlay on top of the target window
	// Using -1 for exitKey since we'll handle it with global hotkeys
	err := overlay.Init(targetWindow, 60, "Drawing Overlay", -1, true)
	if err != nil {
		fmt.Printf("Failed to find window '%s': %v\n", targetWindow, err)
		fmt.Println("Make sure the target application is running!")
		return
	}

	fmt.Println("Overlay initialized!")
	fmt.Println("Press F2 to toggle menu")
	fmt.Println("Press END to exit")

	// Settings with defaults
	settings := Settings{
		DrawBorder:    true,
		DrawCorners:   true,
		DrawCrosshair: true,
		DrawCircles:   true,
		DrawWave:      true,
		DrawBoxes:     true,
		CircleCount:   5,
		CrosshairSize: 20,
		BorderColor:   utils.Cyan,
		CircleSpeed:   1.0,
	}

	// Menu state
	showMenu := false
	menuX, menuY := 50, 50

	// Animation variables
	var time float32 = 0

	// Main loop
	running := true
	for running && overlay.Loop() {
		// Check for global hotkeys (works even when game has focus)

		// END key to exit
		if overlay.IsKeyPressedGlobal(overlay.VK_END) {
			fmt.Println("END pressed - exiting...")
			running = false
			continue
		}

		// F2 to toggle menu
		if overlay.IsKeyPressedGlobal(overlay.VK_F2) {
			showMenu = !showMenu
			overlay.ToggleMouse()
			fmt.Printf("Menu toggled: %v\n", showMenu)
		}

		// Check if target window is focused
		if !overlay.IsWindowFocused(targetWindow) {
			// Window not focused - just do a minimal update and skip drawing
			overlay.BeginDrawing()
			// Draw nothing (transparent)
			overlay.EndDrawing()
			continue
		}

		overlay.BeginDrawing()

		time += 0.016 * settings.CircleSpeed

		// Get overlay dimensions
		w := overlay.GetScreenWidth()
		h := overlay.GetScreenHeight()
		centerX := w / 2
		centerY := h / 2

		// Draw elements based on settings
		if settings.DrawBorder {
			overlay.DrawRectangleLines(5, 5, w-10, h-10, settings.BorderColor, 2)
		}

		if settings.DrawCorners {
			drawCornerMarkers(w, h)
		}

		if settings.DrawCircles {
			drawAnimatedCircles(centerX, centerY, time, settings.CircleCount)
		}

		if settings.DrawCrosshair {
			drawCrosshair(centerX, centerY, settings.CrosshairSize)
		}

		if settings.DrawWave {
			drawSineWave(w, h, time)
		}

		if settings.DrawBoxes {
			drawESPBoxes(w, h)
		}

		// Draw info text (always visible)
		overlay.DrawText("Overlay Active", 20, 20, 20, utils.White)
		overlay.DrawText(fmt.Sprintf("Window: %s", targetWindow), 20, 45, 16, utils.Gray)
		overlay.DrawText("F2: Toggle Menu | END: Exit", 20, h-30, 14, utils.Yellow)
		overlay.DrawFPS(w-100, 20)

		// Draw menu if visible
		if showMenu {
			gui.Begin()
			drawMenu(&menuX, &menuY, &settings, &showMenu)
			gui.End()
		}

		overlay.EndDrawing()
	}

	overlay.Close()
	fmt.Println("Overlay closed")
}

func drawMenu(menuX, menuY *int, settings *Settings, showMenu *bool) {
	var closeClicked bool
	*menuX, *menuY, closeClicked = gui.Window(1, *menuX, *menuY, 280, 380, "Overlay Settings", true)
	if closeClicked {
		*showMenu = false
		overlay.ToggleMouse() // Re-enable mouse passthrough
	}

	x := *menuX + 15
	y := *menuY + 40

	// Checkboxes for features
	settings.DrawBorder = gui.Checkbox(10, x, y, "Draw Border", settings.DrawBorder)
	y += 28

	settings.DrawCorners = gui.Checkbox(11, x, y, "Draw Corners", settings.DrawCorners)
	y += 28

	settings.DrawCrosshair = gui.Checkbox(12, x, y, "Draw Crosshair", settings.DrawCrosshair)
	y += 28

	settings.DrawCircles = gui.Checkbox(13, x, y, "Draw Circles", settings.DrawCircles)
	y += 28

	settings.DrawWave = gui.Checkbox(14, x, y, "Draw Wave", settings.DrawWave)
	y += 28

	settings.DrawBoxes = gui.Checkbox(15, x, y, "Draw ESP Boxes", settings.DrawBoxes)
	y += 35

	// Sliders
	overlay.DrawText("Circle Count:", x, y, 14, utils.White)
	y += 20
	settings.CircleCount = gui.SliderInt(20, x, y, 200, 20, settings.CircleCount, 1, 10)
	overlay.DrawText(fmt.Sprintf("%d", settings.CircleCount), x+210, y, 14, utils.White)
	y += 35

	overlay.DrawText("Crosshair Size:", x, y, 14, utils.White)
	y += 20
	settings.CrosshairSize = gui.SliderInt(21, x, y, 200, 20, settings.CrosshairSize, 5, 50)
	overlay.DrawText(fmt.Sprintf("%d", settings.CrosshairSize), x+210, y, 14, utils.White)
	y += 35

	overlay.DrawText("Animation Speed:", x, y, 14, utils.White)
	y += 20
	settings.CircleSpeed = gui.Slider(22, x, y, 200, 20, settings.CircleSpeed, 0.1, 3.0)
	overlay.DrawText(fmt.Sprintf("%.1f", settings.CircleSpeed), x+210, y, 14, utils.White)
	y += 35

	// Color picker
	overlay.DrawText("Border Color:", x, y, 14, utils.White)
	y += 20
	settings.BorderColor = gui.ColorPickerSimple(30, x, y, settings.BorderColor)
}

func drawCornerMarkers(w, h int) {
	cornerSize := 30
	// Top-left
	overlay.DrawLine(10, 10, 10+cornerSize, 10, utils.Red, 3)
	overlay.DrawLine(10, 10, 10, 10+cornerSize, utils.Red, 3)
	// Top-right
	overlay.DrawLine(w-10, 10, w-10-cornerSize, 10, utils.Red, 3)
	overlay.DrawLine(w-10, 10, w-10, 10+cornerSize, utils.Red, 3)
	// Bottom-left
	overlay.DrawLine(10, h-10, 10+cornerSize, h-10, utils.Red, 3)
	overlay.DrawLine(10, h-10, 10, h-10-cornerSize, utils.Red, 3)
	// Bottom-right
	overlay.DrawLine(w-10, h-10, w-10-cornerSize, h-10, utils.Red, 3)
	overlay.DrawLine(w-10, h-10, w-10, h-10-cornerSize, utils.Red, 3)
}

func drawAnimatedCircles(centerX, centerY int, time float32, count int) {
	for i := 0; i < count; i++ {
		angle := float64(time) + float64(i)*math.Pi*2/float64(count)
		radius := 80.0
		cx := float32(centerX) + float32(math.Cos(angle)*radius)
		cy := float32(centerY) + float32(math.Sin(angle)*radius)

		hue := float32(i) / float32(count)
		color := hsvToColor(hue, 1.0, 1.0)

		overlay.DrawCircle(int(cx), int(cy), 15, color)
		overlay.DrawCircleLines(int(cx), int(cy), 20, utils.White)
	}
}

func drawCrosshair(centerX, centerY, size int) {
	overlay.DrawLine(centerX-size, centerY, centerX+size, centerY, utils.Green, 1)
	overlay.DrawLine(centerX, centerY-size, centerX, centerY+size, utils.Green, 1)
}

func drawSineWave(w, h int, time float32) {
	waveY := h - 100
	prevX, prevY := 0, 0
	for x := 20; x < w-20; x += 5 {
		y := waveY + int(math.Sin(float64(x)*0.02+float64(time)*2)*30)
		if prevX != 0 {
			overlay.DrawLine(prevX, prevY, x, y, utils.Magenta, 2)
		}
		prevX, prevY = x, y
	}
}

func drawESPBoxes(w, h int) {
	boxPositions := []struct{ x, y int }{
		{100, 150},
		{w - 200, 150},
		{100, h - 250},
		{w - 200, h - 250},
	}

	for i, pos := range boxPositions {
		boxW, boxH := 80, 120

		if i%2 == 0 {
			overlay.DrawBox(pos.x, pos.y, boxW, boxH,
				utils.NewColorAlpha(0, 255, 0, 30),
				utils.Green, 1)
		} else {
			overlay.DrawCornerBox(pos.x, pos.y, boxW, boxH, 15, utils.Red, 2)
		}

		health := 0.3 + float32(i)*0.2
		healthColor := utils.Green
		if health < 0.5 {
			healthColor = utils.Yellow
		}
		if health < 0.3 {
			healthColor = utils.Red
		}

		barX := pos.x - 8
		barH := int(float32(boxH) * health)
		overlay.DrawRectangle(barX, pos.y, 4, boxH, utils.NewColorAlpha(0, 0, 0, 150))
		overlay.DrawRectangle(barX, pos.y+boxH-barH, 4, barH, healthColor)

		overlay.DrawText(fmt.Sprintf("Target %d", i+1), pos.x, pos.y-20, 14, utils.White)
	}
}

func hsvToColor(h, s, v float32) utils.Color {
	h = h * 360
	c := v * s
	x := c * (1 - float32(math.Abs(math.Mod(float64(h/60), 2)-1)))
	m := v - c

	var r, g, b float32
	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}

	return utils.Color{
		R: uint8((r + m) * 255),
		G: uint8((g + m) * 255),
		B: uint8((b + m) * 255),
		A: 255,
	}
}
