package cli

import (
	"os"
	"testing"
)

func TestPurge(t *testing.T) {
	// Clean __test__ namespace
	os.RemoveAll(namespace)

	// Initialize CSync
	runInitCommand()

	// Check if directories are created
	for _, dir := range dirs.GetDirs() {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("Directory %s not created", dir)
		}
	}

	// Purge CSync
	runPurgeCommand()

	// Check if directories are purged
	for _, dir := range dirs.GetDirs() {
		if _, err := os.Stat(dir); !os.IsNotExist(err) {
			t.Errorf("Directory %s not purged", dir)
		}
	}
}
