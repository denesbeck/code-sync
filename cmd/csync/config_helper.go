package main

import (
	"encoding/json"
	"os"
)

func GetConfig() *Config {
	Debug("Reading config file")
	config, err := os.ReadFile(dirs.Config)
	if err != nil {
		Debug("Failed to read config file")
		MustSucceed(err, "failed to read config file")
	}

	var content Config
	if err = json.Unmarshal(config, &content); err != nil {
		Debug("Failed to unmarshal config")
		MustSucceed(err, "failed to parse config file")
	}

	Debug("Config retrieved successfully: username=%s, email=%s", content.Username, content.Email)
	return &content
}
