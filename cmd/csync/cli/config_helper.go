package cli

import (
	"encoding/json"
	"log"
	"os"
)

func GetConfig() Config {
	const PATH = "./.csync/config.json"
	config, err := os.ReadFile(PATH)
	if err != nil {
		log.Fatal(err)
	}

	var content Config
	if err = json.Unmarshal(config, &content); err != nil {
		log.Fatal(err)
	}

	return content
}
