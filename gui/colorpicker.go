package gui

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"

	"gomeow/utils"
)

// ColorPicker draws a color picker and returns the selected color
func ColorPicker(id int, x, y, w, h int, color utils.Color) utils.Color {
	const hueBarWidth = 20
	const padding = 8
	const alphaBarHeight = 20

	// Calculate component sizes
	svW := w - hueBarWidth - padding
	svH := h - alphaBarHeight - padding

	// Convert color to HSV for manipulation
	hue, sat, val := rgbToHsv(color.R, color.G, color.B)

	// --- Saturation/Value square ---
	svID := id
	if regionHit(x, y, svW, svH) {
		state.HotItem = svID
		if state.ActiveItem == 0 && state.MouseDown {
			state.ActiveItem = svID
		}
	}

	// Draw SV gradient
	for py := 0; py < svH; py++ {
		for px := 0; px < svW; px++ {
			s := float64(px) / float64(svW-1)
			v := 1.0 - float64(py)/float64(svH-1)
			r, g, b := hsvToRgb(hue, s, v)
			rl.DrawPixel(int32(x+px), int32(y+py), rl.Color{R: r, G: g, B: b, A: 255})
		}
	}
	rl.DrawRectangleLines(int32(x), int32(y), int32(svW), int32(svH), toRLColor(currentTheme.Border))

	// Handle SV interaction
	if state.ActiveItem == svID {
		sat = float64(state.MousePos.X-float32(x)) / float64(svW)
		val = 1.0 - float64(state.MousePos.Y-float32(y))/float64(svH)
		sat = clamp(sat, 0, 1)
		val = clamp(val, 0, 1)
	}

	// Draw SV cursor
	cursorX := x + int(sat*float64(svW))
	cursorY := y + int((1-val)*float64(svH))
	rl.DrawCircleLines(int32(cursorX), int32(cursorY), 5, rl.White)
	rl.DrawCircleLines(int32(cursorX), int32(cursorY), 4, rl.Black)

	// --- Hue bar ---
	hueX := x + svW + padding
	hueID := id + 1
	if regionHit(hueX, y, hueBarWidth, svH) {
		state.HotItem = hueID
		if state.ActiveItem == 0 && state.MouseDown {
			state.ActiveItem = hueID
		}
	}

	// Draw hue gradient
	for py := 0; py < svH; py++ {
		h := float64(py) / float64(svH-1) * 360
		r, g, b := hsvToRgb(h, 1, 1)
		rl.DrawRectangle(int32(hueX), int32(y+py), int32(hueBarWidth), 1, rl.Color{R: r, G: g, B: b, A: 255})
	}
	rl.DrawRectangleLines(int32(hueX), int32(y), int32(hueBarWidth), int32(svH), toRLColor(currentTheme.Border))

	// Handle hue interaction
	if state.ActiveItem == hueID {
		hue = float64(state.MousePos.Y-float32(y)) / float64(svH) * 360
		hue = clamp(hue, 0, 360)
	}

	// Draw hue cursor
	hueCursorY := y + int(hue/360*float64(svH))
	rl.DrawRectangle(int32(hueX-2), int32(hueCursorY-2), int32(hueBarWidth+4), 4, rl.White)
	rl.DrawRectangleLines(int32(hueX-2), int32(hueCursorY-2), int32(hueBarWidth+4), 4, rl.Black)

	// --- Alpha bar ---
	alphaY := y + svH + padding
	alphaID := id + 2
	if regionHit(x, alphaY, w, alphaBarHeight) {
		state.HotItem = alphaID
		if state.ActiveItem == 0 && state.MouseDown {
			state.ActiveItem = alphaID
		}
	}

	// Draw checkerboard background for alpha
	checkerSize := 8
	for py := 0; py < alphaBarHeight; py += checkerSize {
		for px := 0; px < w; px += checkerSize {
			c := rl.LightGray
			if ((px/checkerSize)+(py/checkerSize))%2 == 0 {
				c = rl.White
			}
			rl.DrawRectangle(int32(x+px), int32(alphaY+py), int32(checkerSize), int32(checkerSize), c)
		}
	}

	// Draw alpha gradient
	r, g, b := hsvToRgb(hue, sat, val)
	for px := 0; px < w; px++ {
		a := uint8(float64(px) / float64(w-1) * 255)
		rl.DrawRectangle(int32(x+px), int32(alphaY), 1, int32(alphaBarHeight), rl.Color{R: r, G: g, B: b, A: a})
	}
	rl.DrawRectangleLines(int32(x), int32(alphaY), int32(w), int32(alphaBarHeight), toRLColor(currentTheme.Border))

	// Handle alpha interaction
	alpha := color.A
	if state.ActiveItem == alphaID {
		alpha = uint8(clamp(float64(state.MousePos.X-float32(x))/float64(w)*255, 0, 255))
	}

	// Draw alpha cursor
	alphaCursorX := x + int(float64(alpha)/255*float64(w))
	rl.DrawRectangle(int32(alphaCursorX-2), int32(alphaY-2), 4, int32(alphaBarHeight+4), rl.White)
	rl.DrawRectangleLines(int32(alphaCursorX-2), int32(alphaY-2), 4, int32(alphaBarHeight+4), rl.Black)

	// Convert back to RGB
	finalR, finalG, finalB := hsvToRgb(hue, sat, val)
	return utils.Color{R: finalR, G: finalG, B: finalB, A: alpha}
}

// ColorPickerSimple draws a simplified color picker with preset colors
func ColorPickerSimple(id int, x, y int, color utils.Color) utils.Color {
	const boxSize = 24
	const spacing = 4
	const cols = 8

	presets := []utils.Color{
		utils.Red, utils.Green, utils.Blue, utils.Yellow,
		utils.Orange, utils.Purple, utils.Cyan, utils.Magenta,
		utils.White, utils.Gray, utils.Black, utils.Pink,
		{R: 255, G: 128, B: 0, A: 255},   // Orange
		{R: 0, G: 255, B: 128, A: 255},   // Spring green
		{R: 128, G: 0, B: 255, A: 255},   // Violet
		{R: 255, G: 255, B: 128, A: 255}, // Light yellow
	}

	result := color

	for i, preset := range presets {
		col := i % cols
		row := i / cols
		bx := x + col*(boxSize+spacing)
		by := y + row*(boxSize+spacing)
		boxID := id + i

		if regionHit(bx, by, boxSize, boxSize) {
			state.HotItem = boxID
			if state.ActiveItem == 0 && state.MouseDown {
				state.ActiveItem = boxID
			}
		}

		// Draw color box
		rl.DrawRectangle(int32(bx), int32(by), boxSize, boxSize, toRLColor(preset))

		// Draw border/selection
		borderColor := currentTheme.Border
		if colorsEqual(color, preset) {
			borderColor = utils.White
			rl.DrawRectangleLines(int32(bx-1), int32(by-1), boxSize+2, boxSize+2, toRLColor(borderColor))
		}
		rl.DrawRectangleLines(int32(bx), int32(by), boxSize, boxSize, toRLColor(borderColor))

		// Select on click
		if state.HotItem == boxID && state.ActiveItem == boxID && !state.MouseDown {
			result = preset
		}
	}

	return result
}

// ColorButton draws a button that shows the current color and opens a color picker
func ColorButton(id int, x, y, w, h int, color utils.Color) utils.Color {
	if regionHit(x, y, w, h) {
		state.HotItem = id
		if state.ActiveItem == 0 && state.MouseDown {
			state.ActiveItem = id
		}
	}

	// Draw checkerboard for alpha
	checkerSize := 4
	for py := 0; py < h; py += checkerSize {
		for px := 0; px < w; px += checkerSize {
			c := rl.LightGray
			if ((px/checkerSize)+(py/checkerSize))%2 == 0 {
				c = rl.White
			}
			rl.DrawRectangle(int32(x+px), int32(y+py), int32(checkerSize), int32(checkerSize), c)
		}
	}

	// Draw color
	rl.DrawRectangle(int32(x), int32(y), int32(w), int32(h), toRLColor(color))

	// Draw border
	borderColor := currentTheme.Border
	if state.HotItem == id {
		borderColor = currentTheme.Accent
	}
	rl.DrawRectangleLines(int32(x), int32(y), int32(w), int32(h), toRLColor(borderColor))

	// Draw hex value
	hexStr := fmt.Sprintf("#%02X%02X%02X", color.R, color.G, color.B)
	textWidth := rl.MeasureText(hexStr, 12)
	rl.DrawText(hexStr, int32(x+(w-int(textWidth))/2), int32(y+h+2), 12, toRLColor(currentTheme.Text))

	return color
}

// Helper functions

func rgbToHsv(r, g, b uint8) (h, s, v float64) {
	rf := float64(r) / 255
	gf := float64(g) / 255
	bf := float64(b) / 255

	max := math.Max(rf, math.Max(gf, bf))
	min := math.Min(rf, math.Min(gf, bf))
	delta := max - min

	v = max

	if max == 0 {
		s = 0
	} else {
		s = delta / max
	}

	if delta == 0 {
		h = 0
	} else if max == rf {
		h = 60 * math.Mod((gf-bf)/delta, 6)
	} else if max == gf {
		h = 60 * ((bf-rf)/delta + 2)
	} else {
		h = 60 * ((rf-gf)/delta + 4)
	}

	if h < 0 {
		h += 360
	}

	return h, s, v
}

func hsvToRgb(h, s, v float64) (r, g, b uint8) {
	c := v * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := v - c

	var rf, gf, bf float64

	switch {
	case h < 60:
		rf, gf, bf = c, x, 0
	case h < 120:
		rf, gf, bf = x, c, 0
	case h < 180:
		rf, gf, bf = 0, c, x
	case h < 240:
		rf, gf, bf = 0, x, c
	case h < 300:
		rf, gf, bf = x, 0, c
	default:
		rf, gf, bf = c, 0, x
	}

	r = uint8((rf + m) * 255)
	g = uint8((gf + m) * 255)
	b = uint8((bf + m) * 255)
	return
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func colorsEqual(a, b utils.Color) bool {
	return a.R == b.R && a.G == b.G && a.B == b.B && a.A == b.A
}
