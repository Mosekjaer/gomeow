# pyMeow vs goMeow API Comparison

This document provides a comparison between the original pyMeow (Nim/Python) library and the Go port (goMeow).

## Overview

| Feature | pyMeow (Python) | goMeow (Go) |
|---------|-----------------|-------------|
| Language | Nim → Python Extension | Go |
| Graphics | Raylib | Raylib (raylib-go) |
| Platform | Windows, Linux | Windows, Linux |
| Memory Access | Native | Native (syscall) |

## Process Package

### Python (pyMeow)
```python
import pyMeow as pm

# List processes
processes = pm.enum_processes()

# Open process
proc = pm.open_process("game.exe")
proc = pm.open_process(1234)  # by PID

# Check process
pm.process_exists("game.exe")
pm.pid_exists(1234)

# Get info
pid = pm.get_process_id("game.exe")
name = pm.get_process_name(1234)

# Modules
modules = pm.enum_modules(proc)
module = pm.get_module(proc, "module.dll")
pm.module_exists(proc, "module.dll")

# Close
pm.close_process(proc)
```

### Go (goMeow)
```go
import "gomeow/process"

// List processes
processes, err := process.EnumProcesses()

// Open process
proc, err := process.OpenProcess("game.exe")
proc, err := process.OpenProcessByPID(1234)

// Check process
process.ProcessExists("game.exe")
process.PIDExists(1234)

// Get info
pid, err := process.GetProcessID("game.exe")
name, err := process.GetProcessName(1234)

// Modules
modules, err := proc.EnumModules()
module, err := proc.GetModule("module.dll")
proc.ModuleExists("module.dll")

// Close
proc.Close()
```

## Memory Package

### Python (pyMeow)
```python
# Read operations
value = pm.r_int(proc, address)
value = pm.r_int64(proc, address)
value = pm.r_float(proc, address)
value = pm.r_float64(proc, address)
value = pm.r_bool(proc, address)
value = pm.r_string(proc, address, size)
value = pm.r_bytes(proc, address, size)
vec2 = pm.r_vec2(proc, address)
vec3 = pm.r_vec3(proc, address)

# Write operations
pm.w_int(proc, address, value)
pm.w_float(proc, address, value)
pm.w_bool(proc, address, value)
pm.w_string(proc, address, value)
pm.w_bytes(proc, address, value)

# Pointer chains
addr = pm.pointer_chain_32(proc, base, [offset1, offset2])
addr = pm.pointer_chain_64(proc, base, [offset1, offset2])

# Pattern scanning
addr = pm.aob_scan_module(proc, module, "AA BB ?? CC")
addrs = pm.aob_scan_module_all(proc, module, "AA BB ?? CC")
```

### Go (goMeow)
```go
import "gomeow/memory"

// Read operations
value, err := memory.ReadInt32(proc, address)
value, err := memory.ReadInt64(proc, address)
value, err := memory.ReadFloat32(proc, address)
value, err := memory.ReadFloat64(proc, address)
value, err := memory.ReadBool(proc, address)
value, err := memory.ReadString(proc, address, size)
value, err := memory.ReadBytes(proc, address, size)
vec2, err := memory.ReadVec2(proc, address)
vec3, err := memory.ReadVec3(proc, address)

// Write operations
err := memory.WriteInt32(proc, address, value)
err := memory.WriteFloat32(proc, address, value)
err := memory.WriteBool(proc, address, value)
err := memory.WriteString(proc, address, value)
err := memory.WriteBytes(proc, address, value)

// Pointer chains
addr, err := memory.PointerChain32(proc, base, []uintptr{offset1, offset2})
addr, err := memory.PointerChain64(proc, base, []uintptr{offset1, offset2})

// Pattern scanning
addr, err := memory.AOBScanModule(proc, module, "AA BB ?? CC")
addrs, err := memory.AOBScanModuleAll(proc, module, "AA BB ?? CC")
```

## Overlay Package

### Python (pyMeow)
```python
# Initialize
pm.overlay_init(target="game window", fps=60)
pm.overlay_init()  # fullscreen

# Main loop
while pm.overlay_loop():
    pm.begin_drawing()
    # draw stuff
    pm.end_drawing()

pm.overlay_close()

# Window operations
pm.toggle_mouse()
pm.set_fps(144)
w, h = pm.get_screen_width(), pm.get_screen_height()
```

### Go (goMeow)
```go
import "gomeow/overlay"

// Initialize
overlay.Init("game window", 60, "My Overlay", -1, false)
overlay.Init("Full", 60, "My Overlay", -1, false)  // fullscreen

// Main loop
for overlay.Loop() {
    overlay.BeginDrawing()
    // draw stuff
    overlay.EndDrawing()
}

overlay.Close()

// Window operations
overlay.ToggleMouse()
overlay.SetFPS(144)
w, h := overlay.GetScreenWidth(), overlay.GetScreenHeight()
```

## Drawing Functions

### Python (pyMeow)
```python
# Primitives
pm.draw_pixel(x, y, color)
pm.draw_line(x1, y1, x2, y2, color, thickness)
pm.draw_circle(cx, cy, radius, color)
pm.draw_circle_lines(cx, cy, radius, color)
pm.draw_rectangle(x, y, w, h, color)
pm.draw_rectangle_lines(x, y, w, h, color, thickness)
pm.draw_rectangle_rounded(x, y, w, h, roundness, segments, color)

# Text
pm.draw_text(text, x, y, size, color)
pm.draw_fps(x, y)

# ESP Helpers
pm.draw_box(x, y, w, h, fill_color, outline_color, thickness)
pm.draw_corner_box(x, y, w, h, corner_length, color, thickness)
```

### Go (goMeow)
```go
import "gomeow/overlay"

// Primitives
overlay.DrawPixel(x, y, color)
overlay.DrawLine(x1, y1, x2, y2, color, thickness)
overlay.DrawCircle(cx, cy, radius, color)
overlay.DrawCircleLines(cx, cy, radius, color)
overlay.DrawRectangle(x, y, w, h, color)
overlay.DrawRectangleLines(x, y, w, h, color, thickness)
overlay.DrawRectangleRounded(x, y, w, h, roundness, segments, color)

// Text
overlay.DrawText(text, x, y, size, color)
overlay.DrawFPS(x, y)

// ESP Helpers
overlay.DrawBox(x, y, w, h, fillColor, outlineColor, thickness)
overlay.DrawCornerBox(x, y, w, h, cornerLength, color, thickness)
```

## Input Package

### Python (pyMeow)
```python
# Keyboard
pm.key_pressed(key_code)
pm.press_key(key_code)

# Mouse
pm.mouse_pressed(button)
pm.mouse_move(x, y)
pm.mouse_click(button)
pm.get_mouse_position()
```

### Go (goMeow)
```go
import "gomeow/input"

// Keyboard
input.KeyPressed(keyCode)
input.PressKey(keyCode)

// Mouse
input.MousePressed(button)
input.MouseMove(x, y)
input.MouseClick(button)
input.MousePosition()
```

## Pixel Package

### Python (pyMeow)
```python
# Screen capture
w, h = pm.get_screen_size()
color = pm.get_pixel(x, y)
pm.pixel_search(color, x1, y1, x2, y2)
pm.save_screenshot("file.png", x, y, w, h)
```

### Go (goMeow)
```go
import "gomeow/pixel"

// Screen capture
w, h := pixel.GetScreenSize()
color, err := pixel.GetPixelColor(x, y)
result, found := pixel.PixelSearch(color, x1, y1, x2, y2)
err := pixel.SaveScreenshot("file.png", x, y, w, h)
```

## Vector Math

### Python (pyMeow)
```python
# Vec2
v = pm.vec2(x, y)
v.x, v.y
v.add(other)
v.subtract(other)
v.multiply(scalar)
v.length()
v.distance(other)
v.normalize()

# Vec3
v = pm.vec3(x, y, z)
v.cross(other)
```

### Go (goMeow)
```go
import "gomeow/vec"

// Vec2
v := vec.Vec2{X: x, Y: y}
v.X, v.Y
v.Add(other)
v.Subtract(other)
v.Scale(scalar)
v.Length()
v.Distance(other)
v.Normalize()

// Vec3
v := vec.Vec3{X: x, Y: y, Z: z}
v.Cross(other)
```

## Colors

### Python (pyMeow)
```python
# Predefined colors
pm.RED, pm.GREEN, pm.BLUE, pm.WHITE, pm.BLACK

# Custom colors
color = pm.new_color(r, g, b, a)
color = pm.fade(color, alpha)
```

### Go (goMeow)
```go
import "gomeow/utils"

// Predefined colors
utils.Red, utils.Green, utils.Blue, utils.White, utils.Black

// Custom colors
color := utils.NewColorAlpha(r, g, b, a)
color := utils.Fade(color, alpha)
```

## GUI Package (goMeow exclusive)

goMeow includes a comprehensive immediate-mode GUI system not present in pyMeow:

```go
import "gomeow/gui"

// Frame management
gui.Begin()
defer gui.End()

// Components
clicked := gui.Button(id, x, y, w, h, "Click Me")
checked := gui.Checkbox(id, x, y, "Enable", checked)
value := gui.Slider(id, x, y, w, h, value, min, max)
text := gui.TextBox(id, x, y, w, h, text, maxLen)
selected := gui.Dropdown(id, x, y, w, h, items, selected)
selected := gui.ListView(id, x, y, w, h, items, selected)
color := gui.ColorPicker(id, x, y, w, h, color)

// Containers
x, y, closed := gui.Window(id, x, y, w, h, "Title", closeable)
gui.Panel(x, y, w, h)
scrollY := gui.BeginScrollPanel(id, x, y, w, h, contentH)
gui.EndScrollPanel()

// Tabs
selected := gui.TabBar(id, x, y, w, h, tabs, selected)

// Indicators
gui.ProgressBar(x, y, w, h, value, min, max)
gui.HealthBar(x, y, w, h, health, maxHealth)
gui.LoadingSpinner(x, y, radius)

// Tooltips
gui.Tooltip(id, "Help text")
gui.HelpMarker(id, x, y, "Help text")
```

## Key Differences

1. **Error Handling**: Go uses explicit error returns, Python uses exceptions
2. **Type Safety**: Go is statically typed, pyMeow uses dynamic typing
3. **Method Style**: Go uses methods on structs (e.g., `proc.EnumModules()`), Python uses functions with process parameter
4. **Memory Management**: Go has garbage collection, both handle native resources automatically
5. **GUI Package**: goMeow includes a built-in immediate-mode GUI system
6. **Compilation**: Go compiles to native binary, pyMeow requires Python interpreter

## Migration Tips

1. Replace Python snake_case with Go CamelCase
2. Add error handling for all operations that can fail
3. Use explicit type conversions (Go is strict about types)
4. Initialize structs with `Type{field: value}` syntax
5. Use the gui package for in-overlay menus (new feature in goMeow)
6. Colors require `utils.` prefix instead of `pm.` prefix
