package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
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
	Debug("Logging operation: id=%s, op=%s, path=%s", id, op, path)
	logs, err := os.ReadFile(dirs.StagingLogs)
	if err != nil {
		Debug("Failed to read staging logs")
		log.Fatal(err)
	}
	var content []LogFileEntry
	if len(logs) > 0 {
		if err = json.Unmarshal(logs, &content); err != nil {
			Debug("Failed to unmarshal staging logs")
			log.Fatal(err)
		}
	}
	content = append(content, LogFileEntry{
		Id:   id,
		Op:   op,
		Path: path,
	})
	WriteJson(dirs.StagingLogs, content)
	Debug("Operation logged successfully")
}

func LogEntryLookup(op string, path string) (isLogged bool, logId string, operation string) {
	Debug("Looking up log entry: op=%s, path=%s", op, path)
	logs, err := os.ReadFile(dirs.StagingLogs)
	if err != nil {
		Debug("Failed to read staging logs")
		log.Fatal(err)
	}
	var content []LogFileEntry
	if len(logs) > 0 {
		if err = json.Unmarshal(logs, &content); err != nil {
			Debug("Failed to unmarshal staging logs")
			log.Fatal(err)
		}
		for _, entry := range content {
			// Consider op "*" as a wildcard.
			if op == "*" && entry.Path == path || entry.Op == op && entry.Path == path {
				Debug("Found log entry: id=%s, op=%s", entry.Id, entry.Op)
				return true, entry.Id, entry.Op
			}
		}
	}
	Debug("No matching log entry found")
	return false, "", ""
}

func IsStagingLogsEmpty() bool {
	Debug("Checking if staging logs are empty")
	logs, err := os.ReadFile(dirs.StagingLogs)
	if err != nil {
		Debug("Failed to read staging logs")
		log.Fatal(err)
	}
	var content []LogFileEntry
	if len(logs) > 0 {
		if err = json.Unmarshal(logs, &content); err != nil {
			Debug("Failed to unmarshal staging logs")
			log.Fatal(err)
		}
		if len(content) == 0 {
			Debug("Staging logs are empty")
			return true
		}
		Debug("Staging logs are not empty")
		return false
	}
	Debug("Staging logs file is empty")
	return true
}

func RemoveLogEntry(id string) {
	Debug("Removing log entry: id=%s", id)
	logs, err := os.ReadFile(dirs.StagingLogs)
	if err != nil {
		Debug("Failed to read staging logs")
		log.Fatal(err)
	}
	var content []LogFileEntry
	if len(logs) > 0 {
		if err = json.Unmarshal(logs, &content); err != nil {
			Debug("Failed to unmarshal staging logs")
			log.Fatal(err)
		}
	}
	for i, entry := range content {
		if entry.Id == id {
			Debug("Found and removing log entry: id=%s, op=%s", entry.Id, entry.Op)
			content = slices.Delete(content, i, i+1)
			break
		}
	}
	WriteJson(dirs.StagingLogs, content)
	Debug("Log entry removed successfully")
}

func TruncateLogs() {
	Debug("Truncating staging logs")
	WriteJson(dirs.StagingLogs, []LogFileEntry{})
	Debug("Staging logs truncated successfully")
}

func GetStagingLogsContent() (result *[]LogFileEntry) {
	Debug("Getting staging logs content")
	logs, err := os.ReadFile(dirs.StagingLogs)
	if err != nil {
		Debug("Failed to read staging logs")
		log.Fatal(err)
	}
	var content []LogFileEntry
	if len(logs) > 0 {
		if err = json.Unmarshal(logs, &content); err != nil {
			Debug("Failed to unmarshal staging logs")
			log.Fatal(err)
		}
	} else {
		content = []LogFileEntry{}
		Debug("Staging logs are empty")
		return &content
	}
	Debug("Retrieved %d log entries", len(content))
	return &content
}

func SortByOperationAndPath(content []LogFileEntry) (result *[]LogFileEntry) {
	Debug("Sorting log entries by operation and path")
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
	Debug("Log entries sorted successfully")
	return &content
}

func PrintLogs(content []LogFileEntry) {
	Debug("Printing %d log entries", len(content))
	sortedContent := SortByOperationAndPath(content)
	for _, logEntry := range *sortedContent {
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
	Debug("Log entries printed successfully")
}

// ValidateStagingIntegrity checks if staging logs match actual staged files
// Returns list of orphaned file IDs that should be cleaned up
func ValidateStagingIntegrity() []string {
	Debug("Validating staging integrity")
	logs := GetStagingLogsContent()
	orphanedIds := []string{}

	// Check if all logged files exist in staging
	for _, entry := range *logs {
		var stagingPath string
		switch entry.Op {
		case "ADD":
			stagingPath = dirs.StagingAdded + entry.Id
		case "MOD":
			stagingPath = dirs.StagingModified + entry.Id
		case "REM":
			stagingPath = dirs.StagingRemoved + entry.Id
		}

		if _, err := os.Stat(stagingPath); os.IsNotExist(err) {
			Debug("Found orphaned log entry: %s (path: %s)", entry.Id, entry.Path)
			orphanedIds = append(orphanedIds, entry.Id)
		}
	}

	Debug("Found %d orphaned entries", len(orphanedIds))
	return orphanedIds
}

// CleanOrphanedStagingEntries removes log entries that don't have corresponding staged files
func CleanOrphanedStagingEntries() int {
	Debug("Cleaning orphaned staging entries")
	orphanedIds := ValidateStagingIntegrity()

	for _, id := range orphanedIds {
		RemoveLogEntry(id)
		Debug("Removed orphaned log entry: %s", id)
	}

	Debug("Cleaned %d orphaned entries", len(orphanedIds))
	return len(orphanedIds)
}
