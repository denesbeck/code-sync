package main

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
		t.Errorf("Expected email 'email@email.com', got '%s'", config.Email)
	}

	os.RemoveAll(namespace)
}

func Test_SetDefaultBranch_Errors(t *testing.T) {
	os.RemoveAll(namespace)
	runInitCommand()

	// Test setting default to non-existent branch
	returnCode := setDefaultBranch("nonexistent-branch")
	if returnCode != 216 {
		t.Errorf("Expected return code 216 for non-existent branch, got %d", returnCode)
	}

	// Test setting default branch twice (should return 215)
	runNewCommand("test-branch", "", "")
	setDefaultBranch("test-branch")
	returnCode = setDefaultBranch("test-branch")
	if returnCode != 215 {
		t.Errorf("Expected return code 215 when setting same default branch, got %d", returnCode)
	}

	os.RemoveAll(namespace)
}

func Test_GetConfig_All(t *testing.T) {
	os.RemoveAll(namespace)
	runInitCommand()

	// Set both username and email
	setConfig("username", "testuser")
	setConfig("email", "test@example.com")

	// Get all config
	returnCode, config := getConfig("")
	if returnCode != 604 {
		t.Errorf("Expected return code 604, got %d", returnCode)
	}

	if config.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", config.Username)
	}

	if config.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", config.Email)
	}

	os.RemoveAll(namespace)
}

func Test_GetConfigHelper(t *testing.T) {
	os.RemoveAll(namespace)
	runInitCommand()

	setConfig("username", "testuser")
	setConfig("email", "test@example.com")

	config := GetConfig()

	if config.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", config.Username)
	}

	if config.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", config.Email)
	}

	os.RemoveAll(namespace)
}
