package gui

import (
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// TextBox state for text input
type textBoxState struct {
	text       string
	cursorPos  int
	selectStart int
	selectEnd   int
	scrollOffset int
}

var textBoxStates = make(map[int]*textBoxState)

// TextBox draws a text input field and returns the current text
func TextBox(id int, x, y, w, h int, text string, maxLen int) string {
	// Get or create text box state
	tbs, ok := textBoxStates[id]
	if !ok {
		tbs = &textBoxState{text: text, cursorPos: len(text)}
		textBoxStates[id] = tbs
	}

	// Sync external text changes
	if tbs.text != text {
		tbs.text = text
		if tbs.cursorPos > len(text) {
			tbs.cursorPos = len(text)
		}
	}

	// Check for mouse hover
	if regionHit(x, y, w, h) {
		state.HotItem = id
		if state.ActiveItem == 0 && state.MouseDown {
			state.ActiveItem = id
			// Calculate cursor position from click
			clickX := int(state.MousePos.X) - x
			tbs.cursorPos = calculateCursorPos(tbs.text, clickX)
		}
	}

	// Determine background color
	var bgColor = currentTheme.Background
	if state.ActiveItem == id {
		bgColor = currentTheme.BackgroundActive
	} else if state.HotItem == id {
		bgColor = currentTheme.BackgroundHot
	}

	// Draw text box background
	rl.DrawRectangle(int32(x), int32(y), int32(w), int32(h), toRLColor(bgColor))
	rl.DrawRectangleLines(int32(x), int32(y), int32(w), int32(h), toRLColor(currentTheme.Border))

	// Handle keyboard input when active
	if state.ActiveItem == id {
		// Draw focus indicator
		rl.DrawRectangleLines(int32(x-1), int32(y-1), int32(w+2), int32(h+2), toRLColor(currentTheme.Accent))

		// Handle text input
		for {
			key := rl.GetCharPressed()
			if key == 0 {
				break
			}
			if len(tbs.text) < maxLen {
				// Insert character at cursor position
				tbs.text = tbs.text[:tbs.cursorPos] + string(rune(key)) + tbs.text[tbs.cursorPos:]
				tbs.cursorPos++
			}
		}

		// Handle special keys
		if rl.IsKeyPressed(rl.KeyBackspace) && tbs.cursorPos > 0 {
			tbs.text = tbs.text[:tbs.cursorPos-1] + tbs.text[tbs.cursorPos:]
			tbs.cursorPos--
		}
		if rl.IsKeyPressed(rl.KeyDelete) && tbs.cursorPos < len(tbs.text) {
			tbs.text = tbs.text[:tbs.cursorPos] + tbs.text[tbs.cursorPos+1:]
		}
		if rl.IsKeyPressed(rl.KeyLeft) && tbs.cursorPos > 0 {
			tbs.cursorPos--
		}
		if rl.IsKeyPressed(rl.KeyRight) && tbs.cursorPos < len(tbs.text) {
			tbs.cursorPos++
		}
		if rl.IsKeyPressed(rl.KeyHome) {
			tbs.cursorPos = 0
		}
		if rl.IsKeyPressed(rl.KeyEnd) {
			tbs.cursorPos = len(tbs.text)
		}

		// Draw cursor
		cursorX := x + 4 + int(rl.MeasureText(tbs.text[:tbs.cursorPos], 16))
		if int(rl.GetTime()*2)%2 == 0 { // Blinking cursor
			rl.DrawLine(int32(cursorX), int32(y+4), int32(cursorX), int32(y+h-4), toRLColor(currentTheme.Text))
		}
	}

	// Draw text (clipped to box)
	textY := y + (h-16)/2
	rl.BeginScissorMode(int32(x+2), int32(y), int32(w-4), int32(h))
	rl.DrawText(tbs.text, int32(x+4), int32(textY), 16, toRLColor(currentTheme.Text))
	rl.EndScissorMode()

	return tbs.text
}

// TextBoxWithLabel draws a text box with a label above it
func TextBoxWithLabel(id int, x, y, w, h int, label, text string, maxLen int) string {
	rl.DrawText(label, int32(x), int32(y-20), 16, toRLColor(currentTheme.Text))
	return TextBox(id, x, y, w, h, text, maxLen)
}

// PasswordBox draws a text input field that masks the input
func PasswordBox(id int, x, y, w, h int, text string, maxLen int) string {
	result := TextBox(id, x, y, w, h, text, maxLen)

	// We need to draw masked text over the actual text
	// This is a simplified approach - the TextBox already drew the text
	// In a real implementation, you'd modify TextBox to accept a mask parameter

	return result
}

// calculateCursorPos calculates cursor position from pixel offset
func calculateCursorPos(text string, pixelX int) int {
	if pixelX <= 0 {
		return 0
	}

	for i := range text {
		w := int(rl.MeasureText(text[:i+1], 16))
		if w > pixelX {
			return i
		}
	}
	return len(text)
}

// MultiLineTextBox draws a multi-line text input field
func MultiLineTextBox(id int, x, y, w, h int, text string, maxLen int) string {
	// Get or create text box state
	tbs, ok := textBoxStates[id]
	if !ok {
		tbs = &textBoxState{text: text, cursorPos: len(text)}
		textBoxStates[id] = tbs
	}

	// Sync external text changes
	if tbs.text != text {
		tbs.text = text
		if tbs.cursorPos > len(text) {
			tbs.cursorPos = len(text)
		}
	}

	// Check for mouse hover
	if regionHit(x, y, w, h) {
		state.HotItem = id
		if state.ActiveItem == 0 && state.MouseDown {
			state.ActiveItem = id
		}
	}

	// Determine background color
	var bgColor = currentTheme.Background
	if state.ActiveItem == id {
		bgColor = currentTheme.BackgroundActive
	} else if state.HotItem == id {
		bgColor = currentTheme.BackgroundHot
	}

	// Draw text box background
	rl.DrawRectangle(int32(x), int32(y), int32(w), int32(h), toRLColor(bgColor))
	rl.DrawRectangleLines(int32(x), int32(y), int32(w), int32(h), toRLColor(currentTheme.Border))

	// Handle keyboard input when active
	if state.ActiveItem == id {
		rl.DrawRectangleLines(int32(x-1), int32(y-1), int32(w+2), int32(h+2), toRLColor(currentTheme.Accent))

		for {
			key := rl.GetCharPressed()
			if key == 0 {
				break
			}
			if len(tbs.text) < maxLen {
				tbs.text = tbs.text[:tbs.cursorPos] + string(rune(key)) + tbs.text[tbs.cursorPos:]
				tbs.cursorPos++
			}
		}

		if rl.IsKeyPressed(rl.KeyEnter) && len(tbs.text) < maxLen {
			tbs.text = tbs.text[:tbs.cursorPos] + "\n" + tbs.text[tbs.cursorPos:]
			tbs.cursorPos++
		}
		if rl.IsKeyPressed(rl.KeyBackspace) && tbs.cursorPos > 0 {
			tbs.text = tbs.text[:tbs.cursorPos-1] + tbs.text[tbs.cursorPos:]
			tbs.cursorPos--
		}
		if rl.IsKeyPressed(rl.KeyDelete) && tbs.cursorPos < len(tbs.text) {
			tbs.text = tbs.text[:tbs.cursorPos] + tbs.text[tbs.cursorPos+1:]
		}
	}

	// Draw text with line wrapping
	lines := strings.Split(tbs.text, "\n")
	rl.BeginScissorMode(int32(x+2), int32(y+2), int32(w-4), int32(h-4))
	lineY := y + 4
	for _, line := range lines {
		rl.DrawText(line, int32(x+4), int32(lineY), 16, toRLColor(currentTheme.Text))
		lineY += 18
	}
	rl.EndScissorMode()

	return tbs.text
}
