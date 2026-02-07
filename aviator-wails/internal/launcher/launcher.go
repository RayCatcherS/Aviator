package launcher

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// ProcessInfo tracks running processes
type ProcessInfo struct {
	AppID     string
	AppName   string
	Pid       int
	IsRunning bool
	cmd       *exec.Cmd
}

var (
	runningProcesses = make(map[string]*ProcessInfo)
	processMutex     sync.RWMutex
)

// RunExecutableWithTracking launches the application and tracks its process
func RunExecutableWithTracking(appID, appName, path, args string) (int, error) {
	fmt.Printf("[Launcher] Launching: %s Args: %s\n", path, args)

	// Validate path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return 0, fmt.Errorf("executable not found: %s", path)
	}

	// Parse arguments
	var cmdArgs []string
	if args != "" {
		cmdArgs = strings.Fields(args)
	}

	cmd := exec.Command(path, cmdArgs...)
	cmd.Dir = filepath.Dir(path)

	if err := cmd.Start(); err != nil {
		return 0, err
	}

	pid := 0
	if cmd.Process != nil {
		pid = cmd.Process.Pid

		// Store process info
		processMutex.Lock()
		runningProcesses[appID] = &ProcessInfo{
			AppID:     appID,
			AppName:   appName,
			Pid:       pid,
			IsRunning: true,
			cmd:       cmd,
		}
		processMutex.Unlock()

		// Monitor process in background
		go monitorProcess(appID, cmd)
	}

	return pid, nil
}

// monitorProcess waits for the process to finish and updates status
func monitorProcess(appID string, cmd *exec.Cmd) {
	cmd.Wait() // This blocks until the process finishes

	processMutex.Lock()
	if info, exists := runningProcesses[appID]; exists {
		info.IsRunning = false
	}
	processMutex.Unlock()

	fmt.Printf("[Launcher] Process for app %s has terminated\n", appID)
}

// GetProcessStatus returns the status of a specific app
func GetProcessStatus(appID string) (bool, int) {
	processMutex.RLock()
	defer processMutex.RUnlock()

	if info, exists := runningProcesses[appID]; exists {
		return info.IsRunning, info.Pid
	}
	return false, 0
}

// GetAllProcessStatuses returns status for all apps
func GetAllProcessStatuses() map[string]bool {
	processMutex.RLock()
	defer processMutex.RUnlock()

	statuses := make(map[string]bool)
	for appID, info := range runningProcesses {
		statuses[appID] = info.IsRunning
	}
	return statuses
}

// RunExecutable is the old function, kept for compatibility
func RunExecutable(path string, args string) error {
	fmt.Printf("[Launcher] Launching: %s Args: %s\n", path, args)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("executable not found: %s", path)
	}

	var cmdArgs []string
	if args != "" {
		cmdArgs = strings.Fields(args)
	}

	cmd := exec.Command(path, cmdArgs...)
	cmd.Dir = filepath.Dir(path)

	if err := cmd.Start(); err != nil {
		return err
	}

	if cmd.Process != nil {
		cmd.Process.Release()
	}

	return nil
}
