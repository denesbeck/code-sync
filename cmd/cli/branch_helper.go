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
	metadata := GetBranchesMetadata()
	return metadata.Current
}

func GetDefaultBranchName() string {
	metadata := GetBranchesMetadata()
	return metadata.Default
}

func CreateBranchesMetadata() {
	payload := BranchMetadata{
		Default: InitBranch,
		Current: InitBranch,
	}
	WriteJson(dirs.BranchesMetadata, payload)
}

func GetBranchesMetadata() (m *BranchMetadata) {
	content, err := os.ReadFile(dirs.BranchesMetadata)
	if err != nil {
		log.Fatal(err)
	}
	var metadata BranchMetadata
	if err = json.Unmarshal(content, &metadata); err != nil {
		log.Fatal(err)
	}
	return &metadata
}

func SetBranch(branch string, configParam string) {
	metadata := GetBranchesMetadata()

	if (configParam == DefaultBranch && metadata.Default == branch) || (configParam == CurrentBranch && metadata.Current == branch) {
		color.Red("Branch already set as " + configParam)
		return
	}

	branches := ListBranches()
	if slices.Contains(branches, branch) {
		if configParam == DefaultBranch {
			metadata.Default = branch
		} else {
			metadata.Current = branch
		}
	} else {
		color.Red("Branch does not exist")
	}

	jsonData, err := json.Marshal(metadata)
	if err != nil {
		log.Fatal(err)
	}

	if err = os.WriteFile(dirs.BranchesMetadata, jsonData, 0644); err != nil {
		log.Fatal(err)
	}
}

func ListBranches() []string {
	entries, err := os.ReadDir(dirs.Branches)
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
