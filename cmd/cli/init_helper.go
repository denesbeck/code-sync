package cli

import (
	"os"
)

func IsInitialized() bool {
	Debug("Checking if CSync is initialized")
	if _, err := os.Stat(dirs.Root); !os.IsNotExist(err) {
		Debug("CSync is initialized")
		return true
	}
	Debug("CSync is not initialized")
	return false
}
