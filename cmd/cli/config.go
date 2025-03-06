package cli

import (
	"encoding/json"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.AddCommand(setCmd)
	configCmd.AddCommand(getCmd)

	setCmd.AddCommand(setDefaultBranchCmd)
	setCmd.AddCommand(setUsernameCmd)
	setCmd.AddCommand(setEmailCmd)

	getCmd.AddCommand(getDefaultBranchCmd)
	getCmd.AddCommand(getUsernameCmd)
	getCmd.AddCommand(getEmailCmd)
	getCmd.AddCommand(getUserCmd)
}

type Config struct {
	Username string
	Email    string
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set config values",
	Args:  cobra.ExactArgs(1),
}

var setDefaultBranchCmd = &cobra.Command{
	Use:     "default-branch",
	Short:   "Set default branch",
	Example: "csync config set default-branch <branch-name>",
	Args:    cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		setDefaultBranch(args[0])
	},
}

var setUsernameCmd = &cobra.Command{
	Use:     "username",
	Short:   "Set username",
	Example: "csync config set username <username>",
	Args:    cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		setConfig("username", args[0])
	},
}

var setEmailCmd = &cobra.Command{
	Use:     "email",
	Short:   "Set email",
	Example: "csync config set email <email>",
	Args:    cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		setConfig("email", args[0])
	},
}

var getCmd = &cobra.Command{
	Use:     "get",
	Short:   "Get config values",
	Example: "csync config get username <username>",
	Args:    cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		getConfig(args[0])
	},
}

var getDefaultBranchCmd = &cobra.Command{
	Use:     "default-branch",
	Short:   "Get default branch",
	Example: "csync config get default-branch",
	Args:    cobra.ExactArgs(0),
	Run: func(_ *cobra.Command, args []string) {
		getDefaultBranch()
	},
}

var getUsernameCmd = &cobra.Command{
	Use:     "username",
	Short:   "Get username",
	Example: "csync config get username",
	Args:    cobra.ExactArgs(0),
	Run: func(_ *cobra.Command, args []string) {
		getConfig("username")
	},
}

var getEmailCmd = &cobra.Command{
	Use:     "email",
	Short:   "Get email",
	Example: "csync config get email <email>",
	Args:    cobra.ExactArgs(0),
	Run: func(_ *cobra.Command, args []string) {
		getConfig("email")
	},
}

var getUserCmd = &cobra.Command{
	Use:     "user",
	Short:   "Get username and email",
	Example: "csync config get user",
	Args:    cobra.ExactArgs(0),
	Run: func(_ *cobra.Command, args []string) {
		getConfig("user")
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Config manager",
	Args:  cobra.ExactArgs(1),
}

func setConfig(key string, value string) {
	if initialized := IsInitialized(); !initialized {
		color.Red("CSync not initialized")
	}
	config, err := os.ReadFile(dirs.Config)
	if err != nil {
		log.Fatal(err)
	}

	var content Config
	if err = json.Unmarshal(config, &content); err != nil {
		log.Fatal(err)
	}

	switch key {
	case "username":
		content.Username = value
	case "email":
		content.Email = value
	}

	jsonData, err := json.Marshal(content)
	if err != nil {
		log.Fatal(err)
	}

	if err = os.WriteFile(dirs.Config, jsonData, 0644); err != nil {
		log.Fatal(err)
	}
}

func getConfig(key string) {
	if initialized := IsInitialized(); !initialized {
		color.Red("CSync not initialized")
	}
	config := GetConfig()
	switch key {
	case "username":
		color.Cyan(config.Username)
	case "email":
		color.Cyan(config.Email)
	case "user":
		color.Cyan(config.Username + " <" + config.Email + ">")
	}
}

func setDefaultBranch(branch string) {
	if initialized := IsInitialized(); !initialized {
		color.Red("CSync not initialized")
	}
	SetBranch(branch, "default")
	color.Green("Default branch set to " + branch)
}

func getDefaultBranch() {
	if initialized := IsInitialized(); !initialized {
		color.Red("CSync not initialized")
	}
	branch := GetDefaultBranchName()
	color.Cyan(branch)
}
