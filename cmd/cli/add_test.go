package cli

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
		runAddCommand(namespace + "file" + strconv.Itoa(i) + ".txt")
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
	runAddCommand(file)
	os.Remove(file)
	statusCode := runAddCommand(file)
	if statusCode != 101 {
		t.Errorf("Expected 101, got %d", statusCode)
	}

	os.RemoveAll(namespace)
}

func Test_AddCmdStatusCode102(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	file := namespace + "file.txt"

	os.Create(file)
	runAddCommand(file)
	os.WriteFile(file, []byte("test"), 0644)
	statusCode := runAddCommand(file)
	if statusCode != 102 {
		t.Errorf("Expected 102, got %d", statusCode)
	}

	os.RemoveAll(namespace)
}

func Test_AddCmdStatusCode103(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	os.Create(namespace + "file.txt")
	runAddCommand(namespace + "file.txt")
	statusCode := runAddCommand(namespace + "file.txt")
	if statusCode != 103 {
		t.Errorf("Expected 103, got %d", statusCode)
	}

	os.RemoveAll(namespace)
}

func Test_AddCmdStatusCode104(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	file := namespace + "file.txt"

	hash := GenRandHex(20)
	os.Create(file)
	stageAndLog(hash, file, "modified")

	os.Remove(file)

	statusCode := runAddCommand(file)
	if statusCode != 104 {
		t.Errorf("Expected 104, got %d", statusCode)
	}

	os.RemoveAll(namespace)
}

func Test_AddCmdStatusCode105(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	file := namespace + "file.txt"

	hash := GenRandHex(20)
	os.Create(file)
	stageAndLog(hash, file, "modified")

	os.WriteFile(file, []byte("test"), 0644)

	statusCode := runAddCommand(file)
	if statusCode != 105 {
		t.Errorf("Expected 105, got %d", statusCode)
	}

	os.RemoveAll(namespace)
}

func Test_AddCmdStatusCode106(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	file := namespace + "file.txt"

	hash := GenRandHex(20)
	os.Create(file)
	stageAndLog(hash, file, "modified")

	statusCode := runAddCommand(file)
	if statusCode != 106 {
		t.Errorf("Expected 106, got %d", statusCode)
	}

	os.RemoveAll(namespace)
}

func Test_AddCmdStatusCode107(t *testing.T) {
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
	if statusCode != 107 {
		t.Errorf("Expected 107, got %d", statusCode)
	}

	os.RemoveAll(namespace)
}

func Test_AddCmdStatusCode8(t *testing.T) {
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
	if statusCode != 108 {
		t.Errorf("Expected 108, got %d", statusCode)
	}

	os.RemoveAll(namespace)
}

func Test_AddCmdStatusCode109(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	file := namespace + "file.txt"

	hash := GenRandHex(20)
	os.Create(file)
	stageAndLog(hash, file, "added")
	runCommitCommand("test")

	os.Remove(file)

	statusCode := runAddCommand(file)
	if statusCode != 109 {
		t.Errorf("Expected 109, got %d", statusCode)
	}

	os.RemoveAll(namespace)
}

func Test_AddCmdStatusCode110(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	file := namespace + "file.txt"

	hash := GenRandHex(20)
	os.Create(file)
	stageAndLog(hash, file, "added")
	runCommitCommand("test")

	os.WriteFile(file, []byte("test"), 0644)

	statusCode := runAddCommand(file)
	if statusCode != 110 {
		t.Errorf("Expected 110, got %d", statusCode)
	}

	os.RemoveAll(namespace)
}

func Test_AddCmdStatusCode111(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	file := namespace + "file.txt"

	hash := GenRandHex(20)
	os.Create(file)
	stageAndLog(hash, file, "added")
	runCommitCommand("test")

	statusCode := runAddCommand(file)
	if statusCode != 111 {
		t.Errorf("Expected 111, got %d", statusCode)
	}

	os.RemoveAll(namespace)
}

func Test_AddCmdStatusCode112(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	file := namespace + "file.txt"

	os.Create(file)

	statusCode := runAddCommand(file)
	if statusCode != 112 {
		t.Errorf("Expected 112, got %d", statusCode)
	}

	os.RemoveAll(namespace)
}

func Test_AddCmdStatusCode001(t *testing.T) {
	os.RemoveAll(namespace)

	statusCode := runAddCommand(namespace + "file.txt")
	if statusCode != 001 {
		t.Errorf("Expected 001, got %d", statusCode)
	}

	os.RemoveAll(namespace)
}
