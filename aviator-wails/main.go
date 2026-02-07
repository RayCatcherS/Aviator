package main

import (
	"aviator-wails/internal/config"
	"aviator-wails/internal/discovery"
	"aviator-wails/internal/server"
	"aviator-wails/internal/web"
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	log.Println("Initializing Aviator (Wails)...")

	// 1. Load Config
	cm, err := config.NewConfigManager()
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}
	log.Printf("Loaded %d apps from config.\n", len(cm.GetApps()))

	// 2. Get Static Assets (For HTTP Server)
	webFS, err := web.GetFS()
	if err != nil {
		log.Fatalf("Failed to load embedded assets: %v", err)
	}

	// 3. Initialize Server (but don't start it automatically - user controls via UI)
	srv := server.NewServer(cm, webFS)

	// 4. Discovery service will be started when user starts server
	var ds *discovery.DiscoveryService = nil

	// 5. Create Wails App
	app := NewApp(cm, srv, ds)

	log.Println("Aviator initialized. Use 'Start Server' button to enable web access.")

	// 6. Run Wails Application
	err = wails.Run(&options.App{
		Title:     "Aviator",
		Width:     900,
		Height:    700,
		MinWidth:  750,
		MinHeight: 550,
		Frameless: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 13, G: 17, B: 23, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
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
