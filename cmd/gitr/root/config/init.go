package config

import (
	"github.com/spf13/cobra"
	"github.com/swarupdonepudi/gitr/pkg/config"
	"github.com/swarupdonepudi/gitr/pkg/ui"
)

var Init = &cobra.Command{
	Use:   "init",
	Short: "initialize gitr config",
	Run:   initHandler,
}

func initHandler(cmd *cobra.Command, args []string) {
	if err := config.EnsureInitialConfig(); err != nil {
		ui.GenericError("Configuration Error", "Failed to initialize configuration", err)
	}
	ui.ConfigInitSuccess()
}
