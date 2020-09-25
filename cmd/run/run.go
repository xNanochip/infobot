package run

import (
	"os"
	"os/signal"

	"github.com/urfave/cli/v2"
	"samhofi.us/x/keybase/v2"
)

// Command exports the feeds command set.
var Command = cli.Command{
	Name:   "run",
	Usage:  "run the bot",
	Action: run,
}

func run(c *cli.Context) error {
	var b = bot{
		k: keybase.New(keybase.SetHomePath(c.Path("home"))),
		config: botConfig{
			debug:  c.Bool("debug"),
			json:   c.Bool("json"),
			stdout: c.App.Writer,
			stderr: c.App.ErrWriter,
		},
	}

	b.advertiseCommands()
	defer b.clearCommands()

	// catch ctrl + c
	var trap = make(chan os.Signal, 1)
	signal.Notify(trap, os.Interrupt)
	go func() {
		for _ = range trap {
			b.logDebug("Received interrupt signal")
			b.clearCommands()
			os.Exit(0)
		}
	}()

	b.registerHandlers()
	b.logInfo("Running as user %s", b.k.Username)
	b.k.Run(b.handlers, &b.opts)
	return nil
}
