package cli

import (
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const (
	// Staging directories for added, modified, removed operations.
	stagingAdded    = ".csync/staging/added"
	stagingModified = ".csync/staging/modified"
	stagingRemoved  = ".csync/staging/removed"
	// Log file for tracking staging operations.
	// Format: { id: <hash>, op: add | mod | rem, path: path/to/file }
	stagingLogs = ".csync/staging/logs.json"

	// Commits directory stores directories for each commit hash.
	// "commits/<commit-hash>/file-name" refers to the file in the commit.
	// "commits/<commit-hash>/metadata.json" stores metadata for the commit, e.g. commit message, timestamp.
	commits = ".csync/commits"

	// Initial branch is named "main".
	// "branches/<branch-name>/commits.json" stores commit hashes for the branch.
	defaultBranch = ".csync/branches/main"
	// Whenever a file is added to the project, it is added to the commits.json file.
	// Format: { Commit: <hash>, Path: <path>, Name: <file-name> }
	// Whenever a file is removed from the project, it is removed from the commits.json file.
	// Whenever a file is modified, its commit hash is updated in the commits.json file.
	defaultBranchCommits = ".csync/branches/main/commits.json"
	// "branches/metadata.json" stores default branch and current branch names.
	// Format: { Default: <branch-name>, Current: <branch-name> }
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
