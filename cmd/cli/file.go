package cli

import (
	"bufio"
	"io"
	"log"
	"os"
)

func CopyFile(src, dst string) {
	Debug("Copying file from %s to %s", src, dst)
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		Debug("Source file does not exist: %s", src)
		log.Fatal("Source file does not exist")
	}

	if !sourceFileStat.Mode().IsRegular() {
		Debug("Source is not a regular file: %s", src)
		log.Fatal("Source file is not a regular file")
	}

	source, err := os.Open(src)
	if err != nil {
		Debug("Failed to open source file: %s", src)
		log.Fatal(err)
	}
	defer source.Close()

	if FileExists(dst) {
		Debug("Destination file exists, removing: %s", dst)
		os.Remove(dst)
	}

	path, _ := ParsePath(dst)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		Debug("Creating destination directory: %s", path)
		os.MkdirAll(path, 0700)
	}

	destination, err := os.Create(dst)
	if err != nil {
		Debug("Failed to create destination file: %s", dst)
		log.Fatal(err)
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		Debug("Failed to copy file contents")
		log.Fatal(err)
	}
	Debug("File copied successfully")
}

func RemoveFile(path string) {
	Debug("Removing file/directory: %s", path)
	err := os.RemoveAll(path)
	if err != nil {
		Debug("Failed to remove file/directory: %s", path)
		log.Fatal(err)
	}
	Debug("File/directory removed successfully")
}

func EmptyDir(path string) {
	Debug("Emptying directory: %s", path)
	if err := os.RemoveAll(path); err != nil {
		Debug("Failed to remove directory contents: %s", path)
		log.Fatal((err))
	}
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		Debug("Failed to recreate empty directory: %s", path)
		log.Fatal((err))
	}
	Debug("Directory emptied successfully")
}

func FileExists(path string) bool {
	exists := false
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		exists = true
	}
	Debug("Checking if file exists: %s = %v", path, exists)
	return exists
}

func IsModified(file1, file2 string) bool {
	Debug("Checking if files are modified: %s vs %s", file1, file2)
	f1, err := os.Open(file1)
	if err != nil {
		Debug("Failed to open first file: %s", file1)
		log.Fatal(err)
	}
	defer f1.Close()
	f2, err := os.Open(file2)
	if err != nil {
		Debug("Failed to open second file: %s", file2)
		log.Fatal(err)
	}
	defer f2.Close()
	reader1 := bufio.NewReader(f1)
	reader2 := bufio.NewReader(f2)
	for {
		line1, _, err := reader1.ReadLine()
		line2, _, _ := reader2.ReadLine()
		if err == io.EOF {
			break
		}
		if string(line1) != string(line2) {
			Debug("Files are different")
			return true
		}
	}
	Debug("Files are identical")
	return false
}
