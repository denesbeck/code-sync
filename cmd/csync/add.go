package csync

import (
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:     "add",
	Short:   "This command adds the selected files to the staging area.",
	Example: "csync add",
	RunE: func(_ *cobra.Command, args []string) error {
		if len(args) < 1 {
			color.Red("Please specify a file to add")
			return nil
		}
		return runAddCommand(args[0])
	},
}

func runAddCommand(path string) error {
	// check if csync is initialized
	initialized := IsInitialized()
	if !initialized {
		color.Red("CSync not initialized")
		return nil
	}
	fileInLogs := IsFileListed(path, ".csync/staging/logs.json")
	if fileInLogs {
		// TO BE IMPLEMENTED
		return nil
	} else {
		lastCommit, commitExists := GetLastCommit()
		// there is at least one commit
		if commitExists {
			// file should be deleted?
			isDeleted := IsFileDeleted(lastCommit, path)
			if isDeleted {
				// TO BE IMPLEMENTED: move the file from the appr. commit
				MoveToStaging(path, "removed")
				LogOperation("REM", path)
			} else {
				// new file?
				isNewFile := IsFileListed(path, ".csync/commits/"+lastCommit+"/fileList.json")
				if isNewFile {
					exists := FileExists(path)
					if !exists {
						color.Red("File does not exist")
						return nil
					}
					// add file
					MoveToStaging(path, "added")
					LogOperation("ADD", path)
				}
			}
		} else {
			// check if file exists
			exists := FileExists(path)
			if !exists {
				color.Red("File does not exist")
				return nil
			}
			// add file
			MoveToStaging(path, "added")
			LogOperation("ADD", path)
			// there is a commit
		}
	}

	return nil
}
