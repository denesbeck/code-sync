package csync

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"slices"
	"sort"
	"strings"

	"github.com/fatih/color"
)

type Metadata struct {
	Default string
	Current string
}
type FileListEntry struct {
	CommitId  string
	Path      string
	Timestamp string
}

// read the .csyncignore.json file and return its content
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
	var payload []string
	if err = json.Unmarshal(logs, &payload); err != nil {
		log.Fatal(err)
	}
	return slices.Contains(payload, filePath)
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

// Get the id of the last commit, if there is one
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
func AddToStaging(path string, op string) {
	dirs, file := ParsePath(path)

	fullPath := ".csync/staging/" + op + "/" + dirs

	if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	_, err := CopyFile(path, fullPath+file)
	if err != nil {
		log.Fatal(err)
	}
	color.Green("File added successfully")
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

// Common
func IsInitialized() bool {
	if _, err := os.Stat(".csync"); !os.IsNotExist(err) {
		return true
	}
	return false
}

func CopyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, errors.New("Not a regular file")
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()

	nBytes, err := io.Copy(destination, source)
	return nBytes, err
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
