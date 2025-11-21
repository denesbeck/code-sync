package main

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:     "remove",
	Short:   "Remove the selected files from the staging area",
	Example: "nexio remove <path/to/your/file>",
	Args:    cobra.MinimumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		for _, arg := range args {
			runRemoveCommand(arg)
		}
	},
}

func runRemoveCommand(filePath string) {
	initialized := IsInitialized()
	if !initialized {
		color.Red(COMMON_RETURN_CODES[001])
		return
	}

	isLogged, logId, operation := LogEntryLookup("*", filePath)

	ops := map[string]string{
		"ADD": "added",
		"MOD": "modified",
		"REM": "removed",
	}

	if isLogged {
		RemoveFile(dirs.Staging + ops[operation] + "/" + logId)
		RemoveLogEntry(logId)
		color.Green("File removed from staging")
	} else {
		color.Red("File not staged")
	}
}
