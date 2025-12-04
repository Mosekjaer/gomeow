package process

// Process represents an OS process withrepresents an opened process
type Process struct {
	Name   string
	PID    int
	Base   uintptr
	Debug  bool
	handle uintptr // Windows process handle, unused on Linux
}

// Module represents a loaded module/library in a process
type Module struct {
	Name string
	Base uintptr
	End  uintptr
	Size uintptr
}

// Page represents a memory region/page
type Page struct {
	Start uintptr
	End   uintptr
	Size  uintptr
}

// ProcessInfo is a lightweight struct for enumeration
type ProcessInfo struct {
	Name string
	PID  int
}
