package csync

import (
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const (
	stagingAdded         = ".csync/staging/added"
	stagingModified      = ".csync/staging/modified"
	stagingRemoved       = ".csync/staging/removed"
	commits              = ".csync/commits"
	defaultBranchOrigin  = ".csync/branches/main/original"
	defaultBranchCommits = ".csync/branches/main/commits.txt"
	branchesMetadata     = ".csync/branches/metadata.txt"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:     "init",
	Short:   "This command creates an empty CSync repository.",
	Example: "csync init",
	RunE: func(_ *cobra.Command, _ []string) error {
		return runInitCommand()
	},
}

func runInitCommand() error {
	if _, err := os.Stat(".csync"); !os.IsNotExist(err) {
		color.Red("CSync already initialized")
		log.Fatal(err)
	}
	if err := os.MkdirAll(stagingAdded, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(stagingModified, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(stagingRemoved, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(commits, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(defaultBranchOrigin, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	f, err := os.Create(defaultBranchCommits)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
	f, err = os.Create(branchesMetadata)
	if err != nil {
		log.Fatal(err)
	}
	f.WriteString("default=main")
	f.Close()

	color.Green("CSync initialized")

	return nil
}
