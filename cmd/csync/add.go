package csync

import (
	"encoding/json"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type LogFile struct {
	op   string
	path string
}

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
	// check if csync is initialized
	initialized := IsInitialized()
	if !initialized {
		color.Red("CSync not initialized")
		return nil
	}
	// check if file is already in staging area
	fileInStaging := IsFileListed(path, ".csync/staging/logs.json")
	if fileInStaging {
		// TO BE IMPLEMENTED
		return nil
	} else {
		lastCommit, commitExists := GetLastCommit()
		// there is at least one commit
		if commitExists {
			// file should be deleted?
			isDeleted := isFileDeleted(lastCommit, path)
			if isDeleted {
				// TO BE IMPLEMENTED: move the file from the appr. commit
				MoveToStaging("./.csync/commits/"+lastCommit+"/added/"+path, "removed")
				logOperation("REM", path)
			} else {
				// new file?
				isNewFile := IsFileListed(path, ".csync/commits/"+lastCommit+"/fileList.json")
				if isNewFile {
					exists := FileExists(path)
					if !exists {
						color.Red("File does not exist")
						return nil
					}
					// add file
					MoveToStaging(path, "added")
					logOperation("ADD", path)
				}
			}
		} else {
			// check if file exists
			exists := FileExists(path)
			if !exists {
				color.Red("File does not exist")
				return nil
			}
			// add file
			MoveToStaging(path, "added")
			logOperation("ADD", path)
		}
	}

	return nil
}

// Log changes to the staging/logs.json file
func logOperation(op string, path string) {
	logs, err := os.ReadFile(".csync/staging/logs.json")
	if err != nil {
		log.Fatal(err)
	}
	var payload []LogFile
	if len(logs) > 0 {
		if err = json.Unmarshal(logs, &payload); err != nil {
			log.Fatal(err)
		}
	} else {
		payload = append(payload, LogFile{
			op:   op,
			path: path,
		})
		err = writeJson(".csync/staging/logs.json", payload)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Check if the file exists in the commits/<commit_id>/added directory
// and is missing from the working directory. This means that the file
// should be deleted.
func isFileDeleted(commit string, path string) bool {
	existsInCommits := FileExists("./.csync/commits/" + commit + "/added/" + path)
	existsInWorkdir := FileExists(path)
	return existsInCommits && !existsInWorkdir
}
