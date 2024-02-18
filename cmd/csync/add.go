package csync

import (
	"encoding/json"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type LogFileEntry struct {
	Op   string
	Path string
}

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:     "add",
	Short:   "This command adds the selected files to the staging area",
	Example: "csync add",
	RunE: func(_ *cobra.Command, args []string) error {
		if len(args) < 1 {
			color.Red("Please specify a file to add")
			return nil
		}
		return runAddCommand(args[0])
	},
}

func runAddCommand(filePath string) error {
	// Check if csync is initialized
	initialized := IsInitialized()
	if !initialized {
		color.Red("CSync not initialized")
		return nil
	}
	// Check if file is already in staging area
	fileStaged := IsFileStaged(filePath)
	if fileStaged {
		// TO BE IMPLEMENTED
		return nil
	} else {
		// Check if there is at least one commit registered
		latestCommitId, commitExists := GetLastCommit()

		if commitExists {
			// File should be deleted? Check if it is listed in the latest commit and missing from the working directory
			shouldBeDeleted, srcCommitId := isFileDeleted(filePath, latestCommitId)
			if shouldBeDeleted {
				AddToStaging("./.csync/commits/"+srcCommitId+"/files/"+filePath, "removed")
				logOperation("REM", filePath)
			} else {
				fileCommitted, _ := IsFileCommitted(filePath, latestCommitId)
				// Is it a new file? Check if it was listed in the latest commit
				if !fileCommitted {
					exists := FileExists(filePath)
					if !exists {
						color.Red("File does not exist")
						return nil
					}
					// Add file to staging if it was not listed in the latest commit
					AddToStaging(filePath, "added")
					logOperation("ADD", filePath)
				}
			}
		} else {
			// Check if file exists
			exists := FileExists(filePath)
			if !exists {
				color.Red("File does not exist")
				return nil
			}
			// Add file to staging
			AddToStaging(filePath, "added")
			logOperation("ADD", filePath)
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
	var payload []LogFileEntry
	if len(logs) > 0 {
		if err = json.Unmarshal(logs, &payload); err != nil {
			log.Fatal(err)
		}
	} else {
		payload = append(payload, LogFileEntry{
			Op:   op,
			Path: path,
		})
		err = writeJson(".csync/staging/logs.json", payload)
		if err != nil {
			log.Fatal(err)
		}
	}
}

/*
Check if the file is listed in the commits/<commit_id>/fileList.json
and is missing from the working directory.
This would mean that the file should be deleted.
*/
func isFileDeleted(filePath string, latestCommitId string) (isDeleted bool, srcCommitId string) {
	existsInCommits, sourceCommitId := IsFileCommitted(filePath, latestCommitId)
	existsInWorkdir := FileExists(filePath)
	return existsInCommits && !existsInWorkdir, sourceCommitId
}
