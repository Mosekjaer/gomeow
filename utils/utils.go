package utils

import (
	"gomeow/vec"
)

// WorldToScreen converts a 3D world position to 2D screen coordinates
// viewMatrix should be a 4x4 matrix as a flat array [16]float32
// Returns screen position and whether the position is in front of the camera
func WorldToScreen(pos vec.Vec3, viewMatrix [16]float32, screenWidth, screenHeight int) (vec.Vec2, bool) {
	// Matrix multiplication
	clipX := pos.X*viewMatrix[0] + pos.Y*viewMatrix[4] + pos.Z*viewMatrix[8] + viewMatrix[12]
	clipY := pos.X*viewMatrix[1] + pos.Y*viewMatrix[5] + pos.Z*viewMatrix[9] + viewMatrix[13]
	clipW := pos.X*viewMatrix[3] + pos.Y*viewMatrix[7] + pos.Z*viewMatrix[11] + viewMatrix[15]

	// Check if behind camera
	if clipW < 0.1 {
		return vec.Vec2{}, false
	}

	// Perspective division
	ndcX := clipX / clipW
	ndcY := clipY / clipW

	// Convert to screen coordinates
	screenX := (float32(screenWidth) / 2) * (1 + ndcX)
	screenY := (float32(screenHeight) / 2) * (1 - ndcY)

	return vec.Vec2{X: screenX, Y: screenY}, true
}

// WorldToScreenSimple is a simplified W2S for common game engines
// viewMatrix should be a 4x4 matrix as a flat array [16]float32
func WorldToScreenSimple(pos vec.Vec3, viewMatrix [16]float32, screenWidth, screenHeight int) (vec.Vec2, bool) {
	w := viewMatrix[12]*pos.X + viewMatrix[13]*pos.Y + viewMatrix[14]*pos.Z + viewMatrix[15]

	if w < 0.01 {
		return vec.Vec2{}, false
	}

	x := viewMatrix[0]*pos.X + viewMatrix[1]*pos.Y + viewMatrix[2]*pos.Z + viewMatrix[3]
	y := viewMatrix[4]*pos.X + viewMatrix[5]*pos.Y + viewMatrix[6]*pos.Z + viewMatrix[7]

	screenX := (float32(screenWidth) / 2) + (float32(screenWidth)/2)*x/w
	screenY := (float32(screenHeight) / 2) - (float32(screenHeight)/2)*y/w

	return vec.Vec2{X: screenX, Y: screenY}, true
}

// Clamp clamps a value between min and max
func Clamp(value, min, max float32) float32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// ClampInt clamps an integer value between min and max
func ClampInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Lerp performs linear interpolation between two values
func Lerp(a, b, t float32) float32 {
	return a + (b-a)*t
}

// InBounds checks if a point is within screen bounds
func InBounds(x, y, screenWidth, screenHeight int) bool {
	return x >= 0 && x < screenWidth && y >= 0 && y < screenHeight
}

// InBoundsVec checks if a Vec2 is within screen bounds
func InBoundsVec(pos vec.Vec2, screenWidth, screenHeight int) bool {
	return pos.X >= 0 && pos.X < float32(screenWidth) && pos.Y >= 0 && pos.Y < float32(screenHeight)
}
