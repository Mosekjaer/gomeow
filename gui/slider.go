package gui

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"gomeow/utils"
)

// Slider draws a horizontal slider and returns the new value
func Slider(id int, x, y, w, h int, value, min, max float32) float32 {
	// Calculate thumb position
	thumbW := 10
	range_ := max - min
	normalizedValue := (value - min) / range_
	thumbX := x + int(normalizedValue*float32(w-thumbW))

	// Check for interaction with the entire slider area
	if regionHit(x, y, w, h) {
		state.HotItem = id
		if state.ActiveItem == 0 && state.MouseDown {
			state.ActiveItem = id
		}
	}

	// Draw slider track
	rl.DrawRectangle(int32(x), int32(y+h/2-2), int32(w), 4, toRLColor(currentTheme.Background))
	rl.DrawRectangleLines(int32(x), int32(y+h/2-2), int32(w), 4, toRLColor(currentTheme.Border))

	// Draw filled portion
	fillW := int(normalizedValue * float32(w))
	rl.DrawRectangle(int32(x), int32(y+h/2-2), int32(fillW), 4, toRLColor(currentTheme.SliderFill))

	// Draw thumb
	var thumbColor utils.Color
	if state.ActiveItem == id {
		thumbColor = currentTheme.BackgroundActive
	} else if state.HotItem == id {
		thumbColor = currentTheme.BackgroundHot
	} else {
		thumbColor = currentTheme.Background
	}

	rl.DrawRectangle(int32(thumbX), int32(y), int32(thumbW), int32(h), toRLColor(thumbColor))
	rl.DrawRectangleLines(int32(thumbX), int32(y), int32(thumbW), int32(h), toRLColor(currentTheme.Border))

	// Handle dragging
	if state.ActiveItem == id {
		mouseX := state.MousePos.X
		newValue := min + (mouseX-float32(x))/float32(w)*range_
		return newValue
	}

	return value
}

// SliderInt draws a horizontal slider for integer values
func SliderInt(id int, x, y, w, h int, value, min, max int) int {
	result := Slider(id, x, y, w, h, float32(value), float32(min), float32(max))
	return int(result + 0.5) // Round to nearest int
}

// SliderWithLabel draws a slider with a label and value display
func SliderWithLabel(id int, x, y, w, h int, label string, value, min, max float32) float32 {
	// Draw label
	rl.DrawText(label, int32(x), int32(y-20), 16, toRLColor(currentTheme.Text))

	// Draw value
	valueStr := fmt.Sprintf("%.2f", value)
	valueWidth := rl.MeasureText(valueStr, 16)
	rl.DrawText(valueStr, int32(x+w)-valueWidth, int32(y-20), 16, toRLColor(currentTheme.Text))

	return Slider(id, x, y, w, h, value, min, max)
}

// VerticalSlider draws a vertical slider
func VerticalSlider(id int, x, y, w, h int, value, min, max float32) float32 {
	thumbH := 10
	range_ := max - min
	normalizedValue := (value - min) / range_
	thumbY := y + h - int(normalizedValue*float32(h-thumbH)) - thumbH

	if regionHit(x, y, w, h) {
		state.HotItem = id
		if state.ActiveItem == 0 && state.MouseDown {
			state.ActiveItem = id
		}
	}

	// Draw track
	rl.DrawRectangle(int32(x+w/2-2), int32(y), 4, int32(h), toRLColor(currentTheme.Background))
	rl.DrawRectangleLines(int32(x+w/2-2), int32(y), 4, int32(h), toRLColor(currentTheme.Border))

	// Draw filled portion (from bottom)
	fillH := int(normalizedValue * float32(h))
	rl.DrawRectangle(int32(x+w/2-2), int32(y+h-fillH), 4, int32(fillH), toRLColor(currentTheme.SliderFill))

	// Draw thumb
	var thumbColor utils.Color
	if state.ActiveItem == id {
		thumbColor = currentTheme.BackgroundActive
	} else if state.HotItem == id {
		thumbColor = currentTheme.BackgroundHot
	} else {
		thumbColor = currentTheme.Background
	}

	rl.DrawRectangle(int32(x), int32(thumbY), int32(w), int32(thumbH), toRLColor(thumbColor))
	rl.DrawRectangleLines(int32(x), int32(thumbY), int32(w), int32(thumbH), toRLColor(currentTheme.Border))

	// Handle dragging
	if state.ActiveItem == id {
		mouseY := state.MousePos.Y
		newValue := max - (mouseY-float32(y))/float32(h)*range_
		return newValue
	}

	return value
}
