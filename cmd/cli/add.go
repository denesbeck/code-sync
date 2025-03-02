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
	Args:    cobra.MinimumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		for _, arg := range args {
			runAddCommand(arg)
		}
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
		exists := FileExists(filePath)
		added, id, _ := LogEntryLookup("ADD", filePath)
		if added {
			if !exists {
				removeFileAndLog(id, "added")
				color.Cyan("File removed from staging")
				return
			}
			modified := IsModified(filePath, dirs.StagingAdded+id+"/"+fileName)
			if modified {
				AddToStaging(id, filePath, "added")
				color.Cyan("Staged file updated")
				return
			}
			color.Cyan("File already staged")
			return
		}
		modified, id, _ := LogEntryLookup("MOD", filePath)
		if modified {
			if !exists {
				removeFileAndLog(id, "modified")
				stageAndLog(generatedId, filePath, "removed")
				return
			}
			modified := IsModified(filePath, dirs.StagingModified+id+"/"+fileName)
			if modified {
				AddToStaging(id, filePath, "modified")
				color.Cyan("Staged file updated")
				return
			}
			color.Cyan("File already staged")
			return
		}
		removed, id, _ := LogEntryLookup("REM", filePath)
		if removed {
			if exists {
				removeFileAndLog(id, "removed")
				isCommitted, commitId, fileId := GetFileMetadata(filePath)
				if !isCommitted {
					stageAndLog(generatedId, filePath, "added")
				} else {
					modified := IsModified(filePath, dirs.Commits+commitId+"/"+fileId+"/"+filePath)
					if modified {
						stageAndLog(generatedId, filePath, "modified")
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
			AddToStaging(generatedId, dirs.Commits+commitId+"/"+fileId+"/"+fileName, "removed")
			LogOperation(generatedId, "REM", filePath)
			return
		}

		if isCommitted {
			modified := IsModified(filePath, dirs.Commits+commitId+"/"+fileId+"/"+fileName)
			if modified {
				stageAndLog(generatedId, filePath, "modified")
			} else {
				color.Red("File not modified")
			}
		} else {
			stageAndLog(generatedId, filePath, "added")
		}
	}
}

func removeFileAndLog(id string, op string) {
	RemoveFile(dirs.Staging + op + "/" + id)
	RemoveLogEntry(id)
}

func stageAndLog(id string, path string, op string) {
	logOperations := map[string]string{
		"added":    "ADD",
		"modified": "MOD",
		"removed":  "REM",
	}
	AddToStaging(id, path, op)
	LogOperation(id, logOperations[op], path)
}
