package gui

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"

	"gomeow/utils"
)

// ProgressBar draws a progress bar
func ProgressBar(x, y, w, h int, value, min, max float32) {
	// Normalize value
	normalized := (value - min) / (max - min)
	if normalized < 0 {
		normalized = 0
	}
	if normalized > 1 {
		normalized = 1
	}

	fillW := int(normalized * float32(w))

	// Draw background
	rl.DrawRectangle(int32(x), int32(y), int32(w), int32(h), toRLColor(currentTheme.Background))
	rl.DrawRectangleLines(int32(x), int32(y), int32(w), int32(h), toRLColor(currentTheme.Border))

	// Draw fill
	if fillW > 0 {
		rl.DrawRectangle(int32(x), int32(y), int32(fillW), int32(h), toRLColor(currentTheme.Accent))
	}
}

// ProgressBarEx draws a progress bar with custom colors
func ProgressBarEx(x, y, w, h int, value, min, max float32, bgColor, fillColor utils.Color) {
	normalized := (value - min) / (max - min)
	if normalized < 0 {
		normalized = 0
	}
	if normalized > 1 {
		normalized = 1
	}

	fillW := int(normalized * float32(w))

	rl.DrawRectangle(int32(x), int32(y), int32(w), int32(h), toRLColor(bgColor))
	rl.DrawRectangleLines(int32(x), int32(y), int32(w), int32(h), toRLColor(currentTheme.Border))

	if fillW > 0 {
		rl.DrawRectangle(int32(x), int32(y), int32(fillW), int32(h), toRLColor(fillColor))
	}
}

// ProgressBarWithLabel draws a progress bar with percentage text
func ProgressBarWithLabel(x, y, w, h int, value, min, max float32) {
	ProgressBar(x, y, w, h, value, min, max)

	// Calculate percentage
	normalized := (value - min) / (max - min)
	if normalized < 0 {
		normalized = 0
	}
	if normalized > 1 {
		normalized = 1
	}

	// Draw percentage text
	percentStr := fmt.Sprintf("%.0f%%", normalized*100)
	textWidth := rl.MeasureText(percentStr, 14)
	textX := x + (w-int(textWidth))/2
	textY := y + (h-14)/2
	rl.DrawText(percentStr, int32(textX), int32(textY), 14, toRLColor(currentTheme.Text))
}

// VerticalProgressBar draws a vertical progress bar
func VerticalProgressBar(x, y, w, h int, value, min, max float32) {
	normalized := (value - min) / (max - min)
	if normalized < 0 {
		normalized = 0
	}
	if normalized > 1 {
		normalized = 1
	}

	fillH := int(normalized * float32(h))

	// Draw background
	rl.DrawRectangle(int32(x), int32(y), int32(w), int32(h), toRLColor(currentTheme.Background))
	rl.DrawRectangleLines(int32(x), int32(y), int32(w), int32(h), toRLColor(currentTheme.Border))

	// Draw fill from bottom
	if fillH > 0 {
		rl.DrawRectangle(int32(x), int32(y+h-fillH), int32(w), int32(fillH), toRLColor(currentTheme.Accent))
	}
}

// HealthBar draws a health-style progress bar with gradient coloring
func HealthBar(x, y, w, h int, health, maxHealth float32) {
	ratio := health / maxHealth
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}

	fillW := int(ratio * float32(w))

	// Draw background
	rl.DrawRectangle(int32(x), int32(y), int32(w), int32(h), toRLColor(currentTheme.Background))
	rl.DrawRectangleLines(int32(x), int32(y), int32(w), int32(h), toRLColor(currentTheme.Border))

	// Calculate color based on health percentage
	var fillColor utils.Color
	if ratio > 0.6 {
		fillColor = utils.Green
	} else if ratio > 0.3 {
		fillColor = utils.Yellow
	} else {
		fillColor = utils.Red
	}

	// Draw fill
	if fillW > 0 {
		rl.DrawRectangle(int32(x), int32(y), int32(fillW), int32(h), toRLColor(fillColor))
	}
}

// HealthBarGradient draws a health bar with smooth color gradient
func HealthBarGradient(x, y, w, h int, health, maxHealth float32) {
	ratio := health / maxHealth
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}

	fillW := int(ratio * float32(w))

	// Draw background
	rl.DrawRectangle(int32(x), int32(y), int32(w), int32(h), toRLColor(currentTheme.Background))
	rl.DrawRectangleLines(int32(x), int32(y), int32(w), int32(h), toRLColor(currentTheme.Border))

	// Interpolate color from red to green based on health
	r := uint8(255 * (1 - ratio))
	g := uint8(255 * ratio)
	fillColor := utils.Color{R: r, G: g, B: 0, A: 255}

	// Draw fill
	if fillW > 0 {
		rl.DrawRectangle(int32(x), int32(y), int32(fillW), int32(h), toRLColor(fillColor))
	}
}

// LoadingSpinner draws an animated loading spinner
// Call this every frame for animation
func LoadingSpinner(x, y, radius int) {
	segments := 12
	angleStep := 360.0 / float32(segments)
	time := float32(rl.GetTime())

	for i := 0; i < segments; i++ {
		angle := float32(i)*angleStep + time*360
		alpha := uint8(255 * (float32(i) / float32(segments)))

		// Calculate segment position
		rad := angle * 3.14159 / 180
		segX := float32(x) + float32(radius)*float32(math.Cos(float64(rad)))*0.7
		segY := float32(y) + float32(radius)*float32(math.Sin(float64(rad)))*0.7

		color := utils.Color{
			R: currentTheme.Accent.R,
			G: currentTheme.Accent.G,
			B: currentTheme.Accent.B,
			A: alpha,
		}

		rl.DrawCircle(int32(segX), int32(segY), float32(radius)/6, toRLColor(color))
	}
}

// IndeterminateProgressBar draws an animated indeterminate progress bar
func IndeterminateProgressBar(x, y, w, h int) {
	// Draw background
	rl.DrawRectangle(int32(x), int32(y), int32(w), int32(h), toRLColor(currentTheme.Background))
	rl.DrawRectangleLines(int32(x), int32(y), int32(w), int32(h), toRLColor(currentTheme.Border))

	// Animate a moving segment
	time := float32(rl.GetTime())
	segmentW := w / 4
	pos := int((time * 100)) % (w + segmentW) // Loop position

	// Draw segment with clipping
	rl.BeginScissorMode(int32(x), int32(y), int32(w), int32(h))
	rl.DrawRectangle(int32(x+pos-segmentW), int32(y), int32(segmentW), int32(h), toRLColor(currentTheme.Accent))
	rl.EndScissorMode()
}
