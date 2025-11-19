package main

import (
	"github.com/spf13/cobra"
)

func init() {
	addCmd.Flags().BoolVarP(&Force, "force", "f", false, "Disregard the rules defined in `.csync.rules.yml`")

	rootCmd.AddCommand(addCmd)
}

var Force bool

type AddResult struct {
	FilePath   string
	ReturnCode int
	Message    string
	Success    bool
}

var addCmd = &cobra.Command{
	Use:     "add",
	Short:   "Add the selected files to the staging area",
	Example: "csync add <path/to/your/file>\ncsync add file1 file2 file3\ncsync add .",
	Args:    cobra.MinimumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		Debug("Starting add command with args: %v", args)

		initialized := IsInitialized()
		if !initialized {
			Fail(COMMON_RETURN_CODES[001])
			return
		}

		filePaths, err := ExpandFilePaths(args)
		if err != nil {
			Fail("Failed to expand file paths: " + err.Error())
			return
		}

		Debug("Processing %d files", len(filePaths))

		results := make([]AddResult, 0, len(filePaths))
		for _, filePath := range filePaths {
			result := runAddCommand(filePath, Force)
			results = append(results, result)
		}
		DisplayAddResults(results)
	},
}

func runAddCommand(filePath string, force bool) AddResult {
	result := AddResult{FilePath: filePath}
	returnCode := runAddCommandInternal(filePath, force, &result)
	result.ReturnCode = returnCode
	return result
}

func runAddCommandInternal(filePath string, force bool, _ *AddResult) int {
	Debug("Processing file: %s", filePath)

	if err := ValidatePath(filePath); err != nil {
		Debug("Path is invalid: %s", err.Error())
		return 004
	}

	if !force {
		shouldIgnore := ShouldIgnore(filePath)
		if shouldIgnore {
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
				if err := RemoveFileAndLog(id, "added"); err != nil {
					Debug("Error removing file from staging: %s", err.Error())
					MustSucceed(err, "operation failed")
				}
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
				return 102
			}
			Debug("File was added but not modified")
			return 103
		}
		modified, id, _ := LogEntryLookup("MOD", filePath)
		if modified {
			if !exists {
				Debug("File was modified but no longer exists, removing from staging")
				if err := RemoveFileAndLog(id, "modified"); err != nil {
					Debug("Error removing file from staging: %s", err.Error())
					MustSucceed(err, "operation failed")
				}
				LogOperation(generatedId, "REM", filePath)
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
				return 105
			}
			Debug("File was modified but not changed")
			return 106
		}
		removed, id, _ := LogEntryLookup("REM", filePath)
		if removed {
			if exists {
				Debug("File was removed but exists again, checking modifications")
				if err := RemoveFileAndLog(id, "removed"); err != nil {
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
					if err := StageAndLog(generatedId, filePath, "modified"); err != nil {
						Debug("Error staging file: %s", err.Error())
						MustSucceed(err, "operation failed")
					}
					return 107
				}
				Debug("File was removed but exists again without modifications, removed from staging")
				return 113
			} else {
				Debug("File was removed and still doesn't exist")
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
				if err := StageAndLog(generatedId, filePath, "modified"); err != nil {
					Debug("Error staging file: %s", err.Error())
					MustSucceed(err, "operation failed")
				}
				return 110
			} else {
				Debug("File was committed but not modified")
				return 111
			}
		} else {
			Debug("File is new, staging as added")
			if err := StageAndLog(generatedId, filePath, "added"); err != nil {
				Debug("Error staging file: %s", err.Error())
				MustSucceed(err, "operation failed")
			}
			return 112
		}
	}
	return 100 // Fallback
}
