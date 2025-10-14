package csync

import (
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const (
	// Staging directories for added, modified, removed operations
	stagingAdded    = ".csync/staging/added"
	stagingModified = ".csync/staging/modified"
	stagingRemoved  = ".csync/staging/removed"
	// Log file for tracking staging operations
	// format: { Id: <hash>, Op: ADD | MOD | REM, Path: path/to/file }
	stagingLogs = ".csync/staging/logs.json"
	commits     = ".csync/commits"
	// Initial branch is named "main"
	// "branches/<branch-name>/commits.json" stores commit hashes for the branch
	defaultBranch        = ".csync/branches/main"
	defaultBranchCommits = ".csync/branches/main/commits.json"
	// "branches/metadata.json" stores default branch and current branch names
	// format: { Default: <branch-name>, Current: <branch-name> }
	branchesMetadata = ".csync/branches/metadata.json"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:     "init",
	Short:   "This command creates an empty CSync repository",
	Example: "csync init",
	RunE: func(_ *cobra.Command, _ []string) error {
		return runInitCommand()
	},
}

func runInitCommand() error {
	// check if .csync directory already exists
	if _, err := os.Stat(".csync"); !os.IsNotExist(err) {
		color.Red("CSync already initialized")
		return nil
	}

	// create staging directories: added, modified, removed
	if err := os.MkdirAll(stagingAdded, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(stagingModified, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(stagingRemoved, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	// create staging logs file
	f, err := os.Create(stagingLogs)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()

	// create commits directory
	/*
		Structure of commits directory:
		  commits/
		    |
		    - <commit-hash>/
		      |
		      - added/
		      - modified/
		      - fileList.json
	*/
	if err := os.MkdirAll(commits, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	// create default branch directory and commits file that lists commit hashes
	if err := os.MkdirAll(defaultBranch, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	f, err = os.Create(defaultBranchCommits)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()

	// create branches metadata file which contains default branch and current branch names
	err = CreateBranchesMetadata()
	if err != nil {
		log.Fatal(err)
	}

	color.Green("CSync initialized")

	return nil
}
