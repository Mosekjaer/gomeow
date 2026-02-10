package overlay

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"gomeow/utils"
)

// Font represents a loaded font
type Font struct {
	ID   int
	font rl.Font
}

var fontTable = make(map[int]*Font)

// toRLColor converts utils.Color to raylib Color
func toRLColor(c utils.Color) rl.Color {
	return rl.Color{R: c.R, G: c.G, B: c.B, A: c.A}
}

// DrawFPS draws the current FPS at the specified position
func DrawFPS(x, y int) {
	rl.DrawFPS(int32(x), int32(y))
}

// DrawText draws text at the specified position
func DrawText(text string, x, y int, fontSize int, color utils.Color) {
	rl.DrawText(text, int32(x), int32(y), int32(fontSize), toRLColor(color))
}

// DrawPixel draws a single pixel
func DrawPixel(x, y int, color utils.Color) {
	rl.DrawPixel(int32(x), int32(y), toRLColor(color))
}

// DrawLine draws a line between two points
func DrawLine(startX, startY, endX, endY int, color utils.Color, thickness float32) {
	if thickness <= 1.0 {
		rl.DrawLine(int32(startX), int32(startY), int32(endX), int32(endY), toRLColor(color))
	} else {
		rl.DrawLineEx(
			rl.Vector2{X: float32(startX), Y: float32(startY)},
			rl.Vector2{X: float32(endX), Y: float32(endY)},
			thickness,
			toRLColor(color),
		)
	}
}

// DrawCircle draws a filled circle
func DrawCircle(centerX, centerY int, radius float32, color utils.Color) {
	rl.DrawCircle(int32(centerX), int32(centerY), radius, toRLColor(color))
}

// DrawCircleLines draws a circle outline
func DrawCircleLines(centerX, centerY int, radius float32, color utils.Color) {
	rl.DrawCircleLines(int32(centerX), int32(centerY), radius, toRLColor(color))
}

// DrawCircleSector draws a filled circle sector
func DrawCircleSector(centerX, centerY int, radius, startAngle, endAngle float32, segments int, color utils.Color) {
	rl.DrawCircleSector(
		rl.Vector2{X: float32(centerX), Y: float32(centerY)},
		radius, startAngle, endAngle, int32(segments), toRLColor(color),
	)
}

// DrawCircleSectorLines draws a circle sector outline
func DrawCircleSectorLines(centerX, centerY int, radius, startAngle, endAngle float32, segments int, color utils.Color) {
	rl.DrawCircleSectorLines(
		rl.Vector2{X: float32(centerX), Y: float32(centerY)},
		radius, startAngle, endAngle, int32(segments), toRLColor(color),
	)
}

// DrawRing draws a ring (donut shape)
func DrawRing(centerX, centerY int, innerRadius, outerRadius, startAngle, endAngle float32, segments int, color utils.Color) {
	rl.DrawRing(
		rl.Vector2{X: float32(centerX), Y: float32(centerY)},
		innerRadius, outerRadius, startAngle, endAngle, int32(segments), toRLColor(color),
	)
}

// DrawRingLines draws a ring outline
func DrawRingLines(centerX, centerY int, innerRadius, outerRadius, startAngle, endAngle float32, segments int, color utils.Color) {
	rl.DrawRingLines(
		rl.Vector2{X: float32(centerX), Y: float32(centerY)},
		innerRadius, outerRadius, startAngle, endAngle, int32(segments), toRLColor(color),
	)
}

// DrawEllipse draws a filled ellipse
func DrawEllipse(centerX, centerY int, radiusH, radiusV float32, color utils.Color) {
	rl.DrawEllipse(int32(centerX), int32(centerY), radiusH, radiusV, toRLColor(color))
}

// DrawEllipseLines draws an ellipse outline
func DrawEllipseLines(centerX, centerY int, radiusH, radiusV float32, color utils.Color) {
	rl.DrawEllipseLines(int32(centerX), int32(centerY), radiusH, radiusV, toRLColor(color))
}

// DrawRectangle draws a filled rectangle
func DrawRectangle(x, y, width, height int, color utils.Color) {
	rl.DrawRectangle(int32(x), int32(y), int32(width), int32(height), toRLColor(color))
}

// DrawRectangleLines draws a rectangle outline
func DrawRectangleLines(x, y, width, height int, color utils.Color, lineThickness float32) {
	if lineThickness <= 1.0 {
		rl.DrawRectangleLines(int32(x), int32(y), int32(width), int32(height), toRLColor(color))
	} else {
		rect := rl.Rectangle{X: float32(x), Y: float32(y), Width: float32(width), Height: float32(height)}
		rl.DrawRectangleLinesEx(rect, lineThickness, toRLColor(color))
	}
}

// DrawRectangleRounded draws a filled rounded rectangle
func DrawRectangleRounded(x, y, width, height int, roundness float32, segments int, color utils.Color) {
	rect := rl.Rectangle{X: float32(x), Y: float32(y), Width: float32(width), Height: float32(height)}
	rl.DrawRectangleRounded(rect, roundness, int32(segments), toRLColor(color))
}

// DrawRectangleRoundedLines draws a rounded rectangle outline
func DrawRectangleRoundedLines(x, y, width, height int, roundness float32, segments int, color utils.Color, lineThickness float32) {
	rect := rl.Rectangle{X: float32(x), Y: float32(y), Width: float32(width), Height: float32(height)}
	rl.DrawRectangleRoundedLines(rect, roundness, float32(segments), lineThickness, toRLColor(color))
}

// DrawTriangle draws a filled triangle
func DrawTriangle(x1, y1, x2, y2, x3, y3 int, color utils.Color) {
	rl.DrawTriangle(
		rl.Vector2{X: float32(x1), Y: float32(y1)},
		rl.Vector2{X: float32(x2), Y: float32(y2)},
		rl.Vector2{X: float32(x3), Y: float32(y3)},
		toRLColor(color),
	)
}

// DrawTriangleLines draws a triangle outline
func DrawTriangleLines(x1, y1, x2, y2, x3, y3 int, color utils.Color) {
	rl.DrawTriangleLines(
		rl.Vector2{X: float32(x1), Y: float32(y1)},
		rl.Vector2{X: float32(x2), Y: float32(y2)},
		rl.Vector2{X: float32(x3), Y: float32(y3)},
		toRLColor(color),
	)
}

// DrawPoly draws a filled polygon
func DrawPoly(centerX, centerY int, sides int, radius, rotation float32, color utils.Color) {
	rl.DrawPoly(
		rl.Vector2{X: float32(centerX), Y: float32(centerY)},
		int32(sides), radius, rotation, toRLColor(color),
	)
}

// DrawPolyLines draws a polygon outline
func DrawPolyLines(centerX, centerY int, sides int, radius, rotation, lineThickness float32, color utils.Color) {
	rl.DrawPolyLinesEx(
		rl.Vector2{X: float32(centerX), Y: float32(centerY)},
		int32(sides), radius, rotation, lineThickness, toRLColor(color),
	)
}

// LoadTexture loads a texture from file
func LoadTexture(fileName string) rl.Texture2D {
	return rl.LoadTexture(fileName)
}

// LoadTextureBytes loads a texture from memory
func LoadTextureBytes(fileType string, data []byte) rl.Texture2D {
	img := rl.LoadImageFromMemory(fileType, data, int32(len(data)))
	tex := rl.LoadTextureFromImage(img)
	rl.UnloadImage(img)
	return tex
}

// DrawTexture draws a texture
func DrawTexture(texture rl.Texture2D, x, y int, tint utils.Color, rotation, scale float32) {
	rl.DrawTextureEx(texture, rl.Vector2{X: float32(x), Y: float32(y)}, rotation, scale, toRLColor(tint))
}

// UnloadTexture unloads a texture from memory
func UnloadTexture(texture rl.Texture2D) {
	rl.UnloadTexture(texture)
}

// LoadFont loads a font from file
func LoadFont(fileName string, fontID int) {
	fontTable[fontID] = &Font{
		ID:   fontID,
		font: rl.LoadFont(fileName),
	}
}

// DrawFont draws text using a custom font
func DrawFont(fontID int, text string, x, y int, fontSize, spacing float32, tint utils.Color) {
	if f, ok := fontTable[fontID]; ok {
		rl.DrawTextEx(f.font, text, rl.Vector2{X: float32(x), Y: float32(y)}, fontSize, spacing, toRLColor(tint))
	}
}

// MeasureFont measures text dimensions using a custom font
func MeasureFont(fontID int, text string, fontSize, spacing float32) (width, height float32) {
	if f, ok := fontTable[fontID]; ok {
		size := rl.MeasureTextEx(f.font, text, fontSize, spacing)
		return size.X, size.Y
	}
	return 0, 0
}

// UnloadFont unloads a font
func UnloadFont(fontID int) {
	if f, ok := fontTable[fontID]; ok {
		rl.UnloadFont(f.font)
		delete(fontTable, fontID)
	}
}

// DrawBox draws a 3D-style box (commonly used for ESP)
func DrawBox(x, y, width, height int, color, outlineColor utils.Color, outlineThickness float32) {
	// Fill
	DrawRectangle(x, y, width, height, color)
	// Outline
	DrawRectangleLines(x, y, width, height, outlineColor, outlineThickness)
}

// DrawHealthBar draws a health bar
func DrawHealthBar(x, y, width, height int, health, maxHealth float32, bgColor, fgColor utils.Color) {
	// Background
	DrawRectangle(x, y, width, height, bgColor)
	// Health fill
	healthWidth := int(float32(width) * (health / maxHealth))
	DrawRectangle(x, y, healthWidth, height, fgColor)
}

// DrawCornerBox draws a corner-style box (commonly used for ESP)
func DrawCornerBox(x, y, width, height int, cornerLength int, color utils.Color, thickness float32) {
	// Top-left
	DrawLine(x, y, x+cornerLength, y, color, thickness)
	DrawLine(x, y, x, y+cornerLength, color, thickness)

	// Top-right
	DrawLine(x+width, y, x+width-cornerLength, y, color, thickness)
	DrawLine(x+width, y, x+width, y+cornerLength, color, thickness)

	// Bottom-left
	DrawLine(x, y+height, x+cornerLength, y+height, color, thickness)
	DrawLine(x, y+height, x, y+height-cornerLength, color, thickness)

	// Bottom-right
	DrawLine(x+width, y+height, x+width-cornerLength, y+height, color, thickness)
	DrawLine(x+width, y+height, x+width, y+height-cornerLength, color, thickness)
}
