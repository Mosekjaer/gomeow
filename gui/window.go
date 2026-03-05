package gui

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"gomeow/utils"
)

// Window state for dragging
type windowState struct {
	dragging  bool
	dragOffX  float32
	dragOffY  float32
}

var windowStates = make(map[int]*windowState)

// Window draws a draggable window panel and returns true if the close button was clicked
// Returns new x, y position and whether close was clicked
func Window(id int, x, y, w, h int, title string, closeable bool) (int, int, bool) {
	const titleBarHeight = 24
	const padding = 4

	// Get or create window state
	ws, ok := windowStates[id]
	if !ok {
		ws = &windowState{}
		windowStates[id] = ws
	}

	// Check for title bar interaction (dragging)
	titleBarHit := regionHit(x, y, w-30, titleBarHeight) // Leave space for close button
	if titleBarHit {
		state.HotItem = id
		if state.ActiveItem == 0 && state.MouseDown && !ws.dragging {
			state.ActiveItem = id
			ws.dragging = true
			ws.dragOffX = state.MousePos.X - float32(x)
			ws.dragOffY = state.MousePos.Y - float32(y)
		}
	}

	// Handle dragging
	if ws.dragging {
		if state.MouseDown {
			x = int(state.MousePos.X - ws.dragOffX)
			y = int(state.MousePos.Y - ws.dragOffY)
		} else {
			ws.dragging = false
		}
	}

	// Draw window background
	rl.DrawRectangle(int32(x), int32(y), int32(w), int32(h), toRLColor(currentTheme.Background))
	rl.DrawRectangleLines(int32(x), int32(y), int32(w), int32(h), toRLColor(currentTheme.Border))

	// Draw title bar
	var titleBarColor utils.Color
	if state.HotItem == id || ws.dragging {
		titleBarColor = currentTheme.BackgroundHot
	} else {
		titleBarColor = currentTheme.BackgroundActive
	}
	rl.DrawRectangle(int32(x), int32(y), int32(w), titleBarHeight, toRLColor(titleBarColor))
	rl.DrawLine(int32(x), int32(y+titleBarHeight), int32(x+w), int32(y+titleBarHeight), toRLColor(currentTheme.Border))

	// Draw title
	rl.DrawText(title, int32(x+padding), int32(y+padding), 16, toRLColor(currentTheme.Text))

	// Draw close button if closeable
	closeClicked := false
	if closeable {
		closeX := x + w - 22
		closeY := y + 2
		closeSize := 20

		closeHit := regionHit(closeX, closeY, closeSize, closeSize)
		closeID := id + 10000 // Unique ID for close button

		if closeHit {
			state.HotItem = closeID
			if state.ActiveItem == 0 && state.MouseDown {
				state.ActiveItem = closeID
			}
		}

		var closeBgColor utils.Color
		if state.ActiveItem == closeID {
			closeBgColor = utils.NewColorAlpha(200, 50, 50, 255)
		} else if state.HotItem == closeID {
			closeBgColor = utils.NewColorAlpha(180, 60, 60, 255)
		} else {
			closeBgColor = utils.NewColorAlpha(150, 50, 50, 255)
		}

		rl.DrawRectangle(int32(closeX), int32(closeY), int32(closeSize), int32(closeSize), toRLColor(closeBgColor))
		rl.DrawText("X", int32(closeX+5), int32(closeY+2), 16, toRLColor(utils.White))

		closeClicked = state.HotItem == closeID && state.ActiveItem == closeID && !state.MouseDown
	}

	return x, y, closeClicked
}

// Panel draws a simple panel without title bar
func Panel(x, y, w, h int) {
	rl.DrawRectangle(int32(x), int32(y), int32(w), int32(h), toRLColor(currentTheme.Background))
	rl.DrawRectangleLines(int32(x), int32(y), int32(w), int32(h), toRLColor(currentTheme.Border))
}

// PanelEx draws a panel with custom colors
func PanelEx(x, y, w, h int, bgColor, borderColor utils.Color) {
	rl.DrawRectangle(int32(x), int32(y), int32(w), int32(h), toRLColor(bgColor))
	rl.DrawRectangleLines(int32(x), int32(y), int32(w), int32(h), toRLColor(borderColor))
}

// GroupBox draws a labeled group box
func GroupBox(x, y, w, h int, label string) {
	const labelPadding = 8

	// Draw border with gap for label
	labelWidth := int(rl.MeasureText(label, 14)) + labelPadding*2

	// Top line (with gap)
	rl.DrawLine(int32(x), int32(y), int32(x+labelPadding), int32(y), toRLColor(currentTheme.Border))
	rl.DrawLine(int32(x+labelPadding+labelWidth), int32(y), int32(x+w), int32(y), toRLColor(currentTheme.Border))

	// Other sides
	rl.DrawLine(int32(x+w), int32(y), int32(x+w), int32(y+h), toRLColor(currentTheme.Border))
	rl.DrawLine(int32(x), int32(y+h), int32(x+w), int32(y+h), toRLColor(currentTheme.Border))
	rl.DrawLine(int32(x), int32(y), int32(x), int32(y+h), toRLColor(currentTheme.Border))

	// Draw label
	rl.DrawText(label, int32(x+labelPadding*2), int32(y-7), 14, toRLColor(currentTheme.Text))
}
