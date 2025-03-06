package cli

import (
	"encoding/json"
	"fmt"
	"log"
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

func runHistoryCommand() {
	initialized := IsInitialized()
	if !initialized {
		Debug("CSync not initialized")
		color.Red("CSync not initialized")
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
	for _, commit := range *commits {
		Debug("Processing commit: %s", commit.Id)
		color.Yellow(commit.Id[:40])
		data, err := os.ReadFile(dirs.Commits + commit.Id + "/metadata.json")
		if err != nil {
			Debug("Failed to read commit metadata")
			log.Fatal(err)
		}
		var metadata CommitMetadata
		if err = json.Unmarshal(data, &metadata); err != nil {
			Debug("Failed to unmarshal commit metadata")
			log.Fatal(err)
		}
		color.Cyan("Author:  " + metadata.Author)
		color.Cyan("Date:    " + commit.Timestamp)
		color.Cyan("Message: " + metadata.Message)
		fmt.Println()

		data, err = os.ReadFile(dirs.Commits + commit.Id + "/logs.json")
		if err != nil {
			Debug("Failed to read commit logs")
			log.Fatal(err)
		}
		var logs []LogFileEntry
		if err = json.Unmarshal(data, &logs); err != nil {
			Debug("Failed to unmarshal commit logs")
			log.Fatal(err)
		}
		Debug("Displaying %d log entries for commit", len(logs))
		PrintLogs(logs)
		fmt.Println()
	}
	Debug("History command completed successfully")
}
