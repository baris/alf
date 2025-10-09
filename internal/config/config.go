package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config captures runtime configuration for the alf bot.
type Config struct {
	// AppToken is the Slack App-level token (starts with "xapp").
	AppToken string `yaml:"app_token"`
	// BotToken is the Slack bot token (starts with "xoxb").
	BotToken string `yaml:"bot_token"`
	// SigningSecret validates inbound HTTP requests when using Events API.
	SigningSecret string `yaml:"signing_secret"`
	// LogLevel controls the verbosity of the structured logger.
	LogLevel string `yaml:"log_level"`
	// BrainPath points to the JSON file that persists conversational memory.
	BrainPath string `yaml:"brain_path"`
	// MetricsAddr exposes Prometheus metrics and health probes when set.
	MetricsAddr string `yaml:"metrics_addr"`
	// HTTPAddr optionally enables the Events API receiver when configured.
	HTTPAddr string `yaml:"http_addr"`
	// SocketMode toggles Slack Socket Mode support.
	SocketMode bool `yaml:"socket_mode"`
}

// Defaults returns a Config pre-populated with sensible defaults.
func Defaults() Config {
	return Config{
		LogLevel:    "info",
		BrainPath:   "data/brain.json",
		MetricsAddr: "",
		HTTPAddr:    "",
		SocketMode:  true,
	}
}

// Load returns a Config composed from defaults, an optional YAML file and
// environment variable overrides.
func Load(path string) (Config, error) {
	cfg := Defaults()

	if path != "" {
		if err := fromFile(path, &cfg); err != nil {
			return cfg, err
		}
	}

	overlayEnv(&cfg)

	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}

	if cfg.BrainPath == "" {
		cfg.BrainPath = "data/brain.json"
	}

	if !filepath.IsAbs(cfg.BrainPath) {
		cfg.BrainPath = filepath.Clean(cfg.BrainPath)
	}

	if cfg.SocketMode && cfg.AppToken == "" {
		return cfg, errors.New("app token is required when socket mode is enabled")
	}

	if cfg.BotToken == "" {
		return cfg, errors.New("bot token is required")
	}

	return cfg, nil
}

func fromFile(path string, cfg *Config) error {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	return nil
}

func overlayEnv(cfg *Config) {
	set := func(env string, dest *string) {
		if value := strings.TrimSpace(os.Getenv(env)); value != "" {
			*dest = value
		}
	}

	set("ALF_APP_TOKEN", &cfg.AppToken)
	set("ALF_BOT_TOKEN", &cfg.BotToken)
	set("ALF_SIGNING_SECRET", &cfg.SigningSecret)
	set("ALF_LOG_LEVEL", &cfg.LogLevel)
	set("ALF_BRAIN_PATH", &cfg.BrainPath)
	set("ALF_METRICS_ADDR", &cfg.MetricsAddr)
	set("ALF_HTTP_ADDR", &cfg.HTTPAddr)

	if v := strings.TrimSpace(os.Getenv("ALF_SOCKET_MODE")); v != "" {
		if parsed, err := strconv.ParseBool(v); err == nil {
			cfg.SocketMode = parsed
		}
	}
}
