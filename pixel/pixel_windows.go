//go:build windows

package pixel

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"unsafe"

	"golang.org/x/sys/windows"

	"gomeow/utils"
)

var (
	user32           = windows.NewLazySystemDLL("user32.dll")
	gdi32            = windows.NewLazySystemDLL("gdi32.dll")
	procGetDC        = user32.NewProc("GetDC")
	procReleaseDC    = user32.NewProc("ReleaseDC")
	procGetPixel     = gdi32.NewProc("GetPixel")
	procCreateCompatibleDC = gdi32.NewProc("CreateCompatibleDC")
	procCreateCompatibleBitmap = gdi32.NewProc("CreateCompatibleBitmap")
	procSelectObject = gdi32.NewProc("SelectObject")
	procBitBlt       = gdi32.NewProc("BitBlt")
	procDeleteDC     = gdi32.NewProc("DeleteDC")
	procDeleteObject = gdi32.NewProc("DeleteObject")
	procGetDIBits    = gdi32.NewProc("GetDIBits")
	procGetSystemMetrics = user32.NewProc("GetSystemMetrics")
)

const (
	SM_CXSCREEN = 0
	SM_CYSCREEN = 1
	SRCCOPY     = 0x00CC0020
	DIB_RGB_COLORS = 0
	BI_RGB         = 0
)

type bitmapInfoHeader struct {
	BiSize          uint32
	BiWidth         int32
	BiHeight        int32
	BiPlanes        uint16
	BiBitCount      uint16
	BiCompression   uint32
	BiSizeImage     uint32
	BiXPelsPerMeter int32
	BiYPelsPerMeter int32
	BiClrUsed       uint32
	BiClrImportant  uint32
}

type bitmapInfo struct {
	BmiHeader bitmapInfoHeader
	BmiColors [1]uint32
}

// GetScreenSize returns the screen dimensions
func GetScreenSize() (width, height int) {
	w, _, _ := procGetSystemMetrics.Call(SM_CXSCREEN)
	h, _, _ := procGetSystemMetrics.Call(SM_CYSCREEN)
	return int(w), int(h)
}

// GetPixelColor gets the color of a pixel at the specified screen coordinates
func GetPixelColor(x, y int) utils.Color {
	hdc, _, _ := procGetDC.Call(0)
	defer procReleaseDC.Call(0, hdc)

	colorRef, _, _ := procGetPixel.Call(hdc, uintptr(x), uintptr(y))

	// COLORREF is 0x00BBGGRR
	return utils.Color{
		R: uint8(colorRef & 0xFF),
		G: uint8((colorRef >> 8) & 0xFF),
		B: uint8((colorRef >> 16) & 0xFF),
		A: 255,
	}
}

// PixelSearch searches for a pixel with the specified color in a region
// Returns the coordinates if found, or (-1, -1) if not found
func PixelSearch(x1, y1, x2, y2 int, color utils.Color, tolerance int) (int, int, bool) {
	hdc, _, _ := procGetDC.Call(0)
	defer procReleaseDC.Call(0, hdc)

	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			colorRef, _, _ := procGetPixel.Call(hdc, uintptr(x), uintptr(y))

			r := uint8(colorRef & 0xFF)
			g := uint8((colorRef >> 8) & 0xFF)
			b := uint8((colorRef >> 16) & 0xFF)

			if colorMatches(r, g, b, color.R, color.G, color.B, tolerance) {
				return x, y, true
			}
		}
	}

	return -1, -1, false
}

// PixelSearchAll finds all pixels matching the color in a region
func PixelSearchAll(x1, y1, x2, y2 int, color utils.Color, tolerance int) []image.Point {
	var results []image.Point

	hdc, _, _ := procGetDC.Call(0)
	defer procReleaseDC.Call(0, hdc)

	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			colorRef, _, _ := procGetPixel.Call(hdc, uintptr(x), uintptr(y))

			r := uint8(colorRef & 0xFF)
			g := uint8((colorRef >> 8) & 0xFF)
			b := uint8((colorRef >> 16) & 0xFF)

			if colorMatches(r, g, b, color.R, color.G, color.B, tolerance) {
				results = append(results, image.Point{X: x, Y: y})
			}
		}
	}

	return results
}

// CaptureScreen captures a region of the screen and returns raw pixel data
func CaptureScreen(x, y, width, height int) ([]byte, error) {
	hdcScreen, _, _ := procGetDC.Call(0)
	if hdcScreen == 0 {
		return nil, fmt.Errorf("GetDC failed")
	}
	defer procReleaseDC.Call(0, hdcScreen)

	hdcMem, _, _ := procCreateCompatibleDC.Call(hdcScreen)
	if hdcMem == 0 {
		return nil, fmt.Errorf("CreateCompatibleDC failed")
	}
	defer procDeleteDC.Call(hdcMem)

	hBitmap, _, _ := procCreateCompatibleBitmap.Call(hdcScreen, uintptr(width), uintptr(height))
	if hBitmap == 0 {
		return nil, fmt.Errorf("CreateCompatibleBitmap failed")
	}
	defer procDeleteObject.Call(hBitmap)

	procSelectObject.Call(hdcMem, hBitmap)

	ret, _, _ := procBitBlt.Call(
		hdcMem, 0, 0, uintptr(width), uintptr(height),
		hdcScreen, uintptr(x), uintptr(y),
		SRCCOPY,
	)
	if ret == 0 {
		return nil, fmt.Errorf("BitBlt failed")
	}

	// Prepare bitmap info
	bi := bitmapInfo{
		BmiHeader: bitmapInfoHeader{
			BiSize:        uint32(unsafe.Sizeof(bitmapInfoHeader{})),
			BiWidth:       int32(width),
			BiHeight:      -int32(height), // Negative for top-down
			BiPlanes:      1,
			BiBitCount:    32,
			BiCompression: BI_RGB,
		},
	}

	// Allocate buffer for pixel data (BGRA format)
	bufSize := width * height * 4
	buf := make([]byte, bufSize)

	ret, _, _ = procGetDIBits.Call(
		hdcMem,
		hBitmap,
		0,
		uintptr(height),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&bi)),
		DIB_RGB_COLORS,
	)
	if ret == 0 {
		return nil, fmt.Errorf("GetDIBits failed")
	}

	return buf, nil
}

// SaveScreenshot captures a region and saves it as a PNG file
func SaveScreenshot(filename string, x, y, width, height int) error {
	data, err := CaptureScreen(x, y, width, height)
	if err != nil {
		return err
	}

	// Create image from BGRA data
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for i := 0; i < len(data); i += 4 {
		pixelIdx := i / 4
		py := pixelIdx / width
		px := pixelIdx % width

		// Convert BGRA to RGBA
		img.SetRGBA(px, py, color.RGBA{
			R: data[i+2],
			G: data[i+1],
			B: data[i],
			A: 255,
		})
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}

// EnumPixels calls the callback function for each pixel in the region
func EnumPixels(x1, y1, x2, y2 int, callback func(x, y int, color utils.Color) bool) error {
	width := x2 - x1 + 1
	height := y2 - y1 + 1

	data, err := CaptureScreen(x1, y1, width, height)
	if err != nil {
		return err
	}

	for py := 0; py < height; py++ {
		for px := 0; px < width; px++ {
			idx := (py*width + px) * 4
			color := utils.Color{
				B: data[idx],
				G: data[idx+1],
				R: data[idx+2],
				A: 255,
			}
			if !callback(x1+px, y1+py, color) {
				return nil // Callback requested stop
			}
		}
	}

	return nil
}

func colorMatches(r1, g1, b1, r2, g2, b2 uint8, tolerance int) bool {
	dr := abs(int(r1) - int(r2))
	dg := abs(int(g1) - int(g2))
	db := abs(int(b1) - int(b2))
	return dr <= tolerance && dg <= tolerance && db <= tolerance
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
