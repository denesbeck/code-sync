package main

import (
	"encoding/json"
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
	Username string `json:"username"`
	Email    string `json:"email"`
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

func setConfig(key string, value string) int {
	Debug("Setting config: key=%s, value=%s", key, value)
	if initialized := IsInitialized(); !initialized {
		Fail(COMMON_RETURN_CODES[001])
		return 001
	}
	config, err := os.ReadFile(dirs.Config)
	if err != nil {
		Debug("Failed to read config file")
		MustSucceed(err, "operation failed")
	}

	var content Config
	if err = json.Unmarshal(config, &content); err != nil {
		Debug("Failed to unmarshal config")
		MustSucceed(err, "operation failed")
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
		MustSucceed(err, "operation failed")
	}

	if err = os.WriteFile(dirs.Config, jsonData, 0644); err != nil {
		Debug("Failed to write config file")
		MustSucceed(err, "operation failed")
	}
	Debug("Config updated successfully")
	Info(Capitalize(key) + " set to " + color.BlueString(value) + ".")
	return 603
}

func getConfig(key string) (returnCode int, conf Config) {
	Debug("Getting config: key=%s", key)
	if initialized := IsInitialized(); !initialized {
		Fail(COMMON_RETURN_CODES[001])
		return 001, Config{}
	}
	config := GetConfig()
	switch key {
	case "username":
		Debug("Username: %s", config.Username)
		Info(Capitalize(key) + ": " + color.BlueString(config.Username))
	case "email":
		Debug("Email: %s", config.Email)
		Info(Capitalize(key) + ": " + color.BlueString(config.Email))
	case "user":
		Debug("User: %s <%s>", config.Username, config.Email)
		Info(Capitalize(key) + ": " + color.BlueString(config.Username+" <"+config.Email+">"))
	}
	return 604, *config
}

func setDefaultBranch(branch string) int {
	if initialized := IsInitialized(); !initialized {
		Fail(COMMON_RETURN_CODES[001])
		return 001
	}
	err := SetBranch(branch, "default")
	if err != nil {
		if err.Error() == BRANCH_RETURN_CODES[215] {
			Info("Default branch already set to " + StyledBranch(branch) + ".")
			return 215
		}
		if err.Error() == BRANCH_RETURN_CODES[216] {
			Fail("Branch does not exist: " + StyledBranch(branch))
			return 216
		}
	}
	Debug("Default branch set successfully")
	Success("Default branch set to " + StyledBranch(branch) + ".")
	return 602
}

func getDefaultBranch() (returnCode int, defaultBranch string) {
	Debug("Getting default branch")
	if initialized := IsInitialized(); !initialized {
		Fail(COMMON_RETURN_CODES[001])
		return 001, ""
	}
	branch := GetDefaultBranchName()
	Debug("Default branch: %s", branch)
	Info("Default branch: " + StyledBranch(branch))
	return 601, branch
}
