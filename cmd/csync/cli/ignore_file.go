package cli

import (
	"encoding/json"
	"log"
	"os"

	"github.com/fatih/color"
)

func readCsyncIgnore() []string {
	_, err := os.Stat(".csyncignore.json")
	if os.IsNotExist(err) {
		color.Cyan("INFO: .csyncignore.json not found")
		return []string{}
	}
	if err != nil {
		log.Fatal(err)
		return nil
	}
	var payload []string
	content, err := os.ReadFile(".csyncignore.json")
	if err != nil {
		log.Fatal("Error while opening file: ", err)
	}
	if err = json.Unmarshal(content, &payload); err != nil {
		log.Fatal("Error while parsing data: ", err)
	}
	color.Cyan("INFO: .csyncignore.json found")
	return payload
}
