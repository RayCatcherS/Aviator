package main

import (
	"aviator-wails/internal/config"
	"aviator-wails/internal/discovery"
	"aviator-wails/internal/processmon"
	"aviator-wails/internal/server"
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx             context.Context
	config          *config.ConfigManager
	server          *server.Server
	discovery       *discovery.DiscoveryService
	processMonitor  *processmon.ProcessMonitor
	serverRunning   bool
	isQuitting      bool
	isWindowVisible bool
}

// NewApp creates a new App application struct
func NewApp(cm *config.ConfigManager, srv *server.Server, ds *discovery.DiscoveryService) *App {
	// Use the process monitor from server (shared instance)
	return &App{
		config:          cm,
		server:          srv,
		discovery:       ds,
		processMonitor:  srv.ProcessMonitor,
		serverRunning:   false,
		isQuitting:      false,
		isWindowVisible: true,
	}
}

// IsServerRunning returns true if the server is active
func (a *App) IsServerRunning() bool {
	return a.serverRunning
}

// IsWindowVisible returns true if window is shown
func (a *App) IsWindowVisible() bool {
	return a.isWindowVisible
}

// SetQuitting sets the quit flag
func (a *App) SetQuitting(quitting bool) {
	a.isQuitting = quitting
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	log.Println("Aviator Wails app started")

	// Start background process monitoring
	go a.monitorProcesses()
}

// monitorProcesses polls for running processes every 3 seconds
func (a *App) monitorProcesses() {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-a.ctx.Done():
			return
		case <-ticker.C:
			if err := a.processMonitor.Update(); err != nil {
				log.Printf("Error updating process monitor: %v", err)
			}
		}
	}
}

// shutdown is called when the app closes
func (a *App) shutdown(ctx context.Context) {
	log.Println("Shutting down Aviator")
	if a.server != nil {
		a.server.Stop()
	}
	if a.discovery != nil {
		a.discovery.Shutdown()
	}
}

// GetApps returns all configured applications
func (a *App) GetApps() []config.App {
	return a.config.GetApps()
}

// AddApp adds a new application to the configuration
func (a *App) AddApp(name, path, args string) config.App {
	app := a.config.AddApp(name, path, args)
	a.processMonitor.AddWatch(app.ID, app.Path)
	return app
}

// UpdateApp updates an existing application
func (a *App) UpdateApp(id, name, path, args string) bool {
	success := a.config.UpdateApp(id, name, path, args)
	if success {
		// Update the watch with new path
		a.processMonitor.AddWatch(id, path)
	}
	return success
}

// RemoveApp removes an application from the configuration
func (a *App) RemoveApp(id string) {
	a.config.RemoveApp(id)
	a.processMonitor.RemoveWatch(id)
}

// LaunchApp launches an application by ID
func (a *App) LaunchApp(id string) error {
	_, found := a.config.GetAppByID(id)
	if !found {
		return fmt.Errorf("application not found")
	}

	// The HTTP server handles this via /api/launch endpoint
	return nil
}

// GetServerInfo returns server status information
func (a *App) GetServerInfo() map[string]interface{} {
	localIP := getOutboundIP()
	status := "stopped"
	if a.serverRunning {
		status = "running"
	}
	return map[string]interface{}{
		"localURL":   "http://localhost:8000",
		"networkURL": fmt.Sprintf("http://%s:8000", localIP),
		"status":     status,
		"running":    a.serverRunning,
	}
}

// StartServer starts the HTTP server
func (a *App) StartServer() error {
	if a.serverRunning {
		return fmt.Errorf("server already running")
	}

	// Start server in background
	go func() {
		if err := a.server.Start(8000); err != nil {
			log.Printf("Server error: %v", err)
			a.serverRunning = false
			// Emit event to notify frontend
			if a.ctx != nil {
				runtime.EventsEmit(a.ctx, "server:stopped")
			}
		}
	}()

	// Start discovery service
	if a.discovery != nil {
		a.discovery = nil // Clear old instance if exists
		ds, err := discovery.NewDiscoveryService(8000)
		if err != nil {
			log.Printf("Discovery service error: %v", err)
		} else {
			a.discovery = ds
		}
	}

	a.serverRunning = true
	log.Println("HTTP server started on port 8000")

	// Emit event to notify frontend
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "server:started")
	}

	return nil
}

// StopServer stops the HTTP server
func (a *App) StopServer() error {
	if !a.serverRunning {
		return fmt.Errorf("server not running")
	}

	if a.server != nil {
		if err := a.server.Stop(); err != nil {
			return err
		}
	}

	if a.discovery != nil {
		a.discovery.Shutdown()
	}

	a.serverRunning = false
	log.Println("HTTP server stopped")

	// Emit event to notify frontend
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "server:stopped")
	}

	return nil
}

// SelectFile opens a native file dialog to select an executable
func (a *App) SelectFile() (string, error) {
	result, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Executable",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Executables (*.exe)",
				Pattern:     "*.exe",
			},
			{
				DisplayName: "All Files (*.*)",
				Pattern:     "*.*",
			},
		},
	})

	if err != nil {
		return "", err
	}

	return result, nil
}

// GetProcessStatuses returns the running status of all launched apps
func (a *App) GetProcessStatuses() map[string]bool {
	return a.processMonitor.GetAllStatuses()
}

// Show makes the window visible and focused
func (a *App) Show() {
	if a.ctx != nil {
		runtime.WindowShow(a.ctx)
		a.isWindowVisible = true
	}
}

// Hide hides the window
func (a *App) Hide() {
	if a.ctx != nil {
		runtime.WindowHide(a.ctx)
		a.isWindowVisible = false
	}
}

// OnBeforeClose handles window closing event
func (a *App) AllowClose(ctx context.Context) (prevent bool) {
	if !a.isQuitting {
		// Just hide the window instead of closing
		a.Hide()
		return true // Prevent actual close
	}

	// Attempting to quit fully (from Tray)
	if a.serverRunning {
		// Ask confirmation if server is running
		// Note: We show the window first to ensure dialog is visible
		runtime.WindowShow(ctx)

		res, err := runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
			Type:          runtime.QuestionDialog,
			Title:         "Confirm Exit",
			Message:       "The server is running. Do you want to stop it and exit?",
			Buttons:       []string{"Yes", "No"},
			DefaultButton: "Yes",
			CancelButton:  "No",
		})

		if err != nil || res != "Yes" {
			// User cancelled execution or error occurred
			// Reset quitting flag so next X click just hides
			a.isQuitting = false
			return true // Prevent close
		}

		// User said Yes -> Stop server
		a.StopServer()
	}

	return false // Allow close
}

// GetContext returns the application context
func (a *App) GetContext() context.Context {
	return a.ctx
}

// Helper function to get outbound IP
func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
