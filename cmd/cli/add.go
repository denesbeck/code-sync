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

func runAddCommand(filePath string) int {
	initialized := IsInitialized()
	if !initialized {
		color.Red(COMMON_RETURN_CODES[99])
		return 99
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
				color.Cyan(ADD_RETURN_CODES[1])
				return 1
			}
			modified := IsModified(filePath, dirs.StagingAdded+id+"/"+fileName)
			if modified {
				AddToStaging(id, filePath, "added")
				color.Cyan(ADD_RETURN_CODES[2])
				return 2
			}
			color.Cyan(ADD_RETURN_CODES[3])
			return 3
		}
		modified, id, _ := LogEntryLookup("MOD", filePath)
		if modified {
			if !exists {
				removeFileAndLog(id, "modified")
				/* INFO:
								        * As the file is deleted, we should only log the operation.
												* The file should not be copied to the staging directory.
				                * Therefore calling LogOperation() and skipping AddToStaging().
				*/
				LogOperation(generatedId, "REM", filePath)
				return 4
			}
			modified := IsModified(filePath, dirs.StagingModified+id+"/"+fileName)
			if modified {
				AddToStaging(id, filePath, "modified")
				color.Cyan(ADD_RETURN_CODES[5])
				return 5
			}
			color.Cyan(ADD_RETURN_CODES[6])
			return 6
		}
		removed, id, _ := LogEntryLookup("REM", filePath)
		if removed {
			if exists {
				removeFileAndLog(id, "removed")
				_, commitId, fileId := GetFileMetadata(filePath)
				// If file was staged (REM), it indicates that it was committed, too
				modified := IsModified(filePath, dirs.Commits+commitId+"/"+fileId+"/"+fileName)
				if modified {
					stageAndLog(generatedId, filePath, "modified")
					color.Cyan(ADD_RETURN_CODES[7])
					return 7
				}
				/*
				 If file was staged (REM) and added back to workdir without any modifications, everything should be unchanged.
				*/
			} else {
				color.Cyan(ADD_RETURN_CODES[8])
				return 8
			}
		}
	} else {
		isCommitted, commitId, fileId := GetFileMetadata(filePath)
		isDeleted := IsFileDeleted(filePath)
		if isDeleted {
			AddToStaging(generatedId, dirs.Commits+commitId+"/"+fileId+"/"+fileName, "removed")
			LogOperation(generatedId, "REM", filePath)
			color.Cyan(ADD_RETURN_CODES[9])
			return 9
		}

		if isCommitted {
			modified := IsModified(filePath, dirs.Commits+commitId+"/"+fileId+"/"+fileName)
			if modified {
				stageAndLog(generatedId, filePath, "modified")
				return 10
			} else {
				color.Red(ADD_RETURN_CODES[11])
				return 11
			}
		} else {
			stageAndLog(generatedId, filePath, "added")
			color.Cyan(ADD_RETURN_CODES[12])
			return 12
		}
	}
	return 0
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
