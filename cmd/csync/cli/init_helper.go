package cli

import (
	"os"
)

func IsInitialized() bool {
	if _, err := os.Stat(".csync"); !os.IsNotExist(err) {
		return true
	}
	return false
}

func CreateBranchesMetadata() error {
	branchesMetadata := BranchMetadata{
		Default: "main",
		Current: "main",
	}

	err := WriteJson(".csync/branches/metadata.json", branchesMetadata)
	if err != nil {
		return err
	}
	return nil
}
