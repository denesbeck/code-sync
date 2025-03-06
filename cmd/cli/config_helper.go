package cli

import (
	"encoding/json"
	"log"
	"os"
)

func GetConfig() *Config {
	config, err := os.ReadFile(dirs.Config)
	if err != nil {
		log.Fatal(err)
	}

	var content Config
	if err = json.Unmarshal(config, &content); err != nil {
		log.Fatal(err)
	}

	return &content
}
