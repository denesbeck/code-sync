package cli

import (
	"encoding/json"
	"log"
	"os"
	"sort"
)

type Commit struct {
	Id        string
	Timestamp string
}

type CommitMetadata struct {
	Author  string
	Message string
}

func GetLastCommit() string {
	currentBranchName := GetCurrentBranchName()
	commit := GetLastCommitByBranch(currentBranchName)
	return commit
}

func GetLastCommitByBranch(branch string) string {
	commits, err := os.ReadFile(".csync/branches/" + branch + "/commits.json")
	if err != nil {
		log.Fatal(err)
	}
	var content []Commit
	if err = json.Unmarshal(commits, &content); err != nil {
		log.Fatal(err)
	}
	if len(content) == 0 {
		return ""
	}
	sort.Slice(content, func(i, j int) bool {
		return content[i].Timestamp > content[j].Timestamp
	})
	return content[0].Id
}

func GetCommits() *[]Commit {
	currentBranchName := GetCurrentBranchName()
	commits, err := os.ReadFile(".csync/branches/" + currentBranchName + "/commits.json")
	if err != nil {
		log.Fatal(err)
	}
	var content []Commit
	if err = json.Unmarshal(commits, &content); err != nil {
		log.Fatal(err)
	}
	return &content
}

func GetFileListContent(commitId string) (result *[]FileListEntry) {
	fileList, err := os.ReadFile(".csync/commits/" + commitId + "/fileList.json")
	if err != nil {
		log.Fatal(err)
	}
	var content []FileListEntry
	if len(fileList) > 0 {
		if err = json.Unmarshal(fileList, &content); err != nil {
			log.Fatal(err)
		}
	} else {
		content = []FileListEntry{}
		return &content
	}
	return &content
}

func ProcessFileList(latestCommitId string, newCommitId string) {
	var fileList *[]FileListEntry
	emptyFileList := []FileListEntry{}
	if latestCommitId == "" {
		fileList = &emptyFileList
	} else {
		c := GetFileListContent(latestCommitId)
		if c != nil {
			fileList = c
		} else {
			fileList = &emptyFileList
		}
	}
	stagingLogs := GetStagingLogsContent()

	for _, logEntry := range *stagingLogs {
		switch logEntry.Op {
		case "REM":
			if len(*fileList) == 0 {
				continue
			}
			for i, entry := range *fileList {
				if entry.Path == logEntry.Path {
					*fileList = append((*fileList)[:i], (*fileList)[i+1:]...)
					break
				}
			}
		case "ADD":
			*fileList = append(*fileList, FileListEntry{Id: logEntry.Id, CommitId: newCommitId, Path: logEntry.Path})
			_, fileName := ParsePath(logEntry.Path)
			CopyFile(logEntry.Path, ".csync/commits/"+newCommitId+"/"+logEntry.Id+"/"+fileName)
		case "MOD":
			if len(*fileList) == 0 {
				continue
			}
			for i, entry := range *fileList {
				if logEntry.Path == entry.Path {
					(*fileList)[i].Id = logEntry.Id
					(*fileList)[i].CommitId = newCommitId
				}
			}
			_, fileName := ParsePath(logEntry.Path)
			CopyFile(logEntry.Path, ".csync/commits/"+newCommitId+"/"+logEntry.Id+"/"+fileName)
		}
	}
	WriteJson(".csync/commits/"+newCommitId+"/fileList.json", fileList)
}

func WriteCommitMetadata(commitId string, message string) {
	config := GetConfig()
	WriteJson(".csync/commits/"+commitId+"/metadata.json", CommitMetadata{Author: config.Username + " <" + config.Email + ">", Message: message})
}

func RegisterCommitForBranch(commitId string) {
	currentBranchName := GetCurrentBranchName()
	commits, err := os.ReadFile(".csync/branches/" + currentBranchName + "/commits.json")
	if err != nil {
		log.Fatal(err)
	}
	var content []Commit
	if err = json.Unmarshal(commits, &content); err != nil {
		log.Fatal(err)
	}
	content = append(content, Commit{Id: commitId, Timestamp: getTimestamp()})
	WriteJson(".csync/branches/"+currentBranchName+"/commits.json", content)
}
