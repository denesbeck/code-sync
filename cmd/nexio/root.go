package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nexio",
	Short: "Nexio (Nexio) is a version control system inspired by Git",
}

func Execute() {
	rootCmd.Execute()
}
