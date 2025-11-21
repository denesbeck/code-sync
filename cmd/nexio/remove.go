package main

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:     "remove",
	Short:   "Remove the selected files from the staging area",
	Example: "nexio remove <path/to/your/file>",
	Args:    cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		for _, arg := range args {
			runRemoveCommand(arg)
		}
	},
}

func runRemoveCommand(filePath string) {
	initialized := IsInitialized()
	if !initialized {
		Fail(COMMON_RETURN_CODES[001])
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
		Success(REMOVE_RETURN_CODES[801])
	} else {
		Info(REMOVE_RETURN_CODES[802])
	}
}
