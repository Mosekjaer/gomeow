package gui

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Checkbox draws a checkbox with label and returns the new checked state
func Checkbox(id int, x, y int, label string, checked bool) bool {
	const boxSize = 20
	const spacing = 8

	// Check for mouse hover on the box
	if regionHit(x, y, boxSize, boxSize) {
		state.HotItem = id
		if state.ActiveItem == 0 && state.MouseDown {
			state.ActiveItem = id
		}
	}

	// Determine box color
	var bgColor = currentTheme.Background
	if state.HotItem == id {
		bgColor = currentTheme.BackgroundHot
	}

	// Draw checkbox box
	rl.DrawRectangle(int32(x), int32(y), boxSize, boxSize, toRLColor(bgColor))
	rl.DrawRectangleLines(int32(x), int32(y), boxSize, boxSize, toRLColor(currentTheme.Border))

	// Draw checkmark if checked
	if checked {
		// Draw a simple checkmark
		padding := int32(4)
		rl.DrawRectangle(int32(x)+padding, int32(y)+padding,
			boxSize-padding*2, boxSize-padding*2,
			toRLColor(currentTheme.Accent))
	}

	// Draw label
	rl.DrawText(label, int32(x+boxSize+spacing), int32(y+2), 16, toRLColor(currentTheme.Text))

	// Toggle on click
	if state.HotItem == id && state.ActiveItem == id && !state.MouseDown {
		return !checked
	}
	return checked
}

// CheckboxRight draws a checkbox with label on the left
func CheckboxRight(id int, x, y, labelWidth int, label string, checked bool) bool {
	const boxSize = 20
	const spacing = 8

	boxX := x + labelWidth + spacing

	// Check for mouse hover on the box
	if regionHit(boxX, y, boxSize, boxSize) {
		state.HotItem = id
		if state.ActiveItem == 0 && state.MouseDown {
			state.ActiveItem = id
		}
	}

	// Draw label on left
	rl.DrawText(label, int32(x), int32(y+2), 16, toRLColor(currentTheme.Text))

	// Determine box color
	var bgColor = currentTheme.Background
	if state.HotItem == id {
		bgColor = currentTheme.BackgroundHot
	}

	// Draw checkbox box
	rl.DrawRectangle(int32(boxX), int32(y), boxSize, boxSize, toRLColor(bgColor))
	rl.DrawRectangleLines(int32(boxX), int32(y), boxSize, boxSize, toRLColor(currentTheme.Border))

	// Draw checkmark if checked
	if checked {
		padding := int32(4)
		rl.DrawRectangle(int32(boxX)+padding, int32(y)+padding,
			boxSize-padding*2, boxSize-padding*2,
			toRLColor(currentTheme.Accent))
	}

	// Toggle on click
	if state.HotItem == id && state.ActiveItem == id && !state.MouseDown {
		return !checked
	}
	return checked
}
