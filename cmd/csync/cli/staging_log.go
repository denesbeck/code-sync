package cli

import (
	"encoding/json"
	"log"
	"os"
)

type LogFileEntry struct {
	Id   string
	Op   string
	Path string
}

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
	err = WriteJson(".csync/staging/logs.json", payload)
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
	err = WriteJson(".csync/staging/logs.json", payload)
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
