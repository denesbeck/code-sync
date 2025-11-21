package main

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
	Example: "nexio commit -m <your commit message>",
	Args:    cobra.NoArgs,
	Run: func(_ *cobra.Command, args []string) {
		Debug("Starting commit command with message: %s", Message)
		runCommitCommand(Message)
	},
}

func runCommitCommand(message string) (returnCode int, commitId string) {
	initialized := IsInitialized()
	if !initialized {
		color.Red(COMMON_RETURN_CODES[001])
		return 001, ""
	}

	// Clean up any orphaned staging entries from previous failed operations
	CleanOrphanedStagingEntries()

	empty := IsStagingLogsEmpty()
	if empty {
		Debug("No changes staged for commit")
		color.Red(COMMIT_RETURN_CODES[701])
		return 701, ""
	}

	newCommitId := GenRandHex(20)
	latestCommitId := GetLastCommit().Id
	Debug("Creating new commit: id=%s, parent=%s", newCommitId, latestCommitId)

	ProcessFileList(latestCommitId, newCommitId)
	Debug("Processed file list for commit")

	WriteCommitMetadata(newCommitId, message)
	Debug("Wrote commit metadata")

	if err := CopyFile(dirs.StagingLogs, dirs.Commits+newCommitId+"/logs.json"); err != nil {
		color.Red("Error copying staging logs: " + err.Error())
		return 001, ""
	}
	Debug("Copied staging logs to commit")

	TruncateLogs()
	if err := EmptyDir(dirs.StagingAdded); err != nil {
		color.Red("Error emptying staging added directory: " + err.Error())
		return 001, ""
	}
	if err := EmptyDir(dirs.StagingModified); err != nil {
		color.Red("Error emptying staging modified directory: " + err.Error())
		return 001, ""
	}
	if err := EmptyDir(dirs.StagingRemoved); err != nil {
		color.Red("Error emptying staging removed directory: " + err.Error())
		return 001, ""
	}
	Debug("Cleaned up staging area")

	RegisterCommitForBranch(newCommitId)
	Debug("Registered commit for current branch")

	color.Green("Changes committed successfully")
	return 702, newCommitId
}
