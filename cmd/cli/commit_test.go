package cli

import (
	"encoding/json"
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
		runAddCommand(file)
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
