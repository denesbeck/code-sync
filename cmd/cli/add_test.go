package cli

import (
	"os"
	"strconv"
	"testing"
)

func TestAdd(t *testing.T) {
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

/*
* TODO:
* 1. staged -> add new
* 2. staged -> add modified
* 3. staged -> add removed
* 4. not staged -> deleted
* 5. not staged -> modified
* 6. not staged -> not modified
 */
