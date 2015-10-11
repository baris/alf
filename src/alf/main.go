package main

import (
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

const (
	applicationName        = "alf"
	applicationVersion     = "0.0.1"
	applicationDescription = "The bot from the outer space."
)

var log = logrus.WithFields(logrus.Fields{"app": applicationName})

func initLogrus(logLevel string) {
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RFC822,
		FullTimestamp:   true,
	})

	if level, err := logrus.ParseLevel(logLevel); err == nil {
		logrus.SetLevel(level)
	} else {
		log.Error(err)
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func start(ctx *cli.Context) {
	initLogrus(ctx.String("log_level"))
	c, err := readConfig(ctx.String("config"))
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}

	initAlf(c)
	initBrain(c)

	alf.AddHandler(&AlfHandler{})
	alf.AddHandler(&ScriptsHandler{})
	alf.AddHandler(&QuoteHandler{})
	alf.AddHandler(&WhatisHandler{})
	alf.AddHandler(&MediumHandler{})
	alf.start()
}

func main() {
	app := cli.NewApp()
	app.Name = applicationName
	app.Version = applicationVersion
	app.Usage = applicationDescription
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "log_level, l",
			Value: "info",
			Usage: "Logging level (debug, info, warn, error, fatal, panic) (default=info)",
		},
		cli.StringFlag{
			Name:  "config",
			Value: "",
			Usage: "Configuration file",
		},
	}
	app.Action = start
	app.Run(os.Args)
}
