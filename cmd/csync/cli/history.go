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
		runHistoryCommand()
	},
}

func runHistoryCommand() {
	initialized := IsInitialized()
	if !initialized {
		color.Red("CSync not initialized")
		return
	}

	commits := GetCommits()
	if len(commits) == 0 {
		color.Cyan("No commits found")
		return
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
}
