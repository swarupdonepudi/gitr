package root

import (
	"github.com/spf13/cobra"
	"github.com/swarupdonepudi/gitr/pkg/ui"
)

var VersionLabel = "dev"

var Version = &cobra.Command{
	Use:     "version",
	Short:   "check the version of the cli",
	Aliases: []string{"v"},
	Run:     versionHandler,
}

func versionHandler(cmd *cobra.Command, args []string) {
	ui.Version(VersionLabel)
}
