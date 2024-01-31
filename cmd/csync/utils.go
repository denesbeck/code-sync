package csync

import (
	"errors"
	"io"
	"os"
	"strings"
)

func CheckIfInitialized() bool {
	if _, err := os.Stat(".csync"); !os.IsNotExist(err) {
		return true
	}
	return false
}

func CopyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, errors.New("Not a regular file")
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()

	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func CheckIfFileExists(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}

func ParsePath(path string) (string, string) {
	tmpArr := strings.Split(path, "/")

	dirs := strings.Join(tmpArr[:len(tmpArr)-1], "/")
	file := tmpArr[len(tmpArr)-1]

	if dirs != "" {
		dirs = dirs + "/"
	}

	return dirs, file
}
