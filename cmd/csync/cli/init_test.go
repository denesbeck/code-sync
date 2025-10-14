package cli

import (
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()
	for _, dir := range dirsArr {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("Directory %s not created", dir)
		}
	}
	os.RemoveAll(namespace)
}
