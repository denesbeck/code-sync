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
		color.Red(COMMON_RETURN_CODES[001])
		return 001
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
				color.Cyan(ADD_RETURN_CODES[101])
				return 101
			}
			modified := IsModified(filePath, dirs.StagingAdded+id+"/"+fileName)
			if modified {
				AddToStaging(id, filePath, "added")
				color.Cyan(ADD_RETURN_CODES[102])
				return 102
			}
			color.Cyan(ADD_RETURN_CODES[103])
			return 103
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
				return 104
			}
			modified := IsModified(filePath, dirs.StagingModified+id+"/"+fileName)
			if modified {
				AddToStaging(id, filePath, "modified")
				color.Cyan(ADD_RETURN_CODES[105])
				return 105
			}
			color.Cyan(ADD_RETURN_CODES[106])
			return 106
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
					color.Cyan(ADD_RETURN_CODES[107])
					return 107
				}
				/*
				 If file was staged (REM) and added back to workdir without any modifications, everything should be unchanged.
				*/
			} else {
				color.Cyan(ADD_RETURN_CODES[108])
				return 108
			}
		}
	} else {
		isCommitted, commitId, fileId := GetFileMetadata(filePath)
		isDeleted := IsFileDeleted(filePath)
		if isDeleted {
			AddToStaging(generatedId, dirs.Commits+commitId+"/"+fileId+"/"+fileName, "removed")
			LogOperation(generatedId, "REM", filePath)
			color.Cyan(ADD_RETURN_CODES[109])
			return 109
		}

		if isCommitted {
			modified := IsModified(filePath, dirs.Commits+commitId+"/"+fileId+"/"+fileName)
			if modified {
				stageAndLog(generatedId, filePath, "modified")
				return 110
			} else {
				color.Red(ADD_RETURN_CODES[111])
				return 111
			}
		} else {
			stageAndLog(generatedId, filePath, "added")
			color.Cyan(ADD_RETURN_CODES[112])
			return 112
		}
	}
	return 100
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
