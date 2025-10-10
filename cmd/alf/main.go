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
	var brainPath string
	var user string

	flag.StringVar(&cfgPath, "config", "", "Path to a JSON configuration file")
	flag.StringVar(&brainPath, "brain", "", "Override the path to the brain store")
	flag.StringVar(&user, "user", "", "Name recorded when storing memories")
	flag.Parse()

	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	if brainPath != "" {
		cfg.BrainPath = brainPath
	}

	if user != "" {
		cfg.DefaultUser = user
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
