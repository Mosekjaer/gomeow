package gui

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"gomeow/utils"
)

// ListView draws a scrollable list of items and returns the selected index
// Returns -1 if no item is selected
func ListView(id int, x, y, w, h int, items []string, selected int) int {
	const itemHeight = 28
	const scrollBarWidth = 12

	// Get scroll state
	sps, ok := scrollPanelStates[id]
	if !ok {
		sps = &scrollPanelState{}
		scrollPanelStates[id] = sps
	}

	// Draw background
	rl.DrawRectangle(int32(x), int32(y), int32(w), int32(h), toRLColor(currentTheme.Background))
	rl.DrawRectangleLines(int32(x), int32(y), int32(w), int32(h), toRLColor(currentTheme.Border))

	// Calculate content height
	contentH := len(items) * itemHeight
	maxScroll := float32(contentH - h)
	if maxScroll < 0 {
		maxScroll = 0
	}

	// Handle mouse wheel
	if regionHit(x, y, w, h) {
		wheelMove := rl.GetMouseWheelMove()
		sps.scrollY -= wheelMove * 30
	}

	// Clamp scroll
	if sps.scrollY < 0 {
		sps.scrollY = 0
	}
	if sps.scrollY > maxScroll {
		sps.scrollY = maxScroll
	}

	// Calculate visible range
	startItem := int(sps.scrollY) / itemHeight
	endItem := startItem + (h / itemHeight) + 2
	if endItem > len(items) {
		endItem = len(items)
	}

	// Draw items with clipping
	contentW := w
	if contentH > h {
		contentW = w - scrollBarWidth
	}

	rl.BeginScissorMode(int32(x+1), int32(y+1), int32(contentW-2), int32(h-2))

	newSelected := selected
	for i := startItem; i < endItem; i++ {
		itemY := y + i*itemHeight - int(sps.scrollY)
		itemID := id + 1000 + i

		// Check interaction
		if regionHit(x, itemY, contentW, itemHeight) && itemY >= y && itemY+itemHeight <= y+h {
			state.HotItem = itemID
			if state.ActiveItem == 0 && state.MouseDown {
				state.ActiveItem = itemID
			}
		}

		// Determine item color
		var itemBgColor utils.Color
		if i == selected {
			itemBgColor = currentTheme.Accent
		} else if state.HotItem == itemID {
			itemBgColor = currentTheme.BackgroundHot
		} else {
			itemBgColor = currentTheme.Background
		}

		// Draw item
		rl.DrawRectangle(int32(x+1), int32(itemY), int32(contentW-2), int32(itemHeight), toRLColor(itemBgColor))

		// Draw item text
		textY := itemY + (itemHeight-16)/2
		rl.DrawText(items[i], int32(x+8), int32(textY), 16, toRLColor(currentTheme.Text))

		// Select on click
		if state.HotItem == itemID && state.ActiveItem == itemID && !state.MouseDown {
			newSelected = i
		}
	}

	rl.EndScissorMode()

	// Draw scrollbar if needed
	if contentH > h {
		scrollBarX := x + w - scrollBarWidth
		thumbRatio := float32(h) / float32(contentH)
		thumbH := int(thumbRatio * float32(h))
		if thumbH < 20 {
			thumbH = 20
		}

		scrollRatio := sps.scrollY / maxScroll
		thumbY := y + int(scrollRatio*float32(h-thumbH))

		// Draw track
		rl.DrawRectangle(int32(scrollBarX), int32(y), int32(scrollBarWidth), int32(h), toRLColor(currentTheme.Background))

		// Draw thumb
		thumbID := id + 200000
		if regionHit(scrollBarX, thumbY, scrollBarWidth, thumbH) {
			state.HotItem = thumbID
			if state.ActiveItem == 0 && state.MouseDown {
				state.ActiveItem = thumbID
				sps.dragging = true
				sps.dragStartY = state.MousePos.Y - float32(thumbY)
			}
		}

		if sps.dragging {
			if state.MouseDown {
				newThumbY := state.MousePos.Y - sps.dragStartY - float32(y)
				scrollRatio := newThumbY / float32(h-thumbH)
				sps.scrollY = scrollRatio * maxScroll
				if sps.scrollY < 0 {
					sps.scrollY = 0
				}
				if sps.scrollY > maxScroll {
					sps.scrollY = maxScroll
				}
			} else {
				sps.dragging = false
			}
		}

		thumbColor := currentTheme.Border
		if state.HotItem == thumbID || sps.dragging {
			thumbColor = currentTheme.Accent
		}
		rl.DrawRectangle(int32(scrollBarX+2), int32(thumbY), int32(scrollBarWidth-4), int32(thumbH), toRLColor(thumbColor))
	}

	return newSelected
}

// ListViewMultiSelect draws a list view with multi-selection support
// Returns the updated selection slice
func ListViewMultiSelect(id int, x, y, w, h int, items []string, selected []bool) []bool {
	const itemHeight = 28
	const scrollBarWidth = 12

	// Ensure selected slice is correct size
	if len(selected) != len(items) {
		selected = make([]bool, len(items))
	}

	// Get scroll state
	sps, ok := scrollPanelStates[id]
	if !ok {
		sps = &scrollPanelState{}
		scrollPanelStates[id] = sps
	}

	// Draw background
	rl.DrawRectangle(int32(x), int32(y), int32(w), int32(h), toRLColor(currentTheme.Background))
	rl.DrawRectangleLines(int32(x), int32(y), int32(w), int32(h), toRLColor(currentTheme.Border))

	contentH := len(items) * itemHeight
	maxScroll := float32(contentH - h)
	if maxScroll < 0 {
		maxScroll = 0
	}

	if regionHit(x, y, w, h) {
		wheelMove := rl.GetMouseWheelMove()
		sps.scrollY -= wheelMove * 30
	}

	if sps.scrollY < 0 {
		sps.scrollY = 0
	}
	if sps.scrollY > maxScroll {
		sps.scrollY = maxScroll
	}

	startItem := int(sps.scrollY) / itemHeight
	endItem := startItem + (h / itemHeight) + 2
	if endItem > len(items) {
		endItem = len(items)
	}

	contentW := w
	if contentH > h {
		contentW = w - scrollBarWidth
	}

	rl.BeginScissorMode(int32(x+1), int32(y+1), int32(contentW-2), int32(h-2))

	for i := startItem; i < endItem; i++ {
		itemY := y + i*itemHeight - int(sps.scrollY)
		itemID := id + 1000 + i

		if regionHit(x, itemY, contentW, itemHeight) && itemY >= y && itemY+itemHeight <= y+h {
			state.HotItem = itemID
			if state.ActiveItem == 0 && state.MouseDown {
				state.ActiveItem = itemID
			}
		}

		var itemBgColor utils.Color
		if selected[i] {
			itemBgColor = currentTheme.Accent
		} else if state.HotItem == itemID {
			itemBgColor = currentTheme.BackgroundHot
		} else {
			itemBgColor = currentTheme.Background
		}

		rl.DrawRectangle(int32(x+1), int32(itemY), int32(contentW-2), int32(itemHeight), toRLColor(itemBgColor))

		textY := itemY + (itemHeight-16)/2
		rl.DrawText(items[i], int32(x+8), int32(textY), 16, toRLColor(currentTheme.Text))

		// Toggle selection on click
		if state.HotItem == itemID && state.ActiveItem == itemID && !state.MouseDown {
			selected[i] = !selected[i]
		}
	}

	rl.EndScissorMode()

	// Draw scrollbar if needed
	if contentH > h {
		scrollBarX := x + w - scrollBarWidth
		thumbRatio := float32(h) / float32(contentH)
		thumbH := int(thumbRatio * float32(h))
		if thumbH < 20 {
			thumbH = 20
		}

		scrollRatio := sps.scrollY / maxScroll
		thumbY := y + int(scrollRatio*float32(h-thumbH))

		rl.DrawRectangle(int32(scrollBarX), int32(y), int32(scrollBarWidth), int32(h), toRLColor(currentTheme.Background))

		thumbID := id + 200000
		if regionHit(scrollBarX, thumbY, scrollBarWidth, thumbH) {
			state.HotItem = thumbID
			if state.ActiveItem == 0 && state.MouseDown {
				state.ActiveItem = thumbID
				sps.dragging = true
				sps.dragStartY = state.MousePos.Y - float32(thumbY)
			}
		}

		if sps.dragging {
			if state.MouseDown {
				newThumbY := state.MousePos.Y - sps.dragStartY - float32(y)
				scrollRatio := newThumbY / float32(h-thumbH)
				sps.scrollY = scrollRatio * maxScroll
				if sps.scrollY < 0 {
					sps.scrollY = 0
				}
				if sps.scrollY > maxScroll {
					sps.scrollY = maxScroll
				}
			} else {
				sps.dragging = false
			}
		}

		thumbColor := currentTheme.Border
		if state.HotItem == thumbID || sps.dragging {
			thumbColor = currentTheme.Accent
		}
		rl.DrawRectangle(int32(scrollBarX+2), int32(thumbY), int32(scrollBarWidth-4), int32(thumbH), toRLColor(thumbColor))
	}

	return selected
}
