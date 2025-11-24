package utils

// Color represents an RGBA color
type Color struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

// NewColor creates a new Color with full opacity
func NewColor(r, g, b uint8) Color {
	return Color{R: r, G: g, B: b, A: 255}
}

// NewColorAlpha creates a new Color with specified alpha
func NewColorAlpha(r, g, b, a uint8) Color {
	return Color{R: r, G: g, B: b, A: a}
}

// Predefined colors matching pyMeow
var (
	White      = Color{255, 255, 255, 255}
	Black      = Color{0, 0, 0, 255}
	Blank      = Color{0, 0, 0, 0}
	Red        = Color{255, 0, 0, 255}
	Green      = Color{0, 255, 0, 255}
	Blue       = Color{0, 0, 255, 255}
	Yellow     = Color{255, 255, 0, 255}
	Orange     = Color{255, 165, 0, 255}
	Pink       = Color{255, 192, 203, 255}
	Purple     = Color{128, 0, 128, 255}
	Cyan       = Color{0, 255, 255, 255}
	Magenta    = Color{255, 0, 255, 255}
	Gray       = Color{128, 128, 128, 255}
	DarkGray   = Color{64, 64, 64, 255}
	LightGray  = Color{192, 192, 192, 255}
	Brown      = Color{139, 69, 19, 255}
	Lime       = Color{50, 205, 50, 255}
	Gold       = Color{255, 215, 0, 255}
	SkyBlue    = Color{135, 206, 235, 255}
	DarkBlue   = Color{0, 0, 139, 255}
	DarkGreen  = Color{0, 100, 0, 255}
	DarkRed    = Color{139, 0, 0, 255}
	Maroon     = Color{128, 0, 0, 255}
	Navy       = Color{0, 0, 128, 255}
	Olive      = Color{128, 128, 0, 255}
	Teal       = Color{0, 128, 128, 255}
	Violet     = Color{238, 130, 238, 255}
	Beige      = Color{245, 245, 220, 255}
	RayWhite   = Color{245, 245, 245, 255}
)

// FromHex creates a Color from a hex string (e.g., "#FF0000" or "FF0000")
func FromHex(hex string) Color {
	if len(hex) > 0 && hex[0] == '#' {
		hex = hex[1:]
	}

	if len(hex) < 6 {
		return Black
	}

	r := hexToByte(hex[0:2])
	g := hexToByte(hex[2:4])
	b := hexToByte(hex[4:6])

	a := uint8(255)
	if len(hex) >= 8 {
		a = hexToByte(hex[6:8])
	}

	return Color{R: r, G: g, B: b, A: a}
}

func hexToByte(hex string) uint8 {
	var result uint8
	for _, c := range hex {
		result *= 16
		switch {
		case c >= '0' && c <= '9':
			result += uint8(c - '0')
		case c >= 'a' && c <= 'f':
			result += uint8(c - 'a' + 10)
		case c >= 'A' && c <= 'F':
			result += uint8(c - 'A' + 10)
		}
	}
	return result
}

// WithAlpha returns a copy of the color with a different alpha value
func (c Color) WithAlpha(a uint8) Color {
	return Color{R: c.R, G: c.G, B: c.B, A: a}
}

// Fade returns a copy of the color with alpha multiplied by the given factor (0.0-1.0)
func (c Color) Fade(alpha float32) Color {
	return Color{R: c.R, G: c.G, B: c.B, A: uint8(float32(c.A) * alpha)}
}
