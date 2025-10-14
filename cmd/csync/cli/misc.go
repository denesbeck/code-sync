package cli

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"
	"strings"
)

func WriteJson(path string, data interface{}) error {
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

func GenRandHex(length int) string {
	Rando := rand.Reader
	b := make([]byte, length)
	_, _ = Rando.Read(b)
	return hex.EncodeToString(b)
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
