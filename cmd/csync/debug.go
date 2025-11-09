package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var debugLogger = color.New(color.FgBlue).SprintFunc()

// Debug logs a debug message if DEBUG environment variable is set to "true"
func Debug(format string, args ...interface{}) {
	if os.Getenv("DEBUG") == "true" {
		fmt.Printf("[DEBUG] %s\n", debugLogger(fmt.Sprintf(format, args...)))
	}
}
