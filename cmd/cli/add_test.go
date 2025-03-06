package cli

import (
	"os"
	"strconv"
	"testing"
)

func TestAddNew(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	for i := 1; i < 4; i++ {
		// create files: file1.txt, file2.txt, file3.txt
		os.Create(namespace + "file" + strconv.Itoa(i) + ".txt")
		// add files to staging
		runAddCommand(namespace + "file" + strconv.Itoa(i) + ".txt")
		// check if files are staged
		if IsFileStaged(namespace+"file"+strconv.Itoa(i)+".txt") == false {
			t.Errorf("File not staged")
		}
	}

	os.RemoveAll(namespace)
}

func TestAddStatusCode1(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	file := namespace + "file.txt"

	os.Create(file)
	runAddCommand(file)
	os.Remove(file)
	statusCode := runAddCommand(file)
	if statusCode != 1 {
		t.Errorf("Expected 1, got %d", statusCode)
	}

	os.RemoveAll(namespace)
}

func TestAddStatusCode2(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	file := namespace + "file.txt"

	os.Create(file)
	runAddCommand(file)
	os.WriteFile(file, []byte("test"), 0644)
	statusCode := runAddCommand(file)
	if statusCode != 2 {
		t.Errorf("Expected 2, got %d", statusCode)
	}

	os.RemoveAll(namespace)
}

func TestAddStatusCode3(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	os.Create(namespace + "file.txt")
	runAddCommand(namespace + "file.txt")
	statusCode := runAddCommand(namespace + "file.txt")
	if statusCode != 3 {
		t.Errorf("Expected 3, got %d", statusCode)
	}

	os.RemoveAll(namespace)
}

func TestAddStatusCode4(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	file := namespace + "file.txt"

	hash := GenRandHex(20)
	os.Create(file)
	stageAndLog(hash, file, "modified")

	os.Remove(file)

	statusCode := runAddCommand(file)
	if statusCode != 4 {
		t.Errorf("Expected 4, got %d", statusCode)
	}

	os.RemoveAll(namespace)
}

func TestAddStatusCode5(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	file := namespace + "file.txt"

	hash := GenRandHex(20)
	os.Create(file)
	stageAndLog(hash, file, "modified")

	os.WriteFile(file, []byte("test"), 0644)

	statusCode := runAddCommand(file)
	if statusCode != 5 {
		t.Errorf("Expected 5, got %d", statusCode)
	}

	os.RemoveAll(namespace)
}

func TestAddStatusCode6(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	file := namespace + "file.txt"

	hash := GenRandHex(20)
	os.Create(file)
	stageAndLog(hash, file, "modified")

	statusCode := runAddCommand(file)
	if statusCode != 6 {
		t.Errorf("Expected 6, got %d", statusCode)
	}

	os.RemoveAll(namespace)
}

func TestAddStatusCode7(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	file := namespace + "file.txt"

	hash := GenRandHex(20)
	os.Create(file)
	stageAndLog(hash, file, "added")
	runCommitCommand("test")

	hash = GenRandHex(20)
	LogOperation(hash, "REM", file)

	os.WriteFile(file, []byte("test"), 0644)

	statusCode := runAddCommand(file)
	if statusCode != 7 {
		t.Errorf("Expected 7, got %d", statusCode)
	}

	os.RemoveAll(namespace)
}

func TestAddStatusCode8(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	file := namespace + "file.txt"

	hash := GenRandHex(20)
	os.Create(file)
	stageAndLog(hash, file, "added")
	runCommitCommand("test")

	hash = GenRandHex(20)
	LogOperation(hash, "REM", file)

	os.Remove(file)

	statusCode := runAddCommand(file)
	if statusCode != 8 {
		t.Errorf("Expected 8, got %d", statusCode)
	}

	os.RemoveAll(namespace)
}

func TestAddStatusCode9(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	file := namespace + "file.txt"

	hash := GenRandHex(20)
	os.Create(file)
	stageAndLog(hash, file, "added")
	runCommitCommand("test")

	os.Remove(file)

	statusCode := runAddCommand(file)
	if statusCode != 9 {
		t.Errorf("Expected 9, got %d", statusCode)
	}

	os.RemoveAll(namespace)
}

func TestAddStatusCode10(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	file := namespace + "file.txt"

	hash := GenRandHex(20)
	os.Create(file)
	stageAndLog(hash, file, "added")
	runCommitCommand("test")

	os.WriteFile(file, []byte("test"), 0644)

	statusCode := runAddCommand(file)
	if statusCode != 10 {
		t.Errorf("Expected 10, got %d", statusCode)
	}

	os.RemoveAll(namespace)
}

func TestAddStatusCode11(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	file := namespace + "file.txt"

	hash := GenRandHex(20)
	os.Create(file)
	stageAndLog(hash, file, "added")
	runCommitCommand("test")

	statusCode := runAddCommand(file)
	if statusCode != 11 {
		t.Errorf("Expected 11, got %d", statusCode)
	}

	os.RemoveAll(namespace)
}

func TestAddStatusCode12(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	file := namespace + "file.txt"

	os.Create(file)

	statusCode := runAddCommand(file)
	if statusCode != 12 {
		t.Errorf("Expected 12, got %d", statusCode)
	}

	os.RemoveAll(namespace)
}
