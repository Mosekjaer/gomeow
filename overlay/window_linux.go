//go:build linux

package overlay

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// WindowInfo holds window position and size
type WindowInfo struct {
	X      int
	Y      int
	Width  int
	Height int
}

// GetWindowInfo returns information about a window by name using xwininfo
func GetWindowInfo(name string) (*WindowInfo, error) {
	cmd := exec.Command("xwininfo", "-name", name)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("xwininfo failed (is it installed?): %v", err)
	}

	info := &WindowInfo{}
	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "error") {
			return nil, fmt.Errorf("window '%s' not found", name)
		}

		if strings.Contains(line, "te upper-left X:") {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				info.X, _ = strconv.Atoi(parts[len(parts)-1])
			}
		} else if strings.Contains(line, "te upper-left Y:") {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				info.Y, _ = strconv.Atoi(parts[len(parts)-1])
			}
		} else if strings.Contains(line, "Width:") {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				info.Width, _ = strconv.Atoi(parts[len(parts)-1])
			}
		} else if strings.Contains(line, "Height:") {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				info.Height, _ = strconv.Atoi(parts[len(parts)-1])
			}
		}
	}

	return info, nil
}
