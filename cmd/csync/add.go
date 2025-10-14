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
	// Check if file exists
	exists := FileExists(filePath)
	if !exists {
		color.Red("File does not exist")
		return nil
	}
	// Get file name from file path
	_, file := ParsePath(filePath)
	// Generate a random 32 byte long hex string
	generatedId := GenRandHex(32)
	// Check if there is at least one commit registered
	latestCommitId, commitExists := GetLastCommit()

	// Check if file is already in staging area
	fileStaged := IsFileStaged(filePath)
	if fileStaged {
		added, id := LogEntryLookup("ADD", filePath)
		if added {
			modified := IsModified(filePath, "./.csync/staging/added/"+id+"/"+file)
			if modified {
				AddToStaging(id, filePath, "added")
				return nil
			}
			return nil
		}
		modified, id := LogEntryLookup("MOD", filePath)
		if modified {
			modified := IsModified(filePath, "./.csync/staging/modified/"+id+"/"+file)
			if modified {
				AddToStaging(id, filePath, "modified")
				return nil
			}
		}
		removed, id := LogEntryLookup("REM", filePath)
		if removed {
			fileExists := FileExists("./.csync/staging/removed/" + id + "/" + file)
			if fileExists {
				if !commitExists {
					RemoveFile("./.csync/staging/removed/" + id + "/" + file)
					RemoveLogEntry(id)
					AddToStaging(id, filePath, "added")
					LogOperation(generatedId, "ADD", filePath)
				} else {
					modified := IsModified(filePath, "./.csync/commits/"+latestCommitId+"/"+file)
					if modified {
						RemoveFile("./.csync/staging/removed/" + id + "/" + file)
						RemoveLogEntry(id)
						AddToStaging(id, filePath, "modified")
						LogOperation(generatedId, "MOD", filePath)
					} else {
						RemoveFile("./.csync/staging/removed/" + id + "/" + file)
						RemoveLogEntry(id)
					}
				}
			}
			return nil
		}
	} else {
		if commitExists {
			// File should be deleted? Check if it is listed in the latest commit and missing from the working directory
			shouldBeDeleted, srcCommitId := IsFileDeleted(filePath, latestCommitId)
			if shouldBeDeleted {
				AddToStaging(generatedId, "./.csync/commits/"+srcCommitId+"/files/"+filePath, "removed")
				LogOperation(generatedId, "REM", filePath)
			} else {
				fileCommitted, _ := IsFileCommitted(filePath, latestCommitId)
				// Is it a new file? Check if it was listed in the latest commit
				if !fileCommitted {
					// Add file to staging if it was not listed in the latest commit
					AddToStaging(generatedId, filePath, "added")
					LogOperation(generatedId, "ADD", filePath)
				}
			}
		} else {
			// Add file to staging
			AddToStaging(generatedId, filePath, "added")
			LogOperation(generatedId, "ADD", filePath)
		}
	}
	return nil
}
