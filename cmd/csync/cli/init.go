package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:     "init",
	Short:   "Initialize the CSync version control system",
	Example: "csync init",
	Args:    cobra.NoArgs,
	Run: func(_ *cobra.Command, _ []string) {
		runInitCommand()
	},
}

func runInitCommand() {
	if _, err := os.Stat(namespace + ".csync"); !os.IsNotExist(err) {
		color.Red("CSync already initialized")
		return
	}

	if err := os.MkdirAll(dirs.StagingAdded, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(dirs.StagingModified, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(dirs.StagingRemoved, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(dirs.StagingLogs)
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.WriteString("[]")
	if err != nil {
		log.Fatal(err)
	}
	f.Close()

	if err := os.MkdirAll(dirs.Commits, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(dirs.DefaultBranch, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	f, err = os.Create(dirs.DefaultBranchCommits)
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.WriteString("[]")
	if err != nil {
		log.Fatal(err)
	}
	f.Close()

	CreateBranchesMetadata()

	f, err = os.Create(dirs.Config)
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.WriteString("{ \"Username\": \"\", \"Email\": \"\" }")
	if err != nil {
		log.Fatal(err)
	}
	f.Close()

	color.Green("CSync initialized successfully")
	fmt.Println()
	color.Cyan("Use `csync config set username <your-username>` to set your username and `csync config set email <your-email>` to set your email.")
}
