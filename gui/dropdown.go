package gui

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Dropdown state
type dropdownState struct {
	open bool
}

var dropdownStates = make(map[int]*dropdownState)

// Dropdown draws a dropdown selector and returns the selected index
func Dropdown(id int, x, y, w, h int, items []string, selected int) int {
	const itemHeight = 24
	const arrowSize = 16

	// Get or create dropdown state
	ds, ok := dropdownStates[id]
	if !ok {
		ds = &dropdownState{}
		dropdownStates[id] = ds
	}

	// Check for header interaction
	if regionHit(x, y, w, h) {
		state.HotItem = id
		if state.ActiveItem == 0 && state.MouseDown {
			state.ActiveItem = id
		}
	}

	// Toggle dropdown on click
	if state.HotItem == id && state.ActiveItem == id && !state.MouseDown {
		ds.open = !ds.open
	}

	// Determine header color
	var headerColor = currentTheme.Background
	if state.ActiveItem == id {
		headerColor = currentTheme.BackgroundActive
	} else if state.HotItem == id {
		headerColor = currentTheme.BackgroundHot
	}

	// Draw header
	rl.DrawRectangle(int32(x), int32(y), int32(w), int32(h), toRLColor(headerColor))
	rl.DrawRectangleLines(int32(x), int32(y), int32(w), int32(h), toRLColor(currentTheme.Border))

	// Draw selected item text
	if selected >= 0 && selected < len(items) {
		rl.DrawText(items[selected], int32(x+8), int32(y+(h-16)/2), 16, toRLColor(currentTheme.Text))
	}

	// Draw arrow
	arrowX := x + w - arrowSize - 4
	arrowY := y + (h-arrowSize)/2
	if ds.open {
		// Up arrow
		rl.DrawTriangle(
			rl.Vector2{X: float32(arrowX + arrowSize/2), Y: float32(arrowY + 4)},
			rl.Vector2{X: float32(arrowX + arrowSize - 4), Y: float32(arrowY + arrowSize - 4)},
			rl.Vector2{X: float32(arrowX + 4), Y: float32(arrowY + arrowSize - 4)},
			toRLColor(currentTheme.Text),
		)
	} else {
		// Down arrow
		rl.DrawTriangle(
			rl.Vector2{X: float32(arrowX + 4), Y: float32(arrowY + 4)},
			rl.Vector2{X: float32(arrowX + arrowSize - 4), Y: float32(arrowY + 4)},
			rl.Vector2{X: float32(arrowX + arrowSize/2), Y: float32(arrowY + arrowSize - 4)},
			toRLColor(currentTheme.Text),
		)
	}

	// Draw dropdown list if open
	newSelected := selected
	if ds.open {
		listY := y + h
		listH := len(items) * itemHeight

		// Draw list background
		rl.DrawRectangle(int32(x), int32(listY), int32(w), int32(listH), toRLColor(currentTheme.Background))
		rl.DrawRectangleLines(int32(x), int32(listY), int32(w), int32(listH), toRLColor(currentTheme.Border))

		// Draw items
		for i, item := range items {
			itemY := listY + i*itemHeight
			itemID := id + 1000 + i

			// Check for item interaction
			if regionHit(x, itemY, w, itemHeight) {
				state.HotItem = itemID
				if state.ActiveItem == 0 && state.MouseDown {
					state.ActiveItem = itemID
				}

				// Select item on click
				if state.HotItem == itemID && state.ActiveItem == itemID && !state.MouseDown {
					newSelected = i
					ds.open = false
				}
			}

			// Draw item background
			var itemBgColor = currentTheme.Background
			if i == selected {
				itemBgColor = currentTheme.Accent
			} else if state.HotItem == itemID {
				itemBgColor = currentTheme.BackgroundHot
			}
			rl.DrawRectangle(int32(x+1), int32(itemY), int32(w-2), int32(itemHeight), toRLColor(itemBgColor))

			// Draw item text
			textColor := currentTheme.Text
			rl.DrawText(item, int32(x+8), int32(itemY+(itemHeight-16)/2), 16, toRLColor(textColor))
		}

		// Close dropdown if clicked outside
		if state.MousePressed && !regionHit(x, y, w, h+listH) {
			ds.open = false
		}
	}

	return newSelected
}

// DropdownWithLabel draws a dropdown with a label
func DropdownWithLabel(id int, x, y, w, h int, label string, items []string, selected int) int {
	rl.DrawText(label, int32(x), int32(y-20), 16, toRLColor(currentTheme.Text))
	return Dropdown(id, x, y, w, h, items, selected)
}

// ComboBox is an alias for Dropdown
func ComboBox(id int, x, y, w, h int, items []string, selected int) int {
	return Dropdown(id, x, y, w, h, items, selected)
}
