package cli

import (
	"encoding/json"
	"log"
	"os"
)

type BranchMetadata struct {
	Default string
	Current string
}

func GetCurrentBranchName() string {
	branchesMetadata, err := os.ReadFile(".csync/branches/metadata.json")
	if err != nil {
		log.Fatal(err)
	}
	var metadata BranchMetadata
	if err = json.Unmarshal(branchesMetadata, &metadata); err != nil {
		log.Fatal(err)
	}
	return metadata.Current
}

func GetDefaultBranchName() string {
	branchesMetadata, err := os.ReadFile(".csync/branches/metadata.json")
	if err != nil {
		log.Fatal(err)
	}
	var metadata BranchMetadata
	if err = json.Unmarshal(branchesMetadata, &metadata); err != nil {
		log.Fatal(err)
	}
	return metadata.Default
}

func CreateBranchesMetadata() {
	branchesMetadata := BranchMetadata{
		Default: "main",
		Current: "main",
	}
	WriteJson(".csync/branches/metadata.json", branchesMetadata)
}
