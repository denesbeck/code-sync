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
		Debug("Starting commit command with message: %s", Message)
		runCommitCommand(Message)
	},
}

func runCommitCommand(message string) {
	initialized := IsInitialized()
	if !initialized {
		color.Red(COMMON_RETURN_CODES[001])
		return
	}

	empty := IsStagingLogsEmpty()
	if empty {
		Debug("No changes staged for commit")
		color.Red("Nothing to commit")
		return
	}

	newCommitId := GenRandHex(20)
	latestCommitId := GetLastCommit().Id
	Debug("Creating new commit: id=%s, parent=%s", newCommitId, latestCommitId)

	ProcessFileList(latestCommitId, newCommitId)
	Debug("Processed file list for commit")

	WriteCommitMetadata(newCommitId, message)
	Debug("Wrote commit metadata")

	CopyFile(dirs.StagingLogs, dirs.Commits+newCommitId+"/logs.json")
	Debug("Copied staging logs to commit")

	TruncateLogs()
	EmptyDir(dirs.StagingAdded)
	EmptyDir(dirs.StagingModified)
	EmptyDir(dirs.StagingRemoved)
	Debug("Cleaned up staging area")

	RegisterCommitForBranch(newCommitId)
	Debug("Registered commit for current branch")

	color.Green("Changes committed successfully")
}
