package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(historyCmd)
}

var historyCmd = &cobra.Command{
	Use:     "history",
	Short:   "List all commits for the current branch",
	Example: "csync history",
	Args:    cobra.NoArgs,
	Run: func(_ *cobra.Command, args []string) {
		Debug("Starting history command")
		runHistoryCommand()
	},
}

type History struct {
	author  string
	email   string
	date    string
	message string
	commits []LogFileEntry
}

func runHistoryCommand() (returnCode int, history []History) {
	initialized := IsInitialized()
	if !initialized {
		color.Red(COMMON_RETURN_CODES[001])
		return
	}

	commits := GetCommits()
	if len(*commits) == 0 {
		Debug("No commits found")
		color.Cyan("No commits found")
		return
	}
	if len(*commits) > 20 {
		Debug("Limiting to last 20 commits")
		*commits = (*commits)[:20]
	}

	Debug("Displaying %d commits", len(*commits))

	history = make([]History, 0, len(*commits))

	for _, commit := range *commits {
		Debug("Processing commit: %s", commit.Id)
		color.Yellow(commit.Id[:40])
		data, err := os.ReadFile(dirs.Commits + commit.Id + "/metadata.json")
		if err != nil {
			Debug("Failed to read commit metadata")
			MustSucceed(err, "operation failed")
		}
		var metadata CommitMetadata
		if err = json.Unmarshal(data, &metadata); err != nil {
			Debug("Failed to unmarshal commit metadata")
			MustSucceed(err, "operation failed")
		}
		color.Cyan("Author:  " + metadata.Author)
		color.Cyan("Date:    " + commit.Timestamp)
		color.Cyan("Message: " + metadata.Message)
		fmt.Println()

		data, err = os.ReadFile(dirs.Commits + commit.Id + "/logs.json")
		if err != nil {
			Debug("Failed to read commit logs")
			MustSucceed(err, "operation failed")
		}
		var logs []LogFileEntry
		if err = json.Unmarshal(data, &logs); err != nil {
			Debug("Failed to unmarshal commit logs")
			MustSucceed(err, "operation failed")
		}
		Debug("Displaying %d log entries for commit", len(logs))
		PrintLogs(logs)
		fmt.Println()
		history = append(history, History{
			author:  metadata.Author,
			date:    commit.Timestamp,
			message: metadata.Message,
			commits: logs,
		})
	}
	Debug("History command completed successfully")
	return 401, history
}
