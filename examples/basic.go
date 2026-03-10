package main

import (
	"fmt"
	"log"

	_ "gomeow/memory"  // Used in commented examples
	_ "gomeow/overlay" // Used in commented examples
	"gomeow/process"
	"gomeow/utils"
)

func main() {
	// Example 1: Process and Memory Operations
	fmt.Println("=== goMeow Basic Example ===")
	fmt.Println()

	// List all processes
	fmt.Println("Listing processes:")
	processes, err := process.EnumProcesses()
	if err != nil {
		log.Printf("Failed to enumerate processes: %v", err)
	} else {
		for i, p := range processes[:min(5, len(processes))] {
			fmt.Printf("  %d. %s (PID: %d)\n", i+1, p.Name, p.PID)
		}
		fmt.Printf("  ... and %d more\n", len(processes)-5)
	}
	fmt.Println()

	// Example 2: Opening a process (commented out - needs real process)
	/*
		proc, err := process.OpenProcessByName("notepad.exe")
		if err != nil {
			log.Printf("Failed to open process: %v", err)
		} else {
			defer proc.Close()

			fmt.Printf("Opened: %s (PID: %d)\n", proc.Name, proc.PID)
			fmt.Printf("Base address: 0x%X\n", proc.Base)
			fmt.Printf("Is 64-bit: %v\n", proc.Is64Bit())

			// List modules
			modules, _ := proc.EnumModules()
			fmt.Printf("Loaded modules: %d\n", len(modules))
			for _, m := range modules[:min(3, len(modules))] {
				fmt.Printf("  - %s @ 0x%X (size: %d)\n", m.Name, m.Base, m.Size)
			}

			// Read memory example
			value, err := memory.ReadInt32(proc, proc.Base+0x100)
			if err != nil {
				log.Printf("Failed to read memory: %v", err)
			} else {
				fmt.Printf("Value at base+0x100: %d\n", value)
			}

			// Pattern scan example
			results, err := memory.AOBScanModule(proc, proc.Name, "48 8B ?? 90", false, false)
			if err != nil {
				log.Printf("Pattern scan failed: %v", err)
			} else {
				fmt.Printf("Pattern found at %d locations\n", len(results))
			}
		}
	*/

	// Example 3: Overlay (commented out - needs graphical environment)
	/*
		err = overlay.InitSimple("goMeow Overlay", 60)
		if err != nil {
			log.Fatalf("Failed to init overlay: %v", err)
		}

		for overlay.Loop() {
			overlay.BeginDrawing()

			// Draw some shapes
			overlay.DrawText("goMeow Overlay", 10, 10, 20, utils.White)
			overlay.DrawRectangleLines(100, 100, 200, 150, utils.Red, 2)
			overlay.DrawCircle(200, 175, 30, utils.Green)
			overlay.DrawLine(100, 100, 300, 250, utils.Blue, 1)

			// Draw FPS
			overlay.DrawFPS(10, 40)

			overlay.EndDrawing()
		}

		overlay.Close()
	*/

	fmt.Println("goMeow library loaded successfully!")
	fmt.Println()
	fmt.Println("Available packages:")
	fmt.Println("  - process: Process enumeration, opening, module listing")
	fmt.Println("  - memory:  Read/write memory, pattern scanning, pointer chains")
	fmt.Println("  - overlay: Transparent overlay window with raylib")
	fmt.Println("  - input:   Keyboard/mouse detection and simulation")
	fmt.Println("  - vec:     2D/3D vector mathematics")
	fmt.Println("  - utils:   Colors, world-to-screen, helpers")

	// Demonstrate some utility functions
	fmt.Println()
	fmt.Println("Color examples:")
	fmt.Printf("  Red: %+v\n", utils.Red)
	fmt.Printf("  From hex #00FF00: %+v\n", utils.FromHex("#00FF00"))
	fmt.Printf("  Red with 50%% alpha: %+v\n", utils.Red.Fade(0.5))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
