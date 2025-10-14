package cli

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:     "status",
	Short:   "List the files that are staged for commit",
	Example: "csync status",
	Args:    cobra.NoArgs,
	Run: func(_ *cobra.Command, args []string) {
		runStatusCommand()
	},
}

func runStatusCommand() {
	if initialized := IsInitialized(); !initialized {
		color.Red("CSync not initialized")
		return
	}
	content := GetStagingLogsContent()
	if len(content) == 0 {
		color.Cyan("No files staged for commit")
	} else {
		color.Cyan("Files staged for commit:")
		PrintLogs(content)
	}
}
