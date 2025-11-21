package main

import (
	"os"
)

func IsInitialized() bool {
	Debug("Checking if Nexio is initialized")
	if _, err := os.Stat(dirs.Root); !os.IsNotExist(err) {
		Debug("Nexio is initialized")
		return true
	}
	Debug("Nexio is not initialized")
	return false
}
