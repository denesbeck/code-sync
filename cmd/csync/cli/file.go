package cli

import (
	"errors"
	"io"
	"log"
	"os"
)

func CopyFile(src, dst string) (int64, error) {
	// File exists?
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	// Is it a regular file?
	if !sourceFileStat.Mode().IsRegular() {
		return 0, errors.New("Not a regular file")
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	// Check if the file already exists. If yes, remove it.
	if FileExists(dst) {
		os.Remove(dst)
	}

	// Create the file
	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()

	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func RemoveFile(path string) {
	err := os.Remove(path)
	if err != nil {
		log.Fatal(err)
	}
}

func FileExists(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}

func IsModified(file1, file2 string) bool {
	content1, err := os.ReadFile(file1)
	if err != nil {
		log.Fatal(err)
	}
	content2, err := os.ReadFile(file2)
	if err != nil {
		log.Fatal(err)
	}
	return string(content1) != string(content2)
}
