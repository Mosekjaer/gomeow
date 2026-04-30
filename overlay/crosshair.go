package overlay

import "gomeow/utils"

// CrosshairStyle selects the shape DrawCrosshair renders.
type CrosshairStyle int

const (
	// CrosshairCross is the classic + shape with a center gap.
	CrosshairCross CrosshairStyle = iota

	// CrosshairTShape is a + with the top arm omitted (FPS-style).
	CrosshairTShape

	// CrosshairDot is a single filled dot.
	CrosshairDot

	// CrosshairCircle is a hollow circle outline.
	CrosshairCircle

	// CrosshairCircleDot is a hollow circle with a center dot.
	CrosshairCircleDot
)

// DrawCrosshair renders a crosshair at (x, y) using the existing draw helpers.
//
// size is the half-extent of the cross arms (or the radius for circle styles).
// thickness is the line width; 0 falls back to 1.
// gap is the empty pixel band between the center and where each arm starts;
// pass 0 for no gap.
//
// Call this between BeginDrawing / EndDrawing inside the overlay loop.
func DrawCrosshair(style CrosshairStyle, x, y, size, gap int, thickness float32, color utils.Color) {
	if thickness <= 0 {
		thickness = 1
	}
	if gap < 0 {
		gap = 0
	}

	switch style {
	case CrosshairCross:
		drawCrossArms(x, y, size, gap, thickness, color, true)
	case CrosshairTShape:
		drawCrossArms(x, y, size, gap, thickness, color, false)
	case CrosshairDot:
		r := float32(size) / 2
		if r < 1 {
			r = 1
		}
		DrawCircle(x, y, r, color)
	case CrosshairCircle:
		DrawCircleLines(x, y, float32(size), color)
	case CrosshairCircleDot:
		DrawCircleLines(x, y, float32(size), color)
		DrawCircle(x, y, float32(thickness)+1, color)
	}
}

// drawCrossArms draws the four (or three, if includeTop=false) arms.
func drawCrossArms(x, y, size, gap int, thickness float32, color utils.Color, includeTop bool) {
	// Horizontal arms
	DrawLine(x-size, y, x-gap, y, color, thickness)
	DrawLine(x+gap, y, x+size, y, color, thickness)
	// Vertical arms
	if includeTop {
		DrawLine(x, y-size, x, y-gap, color, thickness)
	}
	DrawLine(x, y+gap, x, y+size, color, thickness)
}
