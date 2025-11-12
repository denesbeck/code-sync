package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "csync",
	Short: "CodeSync (CSync) is a version control system inspired by Git",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		MustSucceed(err, "operation failed")
	}
}
