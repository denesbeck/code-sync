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
	Args:    cobra.NoArgs,
	Run: func(_ *cobra.Command, args []string) {
		runWorkdirCommand()
	},
}

func runWorkdirCommand() {
	if initialized := IsInitialized(); !initialized {
		color.Red("CSync not initialized")
		return
	}
	commitId := GetLastCommit()
	if commitId == "" {
		color.Cyan("No commits yet")
		return
	}
	content := GetFileListContent(commitId)

	sort.Slice(*content, func(i, j int) bool {
		return (*content)[i].Path < (*content)[j].Path
	})

	if len(*content) == 0 {
		color.Cyan("No files committed")
	} else {
		color.Cyan("Files committed:")
		for _, record := range *content {
			fmt.Println("  - " + record.Path)
		}
	}
}
