package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"gomeow/gui"
	"gomeow/overlay"
	"gomeow/utils"
)

func main() {
	// Initialize overlay window
	err := overlay.Init("Full", 60, "GUI Demo", int(rl.KeyEscape), false)
	if err != nil {
		fmt.Println("Failed to init overlay:", err)
		return
	}

	// GUI state
	sliderValue := float32(50)
	checkboxValue := true
	checkbox2Value := false
	textValue := "Hello World"
	dropdownSelected := 0
	dropdownItems := []string{"Option 1", "Option 2", "Option 3", "Option 4"}
	tabSelected := 0
	tabs := []string{"General", "Settings", "Colors"}
	listSelected := 0
	listItems := []string{
		"Item 1", "Item 2", "Item 3", "Item 4", "Item 5",
		"Item 6", "Item 7", "Item 8", "Item 9", "Item 10",
	}
	selectedColor := utils.Red
	progressValue := float32(0)
	windowX, windowY := 500, 100

	// Run main loop
	for overlay.Loop() {
		overlay.BeginDrawing()

		// Dark background for demo (not transparent)
		rl.ClearBackground(rl.Color{R: 30, G: 30, B: 30, A: 255})

		// Begin GUI frame
		gui.Begin()

		// Title
		rl.DrawText("goMeow GUI Demo", 20, 20, 24, rl.White)
		rl.DrawText("Press ESC to exit", 20, 50, 16, rl.Gray)

		// Tab bar
		tabSelected = gui.TabBar(1, 20, 90, 400, 32, tabs, tabSelected)

		// Content based on selected tab
		switch tabSelected {
		case 0: // General
			drawGeneralTab(&sliderValue, &checkboxValue, &checkbox2Value, &textValue, &dropdownSelected, dropdownItems)
		case 1: // Settings
			drawSettingsTab(&listSelected, listItems, &progressValue)
		case 2: // Colors
			drawColorsTab(&selectedColor)
		}

		// Draggable window demo
		var closeClicked bool
		windowX, windowY, closeClicked = gui.Window(100, windowX, windowY, 250, 200, "Draggable Window", true)
		if closeClicked {
			windowX, windowY = 500, 100 // Reset position on close
		}

		// Window content
		gui.Panel(windowX+5, windowY+30, 240, 165)
		rl.DrawText("Drag the title bar!", int32(windowX+15), int32(windowY+50), 16, rl.White)
		rl.DrawText("Click X to reset", int32(windowX+15), int32(windowY+80), 16, rl.Gray)

		if gui.Button(101, windowX+15, windowY+120, 100, 30, "Click Me") {
			fmt.Println("Button in window clicked!")
		}

		// Animate progress
		progressValue += 0.5
		if progressValue > 100 {
			progressValue = 0
		}

		// End GUI frame
		gui.End()

		overlay.EndDrawing()
	}

	overlay.Close()
}

func drawGeneralTab(sliderValue *float32, checkboxValue, checkbox2Value *bool, textValue *string, dropdownSelected *int, dropdownItems []string) {
	y := 140

	// Buttons
	rl.DrawText("Buttons:", 20, int32(y), 16, rl.Gray)
	y += 25

	if gui.Button(10, 20, y, 100, 30, "Button 1") {
		fmt.Println("Button 1 clicked!")
	}

	if gui.Button(11, 130, y, 100, 30, "Button 2") {
		fmt.Println("Button 2 clicked!")
	}

	if gui.ButtonEx(12, 240, y, 100, 30, "Custom", utils.NewColorAlpha(80, 120, 200, 255), utils.White) {
		fmt.Println("Custom button clicked!")
	}

	y += 50

	// Checkboxes
	rl.DrawText("Checkboxes:", 20, int32(y), 16, rl.Gray)
	y += 25

	*checkboxValue = gui.Checkbox(20, 20, y, "Enable feature", *checkboxValue)
	y += 30
	*checkbox2Value = gui.Checkbox(21, 20, y, "Another option", *checkbox2Value)
	y += 40

	// Slider
	rl.DrawText("Slider:", 20, int32(y), 16, rl.Gray)
	y += 25
	*sliderValue = gui.SliderWithLabel(30, 20, y, 200, 20, "Value", *sliderValue, 0, 100)
	y += 50

	// Text input
	rl.DrawText("Text Input:", 20, int32(y), 16, rl.Gray)
	y += 25
	*textValue = gui.TextBox(40, 20, y, 200, 28, *textValue, 50)
	y += 50

	// Dropdown
	rl.DrawText("Dropdown:", 20, int32(y), 16, rl.Gray)
	y += 25
	*dropdownSelected = gui.Dropdown(50, 20, y, 200, 28, dropdownItems, *dropdownSelected)

	// Help marker with tooltip
	gui.HelpMarker(60, 230, y+5, "Select an option from the dropdown list")
}

func drawSettingsTab(listSelected *int, listItems []string, progressValue *float32) {
	y := 140

	// List view
	rl.DrawText("List View:", 20, int32(y), 16, rl.Gray)
	y += 25
	*listSelected = gui.ListView(70, 20, y, 200, 200, listItems, *listSelected)

	// Show selected item
	if *listSelected >= 0 && *listSelected < len(listItems) {
		rl.DrawText(fmt.Sprintf("Selected: %s", listItems[*listSelected]), 230, int32(y), 16, rl.White)
	}

	y += 220

	// Progress bars
	rl.DrawText("Progress Bars:", 20, int32(y), 16, rl.Gray)
	y += 25

	gui.ProgressBar(20, y, 200, 20, *progressValue, 0, 100)
	y += 30

	gui.ProgressBarWithLabel(20, y, 200, 20, *progressValue, 0, 100)
	y += 30

	gui.HealthBar(20, y, 200, 20, *progressValue, 100)
	y += 30

	gui.HealthBarGradient(20, y, 200, 20, *progressValue, 100)
	y += 40

	// Loading spinner
	rl.DrawText("Loading Spinner:", 20, int32(y), 16, rl.Gray)
	gui.LoadingSpinner(100, y+50, 20)
}

func drawColorsTab(selectedColor *utils.Color) {
	y := 140

	// Simple color picker
	rl.DrawText("Color Presets:", 20, int32(y), 16, rl.Gray)
	y += 25
	*selectedColor = gui.ColorPickerSimple(80, 20, y, *selectedColor)

	y += 80

	// Current color display
	rl.DrawText("Selected Color:", 20, int32(y), 16, rl.Gray)
	y += 25
	gui.ColorButton(90, 20, y, 60, 40, *selectedColor)

	// RGB values
	rl.DrawText(fmt.Sprintf("R: %d  G: %d  B: %d  A: %d",
		selectedColor.R, selectedColor.G, selectedColor.B, selectedColor.A),
		100, int32(y+10), 16, rl.White)

	y += 60

	// Full color picker
	rl.DrawText("Full Color Picker:", 20, int32(y), 16, rl.Gray)
	y += 25
	*selectedColor = gui.ColorPicker(91, 20, y, 200, 200, *selectedColor)
}
