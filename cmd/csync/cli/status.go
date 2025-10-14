package cli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:     "status",
	Short:   "This command lists the files that are staged for commit",
	Example: "csync status",
	RunE: func(_ *cobra.Command, args []string) error {
		return runStatusCommand()
	},
}

var (
	add = color.New(color.FgGreen, color.Bold).SprintFunc()
	mod = color.New(color.FgBlue, color.Bold).SprintFunc()
	rem = color.New(color.FgRed, color.Bold).SprintFunc()
)

func runStatusCommand() error {
	if initialized := IsInitialized(); !initialized {
		color.Red("CSync not initialized")
		return nil
	}
	content := GetStagingLogsContent()
	if len(content) == 0 {
		color.Cyan("No files staged for commit")
	} else {

		content = SortByOperationAndPath(content)

		color.Cyan("Files staged for commit:")
		for _, record := range content {
			switch record.Op {
			case "ADD":
				fmt.Println("  " + add(record.Op) + "    " + record.Path)
			case "MOD":
				fmt.Println("  " + mod(record.Op) + "    " + record.Path)
			case "REM":
				fmt.Println("  " + rem(record.Op) + "    " + record.Path)
			default:
				fmt.Println("  " + record.Op + "    " + record.Path)
			}
		}
	}
	return nil
}
