package cli

import (
	"encoding/json"
	"log"
	"os"
)

type FileListEntry struct {
	Id       string
	CommitId string
	Path     string
}

func IsFileStaged(filePath string) bool {
	Debug("Checking if file is staged: %s", filePath)
	logs, err := os.ReadFile(dirs.StagingLogs)
	if err != nil {
		Debug("Failed to read staging logs")
		log.Fatal(err)
	}

	if len(logs) == 0 {
		Debug("No staging logs found")
		return false
	}
	var content []LogFileEntry
	if err = json.Unmarshal(logs, &content); err != nil {
		Debug("Failed to unmarshal staging logs")
		log.Fatal(err)
	}
	for _, entry := range content {
		if entry.Path == filePath {
			Debug("File is staged with operation: %s", entry.Op)
			return true
		}
	}
	Debug("File is not staged")
	return false
}

func GetFileMetadata(filePath string) (isCommitted bool, commitId string, fileId string) {
	Debug("Getting file metadata: %s", filePath)
	latestCommitId := GetLastCommit().Id
	if latestCommitId == "" {
		Debug("No commits found")
		return false, "", ""
	}
	fileList, err := os.ReadFile(dirs.Commits + latestCommitId + "/fileList.json")
	if err != nil {
		Debug("Failed to read file list")
		log.Fatal(err)
	}

	var content []FileListEntry
	if err = json.Unmarshal(fileList, &content); err != nil {
		Debug("Failed to unmarshal file list")
		log.Fatal(err)
	}
	for _, file := range content {
		if file.Path == filePath {
			Debug("File found in commit: id=%s, commitId=%s", file.Id, file.CommitId)
			return true, file.CommitId, file.Id
		}
	}
	Debug("File not found in any commit")
	return false, "", ""
}

func IsFileDeleted(filePath string) bool {
	Debug("Checking if file is deleted: %s", filePath)
	committed, _, _ := GetFileMetadata(filePath)
	existsInWorkdir := FileExists(filePath)
	isDeleted := committed && !existsInWorkdir
	Debug("File deletion status: committed=%v, exists=%v, isDeleted=%v", committed, existsInWorkdir, isDeleted)
	return isDeleted
}
