package run

import (
	"github.com/urfave/cli/v2"
)

// Command exports the feeds command set.
var Command = cli.Command{
	Name:        "run",
	Usage:       "run the bot",
	Subcommands: []*cli.Command{},
}
