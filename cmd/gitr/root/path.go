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
		log.Fatalf("clone url required as argument")
	}
	inputUrl := args[0]
	creDir, err := cmd.PersistentFlags().GetBool(string(cli.CreDir))
	cli.HandleFlagErr(err, cli.CreDir)

	cfg, err := config.NewGitrConfig()
	if err != nil {
		log.Fatalf("failed to get gitr config. err: %v", err)
	}
	repoLocation, err := clone.GetClonePath(cfg, inputUrl, creDir)
	if err != nil {
		log.Fatalf("failed to get clone path. err: %v", err)
	}
	fmt.Println(repoLocation)
	if cfg.CopyRepoPathCdCmdToClipboard {
		if err := clipboard.WriteAll(fmt.Sprintf("cd %s", repoLocation)); err != nil {
			log.Fatalf("err copying repo path to clipboard. %v\n", err)
		}
	}
}
