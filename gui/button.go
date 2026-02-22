package gui

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"gomeow/utils"
)

// Button draws a clickable button and returns true if clicked
func Button(id int, x, y, w, h int, text string) bool {
	// Check for mouse hover
	if regionHit(x, y, w, h) {
		state.HotItem = id
		if state.ActiveItem == 0 && state.MouseDown {
			state.ActiveItem = id
		}
	}

	// Determine button state and color
	var bgColor utils.Color
	if state.ActiveItem == id {
		bgColor = currentTheme.BackgroundActive
	} else if state.HotItem == id {
		bgColor = currentTheme.BackgroundHot
	} else {
		bgColor = currentTheme.Background
	}

	// Draw button background
	rl.DrawRectangle(int32(x), int32(y), int32(w), int32(h), toRLColor(bgColor))
	rl.DrawRectangleLines(int32(x), int32(y), int32(w), int32(h), toRLColor(currentTheme.Border))

	// Draw text centered
	textWidth := rl.MeasureText(text, 16)
	textX := x + (w-int(textWidth))/2
	textY := y + (h-16)/2
	rl.DrawText(text, int32(textX), int32(textY), 16, toRLColor(currentTheme.Text))

	// Return true if button was clicked
	return state.HotItem == id && state.ActiveItem == id && !state.MouseDown
}

// ButtonEx draws a button with custom colors
func ButtonEx(id int, x, y, w, h int, text string, bgColor, textColor utils.Color) bool {
	if regionHit(x, y, w, h) {
		state.HotItem = id
		if state.ActiveItem == 0 && state.MouseDown {
			state.ActiveItem = id
		}
	}

	// Brighten on hover/active
	if state.ActiveItem == id {
		bgColor = brighten(bgColor, 40)
	} else if state.HotItem == id {
		bgColor = brighten(bgColor, 20)
	}

	rl.DrawRectangle(int32(x), int32(y), int32(w), int32(h), toRLColor(bgColor))
	rl.DrawRectangleLines(int32(x), int32(y), int32(w), int32(h), toRLColor(currentTheme.Border))

	textWidth := rl.MeasureText(text, 16)
	textX := x + (w-int(textWidth))/2
	textY := y + (h-16)/2
	rl.DrawText(text, int32(textX), int32(textY), 16, toRLColor(textColor))

	return state.HotItem == id && state.ActiveItem == id && !state.MouseDown
}

// IconButton draws a button with just an icon/symbol
func IconButton(id int, x, y, size int, icon string) bool {
	return Button(id, x, y, size, size, icon)
}

func brighten(c utils.Color, amount uint8) utils.Color {
	r := int(c.R) + int(amount)
	g := int(c.G) + int(amount)
	b := int(c.B) + int(amount)
	if r > 255 {
		r = 255
	}
	if g > 255 {
		g = 255
	}
	if b > 255 {
		b = 255
	}
	return utils.Color{R: uint8(r), G: uint8(g), B: uint8(b), A: c.A}
}
