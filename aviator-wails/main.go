package main

import (
	"aviator-wails/internal/config"
	"aviator-wails/internal/discovery"
	"aviator-wails/internal/processmon"
	"aviator-wails/internal/server"
	"aviator-wails/internal/web"
	"embed"
	"log"
	"time"

	"github.com/energye/systray"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/windows/icon.ico
var iconData []byte

// NewProcessMonitor initializes process monitor with all configured apps
func NewProcessMonitor(cm *config.ConfigManager) *processmon.ProcessMonitor {
	pm := processmon.NewProcessMonitor()
	for _, app := range cm.GetApps() {
		pm.AddWatch(app.ID, app.Path)
	}
	return pm
}

func main() {
	log.Println("Initializing Aviator (Wails)...")

	// 1. Load Config
	cm, err := config.NewConfigManager()
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	// 2. Get Static Assets (For HTTP Server)
	webFS, err := web.GetFS()
	if err != nil {
		log.Fatalf("Failed to load embedded assets: %v", err)
	}

	// 3. Create Process Monitor
	pm := NewProcessMonitor(cm)

	// 4. Initialize Server
	srv := server.NewServer(cm, webFS, pm)

	// 5. Discovery Service
	var ds *discovery.DiscoveryService = nil

	// 6. Create Wails App
	app := NewApp(cm, srv, ds)

	// --- System Tray Logic ---
	go func() {
		systray.Run(func() {
			systray.SetIcon(iconData)
			systray.SetTitle("Aviator")
			systray.SetTooltip("Aviator Launcher")

			// Left Click -> Show Window immediately
			systray.SetOnClick(func(menu systray.IMenu) {
				app.Show()
			})

			// Right Click -> Show Context Menu
			systray.SetOnRClick(func(menu systray.IMenu) {
				menu.ShowMenu()
			})

			// Menu Items
			mShow := systray.AddMenuItem("Minimize Aviator", "Show or Hide the main window")
			mShow.Click(func() {
				if app.IsWindowVisible() {
					app.Hide()
				} else {
					app.Show()
				}
				// Title update is handled by loop
			})

			mServer := systray.AddMenuItem("Start Server", "Start or Stop the server")
			mServer.Click(func() {
				if app.IsServerRunning() {
					app.StopServer()
				} else {
					app.StartServer()
				}
				// Title update is handled by loop
			})

			systray.AddSeparator()

			mQuit := systray.AddMenuItem("Quit", "Quit the application")
			mQuit.Click(func() {
				app.SetQuitting(true)
				if app.GetContext() != nil {
					runtime.Quit(app.GetContext())
				}
				systray.Quit()
			})

			// Update Menu State Loop (Synchronize with App State)
			go func() {
				for {
					// Update Server Status
					if app.IsServerRunning() {
						mServer.SetTitle("Stop Server")
						mServer.SetTooltip("Server is Running")
					} else {
						mServer.SetTitle("Start Server")
						mServer.SetTooltip("Server is Stopped")
					}

					// Update Window Visibility Status
					if app.IsWindowVisible() {
						mShow.SetTitle("Minimize Aviator")
					} else {
						mShow.SetTitle("Show Aviator")
					}

					time.Sleep(200 * time.Millisecond)
				}
			}()

		}, func() {
			// Cleanup
		})
	}()
	// -------------------------

	log.Println("Aviator initialized. Use 'Start Server' button to enable web access.")

	// 7. Run Wails Application
	err = wails.Run(&options.App{
		Title:     "Aviator",
		Width:     700,
		Height:    700,
		MinWidth:  650,
		MinHeight: 550,
		Frameless: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 13, G: 17, B: 23, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		OnBeforeClose:    app.AllowClose, // Intercept close event
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
		},
	})

	if err != nil {
		log.Fatal("Error:", err.Error())
	}
}
