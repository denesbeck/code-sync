package cli

import (
	"os"

	"github.com/fatih/color"
)

func PurgeCSync() {
	initialized := IsInitialized()
	if !initialized {
		color.Red("CSync not initialized")
		return
	}
	os.RemoveAll(".csync")
	color.Green("CSync purged successfully")
}
