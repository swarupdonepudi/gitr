package ui

import (
	"fmt"
	"os"
	"strings"
)

// CloneSuccess displays a success message after cloning a repo
func CloneSuccess(clonePath string, clipboardEnabled bool) {
	// Try to make path more readable with ~ for home dir
	displayPath := clonePath
	if home, err := os.UserHomeDir(); err == nil {
		if strings.HasPrefix(clonePath, home) {
			displayPath = "~" + strings.TrimPrefix(clonePath, home)
		}
	}

	details := []string{
		Path(displayPath),
	}

	if clipboardEnabled {
		details = append(details,
			"",
			Dim("gitr organizes repos by their SCM path structure, so this repo"),
			Dim("was cloned outside your current directory."),
			"",
			fmt.Sprintf("A %s command has been copied to your clipboard.", Cmd("cd")),
			fmt.Sprintf("Press %s to navigate there instantly.", KeyCombo("⌘V")),
		)
	} else {
		details = append(details,
			"",
			fmt.Sprintf("Run %s to navigate to the repo.", Cmd("cd "+displayPath)),
		)
	}

	Success("Repository cloned successfully", details...)
}

// RepoAlreadyExists displays a message when the repo already exists
func RepoAlreadyExists(clonePath string) {
	displayPath := clonePath
	if home, err := os.UserHomeDir(); err == nil {
		if strings.HasPrefix(clonePath, home) {
			displayPath = "~" + strings.TrimPrefix(clonePath, home)
		}
	}

	fmt.Println()
	Info(fmt.Sprintf("Repository already exists at %s", Path(displayPath)))
}

// ConfigInitSuccess displays a success message after config init
func ConfigInitSuccess() {
	configPath := "~/.gitr.yaml"
	Success(
		"Configuration initialized",
		fmt.Sprintf("Config file created at %s", Path(configPath)),
		"",
		fmt.Sprintf("Run %s to customize your settings.", Cmd("gitr config edit")),
	)
}

// PathCopiedToClipboard displays a message when path is copied
func PathCopiedToClipboard(repoPath string) {
	displayPath := repoPath
	if home, err := os.UserHomeDir(); err == nil {
		if strings.HasPrefix(repoPath, home) {
			displayPath = "~" + strings.TrimPrefix(repoPath, home)
		}
	}

	fmt.Println()
	fmt.Printf("%s  %s\n",
		successIcon.Render(iconSuccess),
		successMessage.Render(Path(displayPath)))
	fmt.Println()
	fmt.Printf("   %s\n", Dim("A "+Cmd("cd")+" command has been copied to your clipboard."))
	fmt.Printf("   %s\n", Dim("Press "+KeyCombo("⌘V")+" to navigate there instantly."))
	fmt.Println()
}

// OpeningInBrowser displays a message when opening a URL
func OpeningInBrowser(url string) {
	Info(fmt.Sprintf("Opening %s in browser...", Dim(url)))
}

// WebInfo displays repository web info (for --dry mode)
func WebInfo(provider, hostname, remoteUrl, webUrl, repoPath, repoName, branch string) {
	fmt.Println()
	fmt.Printf("%s  %s\n",
		infoIcon.Render(iconInfo),
		infoMessage.Render("Repository Information"))
	fmt.Println()
	fmt.Printf("   %-12s %s\n", Dim("Provider:"), provider)
	fmt.Printf("   %-12s %s\n", Dim("Hostname:"), hostname)
	fmt.Printf("   %-12s %s\n", Dim("Remote URL:"), Path(remoteUrl))
	fmt.Printf("   %-12s %s\n", Dim("Web URL:"), Path(webUrl))
	fmt.Printf("   %-12s %s\n", Dim("Repo Path:"), repoPath)
	fmt.Printf("   %-12s %s\n", Dim("Repo Name:"), repoName)
	fmt.Printf("   %-12s %s\n", Dim("Branch:"), branch)
	fmt.Println()
}

// Version displays version information
func Version(version string) {
	fmt.Printf("%s %s\n", Dim("gitr version"), version)
}

// ConfigPath displays the config path
func ConfigPath(path string) {
	displayPath := path
	if home, err := os.UserHomeDir(); err == nil {
		if strings.HasPrefix(path, home) {
			displayPath = "~" + strings.TrimPrefix(path, home)
		}
	}
	fmt.Printf("%s  %s\n", infoIcon.Render(iconInfo), Path(displayPath))
}

// ClonePath displays a clone path (for gitr path command)
func ClonePath(path string, clipboardEnabled bool) {
	fmt.Println(path) // Keep raw path for scripts that might parse this

	if clipboardEnabled {
		// Only show the hint if clipboard was used
		fmt.Printf("\n   %s\n", Dim("cd command copied to clipboard. Press "+KeyCombo("⌘V")+" to navigate."))
	}
}

// Cloning displays a message when starting to clone a repository
func Cloning(repoUrl string) {
	fmt.Printf("\n%s  %s %s\n",
		infoIcon.Render("↓"),
		Dim("Cloning"),
		Path(repoUrl))
}
