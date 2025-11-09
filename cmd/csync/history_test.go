package main

import (
	"os"
	"strconv"
	"testing"
)

func Test_History(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	setConfig("username", "testuser")
	setConfig("email", "test@test.com")

	for i := 1; i <= 10; i++ {
		if i > 5 {
			setConfig("username", "testuserX")
			setConfig("email", "testX@test.com")
		}
		fileName := namespace + strconv.Itoa(i) + ".txt"
		os.Create(fileName)
		if i%2 == 0 {
			os.Create(fileName + "a")
			runAddCommand(fileName + "a", false)
		}
		runAddCommand(fileName, false)
		runCommitCommand("Commit " + strconv.Itoa(i))
	}

	statusCode, history := runHistoryCommand()
	if statusCode != 401 {
		t.Errorf("Expected 401, got %d", statusCode)
	}

	if len(history) != 10 {
		t.Errorf("Expected 10 commits, got %d", len(history))
	}

	for i, commit := range history {
		if commit.message != "Commit "+strconv.Itoa(i+1) {
			t.Errorf("Expected commit message 'Commit %d', got '%s'", i+1, commit.message)
		}

		if i < 5 && commit.author != "testuser <test@test.com>" {
			t.Errorf("Expected author 'testuser', got '%s'", commit.author)
		}
		if i > 5 && commit.author != "testuserX <testX@test.com>" {
			t.Errorf("Expected author 'testuser', got '%s'", commit.author)
		}

		if i%2 == 0 && len(commit.commits) != 1 {
			t.Errorf("Expected 1 file in commit %d, got %d", i+1, len(commit.commits))
		}
		if i%2 > 0 && len(commit.commits) != 2 {
			t.Errorf("Expected 2 file in commit %d, got %d", i+1, len(commit.commits))
		}

	}

	os.RemoveAll(namespace)
}
