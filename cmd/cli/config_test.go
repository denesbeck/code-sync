package cli

import (
	"os"
	"testing"
)

func Test_ConfigDefaultBranch(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	returnCode, defaultBranch := getDefaultBranch()

	if returnCode != 601 {
		t.Errorf("Expected return code 601, got %d", returnCode)
	}
	if defaultBranch != "main" {
		t.Errorf("Expected default branch 'main', got '%s'", defaultBranch)
	}

	runNewCommand("test-branch", "", "")
	setDefaultBranch("test-branch")

	returnCode, defaultBranch = getDefaultBranch()

	if returnCode != 601 {
		t.Errorf("Expected return code 601, got %d", returnCode)
	}
	if defaultBranch != "test-branch" {
		t.Errorf("Expected default branch 'test-branch', got '%s'", defaultBranch)
	}

	os.RemoveAll(namespace)
}

func Test_ConfigUsername(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	returnCode := setConfig("username", "testuser")
	if returnCode != 603 {
		t.Errorf("Expected return code 603, got %d", returnCode)
	}
	returnCode, config := getConfig("username")
	if returnCode != 604 {
		t.Errorf("Expected return code 604, got %d", returnCode)
	}
	if config.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", config.Username)
	}
	returnCode = setConfig("email", "email@email.com")
	if returnCode != 603 {
		t.Errorf("Expected return code 603, got %d", returnCode)
	}
	returnCode, config = getConfig("email")
	if returnCode != 604 {
		t.Errorf("Expected return code 604, got %d", returnCode)
	}
	if config.Email != "email@email.com" {
		t.Errorf("Expected email 'email@email.com', got '%s'", config.Username)
	}

	os.RemoveAll(namespace)
}
