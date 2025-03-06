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
		runBranchCommand()
	},
}

var currentCmd = &cobra.Command{
	Use:     "current",
	Short:   "Get current branch",
	Example: "csync branch current",
	Args:    cobra.NoArgs,
	Run: func(_ *cobra.Command, args []string) {
		runCurrentCommand()
	},
}

var defaultCmd = &cobra.Command{
	Use:     "default",
	Short:   "Get default branch",
	Example: "csync branch default",
	Args:    cobra.NoArgs,
	Run: func(_ *cobra.Command, args []string) {
		runDefaultCommand()
	},
}

var newCmd = &cobra.Command{
	Use:     "new",
	Short:   "Create a new branch",
	Example: "csync new <branch-name> --from-commit <commit-id> --from-branch <branch-name>",
	Args:    cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		runNewCommand(args[0], FromCommit, FromBranch)
	},
}

var dropCmd = &cobra.Command{
	Use:     "drop",
	Short:   "Delete a branch",
	Example: "csync drop <branch-name>",
	Args:    cobra.MinimumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {
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
		runSwitchCommand(args[0])
	},
}

func runBranchCommand() {
	initialized := IsInitialized()
	if !initialized {
		color.Red("CSync not initialized")
		return
	}

	branches, err := os.ReadDir("./.csync/branches")
	if err != nil {
		color.Red("No branches found")
		return
	}

	currentBranchName := GetCurrentBranchName()
	defaultBranchName := GetDefaultBranchName()

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
}

func runCurrentCommand() {
	initialized := IsInitialized()
	if !initialized {
		color.Red("CSync not initialized")
		return
	}

	currentBranchName := GetCurrentBranchName()
	fmt.Println(currentBranchName)
}

func runDefaultCommand() {
	initialized := IsInitialized()
	if !initialized {
		color.Red("CSync not initialized")
		return
	}

	defaultBranchName := GetDefaultBranchName()
	fmt.Println(defaultBranchName)
}

func runNewCommand(branchName string, fromCommit string, fromBranch string) {
	initialized := IsInitialized()
	if !initialized {
		color.Red("CSync not initialized")
		return
	}

	if fromCommit != "" && fromBranch != "" {
		color.Red("Cannot create branch from both commit and branch")
		return
	}

	var srcBranch string
	if fromBranch != "" {
		srcBranch = fromBranch
		branches := ListBranches()
		if !slices.Contains(branches, srcBranch) {
			color.Red("Source branch does not exist")
			return
		}
	} else {
		srcBranch = GetCurrentBranchName()
	}

	if fromCommit != "" {
		err := CopyCommitsToBranch(fromCommit, branchName)
		if err != nil {
			color.Red(err.Error())
			return
		}
	} else {
		if err := os.Mkdir(dirs.Branches+branchName, 0755); err != nil {
			color.Red("Branch already exists")
			return
		}

		CopyFile(dirs.Branches+srcBranch+"/commits.json", dirs.Branches+branchName+"/commits.json")

	}
	color.Green("Branch created successfully")
	runSwitchCommand(branchName)
}

func runDropCommand(branchName string) {
	initialized := IsInitialized()
	if !initialized {
		color.Red("CSync not initialized")
		return
	}

	branches := ListBranches()
	if !slices.Contains(branches, branchName) {
		color.Red("Branch does not exist")
		return
	}

	if currentBranchName := GetCurrentBranchName(); currentBranchName == branchName {
		color.Red("Cannot delete current branch")
		return
	}

	if defaultBranchName := GetDefaultBranchName(); defaultBranchName == branchName {
		color.Red("Cannot delete default branch")
		return
	}

	if err := os.RemoveAll(dirs.Branches + branchName); err != nil {
		color.Red("Branch does not exist")
		return
	}
	color.Green("Branch deleted successfully")
}

func runSwitchCommand(branchName string) {
	initialized := IsInitialized()
	if !initialized {
		color.Red("CSync not initialized")
		return
	}

	currentBranch := GetCurrentBranchName()
	if currentBranch == branchName {
		color.Red("You are already on `" + branchName + "` branch")
		return
	}

	branches := ListBranches()
	if !slices.Contains(branches, branchName) {
		color.Red("Branch does not exist")
		return
	}

	commitId := GetLastCommitByBranch(branchName)
	if commitId != "" {
		fileList := GetFileListContent(commitId)
		for _, file := range *fileList {
			_, fileName := ParsePath(file.Path)
			CopyFile(dirs.Commits+file.CommitId+"/"+file.Id+"/"+fileName, "./"+file.Path)
		}
	}
	SetBranch(branchName, "current")
	color.Cyan("Current branch: " + branchName)
}
