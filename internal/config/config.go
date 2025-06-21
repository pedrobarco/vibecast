package config

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var (
	configPath = filepath.Join(os.Getenv("HOME"), ".config/vibecast/config.yaml")
)

type Playlist struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"` // Can be file path or URL
}

type Config struct {
	Playlists []Playlist `yaml:"playlists"`
	// Favourites is a map: playlist name -> list of channel names
	Favourites map[string][]string `yaml:"favourites,omitempty"`
}

func Load() (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Config{
				Favourites: make(map[string][]string),
			}, nil // Return empty config if not found
		}
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.Favourites == nil {
		cfg.Favourites = make(map[string][]string)
	}
	return &cfg, nil
}

func Save(cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0o644)
}
