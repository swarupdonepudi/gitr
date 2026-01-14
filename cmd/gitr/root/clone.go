package root

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"github.com/swarupdonepudi/gitr/internal/cli"
	"github.com/swarupdonepudi/gitr/pkg/clone"
	"github.com/swarupdonepudi/gitr/pkg/config"
	"github.com/swarupdonepudi/gitr/pkg/ui"
	"github.com/swarupdonepudi/gitr/pkg/url"
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
		return
	}
	inputUrl := args[0]

	// Special case: detect "cd <path>" pasted from clipboard that got split by shell
	// When user pastes "cd /path/to/repo", shell splits into args: ["cd", "/path/to/repo"]
	// This happens when another gitr clone finishes and overwrites clipboard with cd command
	if inputUrl == "cd" && len(args) >= 2 {
		reconstructedCmd := "cd " + args[1]
		ui.InvalidCloneURLCdCommand(reconstructedCmd, args[1])
		return
	}

	// Validate the URL before proceeding
	// This catches common mistakes like clipboard being overwritten by another gitr clone
	validation := url.ValidateCloneURL(inputUrl)
	if !validation.Valid {
		switch validation.Reason {
		case url.InvalidURLReasonCdCommand:
			ui.InvalidCloneURLCdCommand(validation.Input, validation.DetectedPath)
		case url.InvalidURLReasonLocalPath:
			ui.InvalidCloneURLLocalPath(validation.Input)
		case url.InvalidURLReasonEmpty:
			ui.CloneURLRequired()
		default:
			ui.InvalidCloneURLFormat(validation.Input)
		}
		return
	}

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
