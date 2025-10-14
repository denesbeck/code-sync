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
	// Check if csync is initialized.
	initialized := IsInitialized()
	if !initialized {
		color.Red("CSync not initialized")
		return nil
	}

	// Get file name from file path.
	_, file := ParsePath(filePath)
	// Generate a random 32 byte long hex string. This will be used as the id for the file in the staging area.
	generatedId := GenRandHex(32)

	// Check if file is already in staging area.
	fileStaged := IsFileStaged(filePath)
	// If it is, we check the log entries to see if the file was added, modified or removed.
	if fileStaged {
		// ADD entry exists? (means the file addition was already staged).
		added, id := LogEntryLookup("ADD", filePath)
		if added {
			// Check if file is modified since ADD entry was created.
			modified := IsModified(filePath, "./.csync/staging/added/"+id+"/"+file)
			if modified {
				// If file is modified, just update the file in staging/added (the id remains the same).
				AddToStaging(id, filePath, "added")
				return nil
			}
			// If file is not modified, do nothing.
			color.Cyan("File already staged")
			return nil
		}
		// MOD entry exists? (means that the change of the file was already staged).
		modified, id := LogEntryLookup("MOD", filePath)
		if modified {
			// Check if file is modified since MOD entry was created.
			modified := IsModified(filePath, "./.csync/staging/modified/"+id+"/"+file)
			if modified {
				// If the file is modified, just update the file in staging/modified
				AddToStaging(id, filePath, "modified")
				color.Cyan("Staged file updated")
				return nil
			}
			// If file is not modified, do nothing
			color.Cyan("File already staged")
			return nil
		}
		// REM entry exists? (means that the removal of the file was already staged).
		removed, id := LogEntryLookup("REM", filePath)
		if removed {
			// If file is removed, check if it was added back.
			exists := FileExists(filePath)
			// If it was...
			if exists {
				// Check if it was committed (exists in the file list of the latest commit).
				isCommitted, commitId, fileId := IsFileCommitted(filePath)
				// If it wasn't committed (THIS SCENARIO SHOULDN'T BE POSSIBLE!), remove the file from staging/removed and add it to staging/added. Log this operation.
				if !isCommitted {
					RemoveFile("./.csync/staging/removed/" + id + file)
					RemoveLogEntry(id)
					AddToStaging(generatedId, filePath, "added")
					LogOperation(generatedId, "ADD", filePath)
					return nil
				} else {
					// If it was committed, check if the file is modified since the last commit.
					modified := IsModified(filePath, "./.csync/commits/"+commitId+"/"+fileId+"/"+filePath)
					// If it was modified...
					if modified {
						// Remove the file from staging/removed and add it to staging/modified. Log this operation.
						RemoveFile("./.csync/staging/removed/" + id + file)
						RemoveLogEntry(id)
						AddToStaging(generatedId, filePath, "modified")
						LogOperation(generatedId, "MOD", filePath)
						return nil
					}
					// If it wasn't modified, remove the file from staging/removed and remove the log entry.
					RemoveFile("./.csync/staging/removed/" + id + file)
					RemoveLogEntry(id)
					return nil
				}
			}
		}
	} else {
		// Check if the file is deleted (exists in the file list of the latest commit but not in the workdir)
		isDeleted := IsFileDeleted(filePath)
		if isDeleted {
			AddToStaging(generatedId, filePath, "removed")
			LogOperation(generatedId, "REM", filePath)
			return nil
		}

		// If it wasn't deleted, check if the file is committed
		isCommitted, commitId, fileId := IsFileCommitted(filePath)
		// If it was committed...
		if isCommitted {
			// Check if the file is modified since the last commit
			modified := IsModified(filePath, "./.csync/commits/"+commitId+"/"+fileId+"/"+filePath)
			// If it was modified...
			if modified {
				// Add the file to staging/modified and log this operation
				AddToStaging(generatedId, filePath, "modified")
				LogOperation(generatedId, "MOD", filePath)
				return nil
			} else {
				// If it wasn't modified, add the file to staging/added and log this operation
				AddToStaging(generatedId, filePath, "added")
				LogOperation(generatedId, "ADD", filePath)
				return nil
			}
		} else {
			// If it wasn't committed, add the file to staging/added and log this operation. This means the file is new.
			AddToStaging(generatedId, filePath, "added")
			LogOperation(generatedId, "ADD", filePath)
			return nil
		}
	}
	return nil
}
