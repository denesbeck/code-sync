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
	RunE: func(_ *cobra.Command, args []string) error {
		return runHistoryCommand()
	},
}

func runHistoryCommand() error {
	initialized := IsInitialized()
	if !initialized {
		color.Red("CSync not initialized")
		return nil
	}

	commits := GetCommits()
	if len(commits) == 0 {
		color.Cyan("No commits found")
		return nil
	}
	if len(commits) > 20 {
		commits = commits[:20]
	}

	for _, commit := range commits {
		color.Yellow(commit.Id[:40])
		data, err := os.ReadFile("./.csync/commits/" + commit.Id + "/metadata.json")
		if err != nil {
			log.Fatal(err)
		}
		var metadata CommitMetadata
		if err = json.Unmarshal(data, &metadata); err != nil {
			log.Fatal(err)
		}
		color.Cyan("Author:  " + metadata.Author)
		color.Cyan("Date:    " + commit.Timestamp)
		color.Cyan("Message: " + metadata.Message)
		fmt.Println()

		data, err = os.ReadFile("./.csync/commits/" + commit.Id + "/logs.json")
		if err != nil {
			log.Fatal(err)
		}
		var logs []LogFileEntry
		if err = json.Unmarshal(data, &logs); err != nil {
			log.Fatal(err)
		}
		PrintLogs(logs)
		fmt.Println()
	}
	return nil
}
