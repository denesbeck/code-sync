package main

import (
	"os"
)

func calculateNamespace() string {
	environment := os.Getenv("NEXIO_ENV")
	if environment == "test" {
		return "__" + environment + "__/"
	}
	return ""
}

var namespace = calculateNamespace()
