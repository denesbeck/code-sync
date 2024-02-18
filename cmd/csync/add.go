package csync

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:     "add",
	Short:   "This command adds the selected files to the staging area",
	Example: "csync add",
	RunE: func(_ *cobra.Command, args []string) error {
		if len(args) < 1 {
			color.Red("Please specify a file to add")
			return nil
		}
		return runAddCommand(args[0])
	},
}

func runAddCommand(filePath string) error {
	// Check if csync is initialized
	initialized := IsInitialized()
	if !initialized {
		color.Red("CSync not initialized")
		return nil
	}
	// Check if file is already in staging area
	fileStaged := IsFileStaged(filePath)
	if fileStaged {
		added := LogEntryLookup("ADD", filePath)
		if added {
			modified := IsModified(filePath, "./.csync/staging/added/"+filePath)
			if modified {
				AddToStaging(filePath, "added")
				return nil
			}
			return nil
		}
		modified := LogEntryLookup("MOD", filePath)
		if modified {
			modified := IsModified(filePath, "./.csync/staging/modified/"+filePath)
			if modified {
				AddToStaging(filePath, "modified")
				return nil
			}
		}
		removed := LogEntryLookup("REM", filePath)
		if removed {
			// TO BE IMPLEMENTED
			return nil
		}
	} else {
		// Check if there is at least one commit registered
		latestCommitId, commitExists := GetLastCommit()

		if commitExists {
			// File should be deleted? Check if it is listed in the latest commit and missing from the working directory
			shouldBeDeleted, srcCommitId := IsFileDeleted(filePath, latestCommitId)
			if shouldBeDeleted {
				AddToStaging("./.csync/commits/"+srcCommitId+"/files/"+filePath, "removed")
				LogOperation("REM", filePath)
			} else {
				fileCommitted, _ := IsFileCommitted(filePath, latestCommitId)
				// Is it a new file? Check if it was listed in the latest commit
				if !fileCommitted {
					exists := FileExists(filePath)
					if !exists {
						color.Red("File does not exist")
						return nil
					}
					// Add file to staging if it was not listed in the latest commit
					AddToStaging(filePath, "added")
					LogOperation("ADD", filePath)
				}
			}
		} else {
			// Check if file exists
			exists := FileExists(filePath)
			if !exists {
				color.Red("File does not exist")
				return nil
			}
			// Add file to staging
			AddToStaging(filePath, "added")
			LogOperation("ADD", filePath)
		}
	}
	return nil
}
