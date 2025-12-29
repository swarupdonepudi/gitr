package config

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/leftbin/go-util/pkg/file"
	"github.com/leftbin/go-util/pkg/shell"
	"github.com/spf13/cobra"
	"github.com/swarupdonepudi/gitr/pkg/ui"
)

var Edit = &cobra.Command{
	Use:   "edit",
	Short: "edit gitr config",
	Run:   editHandler,
}

func editHandler(cmd *cobra.Command, args []string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		ui.GenericError("System Error", "Failed to determine home directory", err)
	}
	gitrConfigPath := filepath.Join(homeDir, ".gitr.yaml")
	if !file.IsFileExists(gitrConfigPath) {
		ui.Error(
			"Configuration Not Found",
			ui.Path(gitrConfigPath)+" does not exist.",
			"Run "+ui.Cmd("gitr config init")+" to create one",
		)
	}
	if err := shell.RunCmd(exec.Command("code", gitrConfigPath)); err != nil {
		ui.FailedToOpenEditor("code", err)
	}
}
