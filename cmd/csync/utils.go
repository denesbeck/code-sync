package csync

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/fatih/color"
)

type File struct {
	Timestamp string
	Path      string
}

func readCsyncIgnore() []string {
	_, err := os.Stat(".csyncignore.json")
	if os.IsNotExist(err) {
		color.Cyan("INFO: .csyncignore.json not found")
		return []string{}
	}
	if err != nil {
		log.Fatal(err)
		return nil
	}
	var payload []string
	content, err := os.ReadFile(".csyncignore.json")
	if err != nil {
		log.Fatal("Error while opening file: ", err)
	}
	if err = json.Unmarshal(content, &payload); err != nil {
		log.Fatal("Error while parsing data: ", err)
	}
	color.Cyan("INFO: .csyncignore.json found")
	return payload
}

func writeJson(path string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}

func CreateFileList() error {
	var fileList []File
	content := readCsyncIgnore()

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// TO BE IMPLEMENTED: wildcards
		isTarget := !info.IsDir() && !strings.HasPrefix(path, ".csync/") && info.Name() != ".csyncignore.json" && !slices.Contains(content, path)

		if isTarget {
			fileList = append(fileList, File{
				Timestamp: info.ModTime().UTC().String(),
				Path:      path,
			})
		}
		err = writeJson(".csync/staging/fileList.json", fileList)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

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
