package csync

import (
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:     "add",
	Short:   "This command adds the selected files to the staging area.",
	Example: "csync add",
	RunE: func(_ *cobra.Command, args []string) error {
		if len(args) < 1 {
			color.Red("Please specify a file to add")
			return nil
		}
		return runAddCommand(args[0])
	},
}

func runAddCommand(path string) error {
	initialized := CheckIfInitialized()
	if !initialized {
		color.Red("CSync not initialized")
		return nil
	}
	exists := CheckIfFileExists(path)
	if !exists {
		color.Red("File does not exist")
		return nil
	}
	fileInLogs := CheckIfFileInLogs(path)
	if fileInLogs {
		color.Red("File already added")
		return nil
	}

	dirs, file := ParsePath(path)

	fullPath := ".csync/staging/added/" + dirs

	if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	_, err := CopyFile(path, fullPath+file)
	if err != nil {
		log.Fatal(err)
	}
	color.Green("File added successfully")
	return nil
}
