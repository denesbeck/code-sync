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
	Short:   "Commit the staged files",
	Example: "csync commit -m <your commit message>",
	RunE: func(_ *cobra.Command, args []string) error {
		return runCommitCommand(Message)
	},
}

func runCommitCommand(message string) error {
	initialized := IsInitialized()
	if !initialized {
		color.Red("CSync not initialized")
	}
	empty := IsLogEntryEmpty()

	if empty {
		color.Red("Nothing to commit")
		return nil
	}

	newCommitId := GenRandHex(32)
	latestCommitId := GetLastCommit()

	ProcessFileList(latestCommitId, newCommitId)

	WriteCommitMetadata(newCommitId, message)

	CopyFile("./.csync/staging/logs.json", "./.csync/commits/"+newCommitId+"/logs.json")

	TruncateLogs()
	EmptyDir("./.csync/staging/added/")
	EmptyDir("./.csync/staging/modified/")
	EmptyDir("./.csync/staging/removed/")

	RegisterCommitForBranch(newCommitId)

	color.Green("Committed successfully")
	return nil
}
