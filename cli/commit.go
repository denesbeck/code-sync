package cli

import (
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
	Short:   "This command commits the staged files",
	Example: "csync commit",
	RunE: func(_ *cobra.Command, args []string) error {
		return runCommitCommand(Message)
	},
}

func runCommitCommand(message string) error {
	// TODO: Implement commit command
	return nil
}
