//go:build ignore

package main

import (
	"gomeow/hotkey"
	"gomeow/overlay"
	"gomeow/utils"
)

// Run with: go run examples/crosshair_example.go
//
// Fullscreen transparent overlay rendering a crosshair in the center of the
// screen. Hotkeys:
//   F1 — cycle crosshair style (cross / T / dot / circle / circle+dot)
//   F2 — toggle HUD visibility
//   F3 — toggle outline color (white / red)
//   End — exit
func main() {
	if err := overlay.InitSimple("crosshair-demo", 60); err != nil {
		panic(err)
	}
	defer overlay.Close()

	style := overlay.CrosshairCross
	showHUD := true
	color := utils.White

	hk := hotkey.New()
	hk.Register(overlay.VK_F1, func() {
		style = (style + 1) % 5
	})
	hk.RegisterToggle(overlay.VK_F2, func(on bool) { showHUD = on })
	hk.RegisterToggle(overlay.VK_F3, func(on bool) {
		if on {
			color = utils.Red
		} else {
			color = utils.White
		}
	})

	// Seed F2 toggle to default ON so the HUD shows up.
	// (RegisterToggle starts at false; flip our local mirror so initial draw matches.)
	showHUD = true

	w := overlay.GetScreenWidth()
	h := overlay.GetScreenHeight()
	cx, cy := w/2, h/2

	labels := map[int]string{
		overlay.VK_F1: "style",
		overlay.VK_F2: "HUD",
		overlay.VK_F3: "red",
	}

	for overlay.Loop() {
		hk.Poll()
		overlay.BeginDrawing()

		overlay.DrawCrosshair(style, cx, cy, 12, 4, 2, color)

		if showHUD {
			overlay.DrawHUD(overlay.HUDOptions{
				ShowFPS:   true,
				ShowFrame: true,
				ShowMouse: true,
				Hotkeys:   hk,
				Toggles:   hk,
				Labels:    labels,
			})
		}

		overlay.EndDrawing()
	}
}
