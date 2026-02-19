package overlay

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"gomeow/vec"
)

// OverlayOptions holds overlay configuration
type OverlayOptions struct {
	ExitKey      int32
	Target       string
	TrackTarget  bool
	TargetX      int
	TargetY      int
	TargetWidth  int
	TargetHeight int
}

var opts OverlayOptions

// Init initializes the overlay window
// target: window name to overlay on, or "Full" for fullscreen
// fps: target framerate (0 for unlimited)
// title: window title
// exitKey: key to close overlay (-1 for default End key)
// trackTarget: whether to follow target window movement
func Init(target string, fps int, title string, exitKey int, trackTarget bool) error {
	var width, height int
	var targetX, targetY int

	// If targeting a specific window, get its info first (before creating our window)
	if target != "Full" {
		winInfo, err := GetWindowInfo(target)
		if err != nil {
			return err
		}

		// Validate dimensions
		if winInfo.Width <= 0 || winInfo.Height <= 0 {
			return fmt.Errorf("target window '%s' has invalid dimensions: %dx%d", target, winInfo.Width, winInfo.Height)
		}

		width = winInfo.Width
		height = winInfo.Height
		targetX = winInfo.X
		targetY = winInfo.Y

		opts.TargetX = winInfo.X
		opts.TargetY = winInfo.Y
		opts.TargetWidth = winInfo.Width
		opts.TargetHeight = winInfo.Height
	} else {
		// For fullscreen, we'll get monitor size after setting config flags
		width = 1920  // Default, will be updated
		height = 1080
	}

	// Suppress raylib logging
	rl.SetTraceLogLevel(rl.LogNone)

	// Set all config flags at once (they are OR'd together)
	flags := uint32(rl.FlagWindowUndecorated | rl.FlagWindowMousePassthrough | rl.FlagWindowTransparent | rl.FlagWindowTopmost | rl.FlagMsaa4xHint)
	rl.SetConfigFlags(flags)

	// Initialize window
	rl.InitWindow(int32(width), int32(height), title)

	// For fullscreen, get actual monitor size now that window exists
	if target == "Full" {
		monitor := rl.GetCurrentMonitor()
		width = int(rl.GetMonitorWidth(monitor))
		height = int(rl.GetMonitorHeight(monitor))
		if width > 0 && height > 0 {
			rl.SetWindowSize(width, height)
		}
	} else {
		// Position overlay on target window
		rl.SetWindowPosition(targetX, targetY)
	}

	// Set FPS after window is created
	if fps > 0 {
		rl.SetTargetFPS(int32(fps))
	}

	// Store our overlay window handle for focus detection
	// Find our own window by title
	if hwnd, err := findWindowByTitle(title); err == nil {
		SetOverlayHandle(hwnd)
	}

	opts.Target = target
	opts.TrackTarget = trackTarget

	if exitKey != -1 {
		opts.ExitKey = int32(exitKey)
	} else {
		opts.ExitKey = 0x23 // VK_END on Windows
	}

	rl.SetExitKey(rl.KeyNull)

	return nil
}

// InitSimple initializes a simple fullscreen overlay
func InitSimple(title string, fps int) error {
	return Init("Full", fps, title, -1, false)
}

// Loop checks if the overlay should continue running
// Returns false when window should close
func Loop() bool {
	rl.ClearBackground(rl.Blank)

	// Check exit key
	if int32(rl.GetKeyPressed()) == opts.ExitKey {
		Close()
		return false
	}

	// Track target window if enabled
	if opts.TrackTarget && opts.Target != "Full" {
		winInfo, err := GetWindowInfo(opts.Target)
		if err == nil {
			if winInfo.X != opts.TargetX || winInfo.Y != opts.TargetY {
				opts.TargetX = winInfo.X
				opts.TargetY = winInfo.Y
				rl.SetWindowPosition(winInfo.X, winInfo.Y)
			}
			if winInfo.Width != opts.TargetWidth || winInfo.Height != opts.TargetHeight {
				opts.TargetWidth = winInfo.Width
				opts.TargetHeight = winInfo.Height
				rl.SetWindowSize(winInfo.Width, winInfo.Height)
			}
		}
	}

	return !rl.WindowShouldClose()
}

// BeginDrawing begins a new drawing frame
func BeginDrawing() {
	rl.BeginDrawing()
}

// EndDrawing ends the drawing frame
func EndDrawing() {
	rl.EndDrawing()
}

// Close closes the overlay window
func Close() {
	rl.CloseWindow()
}

// GetFPS returns the current framerate
func GetFPS() int {
	return int(rl.GetFPS())
}

// SetFPS sets the target framerate
func SetFPS(fps int) {
	rl.SetTargetFPS(int32(fps))
}

// GetScreenWidth returns the overlay width
func GetScreenWidth() int {
	return int(rl.GetScreenWidth())
}

// GetScreenHeight returns the overlay height
func GetScreenHeight() int {
	return int(rl.GetScreenHeight())
}

// GetWindowPosition returns the overlay window position
func GetWindowPosition() vec.Vec2 {
	pos := rl.GetWindowPosition()
	return vec.Vec2{X: pos.X, Y: pos.Y}
}

// SetWindowPosition sets the overlay window position
func SetWindowPosition(x, y int) {
	rl.SetWindowPosition(x, y)
}

// SetWindowSize sets the overlay window size
func SetWindowSize(width, height int) {
	rl.SetWindowSize(width, height)
}

// SetWindowTitle sets the overlay window title
func SetWindowTitle(title string) {
	rl.SetWindowTitle(title)
}

// SetWindowIcon sets the overlay window icon
func SetWindowIcon(imagePath string) {
	img := rl.LoadImage(imagePath)
	rl.SetWindowIcon(*img)
	rl.UnloadImage(img)
}

// ToggleMouse toggles mouse passthrough
func ToggleMouse() {
	if rl.IsWindowState(rl.FlagWindowMousePassthrough) {
		rl.ClearWindowState(rl.FlagWindowMousePassthrough)
	} else {
		rl.SetWindowState(rl.FlagWindowMousePassthrough)
	}
}

// TakeScreenshot saves a screenshot to the given file
func TakeScreenshot(fileName string) {
	rl.TakeScreenshot(fileName)
}
