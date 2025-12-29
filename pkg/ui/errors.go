package ui

import (
	"fmt"
	"os"
)

// Predefined error messages with helpful context

// NotInGitRepo displays an error when the current directory is not a git repo
func NotInGitRepo() {
	Error(
		"Not a Git Repository",
		"The current directory is not inside a git repository.",
		"Navigate to a git repository, or run "+Cmd("git init")+" to create one",
	)
}

// ConfigNotFound displays an error when gitr config is missing
func ConfigNotFound() {
	Error(
		"Configuration Not Found",
		"Could not find or load gitr configuration.",
		"Run "+Cmd("gitr config init")+" to create a configuration file",
	)
}

// ConfigError displays a config-related error
func ConfigError(err error) {
	Error(
		"Configuration Error",
		fmt.Sprintf("Failed to load gitr configuration: %v", err),
		"Run "+Cmd("gitr config init")+" to reset your configuration",
	)
}

// UnknownSCMHost displays an error for unrecognized SCM hosts
func UnknownSCMHost(hostname string) {
	Error(
		"Unknown SCM Host",
		fmt.Sprintf("The hostname %s is not configured in gitr.", Path(hostname)),
		"Add it to your config with "+Cmd("gitr config edit"),
	)
}

// CloneURLRequired displays an error when clone URL is missing
func CloneURLRequired() {
	Error(
		"Clone URL Required",
		"Please provide a repository URL to clone.",
		"Usage: "+Cmd("gitr clone <url>"),
		"Example: "+Cmd("gitr clone https://github.com/owner/repo"),
	)
}

// NoRemotesFound displays an error when git repo has no remotes
func NoRemotesFound() {
	Error(
		"No Remotes Found",
		"This git repository has no remote configured.",
		"Add a remote with "+Cmd("git remote add origin <url>"),
	)
}

// FailedToGetBranch displays an error when branch detection fails
func FailedToGetBranch(err error) {
	Error(
		"Failed to Get Branch",
		fmt.Sprintf("Could not determine the current git branch: %v", err),
		"Make sure you're in a valid git repository with at least one commit",
	)
}

// FailedToClone displays an error when cloning fails
func FailedToClone(err error) {
	Error(
		"Clone Failed",
		fmt.Sprintf("Failed to clone the repository: %v", err),
		"Check your network connection and repository URL",
		"For private repos, ensure you have the correct access token",
	)
}

// ClipboardError displays an error when clipboard operations fail
func ClipboardError(err error) {
	// This is non-fatal, just warn
	Warn(
		"Clipboard Unavailable",
		fmt.Sprintf("Could not copy to clipboard: %v", err),
	)
}

// FlagParseError displays an error for flag parsing issues
func FlagParseError(flag string, err error) {
	Error(
		"Invalid Flag",
		fmt.Sprintf("Failed to parse the %s flag: %v", Cmd("--"+flag), err),
	)
}

// FileNotFound displays an error when a required file is missing
func FileNotFound(path string) {
	Error(
		"File Not Found",
		fmt.Sprintf("The file %s does not exist.", Path(path)),
	)
}

// FailedToOpenEditor displays an error when editor fails to open
func FailedToOpenEditor(editor string, err error) {
	Error(
		"Failed to Open Editor",
		fmt.Sprintf("Could not open %s: %v", Cmd(editor), err),
		"Make sure the editor is installed and in your PATH",
	)
}

// GenericError displays a generic error with custom message
func GenericError(title, message string, err error) {
	if err != nil {
		Error(title, fmt.Sprintf("%s: %v", message, err))
	} else {
		Error(title, message)
	}
}

// Fatal prints an error and exits - direct replacement for log.Fatal
func Fatal(message string) {
	Error("Error", message)
}

// Fatalf prints a formatted error and exits - direct replacement for log.Fatalf
func Fatalf(format string, args ...interface{}) {
	Error("Error", fmt.Sprintf(format, args...))
}

// FatalErr prints an error from an error object and exits
func FatalErr(err error) {
	if err != nil {
		Error("Error", err.Error())
	}
	os.Exit(1)
}
