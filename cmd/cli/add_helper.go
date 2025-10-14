package cli

import (
	"log"
	"os"

	"github.com/fatih/color"
)

func AddToStaging(id string, path string, op string) {
	Debug("Adding file to staging: id=%s, path=%s, op=%s", id, path, op)
	_, file := ParsePath(path)

	if err := os.MkdirAll(dirs.Staging+op+"/"+id, 0755); err != nil {
		Debug("Failed to create staging directory")
		log.Fatal(err)
	}
	CopyFile(path, dirs.Staging+op+"/"+id+"/"+file)
	Debug("File added to staging successfully")
	color.Green("File added to staging")
}
