package main

import (
	"os"
	"strings"
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

func Test_FormatLogs(t *testing.T) {
	// Test with empty log entries
	emptyLogs := []LogFileEntry{}
	result := FormatLogs(emptyLogs)
	if result != "" {
		t.Errorf("Expected empty string for empty logs, got '%s'", result)
	}

	// Test with single log entry
	singleLog := []LogFileEntry{
		{Id: "test1", Op: "ADD", Path: "file1.txt"},
	}
	result = FormatLogs(singleLog)
	if !strings.Contains(result, "└──") {
		t.Errorf("Expected tree format with └── for single entry, got '%s'", result)
	}
	if !strings.Contains(result, "file1.txt") {
		t.Errorf("Expected result to contain 'file1.txt', got '%s'", result)
	}

	// Test with multiple log entries
	multipleLogs := []LogFileEntry{
		{Id: "test1", Op: "ADD", Path: "file1.txt"},
		{Id: "test2", Op: "MOD", Path: "file2.txt"},
		{Id: "test3", Op: "REM", Path: "file3.txt"},
	}
	result = FormatLogs(multipleLogs)
	if !strings.Contains(result, "├──") {
		t.Errorf("Expected tree format with ├── for multiple entries, got '%s'", result)
	}
	if !strings.Contains(result, "└──") {
		t.Errorf("Expected tree format with └── for last entry, got '%s'", result)
	}
	if !strings.Contains(result, "file1.txt") || !strings.Contains(result, "file2.txt") || !strings.Contains(result, "file3.txt") {
		t.Errorf("Expected result to contain all file names, got '%s'", result)
	}

	// Verify sorting (ADD should come before MOD which should come before REM)
	lines := strings.Split(result, "\n")
	addIndex := -1
	modIndex := -1
	remIndex := -1
	for i, line := range lines {
		if strings.Contains(line, "ADD") {
			addIndex = i
		}
		if strings.Contains(line, "MOD") {
			modIndex = i
		}
		if strings.Contains(line, "REM") {
			remIndex = i
		}
	}
	if addIndex > modIndex {
		t.Errorf("Expected ADD to come before MOD in sorted output")
	}
	if modIndex > remIndex {
		t.Errorf("Expected MOD to come before REM in sorted output")
	}
}

func Test_CountOps(t *testing.T) {
	// Test with empty log entries
	emptyLogs := []LogFileEntry{}
	add, mod, rem := CountOps(emptyLogs)
	if add != 0 || mod != 0 || rem != 0 {
		t.Errorf("Expected all counts to be 0 for empty logs, got add=%d, mod=%d, rem=%d", add, mod, rem)
	}

	// Test with mixed operations
	mixedLogs := []LogFileEntry{
		{Id: "test1", Op: "ADD", Path: "file1.txt"},
		{Id: "test2", Op: "ADD", Path: "file2.txt"},
		{Id: "test3", Op: "MOD", Path: "file3.txt"},
		{Id: "test4", Op: "REM", Path: "file4.txt"},
		{Id: "test5", Op: "REM", Path: "file5.txt"},
		{Id: "test6", Op: "REM", Path: "file6.txt"},
	}
	add, mod, rem = CountOps(mixedLogs)
	if add != 2 {
		t.Errorf("Expected 2 ADD operations, got %d", add)
	}
	if mod != 1 {
		t.Errorf("Expected 1 MOD operation, got %d", mod)
	}
	if rem != 3 {
		t.Errorf("Expected 3 REM operations, got %d", rem)
	}

	// Test with only ADD operations
	addOnlyLogs := []LogFileEntry{
		{Id: "test1", Op: "ADD", Path: "file1.txt"},
		{Id: "test2", Op: "ADD", Path: "file2.txt"},
		{Id: "test3", Op: "ADD", Path: "file3.txt"},
	}
	add, mod, rem = CountOps(addOnlyLogs)
	if add != 3 || mod != 0 || rem != 0 {
		t.Errorf("Expected add=3, mod=0, rem=0, got add=%d, mod=%d, rem=%d", add, mod, rem)
	}
}
