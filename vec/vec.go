package vec

import "math"

// Vec2 represents a 2D vector with X and Y components
type Vec2 struct {
	X float32
	Y float32
}

// Vec3 represents a 3D vector with X, Y, and Z components
type Vec3 struct {
	X float32
	Y float32
	Z float32
}

// NewVec2 creates a new Vec2
func NewVec2(x, y float32) Vec2 {
	return Vec2{X: x, Y: y}
}

// NewVec3 creates a new Vec3
func NewVec3(x, y, z float32) Vec3 {
	return Vec3{X: x, Y: y, Z: z}
}

// --- Vec2 Operations ---

// Add adds two Vec2 vectors
func (v Vec2) Add(other Vec2) Vec2 {
	return Vec2{X: v.X + other.X, Y: v.Y + other.Y}
}

// AddValue adds a scalar value to both components
func (v Vec2) AddValue(value float32) Vec2 {
	return Vec2{X: v.X + value, Y: v.Y + value}
}

// Subtract subtracts another Vec2 from this one
func (v Vec2) Subtract(other Vec2) Vec2 {
	return Vec2{X: v.X - other.X, Y: v.Y - other.Y}
}

// SubtractValue subtracts a scalar value from both components
func (v Vec2) SubtractValue(value float32) Vec2 {
	return Vec2{X: v.X - value, Y: v.Y - value}
}

// Multiply multiplies two Vec2 vectors component-wise
func (v Vec2) Multiply(other Vec2) Vec2 {
	return Vec2{X: v.X * other.X, Y: v.Y * other.Y}
}

// Scale multiplies the vector by a scalar
func (v Vec2) Scale(value float32) Vec2 {
	return Vec2{X: v.X * value, Y: v.Y * value}
}

// Divide divides two Vec2 vectors component-wise
func (v Vec2) Divide(other Vec2) Vec2 {
	return Vec2{X: v.X / other.X, Y: v.Y / other.Y}
}

// Length returns the length (magnitude) of the vector
func (v Vec2) Length() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}

// LengthSqr returns the squared length of the vector (faster than Length)
func (v Vec2) LengthSqr() float32 {
	return v.X*v.X + v.Y*v.Y
}

// Distance returns the distance between two Vec2 points
func (v Vec2) Distance(other Vec2) float32 {
	dx := v.X - other.X
	dy := v.Y - other.Y
	return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}

// Closest finds the closest vector from a list of vectors
func (v Vec2) Closest(vectors ...Vec2) Vec2 {
	if len(vectors) == 0 {
		return v
	}

	closest := vectors[0]
	closestDist := v.Distance(closest)

	for _, vec := range vectors[1:] {
		dist := v.Distance(vec)
		if dist < closestDist {
			closest = vec
			closestDist = dist
		}
	}

	return closest
}

// Normalize returns a unit vector in the same direction
func (v Vec2) Normalize() Vec2 {
	length := v.Length()
	if length < 0 {
		return Vec2{}
	}
	return Vec2{X: v.X / length, Y: v.Y / length}
}

// Dot returns the dot product of two vectors
func (v Vec2) Dot(other Vec2) float32 {
	return v.X*other.X + v.Y*other.Y
}

// --- Vec3 Operations ---

// Add adds two Vec3 vectors
func (v Vec3) Add(other Vec3) Vec3 {
	return Vec3{X: v.X + other.X, Y: v.Y + other.Y, Z: v.Z + other.Z}
}

// AddValue adds a scalar value to all components
func (v Vec3) AddValue(value float32) Vec3 {
	return Vec3{X: v.X + value, Y: v.Y + value, Z: v.Z + value}
}

// Subtract subtracts another Vec3 from this one
func (v Vec3) Subtract(other Vec3) Vec3 {
	return Vec3{X: v.X - other.X, Y: v.Y - other.Y, Z: v.Z - other.Z}
}

// SubtractValue subtracts a scalar value from all components
func (v Vec3) SubtractValue(value float32) Vec3 {
	return Vec3{X: v.X - value, Y: v.Y - value, Z: v.Z - value}
}

// Multiply multiplies two Vec3 vectors component-wise
func (v Vec3) Multiply(other Vec3) Vec3 {
	return Vec3{X: v.X * other.X, Y: v.Y * other.Y, Z: v.Z * other.Z}
}

// Scale multiplies the vector by a scalar
func (v Vec3) Scale(value float32) Vec3 {
	return Vec3{X: v.X * value, Y: v.Y * value, Z: v.Z * value}
}

// Divide divides two Vec3 vectors component-wise
func (v Vec3) Divide(other Vec3) Vec3 {
	return Vec3{X: v.X / other.X, Y: v.Y / other.Y, Z: v.Z / other.Z}
}

// Length returns the length (magnitude) of the vector
func (v Vec3) Length() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
}

// LengthSqr returns the squared length of the vector (faster than Length)
func (v Vec3) LengthSqr() float32 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

// Distance returns the distance between two Vec3 points
func (v Vec3) Distance(other Vec3) float32 {
	dx := v.X - other.X
	dy := v.Y - other.Y
	dz := v.Z - other.Z
	return float32(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))
}

// Closest finds the closest vector from a list of vectors
func (v Vec3) Closest(vectors ...Vec3) Vec3 {
	if len(vectors) == 0 {
		return v
	}

	closest := vectors[0]
	closestDist := v.Distance(closest)

	for _, vec := range vectors[1:] {
		dist := v.Distance(vec)
		if dist < closestDist {
			closest = vec
			closestDist = dist
		}
	}

	return closest
}

// Normalize returns a unit vector in the same direction
func (v Vec3) Normalize() Vec3 {
	length := v.Length()
	if length < 0 {
		return Vec3{}
	}
	return Vec3{X: v.X / length, Y: v.Y / length, Z: v.Z / length}
}

// Dot returns the dot product of two vectors
func (v Vec3) Dot(other Vec3) float32 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

// Cross returns the cross product of two vectors
func (v Vec3) Cross(other Vec3) Vec3 {
	return Vec3{
		X: v.Y*other.Z - v.Z*other.Y,
		Y: v.Z*other.X - v.X*other.Z,
		Z: v.X*other.Y - v.Y*other.X,
	}
}

// ToVec2 returns the X and Y components as a Vec2
func (v Vec3) ToVec2() Vec2 {
	return Vec2{X: v.X, Y: v.Y}
}
