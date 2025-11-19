package main

import (
	"os"
	"testing"
)

func Test_LogOperation(t *testing.T) {
	os.RemoveAll(namespace)
	runInitCommand()

	id := GenRandHex(20)
	path := namespace + "test.txt"
	os.WriteFile(path, []byte("test"), 0644)

	// Log an ADD operation
	LogOperation(id, "ADD", path)

	// Check if it was logged using LogEntryLookup
	found, foundId, foundOp := LogEntryLookup("ADD", path)
	if !found {
		t.Errorf("Expected to find logged operation")
	}
	if foundId != id {
		t.Errorf("Expected id %s, got %s", id, foundId)
	}
	if foundOp != "ADD" {
		t.Errorf("Expected operation ADD, got %s", foundOp)
	}

	os.RemoveAll(namespace)
}

func Test_RemoveLogEntry(t *testing.T) {
	os.RemoveAll(namespace)
	runInitCommand()

	id := GenRandHex(20)
	path := namespace + "test.txt"
	os.WriteFile(path, []byte("test"), 0644)

	// Log an operation
	LogOperation(id, "ADD", path)

	// Verify it's there
	found, _, _ := LogEntryLookup("ADD", path)
	if !found {
		t.Errorf("Expected to find log entry before removal")
	}

	// Remove it
	RemoveLogEntry(id)

	// Verify it's gone
	found, _, _ = LogEntryLookup("ADD", path)
	if found {
		t.Errorf("Expected log entry to be removed")
	}

	os.RemoveAll(namespace)
}

func Test_LogEntryLookup(t *testing.T) {
	os.RemoveAll(namespace)
	runInitCommand()

	// Test when entry doesn't exist
	found, _, _ := LogEntryLookup("ADD", "nonexistent.txt")
	if found {
		t.Errorf("Expected not to find nonexistent entry")
	}

	// Add an entry and test
	id := GenRandHex(20)
	path := namespace + "test.txt"
	os.WriteFile(path, []byte("test"), 0644)
	LogOperation(id, "MOD", path)

	found, foundId, foundOp := LogEntryLookup("MOD", path)
	if !found {
		t.Errorf("Expected to find entry")
	}
	if foundId != id {
		t.Errorf("Expected id %s, got %s", id, foundId)
	}
	if foundOp != "MOD" {
		t.Errorf("Expected operation MOD, got %s", foundOp)
	}

	os.RemoveAll(namespace)
}
