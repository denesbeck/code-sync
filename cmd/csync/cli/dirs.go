package cli

import (
	"reflect"
)

type Dirs struct {
	Root                 string
	StagingAdded         string
	StagingModified      string
	StagingRemoved       string
	StagingLogs          string
	Commits              string
	Branches             string
	DefaultBranch        string
	DefaultBranchCommits string
	BranchesMetadata     string
	Config               string
}

var dirs = Dirs{
	Root: namespace + ".csync",
	// Staging directories for `added`, `modified`, `removed` operations.
	StagingAdded:    namespace + ".csync/staging/added",
	StagingModified: namespace + ".csync/staging/modified",
	StagingRemoved:  namespace + ".csync/staging/removed",

	// Log file for tracking staging operations.
	// Format: { Id: <hash>, Op: ADD | MOD | REM, Path: path/to/file }
	StagingLogs: namespace + ".csync/staging/logs.json",

	// Commits directory stores directories for each commit hash.
	// `commits/<commit-hash>/<file-id>/<file-name>`: refers to the file in the commit.
	// `commits/<commit-hash>/logs.json`: copy of the staging logs file at the time of the commit.
	// Format: { Id: <hash>, Op: ADD | MOD | REM, Path: path/to/file }
	// `commits/<commit-hash>/metadata.json` stores metadata for the commit, e.g. commit message, timestamp.
	// Format: { Author: <username <user-email>>, Message: <commit-message> }
	// For each commit hash a file called `commits/<commit-hash>/fileList.json` will be created. It represents the project state at the time of the commit listing all the files with commit hashes.
	// Format: { Id: <hash>, CommitId: <hash>, Path: path/to/file }
	// Before each commit, the `fileList.json` will be copied from the previous commit. This file will be updated according to the changes made in the commit.
	// Whenever a file is added to the project, it is added to the `fileList.json` file.
	// Whenever a file is modified, its commit hash is updated in the fileList.json file with the new commit hash.
	// Whenever a file is removed from the project, it is removed from the fileList.json file.
	Commits: namespace + ".csync/commits",

	Branches: namespace + ".csync/branches",

	// Initial branch is named `main`.
	DefaultBranch: namespace + ".csync/branches/main",

	// "branches/<branch-name>/commits.json" stores commit hashes for the given branch.
	// Format: [ { Id: <commit-hash>, Timestamp: <timestamp> }, ... ]
	DefaultBranchCommits: namespace + ".csync/branches/main/commits.json",

	// "branches/metadata.json" stores default branch and current branch names.
	// Format: { Default: <branch-name>, Current: <branch-name> }
	BranchesMetadata: namespace + ".csync/branches/metadata.json",

	// "config.json" stores CSync config data, e.g. username, email.
	// Format: { Username: <username>, Email: <email> }
	Config: namespace + ".csync/config.json",
}

func (d Dirs) GetDirs() []string {
	fields := reflect.ValueOf(d)
	var dirs []string
	for i := 0; i < fields.NumField(); i++ {
		dirs = append(dirs, fields.Field(i).String())
	}
	return dirs
}
