package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/baris/alf/internal/app"
	"github.com/baris/alf/internal/config"
)

func main() {
	var cfgPath string
	var logLevel string

	flag.StringVar(&cfgPath, "config", "", "Path to a YAML configuration file")
	flag.StringVar(&logLevel, "log-level", "", "Override the configured log level")
	flag.Parse()

	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	if logLevel != "" {
		cfg.LogLevel = logLevel
	}

	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("initialize app: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := application.Run(ctx); err != nil {
		log.Printf("alf stopped: %v", err)
		os.Exit(1)
	}
}
