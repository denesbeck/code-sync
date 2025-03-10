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
	Short:   "Purge CSync and all its data. THIS COMMAND IS IRREVERSIBLE!",
	Example: "csync purge",
	Args:    cobra.NoArgs,
	Run: func(_ *cobra.Command, args []string) {
		runPurgeCommand()
	},
}

func runPurgeCommand() {
	initialized := IsInitialized()
	if !initialized {
		color.Red(COMMON_RETURN_CODES[001])
		return
	}
	if namespace == "" {
		os.RemoveAll(dirs.Root)
	} else {
		os.RemoveAll(namespace)
	}

	color.Green("CSync purged successfully")
}
