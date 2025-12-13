//go:build linux

package memory

import (
	"fmt"
	"syscall"
	"unsafe"

	"gomeow/process"
)

// Read reads raw bytes fromreads memory from the target process using process_vm_readv
func Read(p *process.Process, address uintptr, buffer []byte) error {
	localIov := syscall.Iovec{
		Base: &buffer[0],
		Len:  uint64(len(buffer)),
	}
	remoteIov := syscall.Iovec{
		Base: (*byte)(unsafe.Pointer(address)),
		Len:  uint64(len(buffer)),
	}

	_, _, errno := syscall.Syscall6(
		syscall.SYS_PROCESS_VM_READV,
		uintptr(p.PID),
		uintptr(unsafe.Pointer(&localIov)),
		1,
		uintptr(unsafe.Pointer(&remoteIov)),
		1,
		0,
	)
	if errno != 0 {
		return fmt.Errorf("process_vm_readv failed at 0x%X: %v", address, errno)
	}
	return nil
}

// Write writes memory to the target process using process_vm_writev
func Write(p *process.Process, address uintptr, buffer []byte) error {
	localIov := syscall.Iovec{
		Base: &buffer[0],
		Len:  uint64(len(buffer)),
	}
	remoteIov := syscall.Iovec{
		Base: (*byte)(unsafe.Pointer(address)),
		Len:  uint64(len(buffer)),
	}

	_, _, errno := syscall.Syscall6(
		syscall.SYS_PROCESS_VM_WRITEV,
		uintptr(p.PID),
		uintptr(unsafe.Pointer(&localIov)),
		1,
		uintptr(unsafe.Pointer(&remoteIov)),
		1,
		0,
	)
	if errno != 0 {
		return fmt.Errorf("process_vm_writev failed at 0x%X: %v", address, errno)
	}
	return nil
}

// AllocateMemory is not fully supported on Linux without ptrace
// This is a stub that returns an error
func AllocateMemory(p *process.Process, size int, protection uint32) (uintptr, error) {
	return 0, fmt.Errorf("AllocateMemory requires ptrace implementation on Linux")
}

// FreeMemory is not fully supported on Linux without ptrace
func FreeMemory(p *process.Process, address uintptr) error {
	return fmt.Errorf("FreeMemory requires ptrace implementation on Linux")
}

// PageProtection is not fully supported on Linux without ptrace
func PageProtection(p *process.Process, address uintptr, newProtection uint32) (uint32, error) {
	return 0, fmt.Errorf("PageProtection requires ptrace implementation on Linux")
}

// GetProcAddress is Windows-only
func GetProcAddress(moduleName, functionName string) (uintptr, error) {
	return 0, fmt.Errorf("GetProcAddress is Windows-only")
}

// CreateRemoteThread is Windows-only
func CreateRemoteThread(p *process.Process, startAddress, param uintptr) error {
	return fmt.Errorf("CreateRemoteThread is Windows-only")
}

// InjectLibrary is Windows-only (uses CreateRemoteThread + LoadLibraryA)
func InjectLibrary(p *process.Process, dllPath string) error {
	return fmt.Errorf("InjectLibrary is Windows-only")
}

// InjectShellcode requires ptrace on Linux
func InjectShellcode(p *process.Process, shellcode []byte, param uintptr) error {
	return fmt.Errorf("InjectShellcode requires ptrace implementation on Linux")
}
