//go:build linux

package pixel

/*
#cgo LDFLAGS: -lX11

#include <X11/Xlib.h>
#include <X11/Xutil.h>
#include <stdlib.h>

static Display* display = NULL;
static Window root;
static int screen;

void init_display() {
    if (display == NULL) {
        display = XOpenDisplay(NULL);
        if (display != NULL) {
            screen = DefaultScreen(display);
            root = RootWindow(display, screen);
        }
    }
}

int get_screen_width() {
    init_display();
    if (display == NULL) return 0;
    return DisplayWidth(display, screen);
}

int get_screen_height() {
    init_display();
    if (display == NULL) return 0;
    return DisplayHeight(display, screen);
}

unsigned long get_pixel(int x, int y) {
    init_display();
    if (display == NULL) return 0;

    XImage* image = XGetImage(display, root, x, y, 1, 1, AllPlanes, ZPixmap);
    if (image == NULL) return 0;

    unsigned long pixel = XGetPixel(image, 0, 0);
    XDestroyImage(image);
    return pixel;
}

// Returns pointer to pixel data, caller must free
unsigned char* capture_region(int x, int y, int width, int height) {
    init_display();
    if (display == NULL) return NULL;

    XImage* image = XGetImage(display, root, x, y, width, height, AllPlanes, ZPixmap);
    if (image == NULL) return NULL;

    int size = width * height * 4;
    unsigned char* data = (unsigned char*)malloc(size);
    if (data == NULL) {
        XDestroyImage(image);
        return NULL;
    }

    for (int py = 0; py < height; py++) {
        for (int px = 0; px < width; px++) {
            unsigned long pixel = XGetPixel(image, px, py);
            int idx = (py * width + px) * 4;
            data[idx] = (pixel >> 16) & 0xFF;     // R
            data[idx + 1] = (pixel >> 8) & 0xFF;  // G
            data[idx + 2] = pixel & 0xFF;         // B
            data[idx + 3] = 255;                   // A
        }
    }

    XDestroyImage(image);
    return data;
}
*/
import "C"
import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"unsafe"

	"gomeow/utils"
)

// GetScreenSize returns the screen dimensions
func GetScreenSize() (width, height int) {
	return int(C.get_screen_width()), int(C.get_screen_height())
}

// GetPixelColor gets the color of a pixel at the specified screen coordinates
func GetPixelColor(x, y int) utils.Color {
	pixel := C.get_pixel(C.int(x), C.int(y))
	return utils.Color{
		R: uint8((pixel >> 16) & 0xFF),
		G: uint8((pixel >> 8) & 0xFF),
		B: uint8(pixel & 0xFF),
		A: 255,
	}
}

// PixelSearch searches for a pixel with the specified color in a region
func PixelSearch(x1, y1, x2, y2 int, color utils.Color, tolerance int) (int, int, bool) {
	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			c := GetPixelColor(x, y)
			if colorMatches(c.R, c.G, c.B, color.R, color.G, color.B, tolerance) {
				return x, y, true
			}
		}
	}
	return -1, -1, false
}

// PixelSearchAll finds all pixels matching the color in a region
func PixelSearchAll(x1, y1, x2, y2 int, color utils.Color, tolerance int) []image.Point {
	var results []image.Point

	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			c := GetPixelColor(x, y)
			if colorMatches(c.R, c.G, c.B, color.R, color.G, color.B, tolerance) {
				results = append(results, image.Point{X: x, Y: y})
			}
		}
	}

	return results
}

// CaptureScreen captures a region of the screen and returns raw pixel data (RGBA)
func CaptureScreen(x, y, width, height int) ([]byte, error) {
	cData := C.capture_region(C.int(x), C.int(y), C.int(width), C.int(height))
	if cData == nil {
		return nil, fmt.Errorf("failed to capture screen region")
	}
	defer C.free(unsafe.Pointer(cData))

	size := width * height * 4
	data := make([]byte, size)
	copy(data, (*[1 << 30]byte)(unsafe.Pointer(cData))[:size:size])

	return data, nil
}

// SaveScreenshot captures a region and saves it as a PNG file
func SaveScreenshot(filename string, x, y, width, height int) error {
	data, err := CaptureScreen(x, y, width, height)
	if err != nil {
		return err
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for i := 0; i < len(data); i += 4 {
		pixelIdx := i / 4
		py := pixelIdx / width
		px := pixelIdx % width

		img.SetRGBA(px, py, color.RGBA{
			R: data[i],
			G: data[i+1],
			B: data[i+2],
			A: data[i+3],
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
				R: data[idx],
				G: data[idx+1],
				B: data[idx+2],
				A: data[idx+3],
			}
			if !callback(x1+px, y1+py, color) {
				return nil
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
