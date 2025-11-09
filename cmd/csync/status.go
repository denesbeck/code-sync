package main

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
	currentBranch := GetCurrentBranchName()
	color.Cyan("On branch %s.", currentBranch)
	if len(*content) != 0 {
		Debug("Found %d files staged for commit.", len(*content))
		color.Cyan("\nFiles staged for commit:")
		PrintLogs(*content)
	} else {
		Debug("%s", STATUS_RETURN_CODES[501])
	}

	modified, deleted := GetModifiedOrDeletedFiles()
	if len(modified) > 0 || len(deleted) > 0 {
		Debug("Found %d tracked files that have been modified or deleted.", len(modified)+len(deleted))
		color.Cyan("\nChanges not staged for commit:")
		for _, file := range modified {
			color.Yellow("  modified: " + file)
		}
		for _, file := range deleted {
			color.Yellow("  deleted: " + file)
		}
	} else {
		Debug("%s", STATUS_RETURN_CODES[503])
	}

	untracked := GetUntrackedFiles()
	if len(untracked) != 0 {
		color.Cyan("\nUntracked files:")
		for _, file := range untracked {
			color.Yellow("  " + file)
		}
		color.Cyan("(use \"csync add <file>...\" to track)")
	} else {
		Debug("%s", STATUS_RETURN_CODES[504])
	}

	if len(*content) == 0 && len(modified) == 0 && len(deleted) == 0 && len(untracked) == 0 {
		Debug("%s", STATUS_RETURN_CODES[505])
		color.Cyan("\n" + STATUS_RETURN_CODES[505])
	}
	Debug("Status command completed successfully")
	return 502, *content
}
