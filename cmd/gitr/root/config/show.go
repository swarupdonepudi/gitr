package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/swarupdonepudi/gitr/pkg/config"
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
		log.Fatalf("failed to get gitr config. err: %v", err)
	}
	d, err := yaml.Marshal(&cfg)
	fmt.Printf("")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n%s\n", string(d))
}
