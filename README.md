# goMeow

A Go library for memory manipulation and overlay rendering. Port of [pyMeow](https://github.com/qb-0/pyMeow) (Nim/Python) to Go.

## Features

- **Process Management**: Enumerate, open, and inspect processes and modules
- **Memory Operations**: Read/write memory, pattern scanning, pointer chains
- **Overlay Rendering**: Transparent overlay windows with raylib
- **Input Handling**: Keyboard/mouse detection and simulation
- **Vector Math**: 2D and 3D vector operations
- **Cross-Platform**: Windows and Linux support

## Installation

```bash
go get github.com/mosekjaer/gomeow
```

### Dependencies

- [raylib-go](https://github.com/gen2brain/raylib-go) - For overlay rendering
- [golang.org/x/sys](https://pkg.go.dev/golang.org/x/sys) - For Windows API access

On Linux, you'll also need:
```bash
# For raylib
sudo apt install libgl1-mesa-dev libxi-dev libxcursor-dev libxrandr-dev libxinerama-dev

# For input simulation
sudo apt install libxtst-dev
```

## Quick Start

### Process and Memory

```go
package main

import (
    "fmt"
    "log"

    "gomeow/memory"
    "gomeow/process"
)

func main() {
    // Open a process
    proc, err := process.OpenProcessByName("target.exe")
    if err != nil {
        log.Fatal(err)
    }
    defer proc.Close()

    fmt.Printf("Opened %s (PID: %d)\n", proc.Name, proc.PID)
    fmt.Printf("Base: 0x%X\n", proc.Base)

    // List modules
    modules, _ := proc.EnumModules()
    for _, m := range modules {
        fmt.Printf("  %s @ 0x%X\n", m.Name, m.Base)
    }

    // Read memory
    health, err := memory.ReadInt32(proc, proc.Base+0x1234)
    if err != nil {
        log.Printf("Read failed: %v", err)
    }
    fmt.Printf("Health: %d\n", health)

    // Write memory
    err = memory.WriteInt32(proc, proc.Base+0x1234, 100)

    // Pattern scan
    results, err := memory.AOBScanModule(proc, "target.exe", "48 8B ?? 90 48", false, true)
    if err == nil && len(results) > 0 {
        fmt.Printf("Pattern found at: 0x%X\n", results[0])
    }

    // Pointer chain
    addr, err := memory.PointerChain64(proc, proc.Base+0x1000, 0x10, 0x20, 0x8)
    if err == nil {
        value, _ := memory.ReadFloat32(proc, addr)
        fmt.Printf("Value at chain: %f\n", value)
    }
}
```

### Overlay Rendering

```go
package main

import (
    "gomeow/overlay"
    "gomeow/utils"
)

func main() {
    // Initialize overlay
    err := overlay.Init("Target Window", 60, "My Overlay", -1, true)
    if err != nil {
        panic(err)
    }

    // Main loop
    for overlay.Loop() {
        overlay.BeginDrawing()

        // Draw shapes
        overlay.DrawText("ESP Overlay", 10, 10, 20, utils.White)
        overlay.DrawRectangleLines(100, 100, 200, 300, utils.Red, 2)
        overlay.DrawCircle(200, 250, 5, utils.Green)
        overlay.DrawLine(100, 100, 300, 400, utils.Blue, 1)

        // ESP helpers
        overlay.DrawCornerBox(150, 150, 100, 200, 15, utils.Yellow, 2)
        overlay.DrawHealthBar(150, 355, 100, 5, 75, 100, utils.DarkGray, utils.Green)

        overlay.DrawFPS(10, 40)

        overlay.EndDrawing()
    }

    overlay.Close()
}
```

### Input Detection

```go
package main

import (
    "fmt"
    "time"

    "gomeow/input"
)

func main() {
    fmt.Println("Press keys or move mouse...")

    for {
        // Check key state (Windows VK codes)
        if input.KeyPressed(0x01) { // VK_LBUTTON
            fmt.Println("Left mouse button pressed")
        }

        if input.KeyPressed(0x1B) { // VK_ESCAPE
            fmt.Println("Escape pressed, exiting")
            break
        }

        // Get mouse position
        pos := input.MousePosition()
        fmt.Printf("Mouse: %.0f, %.0f\r", pos.X, pos.Y)

        time.Sleep(16 * time.Millisecond)
    }
}
```

### Vector Math

```go
package main

import (
    "fmt"

    "gomeow/vec"
)

func main() {
    // 3D vectors
    playerPos := vec.NewVec3(100, 50, 200)
    enemyPos := vec.NewVec3(150, 50, 250)

    distance := playerPos.Distance(enemyPos)
    fmt.Printf("Distance: %.2f\n", distance)

    direction := enemyPos.Subtract(playerPos).Normalize()
    fmt.Printf("Direction: %+v\n", direction)

    // 2D vectors
    screenPos := vec.NewVec2(960, 540)
    target := vec.NewVec2(1000, 500)

    delta := target.Subtract(screenPos)
    fmt.Printf("Delta: %+v, Length: %.2f\n", delta, delta.Length())
}
```

### World to Screen

```go
package main

import (
    "gomeow/utils"
    "gomeow/vec"
)

func main() {
    // Read view matrix from game (16 floats)
    var viewMatrix [16]float32
    // ... read from memory ...

    worldPos := vec.NewVec3(1000, 50, 2000)
    screenWidth := 1920
    screenHeight := 1080

    screenPos, visible := utils.WorldToScreen(worldPos, viewMatrix, screenWidth, screenHeight)
    if visible {
        // Draw at screenPos.X, screenPos.Y
    }
}
```

## API Reference

### Process Package

| Function | Description |
|----------|-------------|
| `EnumProcesses()` | List all running processes |
| `OpenProcessByName(name)` | Open process by name |
| `OpenProcessByPID(pid)` | Open process by PID |
| `proc.Close()` | Close process handle |
| `proc.EnumModules()` | List loaded modules |
| `proc.GetModule(name)` | Get module by name |
| `proc.Is64Bit()` | Check if 64-bit process |
| `proc.IsRunning()` | Check if still running |
| `proc.GetPath()` | Get executable path |

### Memory Package

| Function | Description |
|----------|-------------|
| `Read(proc, addr, buf)` | Read raw bytes |
| `Write(proc, addr, buf)` | Write raw bytes |
| `ReadInt32(proc, addr)` | Read signed 32-bit int |
| `WriteInt32(proc, addr, val)` | Write signed 32-bit int |
| `ReadFloat32(proc, addr)` | Read 32-bit float |
| `ReadString(proc, addr, max)` | Read null-terminated string |
| `ReadVec3(proc, addr)` | Read Vec3 (3 floats) |
| `PointerChain64(proc, base, offsets...)` | Follow pointer chain |
| `AOBScanModule(proc, mod, pattern, rel, single)` | Pattern scan module |
| `AOBScanRange(proc, pattern, start, end, rel, single)` | Pattern scan range |
| `AllocateMemory(proc, size, prot)` | Allocate in target (Windows) |
| `InjectLibrary(proc, dllPath)` | Inject DLL (Windows) |

### Overlay Package

| Function | Description |
|----------|-------------|
| `Init(target, fps, title, exitKey, track)` | Initialize overlay |
| `Loop()` | Check if should continue |
| `BeginDrawing()` / `EndDrawing()` | Frame boundaries |
| `Close()` | Close overlay |
| `DrawText(text, x, y, size, color)` | Draw text |
| `DrawLine(x1, y1, x2, y2, color, thick)` | Draw line |
| `DrawRectangle(x, y, w, h, color)` | Draw filled rectangle |
| `DrawRectangleLines(x, y, w, h, color, thick)` | Draw rectangle outline |
| `DrawCircle(x, y, radius, color)` | Draw filled circle |
| `DrawCornerBox(x, y, w, h, len, color, thick)` | Draw corner-style box |
| `DrawHealthBar(x, y, w, h, hp, max, bg, fg)` | Draw health bar |

### Input Package

| Function | Description |
|----------|-------------|
| `KeyPressed(vkey)` | Check if key is pressed |
| `MousePressed(button)` | Check mouse button ("left", "right", "middle") |
| `MousePosition()` | Get cursor position as Vec2 |
| `MouseMove(x, y, relative)` | Move cursor |
| `MouseClick(button)` | Click mouse button |
| `PressKey(vkey)` | Press and release key |

### Vec Package

| Method | Description |
|--------|-------------|
| `Add(other)` | Add vectors |
| `Subtract(other)` | Subtract vectors |
| `Scale(value)` | Multiply by scalar |
| `Length()` | Get magnitude |
| `Distance(other)` | Distance between points |
| `Normalize()` | Get unit vector |
| `Dot(other)` | Dot product |
| `Cross(other)` | Cross product (Vec3 only) |
| `Closest(vectors...)` | Find closest vector |

## Platform Notes

### Windows
- Requires administrator privileges for some operations
- Uses Windows API (kernel32, user32)
- Full feature support

### Linux
- Requires root for memory operations (`process_vm_readv`)
- Uses X11 for input (requires XTest extension)
- DLL injection not available (Windows-only concept)

## License

MIT License - See LICENSE file for details.

## Credits

- Original [pyMeow](https://github.com/qb-0/pyMeow) by qb-0
- [raylib-go](https://github.com/gen2brain/raylib-go) for graphics
