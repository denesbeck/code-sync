package cli

import (
	"log"
	"os"

	"github.com/fatih/color"
)

func AddToStaging(id string, path string, op string) {
	_, file := ParsePath(path)

	if err := os.MkdirAll(".csync/staging/"+op+"/"+id, 0755); err != nil {
		log.Fatal(err)
	}
	_, err := CopyFile(path, ".csync/staging/"+op+"/"+id+"/"+file)
	if err != nil {
		log.Fatal(err)
	}
	color.Green("File added successfully")
}
