package cli

import (
	"log"
	"os"

	"github.com/fatih/color"
)

func AddToStaging(id string, path string, op string) {
	_, file := ParsePath(path)

	if err := os.MkdirAll(dirs.Staging+op+"/"+id, 0755); err != nil {
		log.Fatal(err)
	}
	CopyFile(path, dirs.Staging+op+"/"+id+"/"+file)
	color.Green("File added to staging")
}
