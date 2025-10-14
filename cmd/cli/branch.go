package cli

import (
	"fmt"
	"os"
	"slices"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	newCmd.Flags().StringVarP(&FromCommit, "from-commit", "c", "", "Commit to create branch from")
	newCmd.Flags().StringVarP(&FromBranch, "from-branch", "b", "", "Branch to create branch from")

	rootCmd.AddCommand(branchCmd)

	branchCmd.AddCommand(currentCmd)
	branchCmd.AddCommand(defaultCmd)
	branchCmd.AddCommand(newCmd)
	branchCmd.AddCommand(dropCmd)
	branchCmd.AddCommand(switchCmd)
}

var (
	FromCommit string
	FromBranch string
)

var branchCmd = &cobra.Command{
	Use:     "branch",
	Short:   "Branch management",
	Example: "csync branch",
	Args:    cobra.NoArgs,
	Run: func(_ *cobra.Command, args []string) {
		Debug("Starting branch command")
		runBranchCommand()
	},
}

var currentCmd = &cobra.Command{
	Use:     "current",
	Short:   "Get current branch",
	Example: "csync branch current",
	Args:    cobra.NoArgs,
	Run: func(_ *cobra.Command, args []string) {
		Debug("Starting current branch command")
		runCurrentCommand()
	},
}

var defaultCmd = &cobra.Command{
	Use:     "default",
	Short:   "Get default branch",
	Example: "csync branch default",
	Args:    cobra.NoArgs,
	Run: func(_ *cobra.Command, args []string) {
		Debug("Starting default branch command")
		runDefaultCommand()
	},
}

var newCmd = &cobra.Command{
	Use:     "new",
	Short:   "Create a new branch",
	Example: "csync new <branch-name> --from-commit <commit-id> --from-branch <branch-name>",
	Args:    cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		Debug("Starting new branch command: name=%s, from-commit=%s, from-branch=%s", args[0], FromCommit, FromBranch)
		runNewCommand(args[0], FromCommit, FromBranch)
	},
}

var dropCmd = &cobra.Command{
	Use:     "drop",
	Short:   "Delete a branch",
	Example: "csync drop <branch-name>",
	Args:    cobra.MinimumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		Debug("Starting drop branch command with args: %v", args)
		for _, arg := range args {
			runDropCommand(arg)
		}
	},
}

var switchCmd = &cobra.Command{
	Use:     "switch",
	Short:   "Switch to a branch",
	Example: "csync switch <branch-name>",
	Args:    cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		Debug("Starting switch branch command: branch=%s", args[0])
		runSwitchCommand(args[0])
	},
}

func runBranchCommand() {
	initialized := IsInitialized()
	if !initialized {
		color.Red(COMMON_RETURN_CODES[001])
		return
	}

	branches, err := os.ReadDir(dirs.Branches)
	if err != nil {
		Debug("No branches found")
		color.Red("No branches found")
		return
	}

	currentBranchName := GetCurrentBranchName()
	defaultBranchName := GetDefaultBranchName()
	Debug("Current branch: %s, Default branch: %s", currentBranchName, defaultBranchName)

	for _, branch := range branches {
		if branch.IsDir() {
			branchName := branch.Name()
			if branchName == defaultBranchName {
				branchName = "* " + branchName
			} else {
				branchName = "  " + branchName
			}

			if branch.Name() == currentBranchName {
				color.Green(branchName)
			} else {
				fmt.Println(branchName)
			}
		}
	}
	Debug("Branch command completed successfully")
}

func runCurrentCommand() {
	initialized := IsInitialized()
	if !initialized {
		color.Red(COMMON_RETURN_CODES[001])
		return
	}

	currentBranchName := GetCurrentBranchName()
	Debug("Current branch: %s", currentBranchName)
	fmt.Println(currentBranchName)
}

func runDefaultCommand() {
	initialized := IsInitialized()
	if !initialized {
		color.Red(COMMON_RETURN_CODES[001])
		return
	}

	defaultBranchName := GetDefaultBranchName()
	Debug("Default branch: %s", defaultBranchName)
	fmt.Println(defaultBranchName)
}

func runNewCommand(branchName string, fromCommit string, fromBranch string) int {
	initialized := IsInitialized()
	if !initialized {
		color.Red(COMMON_RETURN_CODES[001])
		return 001
	}

	if !IsValidBranchName(branchName) {
		Debug("Invalid branch name: %s", branchName)
		color.Red(BRANCH_RETURN_CODES[201])
		return 201
	}

	if fromCommit != "" && fromBranch != "" {
		Debug("Cannot create branch from both commit and branch")
		color.Red(BRANCH_RETURN_CODES[202])
		return 202
	}

	var srcBranch string
	if fromBranch != "" {
		srcBranch = fromBranch
		branches := ListBranches()
		if !slices.Contains(branches, srcBranch) {
			Debug("Source branch does not exist: %s", srcBranch)
			color.Red(BRANCH_RETURN_CODES[203])
			return 203
		}
	} else {
		srcBranch = GetCurrentBranchName()
	}

	if fromCommit != "" {
		Debug("Creating branch from commit: %s", fromCommit)
		err := CopyCommitsToBranch(fromCommit, branchName)
		if err != nil {
			Debug("Failed to create branch from commit: %v", err)
			color.Red(BRANCH_RETURN_CODES[204])
			return 204
		}
	} else {
		Debug("Creating branch from branch: %s", srcBranch)
		if err := os.Mkdir(dirs.Branches+branchName, 0755); err != nil {
			Debug("Branch already exists: %s", branchName)
			color.Red(BRANCH_RETURN_CODES[205])
			return 205
		}

		CopyFile(dirs.Branches+srcBranch+"/commits.json", dirs.Branches+branchName+"/commits.json")
	}
	Debug("Branch created successfully: %s", branchName)
	color.Green(BRANCH_RETURN_CODES[206])
	runSwitchCommand(branchName)
	return 206
}

func runDropCommand(branchName string) int {
	initialized := IsInitialized()
	if !initialized {
		color.Red(COMMON_RETURN_CODES[001])
		return 001
	}

	branches := ListBranches()
	if !slices.Contains(branches, branchName) {
		Debug("Branch does not exist: %s", branchName)
		color.Red(BRANCH_RETURN_CODES[207])
		return 207
	}

	if currentBranchName := GetCurrentBranchName(); currentBranchName == branchName {
		Debug("Cannot delete current branch: %s", branchName)
		color.Red(BRANCH_RETURN_CODES[208])
		return 208
	}

	if defaultBranchName := GetDefaultBranchName(); defaultBranchName == branchName {
		Debug("Cannot delete default branch: %s", branchName)
		color.Red(BRANCH_RETURN_CODES[209])
		return 209
	}

	if err := os.RemoveAll(dirs.Branches + branchName); err != nil {
		Debug("Failed to delete branch: %s", branchName)
		color.Red(BRANCH_RETURN_CODES[207])
		return 207
	}
	Debug("Branch deleted successfully: %s", branchName)
	color.Green(BRANCH_RETURN_CODES[210])
	return 210
}

func runSwitchCommand(branchName string) int {
	initialized := IsInitialized()
	if !initialized {
		color.Red(COMMON_RETURN_CODES[001])
		return 001
	}

	currentBranch := GetCurrentBranchName()
	if currentBranch == branchName {
		Debug("Already on branch: %s", branchName)
		color.Red(BRANCH_RETURN_CODES[211])
		return 211
	}

	branches := ListBranches()
	if !slices.Contains(branches, branchName) {
		Debug("Branch does not exist: %s", branchName)
		color.Red(BRANCH_RETURN_CODES[212])
		return 212
	}

	commitId := GetLastCommitByBranch(branchName).Id
	if commitId != "" {
		Debug("Switching to commit: %s", commitId)
		fileList := GetFileListContent(commitId)
		for _, file := range *fileList {
			_, fileName := ParsePath(file.Path)
			CopyFile(dirs.Commits+file.CommitId+"/"+file.Id+"/"+fileName, "./"+file.Path)
		}
	}
	SetBranch(branchName, "current")
	Debug("Switched to branch: %s", branchName)
	color.Cyan(BRANCH_RETURN_CODES[213] + branchName)
	return 213
}
