package csync

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/fatih/color"
)

type Metadata struct {
	Default string
	Current string
}
type LogFileEntry struct {
	Id   string
	Op   string
	Path string
}
type FileListEntry struct {
	Id       string
	CommitId string
	Path     string
}

// ### STAGING LOGS ###
// Log changes to the staging/logs.json file
func LogOperation(id string, op string, path string) {
	logs, err := os.ReadFile(".csync/staging/logs.json")
	if err != nil {
		log.Fatal(err)
	}
	var payload []LogFileEntry
	if len(logs) > 0 {
		if err = json.Unmarshal(logs, &payload); err != nil {
			log.Fatal(err)
		}
	}
	payload = append(payload, LogFileEntry{
		Id:   id,
		Op:   op,
		Path: path,
	})
	err = writeJson(".csync/staging/logs.json", payload)
	if err != nil {
		log.Fatal(err)
	}
}

// Look up the logs.json file for a specific operation and path. It returns a boolean value and the id of the log entry.
func LogEntryLookup(op string, path string) (bool, string) {
	logs, err := os.ReadFile(".csync/staging/logs.json")
	if err != nil {
		log.Fatal(err)
	}
	var payload []LogFileEntry
	if len(logs) > 0 {
		if err = json.Unmarshal(logs, &payload); err != nil {
			log.Fatal(err)
		}
		for _, entry := range payload {
			if entry.Op == op && entry.Path == path {
				return true, entry.Id
			}
		}
	}
	return false, ""
}

func RemoveLogEntry(id string) bool {
	logs, err := os.ReadFile(".csync/staging/logs.json")
	if err != nil {
		log.Fatal(err)
		return false
	}
	var payload []LogFileEntry
	if len(logs) > 0 {
		if err = json.Unmarshal(logs, &payload); err != nil {
			log.Fatal(err)
			return false
		}
	}
	for i, entry := range payload {
		if entry.Id == id {
			payload = append(payload[:i], payload[i+1:]...)
			break
		}
	}
	err = writeJson(".csync/staging/logs.json", payload)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func TruncateLogs() {
	err := os.WriteFile(".csync/staging/logs.json", []byte{}, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// ### STAGING LOGS ENDS ###

/*
Check if the file is listed in the commits/<commit_id>/fileList.json
and is missing from the working directory.
This would mean that the file should be deleted.
*/
func IsFileDeleted(filePath string, latestCommitId string) (isDeleted bool, srcCommitId string) {
	existsInCommits, sourceCommitId := IsFileCommitted(filePath, latestCommitId)
	existsInWorkdir := FileExists(filePath)
	return existsInCommits && !existsInWorkdir, sourceCommitId
}

// Read the .csyncignore.json file and return its content
func readCsyncIgnore() []string {
	_, err := os.Stat(".csyncignore.json")
	if os.IsNotExist(err) {
		color.Cyan("INFO: .csyncignore.json not found")
		return []string{}
	}
	if err != nil {
		log.Fatal(err)
		return nil
	}
	var payload []string
	content, err := os.ReadFile(".csyncignore.json")
	if err != nil {
		log.Fatal("Error while opening file: ", err)
	}
	if err = json.Unmarshal(content, &payload); err != nil {
		log.Fatal("Error while parsing data: ", err)
	}
	color.Cyan("INFO: .csyncignore.json found")
	return payload
}

// Add
// Check if the file is already staged
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
	fileList, err := os.ReadFile(".csync/commits/" + latestCommitId + "/fileList.json")
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

// Get the id of the last commit, if there is one. Otherwise return false.
func GetLastCommit() (string, bool) {
	dirs, err := os.ReadDir(".csync/commits")
	if err != nil {
		log.Fatal(err)
	}
	if len(dirs) > 0 {
		strArr := []string{}
		for _, dir := range dirs {
			if dir.IsDir() {
				strArr = append(strArr, dir.Name())
			}
		}
		sort.Sort(sort.Reverse(sort.StringSlice(strArr)))
		return strArr[0], true
	}
	return "", false
}

// Copies the file to the staging area respecting the operation
func AddToStaging(id string, path string, op string) {
	_, file := ParsePath(path)

	if err := os.MkdirAll(".csync/staging/"+op+"/"+id, 0755); err != nil {
		log.Fatal(err)
	}
	_, err := CopyFile(path, ".csync/staging/"+op+"/"+id+"/"+file)
	if err != nil {
		log.Fatal(err)
	}
	color.Green("File added successfully")
}

// ### Csync Misc ###
func IsInitialized() bool {
	if _, err := os.Stat(".csync"); !os.IsNotExist(err) {
		return true
	}
	return false
}

// Init
func CreateBranchesMetadata() error {
	branchesMetadata := Metadata{
		Default: "main",
		Current: "main",
	}

	err := writeJson(".csync/branches/metadata.json", branchesMetadata)
	if err != nil {
		return err
	}
	return nil
}

// ### Csync Misc ENDS ###

// ### FILE OPERATIONS ###
func CopyFile(src, dst string) (int64, error) {
	// File exists?
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	// Is it a regular file?
	if !sourceFileStat.Mode().IsRegular() {
		return 0, errors.New("Not a regular file")
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	// Check if the file already exists. If yes, remove it.
	if FileExists(dst) {
		os.Remove(dst)
	}

	// Create the file
	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()

	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func RemoveFile(path string) {
	err := os.Remove(path)
	if err != nil {
		log.Fatal(err)
	}
}

func FileExists(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}

func ParsePath(path string) (string, string) {
	tmpArr := strings.Split(path, "/")

	dirs := strings.Join(tmpArr[:len(tmpArr)-1], "/")
	file := tmpArr[len(tmpArr)-1]

	if dirs != "" {
		dirs = dirs + "/"
	}

	return dirs, file
}

func IsModified(file1, file2 string) bool {
	content1, err := os.ReadFile(file1)
	if err != nil {
		log.Fatal(err)
	}
	content2, err := os.ReadFile(file2)
	if err != nil {
		log.Fatal(err)
	}
	return string(content1) != string(content2)
}

// ### FILE OPERATIONS ENDS ###

func writeJson(path string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}

func GenRandHex(length int) string {
	Rando := rand.Reader
	b := make([]byte, length)
	_, _ = Rando.Read(b)
	return hex.EncodeToString(b)
}
