package root

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/swarupdonepudi/gitr/pkg/web"
)

// WebUrlCmd prints the browser URL for the given file inside the current Git repo.
var WebUrlCmd = &cobra.Command{
	Use:   "web-url [file-name]",
	Short: "prints the web url of a file in the repo",
	Args:  cobra.ExactArgs(1),
	Run:   webUrlCmdHandler,
}

func webUrlCmdHandler(cmd *cobra.Command, args []string) {
	url, err := web.FileURLFromPwd(args[0])
	if err != nil {
		log.Fatalf("error: %v\n", err)
		return
	}
	println(url)
}
