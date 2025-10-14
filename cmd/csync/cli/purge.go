package cli

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(purgeCmd)
}

var purgeCmd = &cobra.Command{
	Use:     "purge",
	Short:   "This command purges CSync and all its data. This command is irreversible.",
	Example: "csync purge",
	RunE: func(_ *cobra.Command, args []string) error {
		return runPurgeCommand()
	},
}

func runPurgeCommand() error {
	PurgeCSync()
	return nil
}
