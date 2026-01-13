package upgrade

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/swarupdonepudi/gitr/pkg/ui"
)

// Run executes the upgrade command
func Run(currentVersion string, checkOnly bool, force bool) {
	if currentVersion == "" {
		currentVersion = "dev"
	}

	ui.Info("Checking for updates...")

	latestVersion, err := GetLatestVersion()
	if err != nil {
		ui.Error("Update Check Failed",
			fmt.Sprintf("Failed to check for updates: %v", err),
			"You can manually download from https://github.com/swarupdonepudi/gitr/releases")
		return
	}

	needsUpgrade := CompareVersions(currentVersion, latestVersion)

	fmt.Println()
	if needsUpgrade {
		fmt.Printf("Current version: %s\n", ui.Dim(currentVersion))
		fmt.Printf("Latest version:  %s\n", latestVersion)
	} else {
		fmt.Printf("Current version: %s\n", currentVersion)
		fmt.Printf("Latest version:  %s\n", latestVersion)
	}

	if !needsUpgrade && !force {
		ui.Success(fmt.Sprintf("gitr is already up to date (%s)", currentVersion))
		return
	}

	if checkOnly {
		if needsUpgrade {
			ui.Warn("Update Available",
				fmt.Sprintf("A new version (%s) is available!\n\nRun %s to update.", latestVersion, ui.Cmd("gitr upgrade")))
		}
		return
	}

	if !needsUpgrade && force {
		ui.Info("Forcing upgrade...")
	}

	method := DetectUpgradeMethod()
	ui.Info(fmt.Sprintf("Upgrade method: %s", method.String()))

	var upgradeErr error
	switch method {
	case MethodHomebrew:
		upgradeErr = UpgradeViaHomebrew()
	case MethodDirectDownload:
		upgradeErr = UpgradeViaDirect(latestVersion)
	}

	if upgradeErr != nil {
		handleUpgradeError(upgradeErr, latestVersion)
		os.Exit(1)
	}

	ui.Success(fmt.Sprintf("Successfully upgraded to %s", latestVersion))

	if method == MethodDirectDownload {
		fmt.Println()
		ui.Info("Note: You may need to restart your terminal for changes to take effect.")
	}
}

func handleUpgradeError(err error, latestVersion string) {
	var permErr *PermissionError
	if errors.As(err, &permErr) {
		goos, goarch := GetPlatformInfo()
		downloadURL := BuildDownloadURL(latestVersion, goos, goarch)

		ui.Error("Permission Denied", permErr.Error(),
			"Try running with sudo: sudo gitr upgrade",
			fmt.Sprintf("Or download manually: curl -LO %s", downloadURL))
		return
	}

	goos, goarch := GetPlatformInfo()
	downloadURL := BuildDownloadURL(latestVersion, goos, goarch)

	var manualInstructions string
	if runtime.GOOS == "windows" {
		manualInstructions = fmt.Sprintf("Download from: %s", downloadURL)
	} else {
		manualInstructions = fmt.Sprintf("curl -LO %s && tar -xzf gitr_*.tar.gz && chmod +x gitr && sudo mv gitr /usr/local/bin/", downloadURL)
	}

	ui.Error("Upgrade Failed",
		fmt.Sprintf("Failed to upgrade: %v", err),
		"Visit https://github.com/swarupdonepudi/gitr/releases for manual download",
		manualInstructions)
}
