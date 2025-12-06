//go:build windows

package process

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	kernel32                     = windows.NewLazySystemDLL("kernel32.dll")
	procCreateToolhelp32Snapshot = kernel32.NewProc("CreateToolhelp32Snapshot")
	procProcess32FirstW          = kernel32.NewProc("Process32FirstW")
	procProcess32NextW           = kernel32.NewProc("Process32NextW")
	procModule32FirstW           = kernel32.NewProc("Module32FirstW")
	procModule32NextW            = kernel32.NewProc("Module32NextW")
	procIsWow64Process           = kernel32.NewProc("IsWow64Process")
	procQueryFullProcessImageNameW = kernel32.NewProc("QueryFullProcessImageNameW")
)

const (
	TH32CS_SNAPPROCESS = 0x00000002
	TH32CS_SNAPMODULE  = 0x00000008
	TH32CS_SNAPMODULE32 = 0x00000010
	MAX_PATH           = 260
	PROCESS_ALL_ACCESS = 0x1F0FFF
	STILL_ACTIVE       = 259
)

type processEntry32W struct {
	Size              uint32
	CntUsage          uint32
	ProcessID         uint32
	DefaultHeapID     uintptr
	ModuleID          uint32
	CntThreads        uint32
	ParentProcessID   uint32
	PriorityClassBase int32
	Flags             uint32
	ExeFile           [MAX_PATH]uint16
}

type moduleEntry32W struct {
	Size         uint32
	ModuleID     uint32
	ProcessID    uint32
	GlblcntUsage uint32
	ProccntUsage uint32
	ModBaseAddr  uintptr
	ModBaseSize  uint32
	HModule      uintptr
	Module       [256]uint16
	ExePath      [MAX_PATH]uint16
}

type memoryBasicInformation struct {
	BaseAddress       uintptr
	AllocationBase    uintptr
	AllocationProtect uint32
	RegionSize        uintptr
	State             uint32
	Protect           uint32
	Type              uint32
}

// EnumProcesses returns a list of all running processes
func EnumProcesses() ([]ProcessInfo, error) {
	snapshot, _, err := procCreateToolhelp32Snapshot.Call(TH32CS_SNAPPROCESS, 0)
	if snapshot == uintptr(syscall.InvalidHandle) {
		return nil, fmt.Errorf("CreateToolhelp32Snapshot failed: %v", err)
	}
	defer windows.CloseHandle(windows.Handle(snapshot))

	var processes []ProcessInfo
	var pe processEntry32W
	pe.Size = uint32(unsafe.Sizeof(pe))

	ret, _, _ := procProcess32FirstW.Call(snapshot, uintptr(unsafe.Pointer(&pe)))
	if ret == 0 {
		return nil, fmt.Errorf("Process32First failed")
	}

	for {
		name := windows.UTF16ToString(pe.ExeFile[:])
		processes = append(processes, ProcessInfo{
			Name: name,
			PID:  int(pe.ProcessID),
		})

		ret, _, _ = procProcess32NextW.Call(snapshot, uintptr(unsafe.Pointer(&pe)))
		if ret == 0 {
			break
		}
	}

	return processes, nil
}

// PIDExists checks if a process with the given PID exists
func PIDExists(pid int) bool {
	processes, err := EnumProcesses()
	if err != nil {
		return false
	}
	for _, p := range processes {
		if p.PID == pid {
			return true
		}
	}
	return false
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
	processes, err := EnumProcesses()
	if err != nil {
		return "", err
	}
	for _, p := range processes {
		if p.PID == pid {
			return p.Name, nil
		}
	}
	return "", fmt.Errorf("process with PID %d not found", pid)
}

// OpenProcess opens a process by name or PID
func OpenProcess(identifier interface{}) (*Process, error) {
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

	handle, err := windows.OpenProcess(PROCESS_ALL_ACCESS, false, uint32(pid))
	if err != nil {
		return nil, fmt.Errorf("unable to open process [PID: %d]: %v", pid, err)
	}

	proc := &Process{
		Name:   name,
		PID:    pid,
		handle: uintptr(handle),
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

// Close closes the process handle
func (p *Process) Close() error {
	if p.handle != 0 {
		return windows.CloseHandle(windows.Handle(p.handle))
	}
	return nil
}

// Handle returns the Windows process handle
func (p *Process) Handle() uintptr {
	return p.handle
}

// IsRunning checks if the process is still running
func (p *Process) IsRunning() bool {
	var exitCode uint32
	err := windows.GetExitCodeProcess(windows.Handle(p.handle), &exitCode)
	if err != nil {
		return false
	}
	return exitCode == STILL_ACTIVE
}

// Is64Bit checks if the process is 64-bit
func (p *Process) Is64Bit() bool {
	var isWow64 bool
	ret, _, _ := procIsWow64Process.Call(p.handle, uintptr(unsafe.Pointer(&isWow64)))
	if ret == 0 {
		return false
	}
	// If IsWow64Process returns true, it's a 32-bit process on 64-bit Windows
	// If false, it's either 64-bit on 64-bit Windows, or 32-bit on 32-bit Windows
	return !isWow64
}

// GetPath returns the full path to the process executable
func (p *Process) GetPath() (string, error) {
	var path [MAX_PATH * 2]uint16
	size := uint32(len(path))
	ret, _, err := procQueryFullProcessImageNameW.Call(
		p.handle,
		0,
		uintptr(unsafe.Pointer(&path[0])),
		uintptr(unsafe.Pointer(&size)),
	)
	if ret == 0 {
		return "", fmt.Errorf("QueryFullProcessImageNameW failed: %v", err)
	}
	return windows.UTF16ToString(path[:size]), nil
}

// EnumModules returns a list of all modules loaded in the process
func (p *Process) EnumModules() ([]Module, error) {
	snapshot, _, err := procCreateToolhelp32Snapshot.Call(
		TH32CS_SNAPMODULE|TH32CS_SNAPMODULE32,
		uintptr(p.PID),
	)
	if snapshot == uintptr(syscall.InvalidHandle) {
		return nil, fmt.Errorf("CreateToolhelp32Snapshot failed: %v", err)
	}
	defer windows.CloseHandle(windows.Handle(snapshot))

	var modules []Module
	var me moduleEntry32W
	me.Size = uint32(unsafe.Sizeof(me))

	ret, _, _ := procModule32FirstW.Call(snapshot, uintptr(unsafe.Pointer(&me)))
	if ret == 0 {
		return nil, fmt.Errorf("Module32First failed")
	}

	for {
		name := windows.UTF16ToString(me.Module[:])
		mod := Module{
			Name: name,
			Base: me.ModBaseAddr,
			Size: uintptr(me.ModBaseSize),
			End:  me.ModBaseAddr + uintptr(me.ModBaseSize),
		}
		modules = append(modules, mod)

		ret, _, _ = procModule32NextW.Call(snapshot, uintptr(unsafe.Pointer(&me)))
		if ret == 0 {
			break
		}
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
	var pages []Page
	var mbi memoryBasicInformation
	curAddr := module.Base

	virtualQueryEx := kernel32.NewProc("VirtualQueryEx")

	for curAddr < module.End {
		ret, _, _ := virtualQueryEx.Call(
			p.handle,
			curAddr,
			uintptr(unsafe.Pointer(&mbi)),
			unsafe.Sizeof(mbi),
		)
		if ret == 0 {
			break
		}

		page := Page{
			Start: curAddr,
			End:   curAddr + mbi.RegionSize,
			Size:  mbi.RegionSize,
		}
		pages = append(pages, page)
		curAddr += mbi.RegionSize
	}

	return pages, nil
}
