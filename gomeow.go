// Package gomeow provides a Go library for memory manipulation and overlay rendering.
// It is a port of pyMeow (Nim/Python) to Go.
//
// Main packages:
//   - process: Process enumeration, opening, and module discovery
//   - memory: Memory read/write, pattern scanning, pointer chains
//   - overlay: Transparent overlay window rendering with raylib
//   - input: Keyboard and mouse detection/simulation
//   - pixel: Screen capture and pixel color operations
//   - vec: 2D and 3D vector mathematics
//   - utils: Colors, world-to-screen projection, helpers
//
// Example usage:
//
//	proc, _ := process.OpenProcessByName("game.exe")
//	defer proc.Close()
//
//	health, _ := memory.ReadInt32(proc, healthAddress)
//	fmt.Printf("Health: %d\n", health)
package gomeow

import (
	// Import all subpackages for documentation
	_ "gomeow/input"
	_ "gomeow/memory"
	_ "gomeow/overlay"
	_ "gomeow/pixel"
	_ "gomeow/process"
	_ "gomeow/utils"
	_ "gomeow/vec"
)

// Version of the goMeow library
const Version = "1.0.0"
