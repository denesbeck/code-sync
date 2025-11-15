package main

import (
	"github.com/spf13/cobra"
)

func init() {
	addCmd.Flags().BoolVarP(&Force, "force", "f", false, "Disregard the rules defined in `.csync.rules.yml`")

	rootCmd.AddCommand(addCmd)
}

var Force bool

var addCmd = &cobra.Command{
	Use:     "add",
	Short:   "Add the selected files to the staging area",
	Example: "csync add <path/to/your/file>",
	Args:    cobra.MinimumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		Debug("Starting add command with args: %v", args)
		for _, arg := range args {
			runAddCommand(arg, Force)
		}
	},
}

func runAddCommand(filePath string, force bool) int {
	Debug("Processing file: %s", filePath)
	initialized := IsInitialized()
	if !initialized {
		Fail(COMMON_RETURN_CODES[001])
		return 001
	}

	if err := ValidatePath(filePath); err != nil {
		Debug("Path is invalid: %s", err.Error())
		Fail(COMMON_RETURN_CODES[004])
		return 004
	}

	if !force {
		shouldIgnore := ShouldIgnore(filePath)
		if shouldIgnore {
			Warning(COMMON_RETURN_CODES[002])
			return 002
		}
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
					Debug("Error removing file from staging: %s", err.Error())
					MustSucceed(err, "operation failed")
				}
				Info(ADD_RETURN_CODES[101])
				return 101
			}
			modified, err := IsModified(filePath, dirs.StagingAdded+id+"/"+fileName)
			if err != nil {
				Debug("Error checking if file is modified: %s", err.Error())
				MustSucceed(err, "operation failed")
			}
			if modified {
				Debug("File was added and modified, updating staging")
				if err := AddToStaging(id, filePath, "added"); err != nil {
					Debug("Error adding file to staging: %s", err.Error())
					MustSucceed(err, "operation failed")
				}
				Success(ADD_RETURN_CODES[102])
				return 102
			}
			Debug("File was added but not modified")
			Info(ADD_RETURN_CODES[103])
			return 103
		}
		modified, id, _ := LogEntryLookup("MOD", filePath)
		if modified {
			if !exists {
				Debug("File was modified but no longer exists, removing from staging")
				if err := removeFileAndLog(id, "modified"); err != nil {
					Debug("Error removing file from staging: %s", err.Error())
					MustSucceed(err, "operation failed")
				}
				LogOperation(generatedId, "REM", filePath)
				Info(ADD_RETURN_CODES[104])
				return 104
			}
			modified, err := IsModified(filePath, dirs.StagingModified+id+"/"+fileName)
			if err != nil {
				Debug("Error checking if file is modified: %s", err.Error())
				MustSucceed(err, "operation failed")
			}
			if modified {
				Debug("File was modified and changed, updating staging")
				if err := AddToStaging(id, filePath, "modified"); err != nil {
					Debug("Error adding file to staging: %s", err.Error())
					MustSucceed(err, "operation failed")
				}
				Success(ADD_RETURN_CODES[105])
				return 105
			}
			Debug("File was modified but not changed")
			Info(ADD_RETURN_CODES[106])
			return 106
		}
		removed, id, _ := LogEntryLookup("REM", filePath)
		if removed {
			if exists {
				Debug("File was removed but exists again, checking modifications")
				if err := removeFileAndLog(id, "removed"); err != nil {
					Debug("Error removing file from staging: %s", err.Error())
					MustSucceed(err, "operation failed")
				}
				_, commitId, fileId := GetFileMetadata(filePath)
				modified, err := IsModified(filePath, dirs.Commits+commitId+"/"+fileId+"/"+fileName)
				if err != nil {
					Debug("Error checking if file is modified: %s", err.Error())
					MustSucceed(err, "operation failed")
				}
				if modified {
					Debug("File was removed but modified, adding back as modified")
					if err := stageAndLog(generatedId, filePath, "modified"); err != nil {
						Debug("Error staging file: %s", err.Error())
						MustSucceed(err, "operation failed")
					}
					Success(ADD_RETURN_CODES[107])
					return 107
				}
			} else {
				Debug("File was removed and still doesn't exist")
				Info(ADD_RETURN_CODES[108])
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
				Debug("Error adding file to staging: %s", err.Error())
				MustSucceed(err, "operation failed")
			}
			LogOperation(generatedId, "REM", filePath)
			Success(ADD_RETURN_CODES[109])
			return 109
		}

		if isCommitted {
			modified, err := IsModified(filePath, dirs.Commits+commitId+"/"+fileId+"/"+fileName)
			if err != nil {
				Debug("Error checking if file is modified: %s", err.Error())
				MustSucceed(err, "operation failed")
			}
			if modified {
				Debug("File was committed and modified, staging as modified")
				if err := stageAndLog(generatedId, filePath, "modified"); err != nil {
					Debug("Error staging file: %s", err.Error())
					MustSucceed(err, "operation failed")
				}
				Success(ADD_RETURN_CODES[110])
				return 110
			} else {
				Debug("File was committed but not modified")
				Info(ADD_RETURN_CODES[111])
				return 111
			}
		} else {
			Debug("File is new, staging as added")
			if err := stageAndLog(generatedId, filePath, "added"); err != nil {
				Debug("Error staging file: %s", err.Error())
				MustSucceed(err, "operation failed")
			}
			Success(ADD_RETURN_CODES[112])
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
