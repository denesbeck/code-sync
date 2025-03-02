package cli

import (
	"encoding/json"
	"log"
	"os"

	"github.com/fatih/color"
)

func readCsyncIgnore() *[]string {
	_, err := os.Stat(".csyncignore.json")
	if os.IsNotExist(err) {
		color.Cyan("`.csyncignore.json` not found")
		return nil
	}
	if err != nil {
		log.Fatal(err)
		return nil
	}
	var content []string
	ignoreFile, err := os.ReadFile(".csyncignore.json")
	if err != nil {
		log.Fatal("Error while opening file: ", err)
	}
	if err = json.Unmarshal(ignoreFile, &content); err != nil {
		log.Fatal("Error while parsing data: ", err)
	}
	color.Cyan("`.csyncignore.json` found")
	return &content
}
