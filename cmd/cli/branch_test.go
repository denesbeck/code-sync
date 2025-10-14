package cli

import (
	"os"
	"strconv"
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
	os.RemoveAll(namespace)

	runInitCommand()
	runNewCommand("test-branch", "", "")

	for i := 1; i <= 100; i++ {
		// create 5 test files
		os.Create(namespace + "file" + strconv.Itoa(i) + ".txt")
		// add files to staging
		runAddCommand(namespace + "file" + strconv.Itoa(i) + ".txt")
		// commit files
		runCommitCommand("test commit " + strconv.Itoa(i))
	}

	selectedCommit := GetLastCommit().Id

	for i := 101; i <= 200; i++ {
		// create 5 test files
		os.Create(namespace + "file" + strconv.Itoa(i) + ".txt")
		// add files to staging
		runAddCommand(namespace + "file" + strconv.Itoa(i) + ".txt")
		// commit files
		runCommitCommand("test commit " + strconv.Itoa(i))
	}
	countCommitsOriginalBranch := len(*GetCommits())

	if countCommitsOriginalBranch != 200 {
		t.Errorf("Expected 200 commits, got %d", countCommitsOriginalBranch)
	}

	t.Log(GetCommits())
	t.Log(selectedCommit)
	runNewCommand("test-branch-1", selectedCommit, "")

	countCommitsNewBranch := len(*GetCommits())
	if countCommitsNewBranch != 100 {
		t.Errorf("Expected 100 commits, got %d", countCommitsNewBranch)
	}

	lastCommitNewBranch := GetLastCommit().Id
	if selectedCommit != lastCommitNewBranch {
		t.Errorf("Expected last commit to be %s, got %s", selectedCommit, lastCommitNewBranch)
	}

	os.RemoveAll(namespace)
}

func Test_NewBranchFromBranch(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()
	runNewCommand("test-branch", "", "")
	for i := 1; i < 4; i++ {
		// create test files
		os.Create(namespace + "file" + strconv.Itoa(i) + ".txt")
		// add files to staging
		runAddCommand(namespace + "file" + strconv.Itoa(i) + ".txt")
		// commit files
		runCommitCommand("test commit " + strconv.Itoa(i))
	}
	lastCommitOriginalBranch := GetLastCommit().Id

	runNewCommand("test-branch-1", "", "test-branch")
	commits := GetCommits()
	if len(*commits) != 3 {
		t.Errorf("Expected 3 commits, got %d", len(*commits))
	}
	lastCommitNewBranch := GetLastCommit().Id
	if lastCommitOriginalBranch != lastCommitNewBranch {
		t.Errorf("Expected last commit to be %s, got %s", lastCommitOriginalBranch, lastCommitNewBranch)
	}

	os.RemoveAll(namespace)
}

func Test_SwitchNonExistingBranch(t *testing.T) {
}

func Test_DropBranch(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()
	runNewCommand("test-branch", "", "")
	runNewCommand("test-branch-1", "", "")

	statusCode := runDropCommand("test-branch")

	if statusCode != 210 {
		t.Errorf("Expected 210, got %d", statusCode)
	}

	os.RemoveAll(namespace)
}

func Test_DropDefaultBranch(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()
	runNewCommand("test-branch", "", "")
	runNewCommand("test-branch-1", "", "")

	setDefaultBranch("test-branch")

	runSwitchCommand("test-branch-1")

	statusCode := runDropCommand("test-branch")

	if statusCode != 209 {
		t.Errorf("Expected 209, got %d", statusCode)
	}

	os.RemoveAll(namespace)
}

func Test_DropCurrentBranch(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()
	runNewCommand("test-branch", "", "")

	statusCode := runDropCommand("test-branch")

	if statusCode != 208 {
		t.Errorf("Expected 208, got %d", statusCode)
	}
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

func Test_SwitchBranchAlreadyOnTarget(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	runNewCommand("test-branch", "", "")

	statusCode := runSwitchCommand("test-branch")

	if statusCode != 211 {
		t.Errorf("Expected 211, got %d", statusCode)
	}

	os.RemoveAll(namespace)
}

func Test_SwitchBranchDoesNotExist(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	runNewCommand("test-branch", "", "")

	statusCode := runSwitchCommand("test-branch-1")

	if statusCode != 212 {
		t.Errorf("Expected 212, got %d", statusCode)
	}

	os.RemoveAll(namespace)
}

func Test_SwitchBranchCurrent(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	runNewCommand("test-branch", "", "")
	runNewCommand("test-branch-1", "", "")

	statusCode := runSwitchCommand("test-branch")

	if statusCode != 213 {
		t.Errorf("Expected 213, got %d", statusCode)
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
