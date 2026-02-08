package config

import (
	"aviator-wails/internal/icons"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
)

type App struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
	Args string `json:"args"`
	Icon string `json:"icon,omitempty"` // Base64 encoded PNG icon
}

type Settings struct {
	AutoStart bool `json:"auto_start"`
}

type ConfigManager struct {
	Apps         []App
	Settings     Settings
	FilePath     string // config.json (apps)
	SettingsPath string // settings.json (preferences)
	mu           sync.RWMutex
}

func NewConfigManager() (*ConfigManager, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	// AppData/Local/Aviator usually. On Windows UserConfigDir is AppData/Roaming by default
	// but Python used LOCALAPPDATA. Let's try to match it.
	// os.UserConfigDir() on Windows returns %AppData% (Roaming).
	// We want local AppData.
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		localAppData = configDir // Fallback
	}

	aviatorDir := filepath.Join(localAppData, "Aviator")
	if err := os.MkdirAll(aviatorDir, 0755); err != nil {
		return nil, err
	}

	cm := &ConfigManager{
		Apps:         []App{},
		Settings:     Settings{AutoStart: false},
		FilePath:     filepath.Join(aviatorDir, "config.json"),
		SettingsPath: filepath.Join(aviatorDir, "settings.json"),
	}

	cm.Load()         // Ignore error on load apps
	cm.LoadSettings() // Ignore error on load settings
	return cm, nil
}

func (cm *ConfigManager) Load() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	data, err := os.ReadFile(cm.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			cm.Apps = []App{}
			return nil
		}
		return err
	}

	return json.Unmarshal(data, &cm.Apps)
}

func (cm *ConfigManager) Save() error {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	data, err := json.MarshalIndent(cm.Apps, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(cm.FilePath, data, 0644)
}

func (cm *ConfigManager) AddApp(name, path, args string) App {
	cm.mu.Lock()

	// Extract icon from executable
	iconBase64, err := icons.ExtractIconToBase64(path)
	if err != nil {
		log.Printf("Warning: Could not extract icon from %s: %v", path, err)
		iconBase64 = "" // Empty icon if extraction fails
	}

	app := App{
		ID:   uuid.New().String(),
		Name: name,
		Path: path,
		Args: args,
		Icon: iconBase64,
	}
	cm.Apps = append(cm.Apps, app)
	cm.mu.Unlock()

	cm.Save()
	return app
}

func (cm *ConfigManager) UpdateApp(id, name, path, args string) bool {
	cm.mu.Lock()
	found := false
	for i, app := range cm.Apps {
		if app.ID == id {
			// If path changed, extract new icon
			if app.Path != path {
				iconBase64, err := icons.ExtractIconToBase64(path)
				if err != nil {
					log.Printf("Warning: Could not extract icon from %s: %v", path, err)
					iconBase64 = app.Icon // Keep old icon if extraction fails
				}
				cm.Apps[i].Icon = iconBase64
			}
			cm.Apps[i].Name = name
			cm.Apps[i].Path = path
			cm.Apps[i].Args = args
			found = true
			break
		}
	}
	cm.mu.Unlock()

	if found {
		cm.Save()
	}
	return found
}

func (cm *ConfigManager) RemoveApp(id string) {
	cm.mu.Lock()
	newApps := []App{}
	for _, app := range cm.Apps {
		if app.ID != id {
			newApps = append(newApps, app)
		}
	}
	cm.Apps = newApps
	cm.mu.Unlock()

	cm.Save()
}

func (cm *ConfigManager) GetApps() []App {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	// Return copy to avoid race conditions if caller modifies it
	appsCopy := make([]App, len(cm.Apps))
	copy(appsCopy, cm.Apps)
	return appsCopy
}

func (cm *ConfigManager) GetAppByID(id string) (App, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	for _, app := range cm.Apps {
		if app.ID == id {
			return app, true
		}
	}
	return App{}, false
}

// Settings Management

func (cm *ConfigManager) LoadSettings() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	data, err := os.ReadFile(cm.SettingsPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Start with default settings
			cm.Settings = Settings{AutoStart: IsAutoStartEnabled()}
			return nil
		}
		return err
	}

	err = json.Unmarshal(data, &cm.Settings)

	// Sync with Registry just in case of mismatch on startup
	// If config says true but registry false -> trust registry? Or config?
	// Let's trust registry as source of truth for "active" state,
	// but config file as persistence for user intention.
	// Actually, easier: Read registry state into memory.
	cm.Settings.AutoStart = IsAutoStartEnabled()

	return err
}

func (cm *ConfigManager) SaveSettings() error {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	data, err := json.MarshalIndent(cm.Settings, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(cm.SettingsPath, data, 0644)
}

func (cm *ConfigManager) GetSettings() Settings {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.Settings
}

func (cm *ConfigManager) UpdateSettings(s Settings) error {
	cm.mu.Lock()
	cm.Settings = s
	cm.mu.Unlock()

	// Apply System Changes
	if err := SetAutoStart(s.AutoStart); err != nil {
		return err
	}

	return cm.SaveSettings()
}
