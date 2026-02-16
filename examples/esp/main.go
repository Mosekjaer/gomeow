// ESP Overlay Example
// This example demonstrates how to create an ESP (Extra Sensory Perception) overlay
// that could be used for game cheats/hacks (for educational purposes only).
//
// The example shows:
// - Drawing ESP boxes around entities
// - Drawing health bars
// - Drawing corner boxes
// - Drawing lines to targets
// - Text labels with distance
// - Basic world-to-screen projection simulation
//
// In a real scenario, you would:
// 1. Read entity data from game memory using the memory package
// 2. Transform 3D positions to 2D screen coordinates using WorldToScreen
// 3. Draw ESP elements based on the entity data
//
// Run with: go run examples/esp/main.go

package main

import (
	"fmt"
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"

	"gomeow/gui"
	"gomeow/overlay"
	"gomeow/utils"
	"gomeow/vec"
)

// Simulated entity for demonstration
type Entity struct {
	Position  vec.Vec3
	ScreenX   float32
	ScreenY   float32
	Health    float32
	MaxHealth float32
	Name      string
	Distance  float32
	IsEnemy   bool
	Visible   bool
}

// ESP settings (would typically be configurable via GUI)
type ESPSettings struct {
	DrawBoxes      bool
	DrawCornerBox  bool
	DrawHealthBar  bool
	DrawSnaplines  bool
	DrawNames      bool
	DrawDistance   bool
	EnemyColor     utils.Color
	TeamColor      utils.Color
	SnaplineOrigin int // 0=bottom, 1=center, 2=top
}

func main() {
	// Initialize overlay
	err := overlay.Init("Full", 60, "ESP Overlay Demo", int(rl.KeyEnd), false)
	if err != nil {
		fmt.Println("Failed to init overlay:", err)
		return
	}

	screenW := overlay.GetScreenWidth()
	screenH := overlay.GetScreenHeight()

	// Create simulated entities
	entities := generateEntities(screenW, screenH)

	// ESP settings with defaults
	settings := ESPSettings{
		DrawBoxes:      true,
		DrawCornerBox:  true,
		DrawHealthBar:  true,
		DrawSnaplines:  true,
		DrawNames:      true,
		DrawDistance:   true,
		EnemyColor:     utils.Red,
		TeamColor:      utils.Green,
		SnaplineOrigin: 0,
	}

	// GUI window state
	showMenu := true
	menuX, menuY := 50, 50

	fmt.Println("ESP Overlay Demo Started")
	fmt.Println("Press INSERT to toggle menu")
	fmt.Println("Press END to exit")

	// Main loop
	for overlay.Loop() {
		overlay.BeginDrawing()

		// Toggle menu with INSERT key
		if rl.IsKeyPressed(rl.KeyInsert) {
			showMenu = !showMenu
		}

		// Draw ESP for all entities
		for _, entity := range entities {
			if !entity.Visible {
				continue
			}

			// Choose color based on team
			color := settings.TeamColor
			if entity.IsEnemy {
				color = settings.EnemyColor
			}

			// Calculate box dimensions based on distance
			// In a real scenario, this would be based on the entity's bounding box in world space
			boxH := int(200.0 / (entity.Distance*0.1 + 1))
			boxW := boxH / 2

			x := int(entity.ScreenX) - boxW/2
			y := int(entity.ScreenY) - boxH

			// Draw ESP elements
			if settings.DrawBoxes && !settings.DrawCornerBox {
				// Draw filled box with outline
				overlay.DrawBox(x, y, boxW, boxH, utils.NewColorAlpha(color.R, color.G, color.B, 30), color, 1.0)
			}

			if settings.DrawCornerBox {
				cornerLen := boxW / 4
				if cornerLen < 5 {
					cornerLen = 5
				}
				overlay.DrawCornerBox(x, y, boxW, boxH, cornerLen, color, 1.0)
			}

			if settings.DrawHealthBar {
				healthBarX := x - 5
				healthBarY := y
				healthBarW := 3

				// Get health color
				healthColor := getHealthColor(entity.Health, entity.MaxHealth)

				// Background
				overlay.DrawRectangle(healthBarX, healthBarY, healthBarW, boxH, utils.NewColorAlpha(0, 0, 0, 150))
				// Health fill (from bottom)
				healthH := int(float32(boxH) * (entity.Health / entity.MaxHealth))
				overlay.DrawRectangle(healthBarX, healthBarY+boxH-healthH, healthBarW, healthH, healthColor)
			}

			if settings.DrawSnaplines {
				// Snapline start position
				var startX, startY int
				switch settings.SnaplineOrigin {
				case 0: // Bottom
					startX = screenW / 2
					startY = screenH
				case 1: // Center
					startX = screenW / 2
					startY = screenH / 2
				case 2: // Top
					startX = screenW / 2
					startY = 0
				}

				overlay.DrawLine(startX, startY, int(entity.ScreenX), int(entity.ScreenY), color, 1.0)
			}

			if settings.DrawNames {
				textY := y - 20
				overlay.DrawText(entity.Name, x, textY, 14, color)
			}

			if settings.DrawDistance {
				distText := fmt.Sprintf("%.0fm", entity.Distance)
				textY := y + boxH + 5
				overlay.DrawText(distText, x, textY, 12, utils.White)
			}
		}

		// Draw FPS
		overlay.DrawFPS(10, 10)
		overlay.DrawText("ESP Demo - Press INSERT for menu, END to exit", 10, 30, 14, utils.White)

		// Draw settings menu
		if showMenu {
			gui.Begin()

			var closeClicked bool
			menuX, menuY, closeClicked = gui.Window(1, menuX, menuY, 280, 350, "ESP Settings", true)
			if closeClicked {
				showMenu = false
			}

			// Menu content area
			contentX := menuX + 15
			contentY := menuY + 40

			// Checkboxes for ESP features
			settings.DrawBoxes = gui.Checkbox(10, contentX, contentY, "Draw Boxes", settings.DrawBoxes)
			contentY += 28

			settings.DrawCornerBox = gui.Checkbox(11, contentX, contentY, "Use Corner Box", settings.DrawCornerBox)
			contentY += 28

			settings.DrawHealthBar = gui.Checkbox(12, contentX, contentY, "Draw Health Bars", settings.DrawHealthBar)
			contentY += 28

			settings.DrawSnaplines = gui.Checkbox(13, contentX, contentY, "Draw Snaplines", settings.DrawSnaplines)
			contentY += 28

			settings.DrawNames = gui.Checkbox(14, contentX, contentY, "Draw Names", settings.DrawNames)
			contentY += 28

			settings.DrawDistance = gui.Checkbox(15, contentX, contentY, "Draw Distance", settings.DrawDistance)
			contentY += 35

			// Snapline origin dropdown
			overlay.DrawText("Snapline Origin:", contentX, contentY, 14, utils.White)
			contentY += 20
			snaplineOptions := []string{"Bottom", "Center", "Top"}
			settings.SnaplineOrigin = gui.Dropdown(20, contentX, contentY, 150, 24, snaplineOptions, settings.SnaplineOrigin)
			contentY += 35

			// Color pickers
			overlay.DrawText("Enemy Color:", contentX, contentY, 14, utils.White)
			contentY += 20
			settings.EnemyColor = gui.ColorPickerSimple(30, contentX, contentY, settings.EnemyColor)
			contentY += 60

			overlay.DrawText("Team Color:", contentX, contentY, 14, utils.White)
			contentY += 20
			settings.TeamColor = gui.ColorPickerSimple(31, contentX, contentY, settings.TeamColor)

			gui.End()
		}

		// Animate entities (simulating movement)
		animateEntities(entities, screenW, screenH)

		overlay.EndDrawing()
	}

	overlay.Close()
	fmt.Println("ESP Overlay Demo Closed")
}

// getHealthColor returns a color based on health percentage (green -> yellow -> red)
func getHealthColor(health, maxHealth float32) utils.Color {
	ratio := health / maxHealth
	if ratio > 0.6 {
		return utils.Green
	} else if ratio > 0.3 {
		return utils.Yellow
	}
	return utils.Red
}

// generateEntities creates simulated entities for demonstration
func generateEntities(screenW, screenH int) []Entity {
	names := []string{"Player1", "Enemy1", "Target2", "Bot3", "Sniper4", "Assault5"}
	entities := make([]Entity, 10)

	for i := range entities {
		entities[i] = Entity{
			ScreenX:   float32(rand.Intn(screenW-200) + 100),
			ScreenY:   float32(rand.Intn(screenH-200) + 150),
			Health:    float32(rand.Intn(100) + 1),
			MaxHealth: 100,
			Name:      names[rand.Intn(len(names))],
			Distance:  float32(rand.Intn(200) + 10),
			IsEnemy:   rand.Intn(2) == 0,
			Visible:   true,
		}
	}

	return entities
}

// animateEntities simulates entity movement
func animateEntities(entities []Entity, screenW, screenH int) {
	for i := range entities {
		// Random small movement
		entities[i].ScreenX += float32(rand.Intn(3) - 1)
		entities[i].ScreenY += float32(rand.Intn(3) - 1)

		// Keep within bounds
		if entities[i].ScreenX < 100 {
			entities[i].ScreenX = 100
		}
		if entities[i].ScreenX > float32(screenW-100) {
			entities[i].ScreenX = float32(screenW - 100)
		}
		if entities[i].ScreenY < 150 {
			entities[i].ScreenY = 150
		}
		if entities[i].ScreenY > float32(screenH-50) {
			entities[i].ScreenY = float32(screenH - 50)
		}

		// Random health changes
		if rand.Intn(100) < 2 {
			entities[i].Health = float32(rand.Intn(100) + 1)
		}

		// Simulate distance changes
		entities[i].Distance += float32(math.Sin(float64(rl.GetTime())+float64(i))) * 0.5
		if entities[i].Distance < 10 {
			entities[i].Distance = 10
		}
	}
}
