// Package gui provides GUI components built on top of raylib.
// These components can be used to create interactive overlays with buttons,
// sliders, checkboxes, and other UI elements.
package gui

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"gomeow/utils"
	"gomeow/vec"
)

// State tracks GUI interaction state
type State struct {
	MousePos     vec.Vec2
	MousePressed bool
	MouseDown    bool
	HotItem      int
	ActiveItem   int
}

var state State

// Begin initializes the GUI frame and updatesstarts a new GUI frame - call this before drawing any GUI elements
func Begin() {
	pos := rl.GetMousePosition()
	state.MousePos = vec.Vec2{X: pos.X, Y: pos.Y}
	state.MouseDown = rl.IsMouseButtonDown(rl.MouseLeftButton)
	state.MousePressed = rl.IsMouseButtonPressed(rl.MouseLeftButton)
	state.HotItem = 0
}

// End finishes the GUI frame - call this after drawing all GUI elements
func End() {
	if !state.MouseDown {
		state.ActiveItem = 0
	}
}

// toRLColor converts utils.Color to raylib Color
func toRLColor(c utils.Color) rl.Color {
	return rl.Color{R: c.R, G: c.G, B: c.B, A: c.A}
}

// regionHit checks if mouse is within a rectangle
func regionHit(x, y, w, h int) bool {
	return state.MousePos.X >= float32(x) &&
		state.MousePos.X < float32(x+w) &&
		state.MousePos.Y >= float32(y) &&
		state.MousePos.Y < float32(y+h)
}

// Theme defines colors for GUI elements
type Theme struct {
	Background     utils.Color
	BackgroundHot  utils.Color
	BackgroundActive utils.Color
	Border         utils.Color
	Text           utils.Color
	TextDisabled   utils.Color
	Accent         utils.Color
	SliderFill     utils.Color
}

// DefaultTheme returns a dark theme suitable for overlays
var DefaultTheme = Theme{
	Background:       utils.NewColorAlpha(40, 40, 40, 200),
	BackgroundHot:    utils.NewColorAlpha(60, 60, 60, 200),
	BackgroundActive: utils.NewColorAlpha(80, 80, 80, 200),
	Border:           utils.NewColorAlpha(100, 100, 100, 255),
	Text:             utils.White,
	TextDisabled:     utils.Gray,
	Accent:           utils.NewColorAlpha(0, 150, 255, 255),
	SliderFill:       utils.NewColorAlpha(0, 150, 255, 255),
}

var currentTheme = DefaultTheme

// SetTheme sets the current GUI theme
func SetTheme(theme Theme) {
	currentTheme = theme
}

// GetTheme returns the current GUI theme
func GetTheme() Theme {
	return currentTheme
}
