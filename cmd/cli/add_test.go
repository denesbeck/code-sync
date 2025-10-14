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
		os.Create(namespace + "/file" + strconv.Itoa(i) + ".txt")
		// add files to staging
		runAddCommand(namespace + "/file" + strconv.Itoa(i) + ".txt")
		// check if files are staged
		if IsFileStaged(namespace+"/file"+strconv.Itoa(i)+".txt") == false {
			t.Errorf("File not staged")
		}
	}

	os.RemoveAll(namespace)
}

func TestAddStatusCode1(t *testing.T) {
	os.RemoveAll(namespace)

	runInitCommand()

	file := namespace + "/file.txt"

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

	file := namespace + "/file.txt"

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

	os.Create(namespace + "/file.txt")
	runAddCommand(namespace + "/file.txt")
	statusCode := runAddCommand(namespace + "/file.txt")
	if statusCode != 3 {
		t.Errorf("Expected 3, got %d", statusCode)
	}

	os.RemoveAll(namespace)
}

// func TestAddStatusCode4(t *testing.T) {
// 	os.RemoveAll(namespace)
//
// 	runInitCommand()
//
// 	file := namespace + "/file.txt"
//
// 	hash := GenRandHex(20)
// 	os.Create(file)
// 	stageAndLog(hash, file, "modified")
//
// 	os.Remove(file)
//
// 	statusCode := runAddCommand(file)
// 	if statusCode != 4 {
// 		t.Errorf("Expected 4, got %d", statusCode)
// 	}
//
// 	os.RemoveAll(namespace)
// }
