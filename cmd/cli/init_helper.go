package cli

import (
	"os"
)

func IsInitialized() bool {
	if _, err := os.Stat(dirs.Root); !os.IsNotExist(err) {
		return true
	}
	return false
}
