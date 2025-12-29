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

var Clone = &cobra.Command{
	Use:   "clone",
	Short: "Clone repo to organized, deterministic path (~/scm/{host}/{owner}/{repo})",
	Run:   cloneHandler,
}

func init() {
	Clone.PersistentFlags().BoolP(string(cli.CreDir), "", false, "create full directory hierarchy matching SCM structure")
	Clone.PersistentFlags().StringP(string(cli.Token), "", "", "HTTPS personal access token for authentication")
}

func cloneHandler(cmd *cobra.Command, args []string) {
	if len(args) <= 0 {
		ui.CloneURLRequired()
	}
	inputUrl := args[0]
	dry, err := cmd.InheritedFlags().GetBool(string(cli.Dry))
	cli.HandleFlagErr(err, cli.Dry)
	creDir, err := cmd.PersistentFlags().GetBool(string(cli.CreDir))
	cli.HandleFlagErr(err, cli.CreDir)
	token, err := cmd.PersistentFlags().GetString(string(cli.Token))
	cli.HandleFlagErr(err, cli.Token)
	cfg, err := config.NewGitrConfig()
	if err != nil {
		ui.ConfigError(err)
	}
	clonePath, err := clone.Clone(cfg, inputUrl, token, creDir, dry)
	if err != nil {
		ui.FailedToClone(err)
	}

	// Handle clipboard copy
	clipboardEnabled := cfg.CopyRepoPathCdCmdToClipboard
	if clipboardEnabled {
		if err := clipboard.WriteAll(fmt.Sprintf("cd %s", clonePath)); err != nil {
			ui.ClipboardError(err)
			clipboardEnabled = false
		}
	}

	ui.CloneSuccess(clonePath, clipboardEnabled)
}
