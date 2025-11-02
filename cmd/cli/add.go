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
		Debug("Starting add command with args: %v", args)
		for _, arg := range args {
			runAddCommand(arg)
		}
	},
}

func runAddCommand(filePath string) int {
	Debug("Processing file: %s", filePath)
	initialized := IsInitialized()
	if !initialized {
		color.Red(COMMON_RETURN_CODES[001])
		return 001
	}

	// Validate path to prevent traversal attacks
	if err := ValidatePath(filePath); err != nil {
		color.Red("Invalid file path: " + err.Error())
		return 001
	}

	_, fileName := ParsePath(filePath)
	generatedId := GenRandHex(20)
	Debug("Generated ID for file: %s", generatedId)

	fileStaged := IsFileStaged(filePath)
	if fileStaged {
		Debug("File is already staged: %s", filePath)
		exists := FileExists(filePath)
		added, id, _ := LogEntryLookup("ADD", filePath)
		if added {
			if !exists {
				Debug("File was added but no longer exists, removing from staging")
				if err := removeFileAndLog(id, "added"); err != nil {
					color.Red("Error removing file from staging: " + err.Error())
					return 001
				}
				color.Cyan(ADD_RETURN_CODES[101])
				return 101
			}
			modified, err := IsModified(filePath, dirs.StagingAdded+id+"/"+fileName)
			if err != nil {
				color.Red("Error checking if file is modified: " + err.Error())
				return 001
			}
			if modified {
				Debug("File was added and modified, updating staging")
				if err := AddToStaging(id, filePath, "added"); err != nil {
					color.Red("Error adding file to staging: " + err.Error())
					return 001
				}
				color.Cyan(ADD_RETURN_CODES[102])
				return 102
			}
			Debug("File was added but not modified")
			color.Cyan(ADD_RETURN_CODES[103])
			return 103
		}
		modified, id, _ := LogEntryLookup("MOD", filePath)
		if modified {
			if !exists {
				Debug("File was modified but no longer exists, removing from staging")
				if err := removeFileAndLog(id, "modified"); err != nil {
					color.Red("Error removing file from staging: " + err.Error())
					return 001
				}
				LogOperation(generatedId, "REM", filePath)
				return 104
			}
			modified, err := IsModified(filePath, dirs.StagingModified+id+"/"+fileName)
			if err != nil {
				color.Red("Error checking if file is modified: " + err.Error())
				return 001
			}
			if modified {
				Debug("File was modified and changed, updating staging")
				if err := AddToStaging(id, filePath, "modified"); err != nil {
					color.Red("Error adding file to staging: " + err.Error())
					return 001
				}
				color.Green(ADD_RETURN_CODES[105])
				return 105
			}
			Debug("File was modified but not changed")
			color.Cyan(ADD_RETURN_CODES[106])
			return 106
		}
		removed, id, _ := LogEntryLookup("REM", filePath)
		if removed {
			if exists {
				Debug("File was removed but exists again, checking modifications")
				if err := removeFileAndLog(id, "removed"); err != nil {
					color.Red("Error removing file from staging: " + err.Error())
					return 001
				}
				_, commitId, fileId := GetFileMetadata(filePath)
				modified, err := IsModified(filePath, dirs.Commits+commitId+"/"+fileId+"/"+fileName)
				if err != nil {
					color.Red("Error checking if file is modified: " + err.Error())
					return 001
				}
				if modified {
					Debug("File was removed but modified, adding back as modified")
					if err := stageAndLog(generatedId, filePath, "modified"); err != nil {
						color.Red("Error staging file: " + err.Error())
						return 001
					}
					color.Cyan(ADD_RETURN_CODES[107])
					return 107
				}
			} else {
				Debug("File was removed and still doesn't exist")
				color.Cyan(ADD_RETURN_CODES[108])
				return 108
			}
		}
	} else {
		Debug("File is not staged, checking commit status")
		isCommitted, commitId, fileId := GetFileMetadata(filePath)
		isDeleted := IsFileDeleted(filePath)
		if isDeleted {
			Debug("File was committed but deleted, staging for removal")
			if err := AddToStaging(generatedId, dirs.Commits+commitId+"/"+fileId+"/"+fileName, "removed"); err != nil {
				color.Red("Error adding file to staging: " + err.Error())
				return 001
			}
			LogOperation(generatedId, "REM", filePath)
			color.Green(ADD_RETURN_CODES[109])
			return 109
		}

		if isCommitted {
			modified, err := IsModified(filePath, dirs.Commits+commitId+"/"+fileId+"/"+fileName)
			if err != nil {
				color.Red("Error checking if file is modified: " + err.Error())
				return 001
			}
			if modified {
				Debug("File was committed and modified, staging as modified")
				if err := stageAndLog(generatedId, filePath, "modified"); err != nil {
					color.Red("Error staging file: " + err.Error())
					return 001
				}
				color.Green(ADD_RETURN_CODES[110])
				return 110
			} else {
				Debug("File was committed but not modified")
				color.Red(ADD_RETURN_CODES[111])
				return 111
			}
		} else {
			Debug("File is new, staging as added")
			if err := stageAndLog(generatedId, filePath, "added"); err != nil {
				color.Red("Error staging file: " + err.Error())
				return 001
			}
			color.Green(ADD_RETURN_CODES[112])
			return 112
		}
	}
	return 100
}

func removeFileAndLog(id string, op string) error {
	Debug("Removing file and log entry: id=%s, op=%s", id, op)
	if err := RemoveFile(dirs.Staging + op + "/" + id); err != nil {
		return err
	}
	RemoveLogEntry(id)
	return nil
}

func stageAndLog(id string, path string, op string) error {
	Debug("Staging and logging file: id=%s, path=%s, op=%s", id, path, op)
	logOperations := map[string]string{
		"added":    "ADD",
		"modified": "MOD",
		"removed":  "REM",
	}
	if err := AddToStaging(id, path, op); err != nil {
		return err
	}
	LogOperation(id, logOperations[op], path)
	return nil
}
