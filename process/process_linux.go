//go:build linux

package process

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

// EnumProcesses returns a list of all running processes
func EnumProcesses() ([]ProcessInfo, error) {
	var processes []ProcessInfo

	files, err := os.ReadDir("/proc")
	if err != nil {
		return nil, fmt.Errorf("failed to read /proc: %v", err)
	}

	for _, f := range files {
		if !f.IsDir() {
			continue
		}

		pid, err := strconv.Atoi(f.Name())
		if err != nil {
			continue // Not a PID directory
		}

		commPath := filepath.Join("/proc", f.Name(), "comm")
		name, err := os.ReadFile(commPath)
		if err != nil {
			continue
		}

		processes = append(processes, ProcessInfo{
			Name: strings.TrimSpace(string(name)),
			PID:  pid,
		})
	}

	return processes, nil
}

// PIDExists checks if a process with the given PID exists
func PIDExists(pid int) bool {
	_, err := os.Stat(fmt.Sprintf("/proc/%d", pid))
	return err == nil
}

// ProcessExists checks if a process with the given name exists
func ProcessExists(name string) bool {
	processes, err := EnumProcesses()
	if err != nil {
		return false
	}
	for _, p := range processes {
		if strings.Contains(p.Name, name) {
			return true
		}
	}
	return false
}

// GetProcessID returns the PID of a process by name
func GetProcessID(name string) (int, error) {
	processes, err := EnumProcesses()
	if err != nil {
		return 0, err
	}
	for _, p := range processes {
		if strings.Contains(p.Name, name) {
			return p.PID, nil
		}
	}
	return 0, fmt.Errorf("process '%s' not found", name)
}

// GetProcessName returns the name of a process by PID
func GetProcessName(pid int) (string, error) {
	commPath := fmt.Sprintf("/proc/%d/comm", pid)
	name, err := os.ReadFile(commPath)
	if err != nil {
		return "", fmt.Errorf("process with PID %d not found", pid)
	}
	return strings.TrimSpace(string(name)), nil
}

// OpenProcess opens a process by name or PID
func OpenProcess(identifier interface{}) (*Process, error) {
	// Check for root privileges
	if os.Getuid() != 0 {
		return nil, fmt.Errorf("root access required")
	}

	var pid int
	var name string

	switch v := identifier.(type) {
	case int:
		pid = v
		if !PIDExists(pid) {
			return nil, fmt.Errorf("process with PID %d does not exist", pid)
		}
		var err error
		name, err = GetProcessName(pid)
		if err != nil {
			return nil, err
		}
	case string:
		var err error
		pid, err = GetProcessID(v)
		if err != nil {
			return nil, err
		}
		name = v
	default:
		return nil, fmt.Errorf("identifier must be int (PID) or string (name)")
	}

	proc := &Process{
		Name: name,
		PID:  pid,
	}

	// Get base address from first module
	modules, err := proc.EnumModules()
	if err == nil && len(modules) > 0 {
		proc.Base = modules[0].Base
	}

	return proc, nil
}

// OpenProcessByPID opens a process by its PID
func OpenProcessByPID(pid int) (*Process, error) {
	return OpenProcess(pid)
}

// OpenProcessByName opens a process by its name
func OpenProcessByName(name string) (*Process, error) {
	return OpenProcess(name)
}

// Close closes the process (no-op on Linux)
func (p *Process) Close() error {
	return nil
}

// Handle returns the process handle (PID on Linux)
func (p *Process) Handle() uintptr {
	return uintptr(p.PID)
}

// IsRunning checks if the process is still running
func (p *Process) IsRunning() bool {
	return syscall.Kill(p.PID, 0) == nil
}

// Is64Bit checks if the process is 64-bit
func (p *Process) Is64Bit() bool {
	exePath := fmt.Sprintf("/proc/%d/exe", p.PID)
	f, err := os.Open(exePath)
	if err != nil {
		return false
	}
	defer f.Close()

	// Read ELF header (first 5 bytes)
	header := make([]byte, 5)
	_, err = f.Read(header)
	if err != nil {
		return false
	}

	// ELF magic number: 0x7F 'E' 'L' 'F'
	// 5th byte: 1 = 32-bit, 2 = 64-bit
	if header[0] == 0x7F && header[1] == 'E' && header[2] == 'L' && header[3] == 'F' {
		return header[4] == 2
	}

	return false
}

// GetPath returns the full path to the process executable
func (p *Process) GetPath() (string, error) {
	exePath := fmt.Sprintf("/proc/%d/exe", p.PID)
	path, err := os.Readlink(exePath)
	if err != nil {
		return "", fmt.Errorf("failed to read exe path: %v", err)
	}
	return path, nil
}

// EnumModules returns a list of all modules loaded in the process
func (p *Process) EnumModules() ([]Module, error) {
	mapsPath := fmt.Sprintf("/proc/%d/maps", p.PID)
	f, err := os.Open(mapsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open maps: %v", err)
	}
	defer f.Close()

	moduleMap := make(map[string]*Module)
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) < 6 {
			continue
		}

		// Parse address range
		addrParts := strings.Split(parts[0], "-")
		if len(addrParts) != 2 {
			continue
		}

		startAddr, err := strconv.ParseUint(addrParts[0], 16, 64)
		if err != nil {
			continue
		}
		endAddr, err := strconv.ParseUint(addrParts[1], 16, 64)
		if err != nil {
			continue
		}

		// Get module name (last part, if it's a path)
		name := parts[len(parts)-1]
		if !strings.HasPrefix(name, "/") && !strings.HasPrefix(name, "[") {
			continue
		}

		// Extract just the filename
		if strings.HasPrefix(name, "/") {
			name = filepath.Base(name)
		}

		if existing, ok := moduleMap[name]; ok {
			// Extend existing module
			if uintptr(startAddr) < existing.Base {
				existing.Base = uintptr(startAddr)
			}
			if uintptr(endAddr) > existing.End {
				existing.End = uintptr(endAddr)
			}
			existing.Size = existing.End - existing.Base
		} else {
			// New module
			moduleMap[name] = &Module{
				Name: name,
				Base: uintptr(startAddr),
				End:  uintptr(endAddr),
				Size: uintptr(endAddr - startAddr),
			}
		}
	}

	var modules []Module
	for _, m := range moduleMap {
		modules = append(modules, *m)
	}

	return modules, nil
}

// GetModule returns a module by name
func (p *Process) GetModule(name string) (*Module, error) {
	modules, err := p.EnumModules()
	if err != nil {
		return nil, err
	}
	for _, m := range modules {
		if m.Name == name {
			return &m, nil
		}
	}
	return nil, fmt.Errorf("module '%s' not found", name)
}

// ModuleExists checks if a module is loaded
func (p *Process) ModuleExists(name string) bool {
	_, err := p.GetModule(name)
	return err == nil
}

// EnumMemoryRegions returns memory regions for a module
func (p *Process) EnumMemoryRegions(module *Module) ([]Page, error) {
	mapsPath := fmt.Sprintf("/proc/%d/maps", p.PID)
	f, err := os.Open(mapsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open maps: %v", err)
	}
	defer f.Close()

	var pages []Page
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, module.Name) {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 1 {
			continue
		}

		addrParts := strings.Split(parts[0], "-")
		if len(addrParts) != 2 {
			continue
		}

		startAddr, err := strconv.ParseUint(addrParts[0], 16, 64)
		if err != nil {
			continue
		}
		endAddr, err := strconv.ParseUint(addrParts[1], 16, 64)
		if err != nil {
			continue
		}

		pages = append(pages, Page{
			Start: uintptr(startAddr),
			End:   uintptr(endAddr),
			Size:  uintptr(endAddr - startAddr),
		})
	}

	return pages, nil
}
