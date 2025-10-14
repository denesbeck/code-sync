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
	BranchMetadataPath = "./.csync/branches/metadata.json"
	DefaultBranch      = "default"
	CurrentBranch      = "current"
	InitBranch         = "main"
)

func GetCurrentBranchName() string {
	metadata := GetBranchesMetadata()
	return metadata.Current
}

func GetDefaultBranchName() string {
	metadata := GetBranchesMetadata()
	return metadata.Default
}

func CreateBranchesMetadata() {
	branchesMetadata := BranchMetadata{
		Default: InitBranch,
		Current: InitBranch,
	}
	WriteJson(BranchMetadataPath, branchesMetadata)
}

func GetBranchesMetadata() (m *BranchMetadata) {
	branchesMetadata, err := os.ReadFile(BranchMetadataPath)
	if err != nil {
		log.Fatal(err)
	}
	var metadata BranchMetadata
	if err = json.Unmarshal(branchesMetadata, &metadata); err != nil {
		log.Fatal(err)
	}
	return &metadata
}

func SetBranch(branch string, configParam string) {
	branchesMetadata := GetBranchesMetadata()

	if (configParam == DefaultBranch && branchesMetadata.Default == branch) || (configParam == CurrentBranch && branchesMetadata.Current == branch) {
		color.Red("Branch already set as " + configParam)
		return
	}

	branches := ListBranches()
	if slices.Contains(branches, branch) {
		if configParam == DefaultBranch {
			branchesMetadata.Default = branch
		} else {
			branchesMetadata.Current = branch
		}
	} else {
		color.Red("Branch does not exist")
	}

	jsonData, err := json.Marshal(branchesMetadata)
	if err != nil {
		log.Fatal(err)
	}

	if err = os.WriteFile(BranchMetadataPath, jsonData, 0644); err != nil {
		log.Fatal(err)
	}
}

func ListBranches() []string {
	entries, err := os.ReadDir(".csync/branches")
	if err != nil {
		log.Fatal(err)
	}
	branches := []string{}
	for _, e := range entries {
		if e.IsDir() {
			branches = append(branches, e.Name())
		}
	}
	return branches
}
