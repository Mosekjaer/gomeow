package gui

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Tooltip state
type tooltipState struct {
	hoverTime  float64
	lastHotItem int
}

var ttState = tooltipState{}

const tooltipDelay = 0.5 // seconds before tooltip appears

// Tooltip draws a tooltip if the specified item has been hovered long enough
func Tooltip(id int, text string) {
	if state.HotItem != id {
		if ttState.lastHotItem == id {
			ttState.hoverTime = 0
			ttState.lastHotItem = 0
		}
		return
	}

	// Track hover time
	if ttState.lastHotItem != id {
		ttState.lastHotItem = id
		ttState.hoverTime = rl.GetTime()
	}

	// Check if we've hovered long enough
	elapsed := rl.GetTime() - ttState.hoverTime
	if elapsed < tooltipDelay {
		return
	}

	drawTooltip(text, int(state.MousePos.X)+16, int(state.MousePos.Y)+16)
}

// TooltipImmediate draws a tooltip immediately at the current mouse position
func TooltipImmediate(text string) {
	drawTooltip(text, int(state.MousePos.X)+16, int(state.MousePos.Y)+16)
}

// TooltipAt draws a tooltip at a specific position
func TooltipAt(x, y int, text string) {
	drawTooltip(text, x, y)
}

func drawTooltip(text string, x, y int) {
	const padding = 6
	const fontSize = 14

	// Measure text
	textWidth := int(rl.MeasureText(text, fontSize))
	w := textWidth + padding*2
	h := fontSize + padding*2

	// Adjust position to keep tooltip on screen
	screenW := int(rl.GetScreenWidth())
	screenH := int(rl.GetScreenHeight())

	if x+w > screenW {
		x = screenW - w - 4
	}
	if y+h > screenH {
		y = screenH - h - 4
	}
	if x < 0 {
		x = 4
	}
	if y < 0 {
		y = 4
	}

	// Draw tooltip background
	rl.DrawRectangle(int32(x), int32(y), int32(w), int32(h), toRLColor(currentTheme.Background))
	rl.DrawRectangleLines(int32(x), int32(y), int32(w), int32(h), toRLColor(currentTheme.Border))

	// Draw text
	rl.DrawText(text, int32(x+padding), int32(y+padding), fontSize, toRLColor(currentTheme.Text))
}

// TooltipMultiLine draws a multi-line tooltip
func TooltipMultiLine(lines []string) {
	const padding = 6
	const fontSize = 14
	const lineSpacing = 4

	if len(lines) == 0 {
		return
	}

	// Find widest line
	maxWidth := 0
	for _, line := range lines {
		w := int(rl.MeasureText(line, fontSize))
		if w > maxWidth {
			maxWidth = w
		}
	}

	w := maxWidth + padding*2
	h := len(lines)*(fontSize+lineSpacing) - lineSpacing + padding*2

	x := int(state.MousePos.X) + 16
	y := int(state.MousePos.Y) + 16

	// Adjust position to keep tooltip on screen
	screenW := int(rl.GetScreenWidth())
	screenH := int(rl.GetScreenHeight())

	if x+w > screenW {
		x = screenW - w - 4
	}
	if y+h > screenH {
		y = screenH - h - 4
	}

	// Draw background
	rl.DrawRectangle(int32(x), int32(y), int32(w), int32(h), toRLColor(currentTheme.Background))
	rl.DrawRectangleLines(int32(x), int32(y), int32(w), int32(h), toRLColor(currentTheme.Border))

	// Draw lines
	lineY := y + padding
	for _, line := range lines {
		rl.DrawText(line, int32(x+padding), int32(lineY), fontSize, toRLColor(currentTheme.Text))
		lineY += fontSize + lineSpacing
	}
}

// HelpMarker draws a (?) icon that shows a tooltip on hover
func HelpMarker(id int, x, y int, text string) {
	const size = 16

	if regionHit(x, y, size, size) {
		state.HotItem = id
	}

	// Draw circle with question mark
	bgColor := currentTheme.Background
	if state.HotItem == id {
		bgColor = currentTheme.BackgroundHot
	}

	rl.DrawCircle(int32(x+size/2), int32(y+size/2), float32(size/2), toRLColor(bgColor))
	rl.DrawCircleLines(int32(x+size/2), int32(y+size/2), float32(size/2), toRLColor(currentTheme.Border))
	rl.DrawText("?", int32(x+4), int32(y), 14, toRLColor(currentTheme.Text))

	// Show tooltip when hovering
	Tooltip(id, text)
}
