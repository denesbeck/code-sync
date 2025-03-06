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
		Debug("Setting default branch: %s", args[0])
		setDefaultBranch(args[0])
	},
}

var setUsernameCmd = &cobra.Command{
	Use:     "username",
	Short:   "Set username",
	Example: "csync config set username <username>",
	Args:    cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		Debug("Setting username: %s", args[0])
		setConfig("username", args[0])
	},
}

var setEmailCmd = &cobra.Command{
	Use:     "email",
	Short:   "Set email",
	Example: "csync config set email <email>",
	Args:    cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		Debug("Setting email: %s", args[0])
		setConfig("email", args[0])
	},
}

var getCmd = &cobra.Command{
	Use:     "get",
	Short:   "Get config values",
	Example: "csync config get username <username>",
	Args:    cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		Debug("Getting config value: %s", args[0])
		getConfig(args[0])
	},
}

var getDefaultBranchCmd = &cobra.Command{
	Use:     "default-branch",
	Short:   "Get default branch",
	Example: "csync config get default-branch",
	Args:    cobra.ExactArgs(0),
	Run: func(_ *cobra.Command, args []string) {
		Debug("Getting default branch")
		getDefaultBranch()
	},
}

var getUsernameCmd = &cobra.Command{
	Use:     "username",
	Short:   "Get username",
	Example: "csync config get username",
	Args:    cobra.ExactArgs(0),
	Run: func(_ *cobra.Command, args []string) {
		Debug("Getting username")
		getConfig("username")
	},
}

var getEmailCmd = &cobra.Command{
	Use:     "email",
	Short:   "Get email",
	Example: "csync config get email <email>",
	Args:    cobra.ExactArgs(0),
	Run: func(_ *cobra.Command, args []string) {
		Debug("Getting email")
		getConfig("email")
	},
}

var getUserCmd = &cobra.Command{
	Use:     "user",
	Short:   "Get username and email",
	Example: "csync config get user",
	Args:    cobra.ExactArgs(0),
	Run: func(_ *cobra.Command, args []string) {
		Debug("Getting user info")
		getConfig("user")
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Config manager",
	Args:  cobra.ExactArgs(1),
}

func setConfig(key string, value string) {
	Debug("Setting config: key=%s, value=%s", key, value)
	if initialized := IsInitialized(); !initialized {
		Debug("CSync not initialized")
		color.Red("CSync not initialized")
	}
	config, err := os.ReadFile(dirs.Config)
	if err != nil {
		Debug("Failed to read config file")
		log.Fatal(err)
	}

	var content Config
	if err = json.Unmarshal(config, &content); err != nil {
		Debug("Failed to unmarshal config")
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
		Debug("Failed to marshal config")
		log.Fatal(err)
	}

	if err = os.WriteFile(dirs.Config, jsonData, 0644); err != nil {
		Debug("Failed to write config file")
		log.Fatal(err)
	}
	Debug("Config updated successfully")
}

func getConfig(key string) {
	Debug("Getting config: key=%s", key)
	if initialized := IsInitialized(); !initialized {
		Debug("CSync not initialized")
		color.Red("CSync not initialized")
	}
	config := GetConfig()
	switch key {
	case "username":
		Debug("Username: %s", config.Username)
		color.Cyan(config.Username)
	case "email":
		Debug("Email: %s", config.Email)
		color.Cyan(config.Email)
	case "user":
		Debug("User: %s <%s>", config.Username, config.Email)
		color.Cyan(config.Username + " <" + config.Email + ">")
	}
}

func setDefaultBranch(branch string) {
	Debug("Setting default branch: %s", branch)
	if initialized := IsInitialized(); !initialized {
		Debug("CSync not initialized")
		color.Red("CSync not initialized")
	}
	SetBranch(branch, "default")
	Debug("Default branch set successfully")
	color.Green("Default branch set to " + branch)
}

func getDefaultBranch() {
	Debug("Getting default branch")
	if initialized := IsInitialized(); !initialized {
		Debug("CSync not initialized")
		color.Red("CSync not initialized")
	}
	branch := GetDefaultBranchName()
	Debug("Default branch: %s", branch)
	color.Cyan(branch)
}
