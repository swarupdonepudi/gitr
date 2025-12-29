package root

import (
	"github.com/spf13/cobra"
	"github.com/swarupdonepudi/gitr/cmd/gitr/root/config"
)

var Config = &cobra.Command{
	Use:   "config",
	Short: "see gitr config",
}

func init() {
	Config.AddCommand(config.Init, config.Show, config.Edit)
}
