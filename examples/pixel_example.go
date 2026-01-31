//go:build ignore

package main

import (
	"fmt"
	"time"

	"gomeow/input"
	"gomeow/pixel"
	"gomeow/utils"
)

func main() {
	fmt.Println("=== goMeow Pixel Operations Example ===")
	fmt.Println()

	// Get screen size
	width, height := pixel.GetScreenSize()
	fmt.Printf("Screen size: %dx%d\n", width, height)
	fmt.Println()

	// Get pixel color at center of screen
	centerX, centerY := width/2, height/2
	color := pixel.GetPixelColor(centerX, centerY)
	fmt.Printf("Pixel at center (%d, %d): R=%d G=%d B=%d\n",
		centerX, centerY, color.R, color.G, color.B)
	fmt.Println()

	// Interactive: Show pixel color under mouse cursor
	fmt.Println("Showing pixel color under mouse cursor for 5 seconds...")
	fmt.Println("Move your mouse around to see different colors!")
	fmt.Println()

	ticker := time.NewTicker(100 * time.Millisecond)
	timeout := time.After(5 * time.Second)

	for {
		select {
		case <-timeout:
			ticker.Stop()
			fmt.Println()
			goto done
		case <-ticker.C:
			pos := input.MousePosition()
			x, y := int(pos.X), int(pos.Y)
			c := pixel.GetPixelColor(x, y)
			fmt.Printf("\rMouse (%4d, %4d) -> R=%3d G=%3d B=%3d  ", x, y, c.R, c.G, c.B)
		}
	}
done:
	fmt.Println()

	// Take a screenshot
	fmt.Println("Capturing screenshot of top-left 400x300...")
	err := pixel.SaveScreenshot("screenshot.png", 0, 0, 400, 300)
	if err != nil {
		fmt.Printf("Screenshot failed: %v\n", err)
	} else {
		fmt.Println("Screenshot saved to screenshot.png")
	}
	fmt.Println()

	// Analyze colors in a region
	fmt.Println("Analyzing colors in 100x100 region at (100,100)...")
	colorCounts := make(map[uint32]int)
	startTime := time.Now()

	err = pixel.EnumPixels(100, 100, 199, 199, func(x, y int, c utils.Color) bool {
		key := uint32(c.R)<<16 | uint32(c.G)<<8 | uint32(c.B)
		colorCounts[key]++
		return true
	})

	if err != nil {
		fmt.Printf("Enumeration failed: %v\n", err)
	} else {
		elapsed := time.Since(startTime)
		fmt.Printf("Analyzed 10,000 pixels in %v\n", elapsed)
		fmt.Printf("Found %d unique colors\n", len(colorCounts))

		// Find most common color
		var maxCount int
		var maxColor uint32
		for color, count := range colorCounts {
			if count > maxCount {
				maxCount = count
				maxColor = color
			}
		}

		r := uint8((maxColor >> 16) & 0xFF)
		g := uint8((maxColor >> 8) & 0xFF)
		b := uint8(maxColor & 0xFF)
		fmt.Printf("Most common color: R=%d G=%d B=%d (count: %d)\n", r, g, b, maxCount)
	}
	fmt.Println()

	// Search for the most common color we just found
	fmt.Println("Searching for the most common color in top-left 500x500...")
	c := pixel.GetPixelColor(150, 150) // Get a sample color from the region
	x, y, found := pixel.PixelSearch(0, 0, 500, 500, c, 10)
	if found {
		fmt.Printf("Found matching pixel at (%d, %d)\n", x, y)
	} else {
		fmt.Println("No matching pixels found")
	}
	fmt.Println()

	fmt.Println("Done!")
}
