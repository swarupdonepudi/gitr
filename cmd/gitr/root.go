package gitr

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/swarupdonepudi/gitr/cmd/gitr/root"
	"github.com/swarupdonepudi/gitr/internal/cli"
	"github.com/swarupdonepudi/gitr/pkg/config"
	"os"
	"runtime"
)

var debug bool

const HomebrewAppleSiliconBinPath = "/opt/homebrew/bin"

var rootCmd = &cobra.Command{
	Use:   "gitr",
	Short: "git rapid - the missing link b/w git cli & scm providers",
	Long:  "save time(a ton) by opening git repos on web browser right from the command line",
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, string(cli.Debug), false, "set log level to debug")
	rootCmd.PersistentFlags().BoolP(string(cli.Dry), "", false, "dry run")
	rootCmd.AddCommand(
		root.Version,
		root.Config,
		root.Clone,
		root.Path,
		root.BranchesCmd,
		root.CommitsCmd,
		root.IssuesCmd,
		root.PipelinesCmd,
		root.PrsCmd,
		root.ReleasesCmd,
		root.RemCmd,
		root.TagsCmd,
		root.WebCmd,
		root.WebUrlCmd,
	)
	cobra.OnInitialize(func() {
		if debug {
			log.SetLevel(log.DebugLevel)
			log.Debug("running in debug mode")
		}
		if runtime.GOARCH == "arm64" {
			pathEnvVal := os.Getenv("PATH")
			if err := os.Setenv("PATH", fmt.Sprintf("%s:%s", pathEnvVal, HomebrewAppleSiliconBinPath)); err != nil {
				log.Fatalf("failed to set PATH env. err: %v", err)
			}
		}
	})
	if err := config.EnsureInitialConfig(); err != nil {
		log.Fatalf("failed to initialize config. err %v", err)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("failed to run command. err %v", err)
	}
}
