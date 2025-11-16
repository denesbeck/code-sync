package main

import (
	"os"
	"path/filepath"
	"strings"
)

func AddToStaging(id string, path string, op string) error {
	Debug("Adding file to staging: id=%s, path=%s, op=%s", id, path, op)
	_, file := ParsePath(path)

	if err := os.MkdirAll(dirs.Staging+op+"/"+id, 0755); err != nil {
		Debug("Failed to create staging directory")
		MustSucceed(err, "operation failed")
	}
	if err := CopyFile(path, dirs.Staging+op+"/"+id+"/"+file); err != nil {
		return err
	}
	Debug("File added to staging successfully")
	return nil
}

func DisplayAddResults(results []AddResult) {
	if len(results) == 0 {
		return
	}

	var added, updated, removed, ignored, alreadyStaged, notModified, failed []string

	for _, r := range results {
		switch r.ReturnCode {
		case 112, 110: // File added to staging (new file or modified committed file)
			added = append(added, r.FilePath)
			r.Success = true
			r.Message = "Added"
		case 109: // File deleted from filesystem (committed file)
			removed = append(removed, r.FilePath)
			r.Success = true
			r.Message = "Removed from filesystem"
		case 102, 105, 107: // Staged file updated
			updated = append(updated, r.FilePath)
			r.Success = true
			r.Message = "Updated"
		case 101, 104: // File deleted from filesystem (was staged)
			removed = append(removed, r.FilePath)
			r.Success = true
			r.Message = "Removed from filesystem"
		case 103, 106, 108: // File already staged
			alreadyStaged = append(alreadyStaged, r.FilePath)
			r.Success = true
			r.Message = "Already staged"
		case 111: // File not modified
			notModified = append(notModified, r.FilePath)
			r.Success = true
			r.Message = "Not modified"
		case 002: // Ignored by rules
			ignored = append(ignored, r.FilePath)
			r.Success = false
			r.Message = "Ignored by rules"
		default:
			failed = append(failed, r.FilePath)
			r.Success = false
			r.Message = "Failed"
		}
	}

	if len(added) > 0 {
		BreakLine()
		Success("󰐙 Added to staging: " + FormatFileCount(len(added)))
		Tree(added, true)
	}

	if len(updated) > 0 {
		BreakLine()
		Success("󰓦 Updated in staging: " + FormatFileCount(len(updated)))
		Tree(updated, true)
	}

	if len(removed) > 0 {
		BreakLine()
		Info("󰍷 Removed from filesystem: " + FormatFileCount(len(removed)))
		Tree(removed, true)
	}

	if len(alreadyStaged) > 0 {
		BreakLine()
		Info(" Already staged: " + FormatFileCount(len(alreadyStaged)))
		Tree(alreadyStaged, true)
	}

	if len(notModified) > 0 {
		BreakLine()
		Info(" Not modified: " + FormatFileCount(len(notModified)))
		Tree(notModified, true)
	}

	if len(ignored) > 0 {
		BreakLine()
		Info(" Ignored by rules: " + FormatFileCount(len(ignored)))
		Tree(ignored, true)
	}

	if len(failed) > 0 {
		BreakLine()
		Fail(" Failed: " + FormatFileCount(len(failed)))
		Tree(failed, true)
	}
	BreakLine()
}

func RemoveFileAndLog(id string, op string) error {
	Debug("Removing file and log entry: id=%s, op=%s", id, op)
	RemoveFile(dirs.Staging + op + "/" + id)
	RemoveLogEntry(id)
	return nil
}

func StageAndLog(id string, path string, op string) error {
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

func ExpandFilePaths(args []string) ([]string, error) {
	var filePaths []string

	for _, arg := range args {
		if arg == "." {
			Debug("Expanding current directory recursively")
			err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
				if err != nil {
					Debug("Error walking path %s: %v", path, err)
					MustSucceed(err, "operation failed")
				}

				if info.IsDir() {
					return nil
				}

				if strings.HasPrefix(path, ".csync") {
					return nil
				}

				filePaths = append(filePaths, path)
				return nil
			})
			if err != nil {
				return nil, err
			}

			stagedFiles := GetStagingLogsContent()
			for _, entry := range *stagedFiles {
				if entry.Op == "ADD" || entry.Op == "MOD" {
					if !FileExists(entry.Path) {
						Debug("Found staged file that no longer exists: %s", entry.Path)
						filePaths = append(filePaths, entry.Path)
					}
				}
			}

			_, deletedFiles := GetModifiedOrDeletedFiles()
			for _, deletedFile := range deletedFiles {
				Debug("Found committed file that was deleted: %s", deletedFile)
				filePaths = append(filePaths, deletedFile)
			}
		} else {
			filePaths = append(filePaths, arg)
		}
	}

	Debug("Expanded to %d files", len(filePaths))
	return filePaths, nil
}
