package cli

import (
	"os"
)

func calculateNamespace() string {
	environment := os.Getenv("CSYNC_ENV")
	if environment == "test" {
		return "__" + environment + "__/"
	}
	return ""
}

var namespace = calculateNamespace()
