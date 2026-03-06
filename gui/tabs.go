package gui

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// TabBar draws a horizontal tab bar and returns the selected tab index
func TabBar(id int, x, y, w, h int, tabs []string, selected int) int {
	if len(tabs) == 0 {
		return selected
	}

	tabW := w / len(tabs)

	for i, tab := range tabs {
		tabX := x + i*tabW
		tabID := id + i

		// Check for interaction
		if regionHit(tabX, y, tabW, h) {
			state.HotItem = tabID
			if state.ActiveItem == 0 && state.MouseDown {
				state.ActiveItem = tabID
			}
		}

		// Determine tab color
		var bgColor = currentTheme.Background
		if i == selected {
			bgColor = currentTheme.BackgroundActive
		} else if state.HotItem == tabID {
			bgColor = currentTheme.BackgroundHot
		}

		// Draw tab
		rl.DrawRectangle(int32(tabX), int32(y), int32(tabW), int32(h), toRLColor(bgColor))
		rl.DrawRectangleLines(int32(tabX), int32(y), int32(tabW), int32(h), toRLColor(currentTheme.Border))

		// Draw underline for selected tab
		if i == selected {
			rl.DrawRectangle(int32(tabX), int32(y+h-3), int32(tabW), 3, toRLColor(currentTheme.Accent))
		}

		// Draw text centered
		textWidth := rl.MeasureText(tab, 16)
		textX := tabX + (tabW-int(textWidth))/2
		textY := y + (h-16)/2
		rl.DrawText(tab, int32(textX), int32(textY), 16, toRLColor(currentTheme.Text))

		// Select on click
		if state.HotItem == tabID && state.ActiveItem == tabID && !state.MouseDown {
			selected = i
		}
	}

	return selected
}

// TabBarEx draws a tab bar with icons
func TabBarEx(id int, x, y, w, h int, tabs []string, icons []string, selected int) int {
	if len(tabs) == 0 {
		return selected
	}

	tabW := w / len(tabs)
	const iconSpacing = 6

	for i, tab := range tabs {
		tabX := x + i*tabW
		tabID := id + i

		if regionHit(tabX, y, tabW, h) {
			state.HotItem = tabID
			if state.ActiveItem == 0 && state.MouseDown {
				state.ActiveItem = tabID
			}
		}

		var bgColor = currentTheme.Background
		if i == selected {
			bgColor = currentTheme.BackgroundActive
		} else if state.HotItem == tabID {
			bgColor = currentTheme.BackgroundHot
		}

		rl.DrawRectangle(int32(tabX), int32(y), int32(tabW), int32(h), toRLColor(bgColor))
		rl.DrawRectangleLines(int32(tabX), int32(y), int32(tabW), int32(h), toRLColor(currentTheme.Border))

		if i == selected {
			rl.DrawRectangle(int32(tabX), int32(y+h-3), int32(tabW), 3, toRLColor(currentTheme.Accent))
		}

		// Draw icon + text
		var totalWidth int
		if i < len(icons) && icons[i] != "" {
			iconWidth := int(rl.MeasureText(icons[i], 16))
			textWidth := int(rl.MeasureText(tab, 16))
			totalWidth = iconWidth + iconSpacing + textWidth

			startX := tabX + (tabW-totalWidth)/2
			textY := y + (h-16)/2

			rl.DrawText(icons[i], int32(startX), int32(textY), 16, toRLColor(currentTheme.Text))
			rl.DrawText(tab, int32(startX+iconWidth+iconSpacing), int32(textY), 16, toRLColor(currentTheme.Text))
		} else {
			textWidth := rl.MeasureText(tab, 16)
			textX := tabX + (tabW-int(textWidth))/2
			textY := y + (h-16)/2
			rl.DrawText(tab, int32(textX), int32(textY), 16, toRLColor(currentTheme.Text))
		}

		if state.HotItem == tabID && state.ActiveItem == tabID && !state.MouseDown {
			selected = i
		}
	}

	return selected
}

// VerticalTabs draws a vertical tab list
func VerticalTabs(id int, x, y, w, h int, tabs []string, selected int) int {
	if len(tabs) == 0 {
		return selected
	}

	tabH := 32

	for i, tab := range tabs {
		tabY := y + i*tabH
		if tabY+tabH > y+h {
			break // Don't draw tabs outside bounds
		}

		tabID := id + i

		if regionHit(x, tabY, w, tabH) {
			state.HotItem = tabID
			if state.ActiveItem == 0 && state.MouseDown {
				state.ActiveItem = tabID
			}
		}

		var bgColor = currentTheme.Background
		if i == selected {
			bgColor = currentTheme.BackgroundActive
		} else if state.HotItem == tabID {
			bgColor = currentTheme.BackgroundHot
		}

		rl.DrawRectangle(int32(x), int32(tabY), int32(w), int32(tabH), toRLColor(bgColor))
		rl.DrawRectangleLines(int32(x), int32(tabY), int32(w), int32(tabH), toRLColor(currentTheme.Border))

		// Draw left indicator for selected tab
		if i == selected {
			rl.DrawRectangle(int32(x), int32(tabY), 3, int32(tabH), toRLColor(currentTheme.Accent))
		}

		// Draw text
		textY := tabY + (tabH-16)/2
		rl.DrawText(tab, int32(x+8), int32(textY), 16, toRLColor(currentTheme.Text))

		if state.HotItem == tabID && state.ActiveItem == tabID && !state.MouseDown {
			selected = i
		}
	}

	return selected
}
