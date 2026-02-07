package main

import (
	"aviator/internal/config"
	"aviator/internal/discovery"
	"aviator/internal/gui"
	"aviator/internal/server"
	"aviator/internal/web"
	"fmt"
	"log"
)

func main() {
	log.Println("Initializing Aviator (Go)...")

	// 1. Load Config
	cm, err := config.NewConfigManager()
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}
	fmt.Printf("Loaded %d apps from config.\n", len(cm.GetApps()))

	// 2. Get Static Assets (Embedded)
	webFS, err := web.GetFS()
	if err != nil {
		log.Fatalf("Failed to load embedded assets: %v", err)
	}

	// 3. Start Discovery Service (mDNS)
	ds, err := discovery.NewDiscoveryService(8000)
	if err != nil {
		log.Printf("Failed to start discovery service: %v", err)
	} else {
		defer ds.Shutdown()
	}

	// 4. Initialize Server
	srv := server.NewServer(cm, webFS)
	
	// 5. Initialize GUI
	// We want to open "http://localhost:8000"
	appURL := "http://localhost:8000"
	
	appGUI, err := gui.NewGUI(appURL, cm, srv)
	if err != nil {
		log.Printf("Failed to initialize GUI: %v. Running in console mode.", err)
	}

	// 6. Start Server (in background)
	// GUI assumes it is started, or GUI starts it?
	// In gui.go we set serverRunning = true.
	// So we should start it here.
	go func() {
		if err := srv.Start(8000); err != nil {
			log.Printf("Server stopped or failed: %v", err)
		}
	}()
	log.Println("Server is running on http://localhost:8000")
	
	// 6. Run GUI Loop (Blocking)
	if appGUI != nil {
		appGUI.Run()
	} else {
		// Fallback if GUI failed: block forever
		log.Println("Press Ctrl+C to exit.")
		select {}
	}
}
