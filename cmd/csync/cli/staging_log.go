package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/fatih/color"
)

type LogFileEntry struct {
	Id   string
	Op   string
	Path string
}

var (
	add = color.New(color.FgGreen).SprintFunc()
	mod = color.New(color.FgBlue).SprintFunc()
	rem = color.New(color.FgRed).SprintFunc()
)

func LogOperation(id string, op string, path string) {
	logs, err := os.ReadFile(".csync/staging/logs.json")
	if err != nil {
		log.Fatal(err)
	}
	var content []LogFileEntry
	if len(logs) > 0 {
		if err = json.Unmarshal(logs, &content); err != nil {
			log.Fatal(err)
		}
	}
	content = append(content, LogFileEntry{
		Id:   id,
		Op:   op,
		Path: path,
	})
	WriteJson(".csync/staging/logs.json", content)
}

// Look up the logs.json file for a specific operation and path. It returns a boolean value and the id of the log entry.
func LogEntryLookup(op string, path string) (isLogged bool, logId string, operation string) {
	logs, err := os.ReadFile(".csync/staging/logs.json")
	if err != nil {
		log.Fatal(err)
	}
	var content []LogFileEntry
	if len(logs) > 0 {
		if err = json.Unmarshal(logs, &content); err != nil {
			log.Fatal(err)
		}
		for _, entry := range content {
			// Consider op "*" as a wildcard.
			if op == "*" && entry.Path == path || entry.Op == op && entry.Path == path {
				return true, entry.Id, entry.Op
			}
		}
	}
	return false, "", ""
}

func IsStagingLogsEmpty() bool {
	logs, err := os.ReadFile(".csync/staging/logs.json")
	if err != nil {
		log.Fatal(err)
	}
	var content []LogFileEntry
	if len(logs) > 0 {
		if err = json.Unmarshal(logs, &content); err != nil {
			log.Fatal(err)
		}
		if len(content) == 0 {
			return true
		}
		return false
	}
	return true
}

func RemoveLogEntry(id string) {
	logs, err := os.ReadFile(".csync/staging/logs.json")
	if err != nil {
		log.Fatal(err)
	}
	var content []LogFileEntry
	if len(logs) > 0 {
		if err = json.Unmarshal(logs, &content); err != nil {
			log.Fatal(err)
		}
	}
	for i, entry := range content {
		if entry.Id == id {
			content = append(content[:i], content[i+1:]...)
			break
		}
	}
	WriteJson(".csync/staging/logs.json", content)
}

func TruncateLogs() {
	WriteJson(".csync/staging/logs.json", []LogFileEntry{})
}

func GetStagingLogsContent() (result []LogFileEntry) {
	logs, err := os.ReadFile(".csync/staging/logs.json")
	if err != nil {
		log.Fatal(err)
	}
	var content []LogFileEntry
	if len(logs) > 0 {
		if err = json.Unmarshal(logs, &content); err != nil {
			log.Fatal(err)
		}
	} else {
		content = []LogFileEntry{}
		return content
	}
	return content
}

func SortByOperationAndPath(content []LogFileEntry) (result []LogFileEntry) {
	sort.Slice(content, func(i, j int) bool {
		if content[i].Op == "ADD" && content[j].Op == "MOD" {
			return true
		}
		if content[i].Op == "ADD" && content[j].Op == "REM" {
			return true
		}
		if content[i].Op == "MOD" && content[j].Op == "REM" {
			return true
		}
		if content[i].Op == content[j].Op {
			if content[i].Path < content[j].Path {
				return true
			}
		}
		return false
	})
	return content
}

func PrintLogs(content []LogFileEntry) {
	sortedContent := SortByOperationAndPath(content)
	for _, logEntry := range sortedContent {
		switch logEntry.Op {
		case "ADD":
			fmt.Println("  " + add(logEntry.Op) + "    " + logEntry.Path)
		case "MOD":
			fmt.Println("  " + mod(logEntry.Op) + "    " + logEntry.Path)
		case "REM":
			fmt.Println("  " + rem(logEntry.Op) + "    " + logEntry.Path)
		default:
			fmt.Println("  " + logEntry.Op + "    " + logEntry.Path)
		}
	}
}
