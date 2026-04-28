package overlay

import (
	"fmt"
	"sort"

	rl "github.com/gen2brain/raylib-go/raylib"

	"gomeow/utils"
)

// HotkeySource is anything that can report current key-down state, intended
// to be satisfied by *hotkey.Manager. Defining the interface here keeps the
// overlay package free of a hotkey import.
type HotkeySource interface {
	Snapshot() map[int]bool
}

// ToggleSource exposes the toggle states managed by *hotkey.Manager.
// Optional — pass nil if not used.
type ToggleSource interface {
	ToggleStates() map[int]bool
}

// HUDOptions configures DrawHUD. All fields have sensible zero-value
// defaults; the most common customization is Labels for naming hotkeys.
type HUDOptions struct {
	X, Y     int         // top-left position; default (10, 10)
	FontSize int         // default 14
	LineGap  int         // pixels between lines; default 4
	Color    utils.Color // text color; zero-value → utils.White

	ShowFPS    bool
	ShowFrame  bool
	ShowMouse  bool

	Hotkeys HotkeySource // optional — renders one line per registered key
	Toggles ToggleSource // optional — overlays "ON/OFF" suffix on toggle keys

	// Labels maps key codes to human-readable names (e.g. {VK_F1: "ESP"}).
	// Keys without a label render as "0xNN".
	Labels map[int]string
}

// DrawHUD renders a simple status panel at (X, Y). Call between BeginDrawing
// and EndDrawing. The panel is text-only — overlay it on the transparent
// raylib surface with no background fill.
//
// Typical use:
//
//	hk := hotkey.New()
//	hk.RegisterToggle(overlay.VK_F1, ...)
//
//	for overlay.Loop() {
//	    hk.Poll()
//	    overlay.BeginDrawing()
//	    overlay.DrawHUD(overlay.HUDOptions{
//	        ShowFPS: true, ShowMouse: true,
//	        Hotkeys: hk, Toggles: hk,
//	        Labels: map[int]string{overlay.VK_F1: "ESP"},
//	    })
//	    overlay.EndDrawing()
//	}
func DrawHUD(opts HUDOptions) {
	if opts.FontSize <= 0 {
		opts.FontSize = 14
	}
	if opts.LineGap < 0 {
		opts.LineGap = 0
	}
	if opts.LineGap == 0 {
		opts.LineGap = 4
	}
	if opts.X == 0 && opts.Y == 0 {
		opts.X, opts.Y = 10, 10
	}
	if opts.Color == (utils.Color{}) {
		opts.Color = utils.White
	}

	lineH := opts.FontSize + opts.LineGap
	y := opts.Y
	emit := func(s string) {
		DrawText(s, opts.X, y, opts.FontSize, opts.Color)
		y += lineH
	}

	if opts.ShowFPS {
		emit(fmt.Sprintf("FPS: %d", rl.GetFPS()))
	}
	if opts.ShowFrame {
		emit(fmt.Sprintf("Frame: %.2f ms", rl.GetFrameTime()*1000))
	}
	if opts.ShowMouse {
		mp := rl.GetMousePosition()
		emit(fmt.Sprintf("Mouse: %d, %d", int(mp.X), int(mp.Y)))
	}

	if opts.Hotkeys == nil {
		return
	}
	snap := opts.Hotkeys.Snapshot()
	if len(snap) == 0 {
		return
	}
	var toggles map[int]bool
	if opts.Toggles != nil {
		toggles = opts.Toggles.ToggleStates()
	}

	// Sort for stable output regardless of map iteration order.
	keys := make([]int, 0, len(snap))
	for k := range snap {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, k := range keys {
		label, ok := opts.Labels[k]
		if !ok {
			label = fmt.Sprintf("0x%02X", k)
		}
		state := "off"
		if t, isToggle := toggles[k]; isToggle {
			if t {
				state = "ON"
			}
		} else if snap[k] {
			state = "down"
		}
		emit(fmt.Sprintf("%s: %s", label, state))
	}
}
