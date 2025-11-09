package main

import (
	"os"
	"testing"
)

func TestRemove(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	file := namespace + "file.txt"
	os.Create(file)

	runAddCommand(file)
	isLogged, logId, operation := LogEntryLookup("*", file)
	if !isLogged && operation != "ADD" {
		t.Errorf("Expected log entry to be added, got %s", operation)
	}
	exists := FileExists(dirs.Staging + "added/" + logId + "/file.txt")
	if !exists {
		t.Errorf("Expected file to exist in staging area")
	}

	runRemoveCommand(file)
	isLogged, logId, operation = LogEntryLookup("*", file)
	if isLogged {
		t.Errorf("Expected log entry to be removed, got %s", operation)
	}
	exists = FileExists(dirs.Staging + "added/" + logId + "/file.txt")
	if exists {
		t.Errorf("Expected file to be removed from staging area")
	}

	os.RemoveAll(namespace)
}
