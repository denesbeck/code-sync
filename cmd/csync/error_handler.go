package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

// FatalError prints a user-friendly error message and exits with status code 1
// This centralizes error handling and provides better error messages than log.Fatal
func FatalError(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	Debug("Fatal error: %s", message)
	color.Red("Error: %s", message)
	fmt.Fprintln(os.Stderr)
	color.Yellow("If this is unexpected, try running with DEBUG=true for more information")
	os.Exit(1)
}

// MustSucceed checks for an error and exits with a user-friendly message if found
// This replaces log.Fatal with better error context
func MustSucceed(err error, context string) {
	if err != nil {
		FatalError("%s: %v", context, err)
	}
}
