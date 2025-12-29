package gitr

import (
	"fmt"
	"os"
	"runtime"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/swarupdonepudi/gitr/cmd/gitr/root"
	"github.com/swarupdonepudi/gitr/internal/cli"
	"github.com/swarupdonepudi/gitr/pkg/config"
	"github.com/swarupdonepudi/gitr/pkg/ui"
)

var debug bool

const HomebrewAppleSiliconBinPath = "/opt/homebrew/bin"

var rootCmd = &cobra.Command{
	Use:   "gitr",
	Short: "Clone to organized paths. Open PRs, pipelines, branches instantly.",
	Long: `gitr - Your missing git productivity tool

Clone repos to organized, deterministic paths and navigate to any web page 
(PRs, pipelines, issues, branches) instantly from your terminal.

No more scattered repos. No more browser tab hunting. One CLI, zero friction.

Examples:
  gitr clone https://github.com/owner/repo    # → ~/scm/github.com/owner/repo
  gitr prs                                    # Open PRs in browser
  gitr pipe                                   # Open pipelines/actions
  gitr web                                    # Open repo homepage

Learn more: https://swarupdonepudi.github.io/gitr`,
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
				ui.GenericError("Environment Error", "Failed to configure PATH for Apple Silicon", err)
			}
		}
	})
	if err := config.EnsureInitialConfig(); err != nil {
		ui.GenericError("Configuration Error", "Failed to initialize gitr configuration", err)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// Cobra already prints errors, so we just exit
		os.Exit(1)
	}
}
