package cmd

import (
	"io"

	"samhofi.us/x/infobot/cmd/run"

	"github.com/urfave/cli/v2"
)

var version string

func Run(args []string, stdout io.Writer) error {
	app := cli.NewApp()
	app.Name = "infobot"
	app.Version = version
	app.HideVersion = false
	app.Usage = "A Keybase bot that lets users store information which can be recalled later"
	app.EnableBashCompletion = true
	app.Writer = stdout
	app.HideVersion = false

	app.Commands = []*cli.Command{
		&run.Command,
	}

	if err := app.Run(args); err != nil {
		return err
	}

	return nil
}
