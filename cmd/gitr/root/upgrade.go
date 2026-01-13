package root

import (
	"github.com/spf13/cobra"
	"github.com/swarupdonepudi/gitr/internal/upgrade"
)

var Upgrade = &cobra.Command{
	Use:   "upgrade",
	Short: "upgrade gitr to the latest version",
	Long: `Upgrade gitr to the latest available version.

On macOS, if gitr was installed via Homebrew, this command uses 'brew upgrade --cask'.
On all other platforms, it downloads the latest binary directly from GitHub releases.

Examples:
  # Upgrade to the latest version
  gitr upgrade

  # Check for updates without installing
  gitr upgrade --check

  # Force upgrade even if already on latest version
  gitr upgrade --force`,
	Run: upgradeHandler,
}

func init() {
	Upgrade.Flags().BoolP("check", "c", false, "check for updates without installing")
	Upgrade.Flags().BoolP("force", "f", false, "force upgrade even if already on latest version")
}

func upgradeHandler(cmd *cobra.Command, args []string) {
	checkOnly, _ := cmd.Flags().GetBool("check")
	force, _ := cmd.Flags().GetBool("force")
	upgrade.Run(VersionLabel, checkOnly, force)
}
