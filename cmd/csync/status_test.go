package main

import (
	"os"
	"testing"
)

func Test_IsFileStaged(t *testing.T) {
	os.RemoveAll(namespace)
	runInitCommand()
	file := namespace + "file.txt"
	os.Create(file)
	runAddCommand(file)
	if !IsFileStaged(file) {
		t.Errorf("Expected file %s to be staged", file)
	}
	if IsFileStaged(namespace + "nonexistent.txt") {
		t.Errorf("Expected file %s to not be staged", namespace+"nonexistent.txt")
	}
	os.RemoveAll(namespace)
}

func Test_IsFileDeleted(t *testing.T) {
	os.RemoveAll(namespace)
	runInitCommand()
	file := namespace + "file.txt"
	os.Create(file)
	runAddCommand(file)
	runCommitCommand("test commit")

	os.Remove(file)

	isCommitted, _, _ := GetFileMetadata(file)

	if !isCommitted {
		t.Errorf("Expected file %s to be committed", file)
	}

	if !IsFileDeleted(file) {
		t.Errorf("Expected file %s to be deleted", file)
	}

	if IsFileDeleted(namespace + "nonexistent.txt") {
		t.Errorf("Expected file %s to not be deleted", namespace+"nonexistent.txt")
	}
	os.RemoveAll(namespace)
}

func Test_GetFileMetadata(t *testing.T) {
	os.RemoveAll(namespace)
	runInitCommand()
	file := namespace + "file.txt"
	os.Create(file)
	runAddCommand(file)
	runCommitCommand("test commit")

	isCommitted, commitId, fileId := GetFileMetadata(file)

	if !isCommitted {
		t.Errorf("Expected file %s to be committed", file)
	}

	if commitId == "" {
		t.Errorf("Expected commit ID to not be empty")
	}

	if fileId == "" {
		t.Errorf("Expected file ID to not be empty")
	}

	isCommitted, _, _ = GetFileMetadata(namespace + "nonexistent.txt")
	if isCommitted {
		t.Errorf("Expected file %s to not be committed", namespace+"nonexistent.txt")
	}
	os.RemoveAll(namespace)
}

func Test_StatusCommand(t *testing.T) {
	os.RemoveAll(namespace)
	runInitCommand()
	file := namespace + "file.txt"
	os.Create(file)
	runAddCommand(file)

	returnCode, stagingLogs := runStatusCommand()

	if returnCode != 502 {
		t.Errorf("Expected return code 502, got %d", returnCode)
	}

	if len(stagingLogs) == 0 {
		t.Errorf("Expected staging logs to not be empty")
	}

	runRemoveCommand(file)

	returnCode, stagingLogs = runStatusCommand()

	if returnCode != 501 {
		t.Errorf("Expected return code 501, got %d", returnCode)
	}

	if len(stagingLogs) != 0 {
		t.Errorf("Expected staging logs to be empty")
	}

	os.RemoveAll(namespace)
}
