package cli

import (
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()
	for _, dir := range dirs.GetDirs() {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("Directory %s not created", dir)
		}
	}
	os.RemoveAll(namespace)
}

func TestIsInitialized(t *testing.T) {
	os.RemoveAll(namespace)

	if IsInitialized() {
		t.Errorf("CSync initialized")
	}

	runInitCommand()
	if !IsInitialized() {
		t.Errorf("CSync not initialized")
	}

	os.RemoveAll(namespace)
}
