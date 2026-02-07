package server

import (
	"aviator-wails/internal/config"
	"aviator-wails/internal/launcher"
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Server struct {
	Config     *config.ConfigManager
	FileServer http.Handler
	httpServer *http.Server
}

func NewServer(cm *config.ConfigManager, webFS fs.FS) *Server {
	// Create file server for static files
	fsHandler := http.FileServer(http.FS(webFS))

	return &Server{
		Config:     cm,
		FileServer: fsHandler,
	}
}

func (s *Server) Start(port int) error {
	addr := fmt.Sprintf("0.0.0.0:%d", port)
	log.Printf("Starting server on %s", addr)

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s,
	}

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Stop() error {
	if s.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.httpServer.Shutdown(ctx)
	}
	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// CORS headers for dev
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// API Routing
	if strings.HasPrefix(r.URL.Path, "/api/") {
		s.handleAPI(w, r)
		return
	}

	// Static Files
	s.FileServer.ServeHTTP(w, r)
}

func (s *Server) handleAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch {
	case r.URL.Path == "/api/apps" && r.Method == "GET":
		json.NewEncoder(w).Encode(s.Config.GetApps())

	case strings.HasPrefix(r.URL.Path, "/api/launch/") && r.Method == "POST":
		appID := strings.TrimPrefix(r.URL.Path, "/api/launch/")
		s.handleLaunch(w, appID)

	case r.URL.Path == "/api/info" && r.Method == "GET":
		// Get hostname and username
		hostname, _ := os.Hostname()
		username := os.Getenv("USERNAME") // Windows username

		// Create a friendly display name
		displayName := hostname
		if username != "" {
			displayName = fmt.Sprintf("%s@%s", username, hostname)
		}

		json.NewEncoder(w).Encode(map[string]string{
			"status":   "running",
			"backend":  "go",
			"hostname": displayName,
		})

	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

func (s *Server) handleLaunch(w http.ResponseWriter, appID string) {
	app, found := s.Config.GetAppByID(appID)
	if !found {
		http.Error(w, `{"error": "App not found"}`, http.StatusNotFound)
		return
	}

	if err := launcher.RunExecutable(app.Path, app.Args); err != nil {
		log.Printf("Error launching %s: %v", app.Name, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Launched " + app.Name,
	})
}
