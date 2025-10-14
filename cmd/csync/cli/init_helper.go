package cli

import (
	"os"
)

func IsInitialized() bool {
	if _, err := os.Stat(".csync"); !os.IsNotExist(err) {
		return true
	}
	return false
}
