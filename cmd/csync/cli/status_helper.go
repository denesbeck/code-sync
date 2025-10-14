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
	logs, err := os.ReadFile(".csync/staging/logs.json")
	if err != nil {
		log.Fatal(err)
	}

	if len(logs) == 0 {
		return false
	}
	var payload []LogFileEntry
	if err = json.Unmarshal(logs, &payload); err != nil {
		log.Fatal(err)
	}
	for _, entry := range payload {
		if entry.Path == filePath {
			return true
		}
	}
	return false
}

// Check if the file is already committed, return the commit id where the file was committed the last time
func IsFileCommitted(filePath string, latestCommitId string) (isCommitted bool, srcCommitId string) {
	branchName := GetCurrentBranchName()
	fileList, err := os.ReadFile(".csync/branches/" + branchName + "/fileList.json")
	if err != nil {
		log.Fatal(err)
	}

	var payload []FileListEntry
	if err = json.Unmarshal(fileList, &payload); err != nil {
		log.Fatal(err)
	}
	for _, file := range payload {
		if file.Path == filePath {
			return true, file.CommitId
		}
	}
	return false, ""
}

func IsFileDeleted(filePath string, latestCommitId string) (isDeleted bool, srcCommitId string) {
	existsInCommits, sourceCommitId := IsFileCommitted(filePath, latestCommitId)
	existsInWorkdir := FileExists(filePath)
	return existsInCommits && !existsInWorkdir, sourceCommitId
}
