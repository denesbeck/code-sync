package csync

import (
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const (
	stagingAdded          = ".csync/staging/added"
	stagingModified       = ".csync/staging/modified"
	stagingRemoved        = ".csync/staging/removed"
	stagingLogs           = ".csync/staging/logs.json"
	stagingFileList       = ".csync/staging/filelist.json"
	commits               = ".csync/commits"
	defaultBranchOriginal = ".csync/branches/main/original"
	defaultBranchCommits  = ".csync/branches/main/commits.json"
	branchesMetadata      = ".csync/branches/metadata.json"
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
		return nil
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
	f, err := os.Create(stagingLogs)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()

	err = CreateFileList()
	if err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(commits, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(defaultBranchOriginal, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	f, err = os.Create(defaultBranchCommits)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()

	err = CreateBranchesMetadata()
	if err != nil {
		log.Fatal(err)
	}

	color.Green("CSync initialized")

	return nil
}
