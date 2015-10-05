package main

import (
	"encoding/json"
	"io/ioutil"
)

const (
	defaultName           = "alf"
	defaultUpdateInterval = 60
)

// Config type holds the global configuration values
type Config struct {
	Name           string `json:"name"`
	DatabaseFile   string `json:"databaseFile"`
	SlackToken     string `json:"slackToken"`
	DefaultChannel string `json:"defaultChannel"`
	UpdateInterval int    `json:"updateInterval"`
}

func defaultConfig() Config {
	return Config{
		Name:           defaultName,
		DatabaseFile:   "alf.db",
		SlackToken:     "",
		DefaultChannel: "general",
		UpdateInterval: defaultUpdateInterval,
	}
}

func readConfig(configFile string) (Config, error) {
	c := defaultConfig()
	if configFile == "" {
		log.Error("No configuration file given.")
		return c, nil
	}
	log.Info("Reading configuration file at ", configFile)
	contents, e := ioutil.ReadFile(configFile)
	if e != nil {
		log.Error("Config file error: ", e)
		return c, e
	}
	err := json.Unmarshal(contents, &c)
	if err != nil {
		log.Error("Invalid JSON in config: ", err)
		return c, err
	}
	return c, nil
}
