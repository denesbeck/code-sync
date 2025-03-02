package cli

import (
	"bufio"
	"io"
	"log"
	"os"
)

func CopyFile(src, dst string) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		log.Fatal("Source file does not exist")
	}

	if !sourceFileStat.Mode().IsRegular() {
		log.Fatal("Source file is not a regular file")
	}

	source, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer source.Close()

	if FileExists(dst) {
		os.Remove(dst)
	}

	path, _ := ParsePath(dst)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0700)
	}

	destination, err := os.Create(dst)
	if err != nil {
		log.Fatal(err)
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		log.Fatal(err)
	}
}

func RemoveFile(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		log.Fatal(err)
	}
}

func EmptyDir(path string) {
	if err := os.RemoveAll(path); err != nil {
		log.Fatal((err))
	}
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Fatal((err))
	}
}

func FileExists(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}

func IsModified(file1, file2 string) bool {
	f1, err := os.Open(file1)
	if err != nil {
		log.Fatal(err)
	}
	defer f1.Close()
	f2, err := os.Open(file2)
	if err != nil {
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
			return true
		}
	}
	return false
}
