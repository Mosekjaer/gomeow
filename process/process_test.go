package process

import (
	"os"
	"runtime"
	"testing"
)

// TestEnumProcesses tests that we can enumerate processes
func TestEnumProcesses(t *testing.T) {
	processes, err := EnumProcesses()
	if err != nil {
		t.Fatalf("EnumProcesses failed: %v", err)
	}

	if len(processes) == 0 {
		t.Error("EnumProcesses returned no processes")
	}

	// Verify that each process has valid data
	for _, p := range processes {
		if p.PID == 0 && p.Name != "System Idle Process" && p.Name != "[kernel]" {
			// PID 0 is typically reserved for system processes
			continue
		}

		if p.Name == "" {
			t.Errorf("Process %d has empty name", p.PID)
		}
	}

	t.Logf("Found %d processes", len(processes))
}

// TestPIDExists tests process existence checking by PID
func TestPIDExists(t *testing.T) {
	// Current process should exist
	currentPID := os.Getpid()
	if !PIDExists(currentPID) {
		t.Errorf("PIDExists returned false for current process (PID: %d)", currentPID)
	}

	// Non-existent process (very high PID unlikely to exist)
	if PIDExists(999999999) {
		t.Error("PIDExists returned true for non-existent PID")
	}
}

// TestProcessExists tests process existence checking by name
func TestProcessExists(t *testing.T) {
	// Try to find a common process
	var testProcess string
	if runtime.GOOS == "windows" {
		testProcess = "explorer.exe"
	} else {
		testProcess = "init"
	}

	exists := ProcessExists(testProcess)
	// Just log it - process may or may not exist
	t.Logf("ProcessExists(%s) = %v", testProcess, exists)
}

// TestGetProcessID tests getting PID by name
func TestGetProcessID(t *testing.T) {
	// Try to find a common process
	var testProcess string
	if runtime.GOOS == "windows" {
		testProcess = "explorer.exe"
	} else {
		testProcess = "init"
	}

	pid, err := GetProcessID(testProcess)
	if err != nil {
		t.Logf("GetProcessID(%s) failed: %v - process may not be running", testProcess, err)
	} else {
		t.Logf("GetProcessID(%s) = %d", testProcess, pid)
	}
}

// TestGetProcessName tests getting process name by PID
func TestGetProcessName(t *testing.T) {
	currentPID := os.Getpid()
	name, err := GetProcessName(currentPID)

	if err != nil {
		t.Errorf("GetProcessName failed for current process (PID: %d): %v", currentPID, err)
	} else if name == "" {
		t.Errorf("GetProcessName returned empty string for current process (PID: %d)", currentPID)
	} else {
		t.Logf("Current process name: %s (PID: %d)", name, currentPID)
	}
}

// TestOpenCloseProcess tests opening and closing a process handle
func TestOpenCloseProcess(t *testing.T) {
	// Get list of processes
	processes, err := EnumProcesses()
	if err != nil {
		t.Fatalf("EnumProcesses failed: %v", err)
	}

	if len(processes) == 0 {
		t.Skip("No processes found to test")
	}

	// Try to open the first accessible process
	var opened bool
	for _, p := range processes {
		// Skip system processes on Windows
		if p.PID <= 4 {
			continue
		}

		proc, err := OpenProcess(p.Name)
		if err != nil {
			continue // Process might not be accessible
		}

		// Verify process data
		if proc.PID == 0 {
			t.Errorf("Opened process has PID 0")
		}
		if proc.Handle() == 0 && runtime.GOOS == "windows" {
			t.Errorf("Opened process has nil handle on Windows")
		}

		// Close the process
		proc.Close()
		opened = true
		t.Logf("Successfully opened and closed process: %s (PID: %d)", p.Name, p.PID)
		break
	}

	if !opened {
		t.Log("Could not open any process - may require elevated privileges")
	}
}

// TestOpenProcessByPID tests opening process by PID
func TestOpenProcessByPID(t *testing.T) {
	processes, err := EnumProcesses()
	if err != nil {
		t.Fatalf("EnumProcesses failed: %v", err)
	}

	for _, p := range processes {
		if p.PID <= 4 {
			continue
		}

		proc, err := OpenProcessByPID(p.PID)
		if err != nil {
			continue
		}

		t.Logf("Successfully opened process by PID: %s (PID: %d)", proc.Name, proc.PID)
		proc.Close()
		return
	}

	t.Log("Could not open any process by PID - may require elevated privileges")
}

// TestEnumModules tests module enumeration for an open process
func TestEnumModules(t *testing.T) {
	processes, err := EnumProcesses()
	if err != nil {
		t.Fatalf("EnumProcesses failed: %v", err)
	}

	for _, p := range processes {
		if p.PID <= 4 {
			continue
		}

		proc, err := OpenProcess(p.Name)
		if err != nil {
			continue
		}

		// Try to enumerate modules
		modules, err := proc.EnumModules()
		proc.Close()

		if err != nil {
			continue // Might not have access
		}

		if len(modules) > 0 {
			t.Logf("Process %s has %d modules", p.Name, len(modules))
			t.Logf("First module: %s at 0x%X (size: %d bytes)",
				modules[0].Name, modules[0].Base, modules[0].Size)
			return
		}
	}

	t.Log("Could not enumerate modules for any process - may require elevated privileges")
}

// TestModuleExists tests checking if a module exists
func TestModuleExists(t *testing.T) {
	processes, err := EnumProcesses()
	if err != nil {
		t.Fatalf("EnumProcesses failed: %v", err)
	}

	for _, p := range processes {
		if p.PID <= 4 {
			continue
		}

		proc, err := OpenProcess(p.Name)
		if err != nil {
			continue
		}

		// Get modules first to know what to look for
		modules, err := proc.EnumModules()
		if err != nil || len(modules) == 0 {
			proc.Close()
			continue
		}

		// Check if the first module exists
		exists := proc.ModuleExists(modules[0].Name)
		proc.Close()

		if exists {
			t.Logf("Module %s exists in process %s", modules[0].Name, p.Name)
			return
		}
	}

	t.Log("Could not verify module existence - may require elevated privileges")
}

// TestGetModule tests getting a specific module
func TestGetModule(t *testing.T) {
	processes, err := EnumProcesses()
	if err != nil {
		t.Fatalf("EnumProcesses failed: %v", err)
	}

	for _, p := range processes {
		if p.PID <= 4 {
			continue
		}

		proc, err := OpenProcess(p.Name)
		if err != nil {
			continue
		}

		modules, err := proc.EnumModules()
		if err != nil || len(modules) == 0 {
			proc.Close()
			continue
		}

		// Try to get the first module by name
		module, err := proc.GetModule(modules[0].Name)
		proc.Close()

		if err == nil && module != nil {
			t.Logf("GetModule(%s) returned: Base=0x%X, Size=%d",
				module.Name, module.Base, module.Size)
			return
		}
	}

	t.Log("Could not get module - may require elevated privileges")
}

// TestGetProcessPath tests getting the executable path
func TestGetProcessPath(t *testing.T) {
	processes, err := EnumProcesses()
	if err != nil {
		t.Fatalf("EnumProcesses failed: %v", err)
	}

	for _, p := range processes {
		if p.PID <= 4 {
			continue
		}

		proc, err := OpenProcess(p.Name)
		if err != nil {
			continue
		}

		path, err := proc.GetPath()
		proc.Close()

		if err == nil && path != "" {
			t.Logf("Process %s path: %s", p.Name, path)
			return
		}
	}

	t.Log("Could not get process path - may require elevated privileges")
}

// TestIs64Bit tests architecture detection
func TestIs64Bit(t *testing.T) {
	processes, err := EnumProcesses()
	if err != nil {
		t.Fatalf("EnumProcesses failed: %v", err)
	}

	for _, p := range processes {
		if p.PID <= 4 {
			continue
		}

		proc, err := OpenProcess(p.Name)
		if err != nil {
			continue
		}

		is64 := proc.Is64Bit()
		proc.Close()

		// Just verify it doesn't crash/panic
		t.Logf("Process %s is 64-bit: %v", p.Name, is64)
		return
	}

	t.Log("Could not check architecture - may require elevated privileges")
}

// TestIsRunning tests process running status
func TestIsRunning(t *testing.T) {
	processes, err := EnumProcesses()
	if err != nil {
		t.Fatalf("EnumProcesses failed: %v", err)
	}

	for _, p := range processes {
		if p.PID <= 4 {
			continue
		}

		proc, err := OpenProcess(p.Name)
		if err != nil {
			continue
		}

		isRunning := proc.IsRunning()
		proc.Close()

		if isRunning {
			t.Logf("Process %s (PID: %d) is running", p.Name, p.PID)
			return
		}
	}

	t.Log("Could not verify process running status - may require elevated privileges")
}

// BenchmarkEnumProcesses benchmarks process enumeration
func BenchmarkEnumProcesses(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = EnumProcesses()
	}
}

// BenchmarkPIDExists benchmarks process existence check by PID
func BenchmarkPIDExists(b *testing.B) {
	pid := os.Getpid()
	for i := 0; i < b.N; i++ {
		_ = PIDExists(pid)
	}
}

// BenchmarkProcessExists benchmarks process existence check by name
func BenchmarkProcessExists(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ProcessExists("explorer.exe")
	}
}
