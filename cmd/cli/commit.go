package cli

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	commitCmd.Flags().StringVarP(&Message, "message", "m", "", "Commit message (required)")
	commitCmd.MarkFlagRequired("message")

	rootCmd.AddCommand(commitCmd)
}

var Message string

var commitCmd = &cobra.Command{
	Use:     "commit",
	Short:   "Record changes to the repository",
	Example: "csync commit -m <your commit message>",
	Args:    cobra.NoArgs,
	Run: func(_ *cobra.Command, args []string) {
		runCommitCommand(Message)
	},
}

func runCommitCommand(message string) {
	initialized := IsInitialized()
	if !initialized {
		color.Red("CSync not initialized")
		return
	}

	empty := IsStagingLogsEmpty()
	if empty {
		color.Red("Nothing to commit")
		return
	}

	newCommitId := GenRandHex(20)
	latestCommitId := GetLastCommit()

	ProcessFileList(latestCommitId, newCommitId)

	WriteCommitMetadata(newCommitId, message)

	CopyFile("./.csync/staging/logs.json", "./.csync/commits/"+newCommitId+"/logs.json")

	TruncateLogs()
	EmptyDir("./.csync/staging/added/")
	EmptyDir("./.csync/staging/modified/")
	EmptyDir("./.csync/staging/removed/")

	RegisterCommitForBranch(newCommitId)

	color.Green("Changes committed successfully")
}
