package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/swarupdonepudi/gitr/pkg/config"
	"github.com/swarupdonepudi/gitr/pkg/ui"
	"gopkg.in/yaml.v3"
)

var Show = &cobra.Command{
	Use:   "show",
	Short: "show gitr config",
	Run:   showHandler,
}

func showHandler(cmd *cobra.Command, args []string) {
	cfg, err := config.NewGitrConfig()
	if err != nil {
		ui.ConfigError(err)
	}
	d, err := yaml.Marshal(&cfg)
	if err != nil {
		ui.GenericError("Configuration Error", "Failed to serialize configuration", err)
	}
	fmt.Printf("\n%s\n", string(d))
}
