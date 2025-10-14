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
	Short:   "Add the selected files to the staging area",
	Example: "csync add <path/to/your/file>",
	Args:    cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		runAddCommand(args[0])
	},
}

func runAddCommand(filePath string) {
	initialized := IsInitialized()
	if !initialized {
		color.Red("CSync not initialized")
		return
	}

	_, fileName := ParsePath(filePath)
	generatedId := GenRandHex(20)

	fileStaged := IsFileStaged(filePath)
	if fileStaged {
		added, id, _ := LogEntryLookup("ADD", filePath)
		if added {
			exists := FileExists(filePath)
			if !exists {
				RemoveFile("./.csync/staging/added/" + id)
				RemoveLogEntry(id)
				color.Cyan("File removed from staging")
				return
			}
			modified := IsModified(filePath, "./.csync/staging/added/"+id+"/"+fileName)
			if modified {
				AddToStaging(id, filePath, "added")
			}
			color.Cyan("File already staged")
			return
		}
		modified, id, _ := LogEntryLookup("MOD", filePath)
		if modified {
			exists := FileExists(filePath)
			if !exists {
				RemoveFile("./.csync/staging/modified/" + id)
				RemoveLogEntry(id)
				AddToStaging(generatedId, filePath, "removed")
				LogOperation(generatedId, "REM", filePath)
				return
			}
			modified := IsModified(filePath, "./.csync/staging/modified/"+id+"/"+fileName)
			if modified {
				AddToStaging(id, filePath, "modified")
				color.Cyan("Staged file updated")
			}
			color.Cyan("File already staged")
			return
		}
		removed, id, _ := LogEntryLookup("REM", filePath)
		if removed {
			exists := FileExists(filePath)
			if exists {
				RemoveFile("./.csync/staging/removed/" + id)
				RemoveLogEntry(id)
				isCommitted, commitId, fileId := GetFileMetadata(filePath)
				if !isCommitted {
					AddToStaging(generatedId, filePath, "added")
					LogOperation(generatedId, "ADD", filePath)
				} else {
					modified := IsModified(filePath, "./.csync/commits/"+commitId+"/"+fileId+"/"+filePath)
					if modified {
						AddToStaging(generatedId, filePath, "modified")
						LogOperation(generatedId, "MOD", filePath)
					}
				}
			} else {
				color.Cyan("File already staged")
			}
		}
	} else {
		isCommitted, commitId, fileId := GetFileMetadata(filePath)
		isDeleted := IsFileDeleted(filePath)
		if isDeleted {
			AddToStaging(generatedId, "./.csync/commits/"+commitId+"/"+fileId+"/"+fileName, "removed")
			LogOperation(generatedId, "REM", filePath)
			return
		}

		if isCommitted {
			modified := IsModified(filePath, "./.csync/commits/"+commitId+"/"+fileId+"/"+fileName)
			if modified {
				AddToStaging(generatedId, filePath, "modified")
				LogOperation(generatedId, "MOD", filePath)
			} else {
				color.Red("File not modified")
			}
		} else {
			AddToStaging(generatedId, filePath, "added")
			LogOperation(generatedId, "ADD", filePath)
		}
	}
}
