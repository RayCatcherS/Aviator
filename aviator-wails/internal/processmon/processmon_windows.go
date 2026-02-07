package processmon

import (
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"unsafe"
)

var (
	kernel32                     = syscall.NewLazyDLL("kernel32.dll")
	procCreateToolhelp32Snapshot = kernel32.NewProc("CreateToolhelp32Snapshot")
	procProcess32First           = kernel32.NewProc("Process32FirstW")
	procProcess32Next            = kernel32.NewProc("Process32NextW")
	procCloseHandle              = kernel32.NewProc("CloseHandle")
)

const (
	TH32CS_SNAPPROCESS = 0x00000002
	MAX_PATH           = 260
)

type PROCESSENTRY32 struct {
	dwSize              uint32
	cntUsage            uint32
	th32ProcessID       uint32
	th32DefaultHeapID   uintptr
	th32ModuleID        uint32
	cntThreads          uint32
	th32ParentProcessID uint32
	pcPriClassBase      int32
	dwFlags             uint32
	szExeFile           [MAX_PATH]uint16
}

// ProcessMonitor monitors running processes
type ProcessMonitor struct {
	watchedProcesses map[string]string // appID -> exe filename (lowercase)
	runningStatus    map[string]bool   // appID -> isRunning
	mu               sync.RWMutex
}

// NewProcessMonitor creates a new process monitor
func NewProcessMonitor() *ProcessMonitor {
	return &ProcessMonitor{
		watchedProcesses: make(map[string]string),
		runningStatus:    make(map[string]bool),
	}
}

// AddWatch adds an application to watch
func (pm *ProcessMonitor) AddWatch(appID, exePath string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Extract just the executable filename (case-insensitive)
	exeName := filepath.Base(exePath)
	pm.watchedProcesses[appID] = strings.ToLower(exeName)
}

// RemoveWatch removes an application from watching
func (pm *ProcessMonitor) RemoveWatch(appID string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	delete(pm.watchedProcesses, appID)
	delete(pm.runningStatus, appID)
}

// Update scans all running processes and updates status
func (pm *ProcessMonitor) Update() error {
	// 1. Get watched apps snapshot (read-only lock)
	pm.mu.RLock()
	watched := make(map[string]string) // appID -> exeName
	for id, exe := range pm.watchedProcesses {
		watched[id] = exe
	}
	pm.mu.RUnlock()

	// 2. Scan processes (NO LOCK held here, expensive operation)
	snapshot, _, _ := procCreateToolhelp32Snapshot.Call(
		uintptr(TH32CS_SNAPPROCESS),
		0,
	)
	if snapshot == 0 || snapshot == uintptr(syscall.InvalidHandle) {
		return syscall.GetLastError()
	}
	defer procCloseHandle.Call(snapshot)

	var pe PROCESSENTRY32
	pe.dwSize = uint32(unsafe.Sizeof(pe))

	// Local map for results
	currentRunning := make(map[string]bool)
	for id := range watched {
		currentRunning[id] = false
	}

	// Get first process
	ret, _, _ := procProcess32First.Call(snapshot, uintptr(unsafe.Pointer(&pe)))
	if ret == 0 {
		return syscall.GetLastError()
	}

	// Iterate through processes
	for {
		exeName := syscall.UTF16ToString(pe.szExeFile[:])
		exeNameLower := strings.ToLower(exeName)

		// Check if this process matches any watched app
		for appID, watchedExe := range watched {
			if exeNameLower == watchedExe {
				currentRunning[appID] = true
			}
		}

		// Get next process
		ret, _, _ := procProcess32Next.Call(snapshot, uintptr(unsafe.Pointer(&pe)))
		if ret == 0 {
			break
		}
	}

	// 3. Update status (Write lock, very fast)
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.runningStatus = currentRunning

	return nil
}

// GetStatus returns the running status of a specific app
func (pm *ProcessMonitor) GetStatus(appID string) bool {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.runningStatus[appID]
}

// GetAllStatuses returns the running status of all watched apps
func (pm *ProcessMonitor) GetAllStatuses() map[string]bool {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	result := make(map[string]bool)
	for appID, status := range pm.runningStatus {
		result[appID] = status
	}
	return result
}
