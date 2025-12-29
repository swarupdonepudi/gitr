package root

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"github.com/swarupdonepudi/gitr/internal/cli"
	"github.com/swarupdonepudi/gitr/pkg/clone"
	"github.com/swarupdonepudi/gitr/pkg/config"
	"github.com/swarupdonepudi/gitr/pkg/ui"
)

var Path = &cobra.Command{
	Use:   "path",
	Short: "prints the path to which the repo is cloned/will be cloned",
	Run:   pathHandler,
}

func init() {
	Path.PersistentFlags().BoolP(string(cli.CreDir), "", false, "cre folders to mimic repo path on scm")
}

func pathHandler(cmd *cobra.Command, args []string) {
	if len(args) <= 0 {
		ui.CloneURLRequired()
	}
	inputUrl := args[0]
	creDir, err := cmd.PersistentFlags().GetBool(string(cli.CreDir))
	cli.HandleFlagErr(err, cli.CreDir)

	cfg, err := config.NewGitrConfig()
	if err != nil {
		ui.ConfigError(err)
	}
	repoLocation, err := clone.GetClonePath(cfg, inputUrl, creDir)
	if err != nil {
		ui.GenericError("Failed to Get Path", "Could not determine clone path for the repository", err)
	}

	clipboardEnabled := cfg.CopyRepoPathCdCmdToClipboard
	if clipboardEnabled {
		if err := clipboard.WriteAll(fmt.Sprintf("cd %s", repoLocation)); err != nil {
			ui.ClipboardError(err)
			clipboardEnabled = false
		}
	}

	ui.ClonePath(repoLocation, clipboardEnabled)
}
