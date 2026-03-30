package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Config holds all user preferences persisted to disk.
type Config struct {
	APIKey string `json:"api_key"`
	Units  string `json:"units"` // metric | imperial | standard
}

// Returns .weather-cli/config.json on all platforms.
func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".weather-cli", "config.json"), nil
}

// Load reads the config file from disk.
// Returns a default Config (units: metric) if the file does not exist yet.
func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Config{Units: "metric"}, nil
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.Units == "" {
		cfg.Units = "metric"
	}
	return &cfg, nil
}

// Save writes the config to disk with 0600 permissions.
func Save(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}
