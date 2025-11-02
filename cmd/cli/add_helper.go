package cli

import (
	"os"
)

func AddToStaging(id string, path string, op string) error {
	Debug("Adding file to staging: id=%s, path=%s, op=%s", id, path, op)
	_, file := ParsePath(path)

	if err := os.MkdirAll(dirs.Staging+op+"/"+id, 0755); err != nil {
		Debug("Failed to create staging directory")
		return err
	}
	if err := CopyFile(path, dirs.Staging+op+"/"+id+"/"+file); err != nil {
		return err
	}
	Debug("File added to staging successfully")
	return nil
}
