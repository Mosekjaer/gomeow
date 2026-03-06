package gui

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// ScrollPanel state
type scrollPanelState struct {
	scrollY     float32
	scrollX     float32
	dragging    bool
	dragStartY  float32
}

var scrollPanelStates = make(map[int]*scrollPanelState)

// BeginScrollPanel starts a scrollable panel area
// Returns the scroll offset that should be applied to content
func BeginScrollPanel(id int, x, y, w, h, contentH int) float32 {
	const scrollBarWidth = 12

	// Get or create state
	sps, ok := scrollPanelStates[id]
	if !ok {
		sps = &scrollPanelState{}
		scrollPanelStates[id] = sps
	}

	// Draw panel background
	rl.DrawRectangle(int32(x), int32(y), int32(w), int32(h), toRLColor(currentTheme.Background))
	rl.DrawRectangleLines(int32(x), int32(y), int32(w), int32(h), toRLColor(currentTheme.Border))

	// Calculate scroll limits
	maxScroll := float32(contentH - h)
	if maxScroll < 0 {
		maxScroll = 0
	}

	// Handle mouse wheel scrolling when hovering over panel
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

	// Draw scrollbar if needed
	if contentH > h {
		scrollBarX := x + w - scrollBarWidth
		scrollBarH := h

		// Calculate thumb size and position
		thumbRatio := float32(h) / float32(contentH)
		thumbH := int(thumbRatio * float32(scrollBarH))
		if thumbH < 20 {
			thumbH = 20
		}

		scrollRatio := sps.scrollY / maxScroll
		thumbY := y + int(scrollRatio*float32(scrollBarH-thumbH))

		// Draw scrollbar track
		rl.DrawRectangle(int32(scrollBarX), int32(y), int32(scrollBarWidth), int32(scrollBarH), toRLColor(currentTheme.Background))

		// Check thumb interaction
		thumbID := id + 100000
		if regionHit(scrollBarX, thumbY, scrollBarWidth, thumbH) {
			state.HotItem = thumbID
			if state.ActiveItem == 0 && state.MouseDown {
				state.ActiveItem = thumbID
				sps.dragging = true
				sps.dragStartY = state.MousePos.Y - float32(thumbY)
			}
		}

		// Handle thumb dragging
		if sps.dragging {
			if state.MouseDown {
				newThumbY := state.MousePos.Y - sps.dragStartY - float32(y)
				scrollRatio := newThumbY / float32(scrollBarH-thumbH)
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

		// Draw thumb
		thumbColor := currentTheme.Border
		if state.HotItem == thumbID || sps.dragging {
			thumbColor = currentTheme.Accent
		}
		rl.DrawRectangle(int32(scrollBarX+2), int32(thumbY), int32(scrollBarWidth-4), int32(thumbH), toRLColor(thumbColor))
	}

	// Begin scissor mode for clipped content
	rl.BeginScissorMode(int32(x+1), int32(y+1), int32(w-scrollBarWidth-2), int32(h-2))

	return sps.scrollY
}

// EndScrollPanel ends the scrollable panel area
func EndScrollPanel() {
	rl.EndScissorMode()
}

// ScrollToTop scrolls a panel to the top
func ScrollToTop(id int) {
	if sps, ok := scrollPanelStates[id]; ok {
		sps.scrollY = 0
	}
}

// ScrollToBottom scrolls a panel to show the bottom content
func ScrollToBottom(id int, h, contentH int) {
	if sps, ok := scrollPanelStates[id]; ok {
		maxScroll := float32(contentH - h)
		if maxScroll < 0 {
			maxScroll = 0
		}
		sps.scrollY = maxScroll
	}
}

// GetScrollOffset returns the current scroll offset for a panel
func GetScrollOffset(id int) float32 {
	if sps, ok := scrollPanelStates[id]; ok {
		return sps.scrollY
	}
	return 0
}

// SetScrollOffset sets the scroll offset for a panel
func SetScrollOffset(id int, offset float32) {
	sps, ok := scrollPanelStates[id]
	if !ok {
		sps = &scrollPanelState{}
		scrollPanelStates[id] = sps
	}
	sps.scrollY = offset
}
