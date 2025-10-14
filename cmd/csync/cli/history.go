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
	if len(commits) > 10 {
		commits = commits[:10]
	}

	color.Cyan("Commits registered:")
	for i, commit := range commits {
		color.Blue("(" + fmt.Sprint(i+1) + ") " + commit[:8] + "..." + commit[len(commit)-8:])
		logs, err := os.ReadFile("./.csync/commits/" + commit + "/logs.json")
		if err != nil {
			log.Fatal(err)
		}
		var content []LogFileEntry
		if err = json.Unmarshal(logs, &content); err != nil {
			log.Fatal(err)
		}
		PrintLogs(content)
	}
	return nil
}
