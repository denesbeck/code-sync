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
	RunE: func(_ *cobra.Command, args []string) error {
		return runPurgeCommand()
	},
}

func runPurgeCommand() error {
	initialized := IsInitialized()
	if !initialized {
		color.Red("CSync not initialized")
		return nil
	}
	os.RemoveAll(".csync")
	color.Green("CSync purged successfully")
	return nil
}
