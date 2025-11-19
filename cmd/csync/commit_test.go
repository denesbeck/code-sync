package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"testing"
)

func TestCommit(t *testing.T) {
	os.RemoveAll(namespace)
	runInitCommand()

	setConfig("username", "test user")
	setConfig("email", "test@test.com")

	for i := range 10 {
		file := namespace + "file" + strconv.Itoa(i) + ".txt"
		os.Create(file)
		runAddCommand(file, false)
		returnCode, commitId := runCommitCommand("test commit " + strconv.Itoa(i))
		if returnCode != 702 {
			t.Errorf("Expected 702, got %d", returnCode)
		}
		commits := GetCommits()
		if len(*commits) == 0 {
			t.Errorf("Expected at least one commit, got %d", len(*commits))
		}
		lastCommit := GetLastCommit()
		if lastCommit.Id != commitId {
			t.Errorf("Expected commit ID %s, got %s", commitId, lastCommit.Id)
		}
		if len(*commits) == 1 {
			if lastCommit.Next != "" {
				t.Errorf("Expected no next commit, got %s", lastCommit.Next)
			}
		}
		if len(*commits) > 1 {
			if (*commits)[i-1].Next != commitId {
				t.Errorf("Expected next commit ID %s, got %s", commitId, (*commits)[i-1].Next)
			}
		}
		metadata, err := os.ReadFile(dirs.Commits + commitId + "/metadata.json")
		if err != nil {
			t.Errorf("Failed to read metadata file: %v", err)
		}
		var content CommitMetadata
		if err = json.Unmarshal(metadata, &content); err != nil {
			t.Errorf("Failed to unmarshal metadata: %v", err)
		}
		if content.Message != "test commit "+strconv.Itoa(i) {
			t.Errorf("Expected commit message 'test commit %s', got '%s'", strconv.Itoa(i), content.Message)
		}
		if content.Author != "test user <test@test.com>" {
			t.Errorf("Expected commit author `test user <test@test.com>`, got '%s'", content.Author)
		}
	}
}

func Test_CountCommits(t *testing.T) {
	os.RemoveAll(namespace)
	runInitCommand()

	// Initially should have 0 commits
	count := CountCommits()
	if count != 0 {
		t.Errorf("Expected 0 commits initially, got %d", count)
	}

	// Make a commit
	file := namespace + "test.txt"
	os.WriteFile(file, []byte("content"), 0644)
	runAddCommand(file, false)
	runCommitCommand("First commit")

	count = CountCommits()
	if count != 1 {
		t.Errorf("Expected 1 commit, got %d", count)
	}

	// Make another commit
	os.WriteFile(file, []byte("updated"), 0644)
	runAddCommand(file, false)
	runCommitCommand("Second commit")

	count = CountCommits()
	if count != 2 {
		t.Errorf("Expected 2 commits, got %d", count)
	}

	os.RemoveAll(namespace)
}

func Test_GetCommits_Simple(t *testing.T) {
	os.RemoveAll(namespace)
	runInitCommand()

	// Make some commits
	for i := 1; i <= 3; i++ {
		file := namespace + "file" + fmt.Sprintf("%d", i) + ".txt"
		os.WriteFile(file, []byte("content"), 0644)
		runAddCommand(file, false)
		runCommitCommand(fmt.Sprintf("Commit %d", i))
	}

	commits := GetCommits()
	if len(*commits) != 3 {
		t.Errorf("Expected 3 commits, got %d", len(*commits))
	}

	os.RemoveAll(namespace)
}
