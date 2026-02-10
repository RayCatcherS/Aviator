package server

import (
	"aviator-wails/internal/config"
	"aviator-wails/internal/launcher"
	"aviator-wails/internal/processmon"
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Server struct {
	Config         *config.ConfigManager
	ProcessMonitor *processmon.ProcessMonitor
	FileServer     http.Handler
	httpServer     *http.Server

	// Key Bucket (Session Pool)
	keyBucket   map[string]time.Time
	bucketMutex sync.RWMutex

	// Cache for process statuses
	statusCache      map[string]bool
	statusCacheTime  time.Time
	statusCacheMutex sync.RWMutex
}

func NewServer(cm *config.ConfigManager, webFS fs.FS, pm *processmon.ProcessMonitor) *Server {
	// Create file server for static files
	fsHandler := http.FileServer(http.FS(webFS))

	return &Server{
		Config:         cm,
		ProcessMonitor: pm,
		FileServer:     fsHandler,
		statusCache:    make(map[string]bool),
		keyBucket:      make(map[string]time.Time),
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
	log.Printf("[REQ] Incoming request: %s %s", r.Method, r.URL.Path)

	// CORS headers for dev
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	// Disabilita cache per forzare aggiornamenti su mobile
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Public check for server health
	if r.URL.Path == "/ping" {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("PONG"))
		return
	}

	// Real logic but at top level to avoid handleAPI/switch potential issues
	if r.URL.Path == "/api/info" {
		w.Header().Set("Content-Type", "application/json")
		settings := s.Config.GetSettings()
		authorized := s.isAuthorized(r)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":        "running",
			"backend":       "go",
			"version":       "v2.8.2",
			"hostname":      "Aviator Desktop",
			"auth_required": settings.AuthEnabled,
			"is_authorized": authorized,
		})
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
		if !s.isAuthorized(r) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		json.NewEncoder(w).Encode(s.Config.GetApps())

	case strings.HasPrefix(r.URL.Path, "/api/launch/") && r.Method == "POST":
		if !s.isAuthorized(r) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		appID := strings.TrimPrefix(r.URL.Path, "/api/launch/")
		s.handleLaunch(w, appID)

	case r.URL.Path == "/api/process-statuses" && r.Method == "GET":
		if !s.isAuthorized(r) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		s.statusCacheMutex.RLock()
		cacheAge := time.Since(s.statusCacheTime)
		if cacheAge < time.Second && len(s.statusCache) > 0 {
			cached := make(map[string]bool)
			for k, v := range s.statusCache {
				cached[k] = v
			}
			s.statusCacheMutex.RUnlock()
			json.NewEncoder(w).Encode(cached)
			return
		}
		s.statusCacheMutex.RUnlock()

		// Update cache
		statuses := s.ProcessMonitor.GetAllStatuses()
		s.statusCacheMutex.Lock()
		s.statusCache = statuses
		s.statusCacheTime = time.Now()
		s.statusCacheMutex.Unlock()

		json.NewEncoder(w).Encode(statuses)

	case r.URL.Path == "/api/auth" && r.Method == "POST":
		var authData struct {
			PIN string `json:"pin"`
		}
		if err := json.NewDecoder(r.Body).Decode(&authData); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		if s.Config.VerifyWebPIN(authData.PIN) {
			// Generate unique access key
			key := strings.ReplaceAll(uuid.New().String(), "-", "")
			s.bucketMutex.Lock()
			s.keyBucket[key] = time.Now()
			s.bucketMutex.Unlock()

			// Set HttpOnly Cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "aviator_key",
				Value:    key,
				Path:     "/",
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
				MaxAge:   86400, // 24 hours
			})

			json.NewEncoder(w).Encode(map[string]string{
				"status": "success",
			})
		} else {
			http.Error(w, "Invalid PIN", http.StatusUnauthorized)
		}

	case r.URL.Path == "/api/logout" && r.Method == "POST":
		cookie, err := r.Cookie("aviator_key")
		if err == nil {
			s.bucketMutex.Lock()
			delete(s.keyBucket, cookie.Value)
			s.bucketMutex.Unlock()
		}

		// Clear cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "aviator_key",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1,
		})
		w.WriteHeader(http.StatusOK)

	case r.URL.Path == "/api/process-statuses" && r.Method == "GET":
		if !s.isAuthorized(r) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// Return cached statuses if less than 1 second old
		s.statusCacheMutex.RLock()
		cacheAge := time.Since(s.statusCacheTime)
		if cacheAge < time.Second && len(s.statusCache) > 0 {
			cached := make(map[string]bool)
			for k, v := range s.statusCache {
				cached[k] = v
			}
			s.statusCacheMutex.RUnlock()
			json.NewEncoder(w).Encode(cached)
			return
		}
		s.statusCacheMutex.RUnlock()

		// Update cache
		statuses := s.ProcessMonitor.GetAllStatuses()
		s.statusCacheMutex.Lock()
		s.statusCache = statuses
		s.statusCacheTime = time.Now()
		s.statusCacheMutex.Unlock()

		json.NewEncoder(w).Encode(statuses)

	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

func (s *Server) isAuthorized(r *http.Request) bool {
	if !s.Config.GetSettings().AuthEnabled {
		return true
	}

	cookie, err := r.Cookie("aviator_key")
	if err != nil {
		return false
	}

	key := cookie.Value

	s.bucketMutex.RLock()
	lastSeen, exists := s.keyBucket[key]
	s.bucketMutex.RUnlock()

	if !exists {
		return false
	}

	// Session valid for 24 hours
	if time.Since(lastSeen) > 24*time.Hour {
		s.bucketMutex.Lock()
		delete(s.keyBucket, key)
		s.bucketMutex.Unlock()
		return false
	}

	// Update last seen
	s.bucketMutex.Lock()
	s.keyBucket[key] = time.Now()
	s.bucketMutex.Unlock()

	return true
}

func (s *Server) handleLaunch(w http.ResponseWriter, appID string) {
	app, found := s.Config.GetAppByID(appID)
	if !found {
		http.Error(w, `{"error": "App not found"}`, http.StatusNotFound)
		return
	}

	pid, err := launcher.RunExecutableWithTracking(app.ID, app.Name, app.Path, app.Args)
	if err != nil {
		log.Printf("Error launching %s: %v", app.Name, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Launched " + app.Name,
		"pid":     pid,
	})
}
