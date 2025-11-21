package main

import (
	"fmt"
	"os"
)

func FatalError(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	Debug("Fatal error: %s", message)
	BreakLine()
	Fail(fmt.Sprintf("Error: %s", message))
	fmt.Fprintln(os.Stderr)
	Fail("Try running with " + Code("DEBUG=true") + " for more information.")
	BreakLine()
	os.Exit(1)
}

func MustSucceed(err error, context string) {
	if err != nil {
		FatalError("%s -- %v", context, ErrorMsg(err.Error()))
	}
}
