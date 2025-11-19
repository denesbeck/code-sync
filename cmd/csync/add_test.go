package main

import (
	"os"
	"strconv"
	"testing"
)

func Test_AddToStaging(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	for i := 1; i < 4; i++ {
		// create files: file1.txt, file2.txt, file3.txt
		os.Create(namespace + "file" + strconv.Itoa(i) + ".txt")
		// add files to staging
		runAddCommand(namespace+"file"+strconv.Itoa(i)+".txt", false)
		// check if files are staged
		if IsFileStaged(namespace+"file"+strconv.Itoa(i)+".txt") == false {
			t.Errorf("File not staged")
		}
	}

	os.RemoveAll(namespace)
}

func Test_AddCmdStatusCode101(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	file := namespace + "file.txt"

	os.Create(file)
	runAddCommand(file, false)
	os.Remove(file)
	result := runAddCommand(file, false)
	if result.ReturnCode != 101 {
		t.Errorf("Expected 101, got %d", result.ReturnCode)
	}

	os.RemoveAll(namespace)
}

func Test_AddCmdStatusCode102(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	file := namespace + "file.txt"

	os.Create(file)
	runAddCommand(file, false)
	os.WriteFile(file, []byte("test"), 0644)
	result := runAddCommand(file, false)
	if result.ReturnCode != 102 {
		t.Errorf("Expected 102, got %d", result.ReturnCode)
	}

	os.RemoveAll(namespace)
}

func Test_AddCmdStatusCode103(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	os.Create(namespace + "file.txt")
	runAddCommand(namespace+"file.txt", false)
	result := runAddCommand(namespace+"file.txt", false)
	if result.ReturnCode != 103 {
		t.Errorf("Expected 103, got %d", result.ReturnCode)
	}

	os.RemoveAll(namespace)
}

func Test_AddCmdStatusCode104(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	file := namespace + "file.txt"

	hash := GenRandHex(20)
	os.Create(file)
	StageAndLog(hash, file, "modified")

	os.Remove(file)

	result := runAddCommand(file, false)
	if result.ReturnCode != 104 {
		t.Errorf("Expected 104, got %d", result.ReturnCode)
	}

	os.RemoveAll(namespace)
}

func Test_AddCmdStatusCode105(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	file := namespace + "file.txt"

	hash := GenRandHex(20)
	os.Create(file)
	StageAndLog(hash, file, "modified")

	os.WriteFile(file, []byte("test"), 0644)

	result := runAddCommand(file, false)
	if result.ReturnCode != 105 {
		t.Errorf("Expected 105, got %d", result.ReturnCode)
	}

	os.RemoveAll(namespace)
}

func Test_AddCmdStatusCode106(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	file := namespace + "file.txt"

	hash := GenRandHex(20)
	os.Create(file)
	StageAndLog(hash, file, "modified")

	result := runAddCommand(file, false)
	if result.ReturnCode != 106 {
		t.Errorf("Expected 106, got %d", result.ReturnCode)
	}

	os.RemoveAll(namespace)
}

func Test_AddCmdStatusCode107(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	file := namespace + "file.txt"

	hash := GenRandHex(20)
	os.Create(file)
	StageAndLog(hash, file, "added")
	runCommitCommand("test")

	hash = GenRandHex(20)
	LogOperation(hash, "REM", file)

	os.WriteFile(file, []byte("test"), 0644)

	result := runAddCommand(file, false)
	if result.ReturnCode != 107 {
		t.Errorf("Expected 107, got %d", result.ReturnCode)
	}

	os.RemoveAll(namespace)
}

func Test_AddCmdStatusCode8(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	file := namespace + "file.txt"

	hash := GenRandHex(20)
	os.Create(file)
	StageAndLog(hash, file, "added")
	runCommitCommand("test")

	hash = GenRandHex(20)
	LogOperation(hash, "REM", file)

	os.Remove(file)

	result := runAddCommand(file, false)
	if result.ReturnCode != 108 {
		t.Errorf("Expected 108, got %d", result.ReturnCode)
	}

	os.RemoveAll(namespace)
}

func Test_AddCmdStatusCode109(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	file := namespace + "file.txt"

	hash := GenRandHex(20)
	os.Create(file)
	StageAndLog(hash, file, "added")
	runCommitCommand("test")

	os.Remove(file)

	result := runAddCommand(file, false)
	if result.ReturnCode != 109 {
		t.Errorf("Expected 109, got %d", result.ReturnCode)
	}

	os.RemoveAll(namespace)
}

func Test_AddCmdStatusCode110(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	file := namespace + "file.txt"

	hash := GenRandHex(20)
	os.Create(file)
	StageAndLog(hash, file, "added")
	runCommitCommand("test")

	os.WriteFile(file, []byte("test"), 0644)

	result := runAddCommand(file, false)
	if result.ReturnCode != 110 {
		t.Errorf("Expected 110, got %d", result.ReturnCode)
	}

	os.RemoveAll(namespace)
}

func Test_AddCmdStatusCode111(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	file := namespace + "file.txt"

	hash := GenRandHex(20)
	os.Create(file)
	StageAndLog(hash, file, "added")
	runCommitCommand("test")

	result := runAddCommand(file, false)
	if result.ReturnCode != 111 {
		t.Errorf("Expected 111, got %d", result.ReturnCode)
	}

	os.RemoveAll(namespace)
}

func Test_AddCmdStatusCode112(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	file := namespace + "file.txt"

	os.Create(file)

	result := runAddCommand(file, false)
	if result.ReturnCode != 112 {
		t.Errorf("Expected 112, got %d", result.ReturnCode)
	}

	os.RemoveAll(namespace)
}
