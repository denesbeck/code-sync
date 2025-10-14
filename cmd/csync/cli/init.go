package cli

import (
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const (
	// Staging directories for `added`, `modified`, `removed` operations.
	stagingAdded    = ".csync/staging/added"
	stagingModified = ".csync/staging/modified"
	stagingRemoved  = ".csync/staging/removed"
	// Log file for tracking staging operations.
	// Format: { Id: <hash>, Op: ADD | MOD | REM, Path: path/to/file }
	stagingLogs = ".csync/staging/logs.json"

	// Commits directory stores directories for each commit hash.
	// `commits/<commit-hash>/<file-id>/<file-name>` refers to the file in the commit.
	// `commits/<commit-hash>/logs.json` is a copy of the staging logs file at the time of the commit.
	// `commits/<commit-hash>/metadata.json` stores metadata for the commit, e.g. commit message, timestamp.
	commits = ".csync/commits"
	// For each commit hash a file called `commits/<commit-hash>/fileList.json` will be created. It represents the project state at the time of the commit listing all the files with commit hashes.
	// Format: { Id: <hash>, CommitId: <hash>, Path: <base64-encoded path> }
	// Before each commit, the `fileList.json` will be copied from the previous commit. This file will be updated according to the changes made in the commit.
	// Whenever a file is added to the project, it is added to the `fileList.json` file.
	// Whenever a file is modified, its commit hash is updated in the fileList.json file with the new commit hash.
	// Whenever a file is removed from the project, it is removed from the fileList.json file.

	// Initial branch is named `main`.
	// "branches/<branch-name>/commits.json" stores commit hashes for the branch.
	defaultBranch = ".csync/branches/main"

	// "branches/<branch-name>/commits.json" stores commit hashes for the given branch.
	// Format: [ { Id: <commit-hash>, Timestamp: <timestamp> }, ... ]
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
	Short:   "Initialize the CSync version control system",
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
				          - <file-name> (file in the commit)
				          - metadata.json
		              - logs.json
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
	_, err = f.WriteString("[]")
	if err != nil {
		log.Fatal(err)
	}
	f.Close()

	// create branches metadata file which contains default branch and current branch names
	CreateBranchesMetadata()

	color.Green("CSync initialized")

	return nil
}
