package cli

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(purgeCmd)
}

var purgeCmd = &cobra.Command{
	Use:     "purge",
	Short:   "Purge CSync and all its data. This command is irreversible!",
	Example: "csync purge",
	Args:    cobra.NoArgs,
	Run: func(_ *cobra.Command, args []string) {
		runPurgeCommand()
	},
}

func runPurgeCommand() {
	initialized := IsInitialized()
	if !initialized {
		color.Red("CSync not initialized")
		return
	}
	os.RemoveAll(".csync")
	color.Green("CSync purged successfully")
}
