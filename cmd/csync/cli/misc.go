package cli

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"os"
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

func getTimestamp() string {
	return time.Now().Format(time.RFC3339)
}
