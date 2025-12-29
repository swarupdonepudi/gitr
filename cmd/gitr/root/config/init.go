package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/swarupdonepudi/gitr/pkg/config"
)

var Init = &cobra.Command{
	Use:   "init",
	Short: "initialize gitr config",
	Run:   initHandler,
}

func initHandler(cmd *cobra.Command, args []string) {
	if err := config.EnsureInitialConfig(); err != nil {
		log.Fatalf("%v", err)
	}
	log.Info("success!")
}
