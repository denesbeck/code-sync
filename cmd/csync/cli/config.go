package cli

import (
	"encoding/json"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configCmd)
}

type Config struct {
	Username string
	Email    string
}

var configCmd = &cobra.Command{
	Use:     "config",
	Short:   "Config manager",
	Example: "csync config",
	RunE: func(_ *cobra.Command, args []string) error {
		if args[0] == "set" {
			setConfig(args[1], args[2])
			return nil
		}
		if args[0] == "get" {
			getConfig(args[1])
			return nil
		}
		color.Red("Invalid command")
		return nil
	},
}

func setConfig(key string, value string) {
	if initialized := IsInitialized(); !initialized {
		color.Red("CSync not initialized")
	}
	const PATH = "./.csync/config.json"
	config, err := os.ReadFile(PATH)
	if err != nil {
		log.Fatal(err)
	}

	var content Config
	if err = json.Unmarshal(config, &content); err != nil {
		log.Fatal(err)
	}

	switch key {
	case "username":
		content.Username = value
	case "email":
		content.Email = value
	}

	jsonData, err := json.Marshal(content)
	if err != nil {
		log.Fatal(err)
	}

	if err = os.WriteFile(PATH, jsonData, 0644); err != nil {
		log.Fatal(err)
	}
}

func getConfig(key string) {
	if initialized := IsInitialized(); !initialized {
		color.Red("CSync not initialized")
	}
	config := GetConfig()
	switch key {
	case "username":
		color.Cyan(config.Username)
	case "email":
		color.Cyan(config.Email)
	case "user":
		color.Cyan(config.Username + " <" + config.Email + ">")
	}
}
