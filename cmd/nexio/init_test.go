package main

import (
	"os"
	"testing"
)

func Test_Init(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()
	for _, dir := range dirs.GetDirs() {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("Directory %s not created", dir)
		}
	}
	os.RemoveAll(namespace)
}

func Test_IsInitialized(t *testing.T) {
	os.RemoveAll(namespace)

	if IsInitialized() {
		t.Errorf("Nexio initialized")
	}

	runInitCommand()
	if !IsInitialized() {
		t.Errorf("Nexio not initialized")
	}

	os.RemoveAll(namespace)
}
