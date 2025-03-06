package cli

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"slices"
	"sort"
)

type Commit struct {
	Id           string
	Timestamp    string
	PrevCommitId string
}

type CommitMetadata struct {
	Author  string
	Message string
}

func GetLastCommit() string {
	Debug("Getting last commit")
	currentBranchName := GetCurrentBranchName()
	commit := GetLastCommitByBranch(currentBranchName)
	Debug("Last commit: %s", commit)
	return commit
}

func GetLastCommitByBranch(branch string) string {
	Debug("Getting last commit for branch: %s", branch)
	commits, err := os.ReadFile(dirs.Branches + branch + "/commits.json")
	if err != nil {
		Debug("Failed to read commits file")
		log.Fatal(err)
	}
	var content []Commit
	if err = json.Unmarshal(commits, &content); err != nil {
		Debug("Failed to unmarshal commits")
		log.Fatal(err)
	}
	if len(content) == 0 {
		Debug("No commits found for branch")
		return ""
	}
	sort.Slice(content, func(i, j int) bool {
		return content[i].Timestamp > content[j].Timestamp
	})
	Debug("Last commit for branch: %s", content[0].Id)
	return content[0].Id
}

func GetCommits() *[]Commit {
	Debug("Getting all commits")
	currentBranchName := GetCurrentBranchName()
	commits, err := os.ReadFile(dirs.Branches + currentBranchName + "/commits.json")
	if err != nil {
		Debug("Failed to read commits file")
		log.Fatal(err)
	}
	var content []Commit
	if err = json.Unmarshal(commits, &content); err != nil {
		Debug("Failed to unmarshal commits")
		log.Fatal(err)
	}
	Debug("Retrieved %d commits", len(content))
	return &content
}

func GetFileListContent(commitId string) (result *[]FileListEntry) {
	Debug("Getting file list for commit: %s", commitId)
	fileList, err := os.ReadFile(dirs.Commits + commitId + "/fileList.json")
	if err != nil {
		Debug("Failed to read file list")
		log.Fatal(err)
	}
	var content []FileListEntry
	if len(fileList) > 0 {
		if err = json.Unmarshal(fileList, &content); err != nil {
			Debug("Failed to unmarshal file list")
			log.Fatal(err)
		}
	} else {
		content = []FileListEntry{}
		Debug("Empty file list")
		return &content
	}
	Debug("Retrieved %d files from commit", len(content))
	return &content
}

func ProcessFileList(latestCommitId string, newCommitId string) {
	Debug("Processing file list: latest=%s, new=%s", latestCommitId, newCommitId)
	var fileList *[]FileListEntry
	emptyFileList := []FileListEntry{}
	if latestCommitId == "" {
		Debug("No previous commit, starting with empty file list")
		fileList = &emptyFileList
	} else {
		c := GetFileListContent(latestCommitId)
		if c != nil {
			fileList = c
			Debug("Using file list from previous commit")
		} else {
			fileList = &emptyFileList
			Debug("No file list found in previous commit")
		}
	}
	stagingLogs := GetStagingLogsContent()

	for _, logEntry := range *stagingLogs {
		Debug("Processing staging log entry: op=%s, path=%s", logEntry.Op, logEntry.Path)
		switch logEntry.Op {
		case "REM":
			if len(*fileList) == 0 {
				Debug("Skipping REM operation - no files in list")
				continue
			}
			for i, entry := range *fileList {
				if entry.Path == logEntry.Path {
					Debug("Removing file from list: %s", entry.Path)
					*fileList = append((*fileList)[:i], (*fileList)[i+1:]...)
					break
				}
			}
		case "ADD":
			Debug("Adding new file to list: %s", logEntry.Path)
			*fileList = append(*fileList, FileListEntry{Id: logEntry.Id, CommitId: newCommitId, Path: logEntry.Path})
			_, fileName := ParsePath(logEntry.Path)
			CopyFile(logEntry.Path, dirs.Commits+newCommitId+"/"+logEntry.Id+"/"+fileName)
		case "MOD":
			if len(*fileList) == 0 {
				Debug("Skipping MOD operation - no files in list")
				continue
			}
			for i, entry := range *fileList {
				if logEntry.Path == entry.Path {
					Debug("Updating file in list: %s", entry.Path)
					(*fileList)[i].Id = logEntry.Id
					(*fileList)[i].CommitId = newCommitId
				}
			}
			_, fileName := ParsePath(logEntry.Path)
			CopyFile(logEntry.Path, dirs.Commits+newCommitId+"/"+logEntry.Id+"/"+fileName)
		}
	}
	WriteJson(dirs.Commits+newCommitId+"/fileList.json", fileList)
	Debug("File list processed successfully")
}

func WriteCommitMetadata(commitId string, message string) {
	Debug("Writing commit metadata: id=%s, message=%s", commitId, message)
	config := GetConfig()
	WriteJson(dirs.Commits+commitId+"/metadata.json", CommitMetadata{Author: config.Username + " <" + config.Email + ">", Message: message})
	Debug("Commit metadata written successfully")
}

func RegisterCommitForBranch(commitId string) {
	Debug("Registering commit for branch: %s", commitId)
	currentBranchName := GetCurrentBranchName()
	commits, err := os.ReadFile(dirs.Branches + currentBranchName + "/commits.json")
	if err != nil {
		Debug("Failed to read commits file")
		log.Fatal(err)
	}
	var content []Commit
	if err = json.Unmarshal(commits, &content); err != nil {
		Debug("Failed to unmarshal commits")
		log.Fatal(err)
	}
	lastCommit := GetLastCommit()
	content = append(content, Commit{Id: commitId, Timestamp: GetTimestamp(), PrevCommitId: lastCommit})
	WriteJson(dirs.Branches+currentBranchName+"/commits.json", content)
	Debug("Commit registered successfully")
}

func CopyCommitsToBranch(commitId string, targetBranch string) error {
	Debug("Copying commits to branch: commit=%s, branch=%s", commitId, targetBranch)
	commits := GetCommits()
	commitIds := []string{}
	for _, commit := range *commits {
		commitIds = append(commitIds, commit.Id)
	}
	if !slices.Contains(commitIds, commitId) {
		Debug("Commit does not exist: %s", commitId)
		return errors.New("Commit does not exist")
	}

	if err := os.Mkdir(dirs.Branches+targetBranch, 0755); err != nil {
		Debug("Branch already exists: %s", targetBranch)
		return errors.New("Branch already exists")
	}

	index := FindIndex(commitIds, commitId)
	*commits = (*commits)[:(index + 1)]
	WriteJson(dirs.Branches+targetBranch+"/commits.json", *commits)
	Debug("Commits copied successfully")
	return nil
}
