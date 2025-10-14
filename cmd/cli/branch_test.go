package cli

import (
	"os"
	"testing"
)

func Test_NewBranchCmd(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	runNewCommand("test-branch", "", "")

	branches, err := os.ReadDir(dirs.Branches)
	if err != nil {
		t.Error(err)
	}

	found := false

	for _, branch := range branches {
		if branch.IsDir() {
			if branch.Name() == "test-branch" {
				found = true
				break
			}
		}
	}

	if !found {
		t.Error("Branch `test-branch` not found")
	}

	os.RemoveAll(namespace)
}

func Test_NewBranchFromCommit(t *testing.T) {
	// TODO: Implement
}

func Test_NewBranchFromBranch(t *testing.T) {
	// TODO: Implement
}

func Test_DropDefaultBranch(t *testing.T) {
	// TODO: Implement
}

func Test_DropCurrentBranch(t *testing.T) {
	// TODO: Implement
}

func Test_NewBranchAlreadyExists(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	statusCode := runNewCommand("test-branch", "", "")
	if statusCode != 206 {
		t.Errorf("Expected 206, got %d", statusCode)
	}

	statusCode = runNewCommand("test-branch", "", "")
	if statusCode != 205 {
		t.Errorf("Expected 205, got %d", statusCode)
	}

	os.RemoveAll(namespace)
}

func Test_NewBranchInvalidName(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	statusCode := runNewCommand("test-branch", "", "")
	if statusCode != 206 {
		t.Errorf("Expected 206, got %d", statusCode)
	}

	statusCode = runNewCommand("test branch %#^@#&", "", "")
	if statusCode != 201 {
		t.Errorf("Expected 201, got %d", statusCode)
	}

	os.RemoveAll(namespace)
}

// Test: new, switch, drop, current, default
func Test_Branching(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	runNewCommand("test-branch1", "", "")
	runNewCommand("test-branch2", "", "")
	runNewCommand("test-branch3", "", "")

	runSwitchCommand("test-branch2")

	currentBranchName := GetCurrentBranchName()

	if currentBranchName != "test-branch2" {
		t.Errorf("Expected current branch to be `test-branch2`, got '%s'", currentBranchName)
	}

	runSwitchCommand("test-branch1")

	currentBranchName = GetCurrentBranchName()

	if currentBranchName != "test-branch1" {
		t.Errorf("Expected current branch to be `test-branch1`, got '%s'", currentBranchName)
	}

	runSwitchCommand("test-branch3")

	currentBranchName = GetCurrentBranchName()

	if currentBranchName != "test-branch3" {
		t.Errorf("Expected current branch to be `test-branch3`, got '%s'", currentBranchName)
	}

	defaultBranchName := GetDefaultBranchName()
	if defaultBranchName != "main" {
		t.Errorf("Expected default branch to be `main`, got '%s'", defaultBranchName)
	}

	SetBranch("test-branch2", "default")

	defaultBranchName = GetDefaultBranchName()
	if defaultBranchName != "test-branch2" {
		t.Errorf("Expected default branch to be `test-branch2`, got '%s'", defaultBranchName)
	}

	runDropCommand("test-branch1")

	branches, err := os.ReadDir(dirs.Branches)
	if err != nil {
		t.Error(err)
	}

	found := false

	for _, branch := range branches {
		if branch.IsDir() {
			if branch.Name() == "test-branch1" {
				found = true
				break
			}
		}
	}

	if found {
		t.Error("branch `test-branch1` not deleted")
	}

	os.RemoveAll(namespace)
}
