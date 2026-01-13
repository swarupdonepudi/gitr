package upgrade

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/swarupdonepudi/gitr/pkg/ui"
)

// UpgradeViaHomebrew upgrades the CLI using Homebrew cask
func UpgradeViaHomebrew() error {
	fmt.Println()
	ui.Info("Updating Homebrew...")

	updateCmd := exec.Command("brew", "update")
	updateCmd.Stdout = os.Stdout
	updateCmd.Stderr = os.Stderr
	if err := updateCmd.Run(); err != nil {
		return fmt.Errorf("failed to update Homebrew: %w", err)
	}

	fmt.Println()
	ui.Info("Upgrading gitr...")

	upgradeCmd := exec.Command("brew", "upgrade", "--cask", "gitr")
	upgradeCmd.Stdout = os.Stdout
	upgradeCmd.Stderr = os.Stderr
	if err := upgradeCmd.Run(); err != nil {
		return fmt.Errorf("Homebrew upgrade failed: %w", err)
	}

	return nil
}
