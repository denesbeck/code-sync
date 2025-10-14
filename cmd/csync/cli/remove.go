package cli

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:     "rm",
	Short:   "This command removes the selected files from the staging area",
	Example: "csync rm <path/to/your/file>",
	RunE: func(_ *cobra.Command, args []string) error {
		if len(args) < 1 {
			color.Red("Please specify a file to remove")
			return nil
		}
		return runRemoveCommand(args[0])
	},
}

func runRemoveCommand(filePath string) error {
	initialized := IsInitialized()
	if !initialized {
		color.Red("CSync not initialized")
		return nil
	}

	isLogged, logId, operation := LogEntryLookup("*", filePath)

	// create ops struct to translate operation to human readable string
	ops := map[string]string{
		"ADD": "added",
		"MOD": "modified",
		"REM": "removed",
	}

	if isLogged {
		RemoveFile("./.csync/staging/" + ops[operation] + "/" + logId)
		RemoveLogEntry(logId)
		color.Green("File removed from staging")
		return nil
	} else {
		color.Red("File not staged")
		return nil
	}
}
