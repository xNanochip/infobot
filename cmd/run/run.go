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
			stdout: c.App.Writer,
			stderr: c.App.ErrWriter,
		},
	}

	clearCommands := b.advertiseCommands()
	defer clearCommands()

	// catch ctrl + c
	var trap = make(chan os.Signal, 1)
	signal.Notify(trap, os.Interrupt)
	go func() {
		for _ = range trap {
			b.log_debug("Received interrupt signal. Clearing commands and exiting...")
			clearCommands()
			os.Exit(0)
		}
	}()

	b.registerHandlers()
	b.log_info("Running as user %s", b.k.Username)
	b.k.Run(b.handlers, &b.opts)
	return nil
}
