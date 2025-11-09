package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func WriteJson(fullPath string, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	if _, err = os.Stat(fullPath); os.IsNotExist(err) {
		path, _ := ParsePath(fullPath)
		os.MkdirAll(path, 0700)
	}
	err = os.WriteFile(fullPath, jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func GenRandHex(length int) string {
	Rando := rand.Reader
	b := make([]byte, length)
	_, _ = Rando.Read(b)
	return hex.EncodeToString(b)
}

func ParsePath(fullPath string) (path string, fileName string) {
	tmpArr := strings.Split(fullPath, "/")

	dirs := strings.Join(tmpArr[:len(tmpArr)-1], "/")
	file := tmpArr[len(tmpArr)-1]

	if dirs != "" {
		dirs = dirs + "/"
	}

	return dirs, file
}

// ValidatePath ensures a file path doesn't escape the working directory
func ValidatePath(userPath string) error {
	// Get the working directory
	workdir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Clean and make the user path absolute
	cleanPath := filepath.Clean(userPath)
	var absPath string
	if filepath.IsAbs(cleanPath) {
		absPath = cleanPath
	} else {
		absPath = filepath.Join(workdir, cleanPath)
	}

	// Check if the absolute path is within the working directory
	relPath, err := filepath.Rel(workdir, absPath)
	if err != nil {
		return err
	}

	// If the relative path starts with "..", it's trying to escape
	if strings.HasPrefix(relPath, "..") {
		return errors.New("path traversal detected: path escapes working directory")
	}

	return nil
}

func GetTimestamp() string {
	return time.Now().Format(time.RFC3339)
}

func FindIndex(arr []string, val string) int {
	for i, v := range arr {
		if v == val {
			return i
		}
	}
	return -1
}

/*
* Checks if a branch name follows GitHub's naming conventions:
*
* - Can contain alphanumeric characters, hyphens (-), underscores (_), and forward slashes (/)
* - Cannot start with a hyphen, underscore, or forward slash
* - Cannot have consecutive hyphens, underscores, or forward slashes
* - Cannot end with a forward slash
* - Cannot contain control characters or spaces
* - Cannot be empty
 */
func IsValidBranchName(name string) bool {
	Debug("Validating branch name: %s", name)

	// Empty check
	if name == "" {
		Debug("Branch name cannot be empty")
		return false
	}

	// Check if starts with invalid characters
	if strings.HasPrefix(name, "-") || strings.HasPrefix(name, "_") || strings.HasPrefix(name, "/") {
		Debug("Branch name cannot start with -, _, or /")
		return false
	}

	// Check for consecutive special characters
	if strings.Contains(name, "--") || strings.Contains(name, "__") || strings.Contains(name, "//") {
		Debug("Branch name cannot contain consecutive -, _, or /")
		return false
	}

	// Check if ends with forward slash
	if strings.HasSuffix(name, "/") {
		Debug("Branch name cannot end with /")
		return false
	}

	// Main regex pattern:
	// ^[a-zA-Z0-9] - Must start with alphanumeric
	// [a-zA-Z0-9\-_/]* - Can contain alphanumeric, hyphens, underscores, and forward slashes
	// $ - End of string
	pattern := `^[a-zA-Z0-9][a-zA-Z0-9\-_/]*$`
	matched, err := regexp.MatchString(pattern, name)
	if err != nil {
		Debug("Error validating branch name: %v", err)
		return false
	}

	Debug("Branch name validation result: %v", matched)
	return matched
}
