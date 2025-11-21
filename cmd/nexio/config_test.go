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

func Test_ConfigName(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	returnCode := setConfig("name", "testuser")
	if returnCode != 603 {
		t.Errorf("Expected return code 603, got %d", returnCode)
	}
	returnCode, config := getConfig("name")
	if returnCode != 604 {
		t.Errorf("Expected return code 604, got %d", returnCode)
	}
	if config.Name != "testuser" {
		t.Errorf("Expected name 'testuser', got '%s'", config.Name)
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

	// Set both name and email
	setConfig("name", "testuser")
	setConfig("email", "test@example.com")

	// Get all config
	returnCode, config := getConfig("")
	if returnCode != 604 {
		t.Errorf("Expected return code 604, got %d", returnCode)
	}

	if config.Name != "testuser" {
		t.Errorf("Expected name 'testuser', got '%s'", config.Name)
	}

	if config.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", config.Email)
	}

	os.RemoveAll(namespace)
}

func Test_GetConfigHelper(t *testing.T) {
	os.RemoveAll(namespace)
	runInitCommand()

	setConfig("name", "testuser")
	setConfig("email", "test@example.com")

	config := GetConfig()

	if config.Name != "testuser" {
		t.Errorf("Expected name 'testuser', got '%s'", config.Name)
	}

	if config.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", config.Email)
	}

	os.RemoveAll(namespace)
}
func Test_GetConfig_NameNotSet(t *testing.T) {
	os.RemoveAll(namespace)
	runInitCommand()

	// Don't set name, try to get it
	returnCode, config := getConfig("name")
	if returnCode != 605 {
		t.Errorf("Expected return code 605 for unset name, got %d", returnCode)
	}
	if config.Name != "" {
		t.Errorf("Expected empty name, got '%s'", config.Name)
	}

	os.RemoveAll(namespace)
}

func Test_GetConfig_EmailNotSet(t *testing.T) {
	os.RemoveAll(namespace)
	runInitCommand()

	// Don't set email, try to get it
	returnCode, config := getConfig("email")
	if returnCode != 606 {
		t.Errorf("Expected return code 606 for unset email, got %d", returnCode)
	}
	if config.Email != "" {
		t.Errorf("Expected empty email, got '%s'", config.Email)
	}

	os.RemoveAll(namespace)
}

func Test_GetConfig_UserNotSet(t *testing.T) {
	os.RemoveAll(namespace)
	runInitCommand()

	// Test with neither name nor email set
	returnCode, config := getConfig("user")
	if returnCode != 607 {
		t.Errorf("Expected return code 607 when neither name nor email is set, got %d", returnCode)
	}

	// Test with only name set (email missing)
	setConfig("name", "testuser")
	returnCode, config = getConfig("user")
	if returnCode != 607 {
		t.Errorf("Expected return code 607 when email is not set, got %d", returnCode)
	}

	// Test with only email set (name missing)
	os.RemoveAll(namespace)
	runInitCommand()
	setConfig("email", "test@example.com")
	returnCode, config = getConfig("user")
	if returnCode != 607 {
		t.Errorf("Expected return code 607 when name is not set, got %d", returnCode)
	}

	// Verify that when both are set, it returns 604
	setConfig("name", "testuser")
	returnCode, config = getConfig("user")
	if returnCode != 604 {
		t.Errorf("Expected return code 604 when both name and email are set, got %d", returnCode)
	}
	if config.Name != "testuser" || config.Email != "test@example.com" {
		t.Errorf("Expected name 'testuser' and email 'test@example.com', got '%s' and '%s'", config.Name, config.Email)
	}

	os.RemoveAll(namespace)
}
