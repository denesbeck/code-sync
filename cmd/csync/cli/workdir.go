package cli

import (
	"fmt"
	"sort"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(workdirCmd)
}

var workdirCmd = &cobra.Command{
	Use:     "workdir",
	Short:   "List the files that are committed",
	Example: "csync workdir",
	RunE: func(_ *cobra.Command, args []string) error {
		return runWorkdirCommand()
	},
}

func runWorkdirCommand() error {
	if initialized := IsInitialized(); !initialized {
		color.Red("CSync not initialized")
		return nil
	}
	commitId := GetLastCommit()
	content := GetFileListContent(commitId)

	sort.Slice(content, func(i, j int) bool {
		return content[i].Path < content[j].Path
	})

	if len(content) == 0 {
		color.Cyan("No files committed")
	} else {
		color.Cyan("Files committed:")
		for _, record := range content {
			fmt.Println("  - " + record.Path)
		}
	}
	return nil
}
