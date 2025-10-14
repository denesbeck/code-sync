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
		Debug("Starting init command")
		runInitCommand()
	},
}

func runInitCommand() {
	if _, err := os.Stat(dirs.Root); !os.IsNotExist(err) {
		Debug("CSync already initialized")
		color.Red("CSync already initialized")
		return
	}

	Debug("Creating staging directories")
	if err := os.MkdirAll(dirs.StagingAdded, os.ModePerm); err != nil {
		Debug("Failed to create staging/added directory")
		log.Fatal(err)
	}
	if err := os.MkdirAll(dirs.StagingModified, os.ModePerm); err != nil {
		Debug("Failed to create staging/modified directory")
		log.Fatal(err)
	}
	if err := os.MkdirAll(dirs.StagingRemoved, os.ModePerm); err != nil {
		Debug("Failed to create staging/removed directory")
		log.Fatal(err)
	}

	Debug("Creating staging logs file")
	f, err := os.Create(dirs.StagingLogs)
	if err != nil {
		Debug("Failed to create staging logs file")
		log.Fatal(err)
	}
	_, err = f.WriteString("[]")
	if err != nil {
		Debug("Failed to write initial staging logs")
		log.Fatal(err)
	}
	f.Close()

	Debug("Creating commits directory")
	if err := os.MkdirAll(dirs.Commits, os.ModePerm); err != nil {
		Debug("Failed to create commits directory")
		log.Fatal(err)
	}

	Debug("Creating default branch directory")
	if err := os.MkdirAll(dirs.DefaultBranch, os.ModePerm); err != nil {
		Debug("Failed to create default branch directory")
		log.Fatal(err)
	}
	Debug("Creating default branch commits file")
	f, err = os.Create(dirs.DefaultBranchCommits)
	if err != nil {
		Debug("Failed to create default branch commits file")
		log.Fatal(err)
	}
	_, err = f.WriteString("[]")
	if err != nil {
		Debug("Failed to write initial default branch commits")
		log.Fatal(err)
	}
	f.Close()

	Debug("Creating branches metadata")
	CreateBranchesMetadata()

	Debug("Creating config file")
	f, err = os.Create(dirs.Config)
	if err != nil {
		Debug("Failed to create config file")
		log.Fatal(err)
	}
	_, err = f.WriteString("{ \"Username\": \"\", \"Email\": \"\" }")
	if err != nil {
		Debug("Failed to write initial config")
		log.Fatal(err)
	}
	f.Close()

	Debug("CSync initialized successfully")
	color.Green("CSync initialized successfully")
	fmt.Println()
	color.Cyan("Use `csync config set username <your-username>` to set your username and `csync config set email <your-email>` to set your email.")
}
