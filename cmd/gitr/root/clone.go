package root

import (
	"fmt"
	"github.com/atotto/clipboard"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/swarupdonepudi/gitr/internal/cli"
	"github.com/swarupdonepudi/gitr/pkg/clone"
	"github.com/swarupdonepudi/gitr/pkg/config"
)

var Clone = &cobra.Command{
	Use:   "clone",
	Short: "clones repo to mimic folder structure to the scm repo hierarchy",
	Run:   cloneHandler,
}

func init() {
	Clone.PersistentFlags().BoolP(string(cli.CreDir), "", false, "cre folders to mimic repo path on scm")
	Clone.PersistentFlags().StringP(string(cli.Token), "", "", "https personal access token help: https://docs.gitlab.com/12.10/ee/user/profile/personal_access_tokens.html")
}

func cloneHandler(cmd *cobra.Command, args []string) {
	if len(args) <= 0 {
		log.Fatalf("clone url required as argument")
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
		log.Fatalf("failed to get gitr config. err: %v", err)
	}
	clonePath, err := clone.Clone(cfg, inputUrl, token, creDir, dry)
	if err != nil {
		log.Fatalf("failed to clone repo. err: %v", err)
	}
	log.Infof("repo path: %s", clonePath)
	if !cfg.CopyRepoPathCdCmdToClipboard {
		log.Infof("run this command to navigate to repo path: cd %s", clonePath)
	}
	if err := clipboard.WriteAll(fmt.Sprintf("cd %s", clonePath)); err != nil {
		log.Fatalf("failed to copying repo path to clipboard. err: %v", err)
	}
	log.Info("command to navigate to repo path has been added to clipboard. run cmd+v to paste the command")
}
