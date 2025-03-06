package cli

import (
	"encoding/json"
	"log"
	"os"
	"slices"

	"github.com/fatih/color"
)

type BranchMetadata struct {
	Default string
	Current string
}

const (
	DefaultBranch = "default"
	CurrentBranch = "current"
	InitBranch    = "main"
)

func GetCurrentBranchName() string {
	Debug("Getting current branch name")
	metadata := GetBranchesMetadata()
	Debug("Current branch: %s", metadata.Current)
	return metadata.Current
}

func GetDefaultBranchName() string {
	Debug("Getting default branch name")
	metadata := GetBranchesMetadata()
	Debug("Default branch: %s", metadata.Default)
	return metadata.Default
}

func CreateBranchesMetadata() {
	Debug("Creating branches metadata")
	payload := BranchMetadata{
		Default: InitBranch,
		Current: InitBranch,
	}
	WriteJson(dirs.BranchesMetadata, payload)
	Debug("Branches metadata created with initial branch: %s", InitBranch)
}

func GetBranchesMetadata() (m *BranchMetadata) {
	Debug("Reading branches metadata")
	content, err := os.ReadFile(dirs.BranchesMetadata)
	if err != nil {
		Debug("Failed to read branches metadata")
		log.Fatal(err)
	}
	var metadata BranchMetadata
	if err = json.Unmarshal(content, &metadata); err != nil {
		Debug("Failed to unmarshal branches metadata")
		log.Fatal(err)
	}
	Debug("Branches metadata retrieved successfully")
	return &metadata
}

func SetBranch(branch string, configParam string) {
	Debug("Setting branch: branch=%s, config=%s", branch, configParam)
	metadata := GetBranchesMetadata()

	if (configParam == DefaultBranch && metadata.Default == branch) || (configParam == CurrentBranch && metadata.Current == branch) {
		Debug("Branch already set as %s", configParam)
		color.Red("Branch already set as " + configParam)
		return
	}

	branches := ListBranches()
	if slices.Contains(branches, branch) {
		if configParam == DefaultBranch {
			metadata.Default = branch
			Debug("Setting default branch to: %s", branch)
		} else {
			metadata.Current = branch
			Debug("Setting current branch to: %s", branch)
		}
	} else {
		Debug("Branch does not exist: %s", branch)
		color.Red("Branch does not exist")
	}

	jsonData, err := json.Marshal(metadata)
	if err != nil {
		Debug("Failed to marshal branch metadata")
		log.Fatal(err)
	}

	if err = os.WriteFile(dirs.BranchesMetadata, jsonData, 0644); err != nil {
		Debug("Failed to write branch metadata")
		log.Fatal(err)
	}
	Debug("Branch metadata updated successfully")
}

func ListBranches() []string {
	Debug("Listing all branches")
	entries, err := os.ReadDir(dirs.Branches)
	if err != nil {
		Debug("Failed to read branches directory")
		log.Fatal(err)
	}
	branches := []string{}
	for _, e := range entries {
		if e.IsDir() {
			branches = append(branches, e.Name())
		}
	}
	Debug("Found %d branches: %v", len(branches), branches)
	return branches
}
