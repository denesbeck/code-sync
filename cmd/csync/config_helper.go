package main

import (
	"encoding/json"
	"log"
	"os"
)

func GetConfig() *Config {
	Debug("Reading config file")
	config, err := os.ReadFile(dirs.Config)
	if err != nil {
		Debug("Failed to read config file")
		log.Fatal(err)
	}

	var content Config
	if err = json.Unmarshal(config, &content); err != nil {
		Debug("Failed to unmarshal config")
		log.Fatal(err)
	}

	Debug("Config retrieved successfully: username=%s, email=%s", content.Username, content.Email)
	return &content
}
