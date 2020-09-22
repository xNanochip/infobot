package cmd

import (
	"io"

	"samhofi.us/x/infobot/cmd/run"
	"samhofi.us/x/infobot/cmd/team"

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

	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:    "debug",
			Aliases: []string{"d"},
			EnvVars: []string{"INFOBOT_DEBUG"},
			Usage:   "Enable debug mode",
		},
		&cli.PathFlag{
			Name:    "home",
			Aliases: []string{"H"},
			EnvVars: []string{"INFOBOT_HOME"},
			Usage:   "Set an alternate home directory for the Keybase client",
		},
		&cli.BoolFlag{
			Name:    "json",
			Aliases: []string{"j"},
			EnvVars: []string{"INFOBOT_JSON"},
			Usage:   "Output logs in JSON format",
		},
	}

	app.Commands = []*cli.Command{
		&run.Command,
		&team.Command,
	}

	if err := app.Run(args); err != nil {
		return err
	}

	return nil
}
