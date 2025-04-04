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
		Debug("Starting status command")
		runStatusCommand()
	},
}

func runStatusCommand() (returnCode int, stagingLogs []LogFileEntry) {
	if initialized := IsInitialized(); !initialized {
		color.Red(COMMON_RETURN_CODES[001])
		return 001, nil
	}
	content := GetStagingLogsContent()
	if len(*content) == 0 {
		Debug("No files staged for commit")
		color.Cyan(STATUS_RETURN_CODES[501])
		return 501, nil
	} else {
		Debug("Found %d files staged for commit", len(*content))
		color.Cyan("Files staged for commit:")
		PrintLogs(*content)
	}
	Debug("Status command completed successfully")
	return 502, *content
}
