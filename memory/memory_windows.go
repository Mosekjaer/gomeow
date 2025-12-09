//go:build windows

package memory

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"

	"gomeow/process"
)

var (
	kernel32                  = windows.NewLazySystemDLL("kernel32.dll")
	procReadProcessMemory     = kernel32.NewProc("ReadProcessMemory")
	procWriteProcessMemory    = kernel32.NewProc("WriteProcessMemory")
	procVirtualAllocEx        = kernel32.NewProc("VirtualAllocEx")
	procVirtualFreeEx         = kernel32.NewProc("VirtualFreeEx")
	procVirtualProtectEx      = kernel32.NewProc("VirtualProtectEx")
	procVirtualQueryEx        = kernel32.NewProc("VirtualQueryEx")
	procGetProcAddress        = kernel32.NewProc("GetProcAddress")
	procGetModuleHandleA      = kernel32.NewProc("GetModuleHandleA")
	procCreateRemoteThread    = kernel32.NewProc("CreateRemoteThread")
)

const (
	MEM_COMMIT             = 0x1000
	MEM_RESERVE            = 0x2000
	MEM_RELEASE            = 0x8000
	PAGE_EXECUTE_READWRITE = 0x40
	PAGE_READWRITE         = 0x04
	INFINITE               = 0xFFFFFFFF
)

// Read reads memory from the target process
func Read(p *process.Process, address uintptr, buffer []byte) error {
	var bytesRead uintptr
	ret, _, err := procReadProcessMemory.Call(
		p.Handle(),
		address,
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(len(buffer)),
		uintptr(unsafe.Pointer(&bytesRead)),
	)
	if ret == 0 {
		return fmt.Errorf("ReadProcessMemory failed at 0x%X: %v", address, err)
	}
	return nil
}

// Write writes memory to the target process
func Write(p *process.Process, address uintptr, buffer []byte) error {
	var bytesWritten uintptr
	ret, _, err := procWriteProcessMemory.Call(
		p.Handle(),
		address,
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(len(buffer)),
		uintptr(unsafe.Pointer(&bytesWritten)),
	)
	if ret == 0 {
		return fmt.Errorf("WriteProcessMemory failed at 0x%X: %v", address, err)
	}
	return nil
}

// AllocateMemory allocates memory in the target process
func AllocateMemory(p *process.Process, size int, protection uint32) (uintptr, error) {
	if protection == 0 {
		protection = PAGE_EXECUTE_READWRITE
	}

	addr, _, err := procVirtualAllocEx.Call(
		p.Handle(),
		0,
		uintptr(size),
		MEM_COMMIT|MEM_RESERVE,
		uintptr(protection),
	)
	if addr == 0 {
		return 0, fmt.Errorf("VirtualAllocEx failed: %v", err)
	}
	return addr, nil
}

// FreeMemory frees memory in the target process
func FreeMemory(p *process.Process, address uintptr) error {
	ret, _, err := procVirtualFreeEx.Call(
		p.Handle(),
		address,
		0,
		MEM_RELEASE,
	)
	if ret == 0 {
		return fmt.Errorf("VirtualFreeEx failed: %v", err)
	}
	return nil
}

// PageProtection changes memory protection and returns the old protection
func PageProtection(p *process.Process, address uintptr, newProtection uint32) (uint32, error) {
	var oldProtection uint32

	// First query the region size
	type memoryBasicInfo struct {
		BaseAddress       uintptr
		AllocationBase    uintptr
		AllocationProtect uint32
		RegionSize        uintptr
		State             uint32
		Protect           uint32
		Type              uint32
	}

	var mbi memoryBasicInfo
	ret, _, _ := procVirtualQueryEx.Call(
		p.Handle(),
		address,
		uintptr(unsafe.Pointer(&mbi)),
		unsafe.Sizeof(mbi),
	)
	if ret == 0 {
		return 0, fmt.Errorf("VirtualQueryEx failed")
	}

	ret, _, err := procVirtualProtectEx.Call(
		p.Handle(),
		address,
		mbi.RegionSize,
		uintptr(newProtection),
		uintptr(unsafe.Pointer(&oldProtection)),
	)
	if ret == 0 {
		return 0, fmt.Errorf("VirtualProtectEx failed: %v", err)
	}

	return oldProtection, nil
}

// GetProcAddress gets the address of a function in a module
func GetProcAddress(moduleName, functionName string) (uintptr, error) {
	moduleNamePtr, err := windows.BytePtrFromString(moduleName)
	if err != nil {
		return 0, err
	}
	functionNamePtr, err := windows.BytePtrFromString(functionName)
	if err != nil {
		return 0, err
	}

	hModule, _, _ := procGetModuleHandleA.Call(uintptr(unsafe.Pointer(moduleNamePtr)))
	if hModule == 0 {
		return 0, fmt.Errorf("GetModuleHandleA failed for %s", moduleName)
	}

	addr, _, _ := procGetProcAddress.Call(hModule, uintptr(unsafe.Pointer(functionNamePtr)))
	if addr == 0 {
		return 0, fmt.Errorf("GetProcAddress failed for %s", functionName)
	}

	return addr, nil
}

// CreateRemoteThread creates a thread in the target process
func CreateRemoteThread(p *process.Process, startAddress, param uintptr) error {
	var threadID uint32
	hThread, _, err := procCreateRemoteThread.Call(
		p.Handle(),
		0,
		0,
		startAddress,
		param,
		0,
		uintptr(unsafe.Pointer(&threadID)),
	)
	if hThread == 0 {
		return fmt.Errorf("CreateRemoteThread failed: %v", err)
	}
	defer windows.CloseHandle(windows.Handle(hThread))

	// Wait for thread to complete
	windows.WaitForSingleObject(windows.Handle(hThread), INFINITE)

	return nil
}

// InjectLibrary injects a DLL into the target process
func InjectLibrary(p *process.Process, dllPath string) error {
	// Allocate memory for DLL path
	pathBytes := append([]byte(dllPath), 0)
	allocAddr, err := AllocateMemory(p, len(pathBytes), PAGE_READWRITE)
	if err != nil {
		return fmt.Errorf("failed to allocate memory: %v", err)
	}

	// Write DLL path to allocated memory
	if err := Write(p, allocAddr, pathBytes); err != nil {
		FreeMemory(p, allocAddr)
		return fmt.Errorf("failed to write DLL path: %v", err)
	}

	// Get LoadLibraryA address
	loadLibAddr, err := GetProcAddress("kernel32.dll", "LoadLibraryA")
	if err != nil {
		FreeMemory(p, allocAddr)
		return fmt.Errorf("failed to get LoadLibraryA address: %v", err)
	}

	// Create remote thread to call LoadLibraryA
	if err := CreateRemoteThread(p, loadLibAddr, allocAddr); err != nil {
		FreeMemory(p, allocAddr)
		return fmt.Errorf("failed to create remote thread: %v", err)
	}

	return nil
}

// InjectShellcode injects and executes shellcode in the target process
func InjectShellcode(p *process.Process, shellcode []byte, param uintptr) error {
	// Allocate memory for shellcode
	allocAddr, err := AllocateMemory(p, len(shellcode), PAGE_EXECUTE_READWRITE)
	if err != nil {
		return fmt.Errorf("failed to allocate memory: %v", err)
	}

	// Write shellcode to allocated memory
	if err := Write(p, allocAddr, shellcode); err != nil {
		FreeMemory(p, allocAddr)
		return fmt.Errorf("failed to write shellcode: %v", err)
	}

	// Create remote thread to execute shellcode
	if err := CreateRemoteThread(p, allocAddr, param); err != nil {
		FreeMemory(p, allocAddr)
		return fmt.Errorf("failed to create remote thread: %v", err)
	}

	return nil
}
