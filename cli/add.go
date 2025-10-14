package cli

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
		// ADD entry exists?
		added, id := LogEntryLookup("ADD", filePath)
		if added {
			// Check if file is modified since ADD entry was created
			modified := IsModified(filePath, "./.csync/staging/added/"+id+"/"+file)
			if modified {
				// If file is modified, just update the file in staging/added
				AddToStaging(id, filePath, "added")
				return nil
			}
			// If file is not modified, do nothing
			return nil
		}
		// MOD entry exists?
		modified, id := LogEntryLookup("MOD", filePath)
		if modified {
			// Check if file is modified since MOD entry was created
			modified := IsModified(filePath, "./.csync/staging/modified/"+id+"/"+file)
			if modified {
				// If file is modified, just update the file in staging/modified
				AddToStaging(id, filePath, "modified")
				return nil
			}
			// If file is not modified, do nothing
			return nil
		}
		// REM entry exists?
		removed, id := LogEntryLookup("REM", filePath)
		if removed {
			// Check if file is still removed
			fileExists := FileExists("./.csync/staging/removed/" + id + "/" + file)
			// If file is not removed, it means that it has been added back, therefore...
			if fileExists {
				// Check if file is listed in any of the commits
				if !commitExists {
					// If the file is not listed in any of the commits, it means that it's a new file, therefore...
					//... remove the file from staging/removed and add it to staging/added
					RemoveFile("./.csync/staging/removed/" + id + "/" + file)
					// ... remove the log entry
					RemoveLogEntry(id)
					// ... add the file to staging/added
					AddToStaging(id, filePath, "added")
					// ... log the ADD operation
					LogOperation(generatedId, "ADD", filePath)
				} else {
					// If the file is listed in any of the commits, it means that it's already tracked, therefore...
					modified := IsModified(filePath, "./.csync/commits/"+latestCommitId+"/"+file)
					// ... check if the file is modified
					if modified {
						// If the file was modified, remove it from staging/removed and add it to staging/modified
						// Also, remove the REM log entry and add a new MOD log entry
						RemoveFile("./.csync/staging/removed/" + id + "/" + file)
						RemoveLogEntry(id)
						AddToStaging(id, filePath, "modified")
						LogOperation(generatedId, "MOD", filePath)
					} else {
						// If the file was not modified, remove it from staging/removed and remove the REM log entry, because
						// the file is not removed anymore and nothing has changed
						RemoveFile("./.csync/staging/removed/" + id + "/" + file)
						RemoveLogEntry(id)
					}
				}
			}
			// If file is still removed, do nothing
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
