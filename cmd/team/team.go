package team

import "github.com/urfave/cli/v2"

// Command exports the team command set.
var Command = cli.Command{
	Name:  "team",
	Usage: "manage a keys and settings for a team",
	Subcommands: []*cli.Command{
		&teamListKeysCmd,
		&teamAddKeyCmd,
		&teamEditKeyCmd,
		&teamAppendKeyCmd,
		&teamLockKeyCmd,
		&teamUnlockKeyCmd,
		&teamReadKeyCmd,
		&teamDeleteKeyCmd,
		&teamAuditKeyCmd,
		&teamGetSettingsCmd,
		&teamSetCmd,
	},
}
