package vec

import (
	"math"
	"testing"
)

const epsilon = 0.0001

func floatEquals(a, b float32) bool {
	return math.Abs(float64(a-b)) < epsilon
}

// Vec2 Tests

func TestNewVec2(t *testing.T) {
	v := NewVec2(3.0, 4.0)
	if v.X != 3.0 || v.Y != 4.0 {
		t.Errorf("NewVec2(3, 4) = %v; want {3, 4}", v)
	}
}

func TestVec2Add(t *testing.T) {
	v1 := NewVec2(1, 2)
	v2 := NewVec2(3, 4)
	result := v1.Add(v2)
	if result.X != 4 || result.Y != 6 {
		t.Errorf("Vec2.Add: got %v; want {4, 6}", result)
	}
}

func TestVec2AddValue(t *testing.T) {
	v := NewVec2(1, 2)
	result := v.AddValue(5)
	if result.X != 6 || result.Y != 7 {
		t.Errorf("Vec2.AddValue: got %v; want {6, 7}", result)
	}
}

func TestVec2Subtract(t *testing.T) {
	v1 := NewVec2(5, 7)
	v2 := NewVec2(2, 3)
	result := v1.Subtract(v2)
	if result.X != 3 || result.Y != 4 {
		t.Errorf("Vec2.Subtract: got %v; want {3, 4}", result)
	}
}

func TestVec2Scale(t *testing.T) {
	v := NewVec2(2, 3)
	result := v.Scale(2)
	if result.X != 4 || result.Y != 6 {
		t.Errorf("Vec2.Scale: got %v; want {4, 6}", result)
	}
}

func TestVec2Length(t *testing.T) {
	v := NewVec2(3, 4)
	length := v.Length()
	if !floatEquals(length, 5) {
		t.Errorf("Vec2.Length: got %f; want 5", length)
	}
}

func TestVec2LengthSqr(t *testing.T) {
	v := NewVec2(3, 4)
	lengthSqr := v.LengthSqr()
	if !floatEquals(lengthSqr, 25) {
		t.Errorf("Vec2.LengthSqr: got %f; want 25", lengthSqr)
	}
}

func TestVec2Distance(t *testing.T) {
	v1 := NewVec2(0, 0)
	v2 := NewVec2(3, 4)
	dist := v1.Distance(v2)
	if !floatEquals(dist, 5) {
		t.Errorf("Vec2.Distance: got %f; want 5", dist)
	}
}

func TestVec2Normalize(t *testing.T) {
	v := NewVec2(3, 4)
	result := v.Normalize()
	if !floatEquals(result.X, 0.6) || !floatEquals(result.Y, 0.8) {
		t.Errorf("Vec2.Normalize: got %v; want {0.6, 0.8}", result)
	}
	if !floatEquals(result.Length(), 1) {
		t.Errorf("Normalized vector length should be 1, got %f", result.Length())
	}
}

func TestVec2NormalizeZero(t *testing.T) {
	v := NewVec2(0, 0)
	result := v.Normalize()
	if result.X != 0 || result.Y != 0 {
		t.Errorf("Vec2.Normalize(zero): got %v; want {0, 0}", result)
	}
}

func TestVec2Dot(t *testing.T) {
	v1 := NewVec2(1, 2)
	v2 := NewVec2(3, 4)
	dot := v1.Dot(v2)
	if !floatEquals(dot, 11) { // 1*3 + 2*4 = 11
		t.Errorf("Vec2.Dot: got %f; want 11", dot)
	}
}

func TestVec2Closest(t *testing.T) {
	origin := NewVec2(0, 0)
	v1 := NewVec2(10, 0)
	v2 := NewVec2(3, 4) // Distance 5
	v3 := NewVec2(20, 20)

	closest := origin.Closest(v1, v2, v3)
	if closest.X != v2.X || closest.Y != v2.Y {
		t.Errorf("Vec2.Closest: got %v; want %v", closest, v2)
	}
}

func TestVec2ClosestEmpty(t *testing.T) {
	v := NewVec2(5, 5)
	result := v.Closest()
	if result.X != v.X || result.Y != v.Y {
		t.Errorf("Vec2.Closest(empty): got %v; want %v", result, v)
	}
}

// Vec3 Tests

func TestNewVec3(t *testing.T) {
	v := NewVec3(1, 2, 3)
	if v.X != 1 || v.Y != 2 || v.Z != 3 {
		t.Errorf("NewVec3(1, 2, 3) = %v; want {1, 2, 3}", v)
	}
}

func TestVec3Add(t *testing.T) {
	v1 := NewVec3(1, 2, 3)
	v2 := NewVec3(4, 5, 6)
	result := v1.Add(v2)
	if result.X != 5 || result.Y != 7 || result.Z != 9 {
		t.Errorf("Vec3.Add: got %v; want {5, 7, 9}", result)
	}
}

func TestVec3Subtract(t *testing.T) {
	v1 := NewVec3(5, 7, 9)
	v2 := NewVec3(1, 2, 3)
	result := v1.Subtract(v2)
	if result.X != 4 || result.Y != 5 || result.Z != 6 {
		t.Errorf("Vec3.Subtract: got %v; want {4, 5, 6}", result)
	}
}

func TestVec3Scale(t *testing.T) {
	v := NewVec3(1, 2, 3)
	result := v.Scale(2)
	if result.X != 2 || result.Y != 4 || result.Z != 6 {
		t.Errorf("Vec3.Scale: got %v; want {2, 4, 6}", result)
	}
}

func TestVec3Length(t *testing.T) {
	v := NewVec3(2, 3, 6) // sqrt(4 + 9 + 36) = sqrt(49) = 7
	length := v.Length()
	if !floatEquals(length, 7) {
		t.Errorf("Vec3.Length: got %f; want 7", length)
	}
}

func TestVec3Distance(t *testing.T) {
	v1 := NewVec3(0, 0, 0)
	v2 := NewVec3(2, 3, 6)
	dist := v1.Distance(v2)
	if !floatEquals(dist, 7) {
		t.Errorf("Vec3.Distance: got %f; want 7", dist)
	}
}

func TestVec3Normalize(t *testing.T) {
	v := NewVec3(0, 3, 4)
	result := v.Normalize()
	if !floatEquals(result.Length(), 1) {
		t.Errorf("Normalized vector length should be 1, got %f", result.Length())
	}
}

func TestVec3Dot(t *testing.T) {
	v1 := NewVec3(1, 2, 3)
	v2 := NewVec3(4, 5, 6)
	dot := v1.Dot(v2)
	if !floatEquals(dot, 32) { // 1*4 + 2*5 + 3*6 = 32
		t.Errorf("Vec3.Dot: got %f; want 32", dot)
	}
}

func TestVec3Cross(t *testing.T) {
	v1 := NewVec3(1, 0, 0)
	v2 := NewVec3(0, 1, 0)
	result := v1.Cross(v2)
	// i x j = k
	if result.X != 0 || result.Y != 0 || result.Z != 1 {
		t.Errorf("Vec3.Cross(i, j): got %v; want {0, 0, 1}", result)
	}
}

func TestVec3CrossAntiCommutative(t *testing.T) {
	v1 := NewVec3(1, 2, 3)
	v2 := NewVec3(4, 5, 6)
	cross1 := v1.Cross(v2)
	cross2 := v2.Cross(v1)
	// a x b = -(b x a)
	if !floatEquals(cross1.X, -cross2.X) || !floatEquals(cross1.Y, -cross2.Y) || !floatEquals(cross1.Z, -cross2.Z) {
		t.Errorf("Cross product should be anti-commutative")
	}
}

func TestVec3ToVec2(t *testing.T) {
	v3 := NewVec3(1, 2, 3)
	v2 := v3.ToVec2()
	if v2.X != 1 || v2.Y != 2 {
		t.Errorf("Vec3.ToVec2: got %v; want {1, 2}", v2)
	}
}

func TestVec3Closest(t *testing.T) {
	origin := NewVec3(0, 0, 0)
	v1 := NewVec3(10, 0, 0)
	v2 := NewVec3(2, 3, 6) // Distance 7
	v3 := NewVec3(1, 1, 1) // Distance sqrt(3) ≈ 1.73

	closest := origin.Closest(v1, v2, v3)
	if closest.X != v3.X || closest.Y != v3.Y || closest.Z != v3.Z {
		t.Errorf("Vec3.Closest: got %v; want %v", closest, v3)
	}
}

// Benchmark tests

func BenchmarkVec2Add(b *testing.B) {
	v1 := NewVec2(1, 2)
	v2 := NewVec2(3, 4)
	for i := 0; i < b.N; i++ {
		_ = v1.Add(v2)
	}
}

func BenchmarkVec2Length(b *testing.B) {
	v := NewVec2(3, 4)
	for i := 0; i < b.N; i++ {
		_ = v.Length()
	}
}

func BenchmarkVec2Normalize(b *testing.B) {
	v := NewVec2(3, 4)
	for i := 0; i < b.N; i++ {
		_ = v.Normalize()
	}
}

func BenchmarkVec3Cross(b *testing.B) {
	v1 := NewVec3(1, 2, 3)
	v2 := NewVec3(4, 5, 6)
	for i := 0; i < b.N; i++ {
		_ = v1.Cross(v2)
	}
}
