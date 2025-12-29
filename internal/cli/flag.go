package cli

import "github.com/swarupdonepudi/gitr/pkg/ui"

type Flag string

const (
	Dry    Flag = "dry"
	CreDir Flag = "create-dir"
	Debug  Flag = "debug"
	Token  Flag = "token"
)

func HandleFlagErr(err error, flag Flag) {
	if err != nil {
		ui.FlagParseError(string(flag), err)
	}
}
