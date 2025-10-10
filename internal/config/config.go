package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config captures runtime configuration for the alf console bot.
type Config struct {
	BrainPath   string `json:"brain_path"`
	DefaultUser string `json:"default_user"`
	Greeting    string `json:"greeting"`
	Prompt      string `json:"prompt"`
}

// Defaults returns a Config pre-populated with sensible defaults.
func Defaults() Config {
	return Config{
		BrainPath:   "data/brain.json",
		DefaultUser: "operator",
		Greeting:    "alf ready. type 'help' to explore commands or 'exit' to quit.",
		Prompt:      "> ",
	}
}

// Load returns a Config composed from defaults, an optional JSON file, and
// environment variable overrides.
func Load(path string) (Config, error) {
	cfg := Defaults()

	if path != "" {
		if err := fromFile(path, &cfg); err != nil {
			return cfg, err
		}
	}

	overlayEnv(&cfg)

	if cfg.BrainPath == "" {
		cfg.BrainPath = Defaults().BrainPath
	}

	if !filepath.IsAbs(cfg.BrainPath) {
		cfg.BrainPath = filepath.Clean(cfg.BrainPath)
	}

	if cfg.DefaultUser == "" {
		cfg.DefaultUser = Defaults().DefaultUser
	}

	if cfg.Prompt == "" {
		cfg.Prompt = Defaults().Prompt
	}

	return cfg, nil
}

func fromFile(path string, cfg *Config) error {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	return nil
}

func overlayEnv(cfg *Config) {
	if value := strings.TrimSpace(os.Getenv("ALF_BRAIN_PATH")); value != "" {
		cfg.BrainPath = value
	}

	if value := strings.TrimSpace(os.Getenv("ALF_DEFAULT_USER")); value != "" {
		cfg.DefaultUser = value
	}

	if value := os.Getenv("ALF_GREETING"); strings.TrimSpace(value) != "" {
		cfg.Greeting = value
	}

	if value := os.Getenv("ALF_PROMPT"); strings.TrimSpace(value) != "" {
		cfg.Prompt = value
	}
}
