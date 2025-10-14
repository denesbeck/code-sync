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
	var content []LogFileEntry
	if err = json.Unmarshal(logs, &content); err != nil {
		log.Fatal(err)
	}
	for _, entry := range content {
		if entry.Path == filePath {
			return true
		}
	}
	return false
}

func GetFileMetadata(filePath string) (isCommitted bool, commitId string, fileId string) {
	latestCommitId := GetLastCommit()
	if latestCommitId == "" {
		return false, "", ""
	}
	fileList, err := os.ReadFile(".csync/commits/" + latestCommitId + "/fileList.json")
	if err != nil {
		log.Fatal(err)
	}

	var content []FileListEntry
	if err = json.Unmarshal(fileList, &content); err != nil {
		log.Fatal(err)
	}
	for _, file := range content {
		if file.Path == filePath {
			return true, file.CommitId, file.Id
		}
	}
	return false, "", ""
}

func IsFileDeleted(filePath string) bool {
	committed, _, _ := GetFileMetadata(filePath)
	existsInWorkdir := FileExists(filePath)
	return committed && !existsInWorkdir
}
