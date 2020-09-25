package team

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
	"samhofi.us/x/infobot/pkg/infobot"
	"samhofi.us/x/infobot/pkg/utils"
	"samhofi.us/x/keybase/v2"
)

var teamAddKeyCmd = cli.Command{
	Name:      "add-key",
	Usage:     "Add a key",
	ArgsUsage: "<team name> <key> <value>",
	Action:    teamAddKey,
}

func teamAddKey(c *cli.Context) (err error) {
	var args = c.Args()
	if args.Len() == 0 {
		return fmt.Errorf("missing team name, key, and value")
	}
	if args.Len() == 1 {
		return fmt.Errorf("missing key and value")
	}
	if args.Len() == 2 {
		return fmt.Errorf("missing value")
	}

	kb := keybase.New(keybase.SetHomePath(c.Path("home")))

	var (
		teamName = strings.ToLower(args.Get(0))
		key      = args.Get(1)
		value    = strings.Join(args.Slice()[2:args.Len()], " ")
	)

	// make sure key doesn't already exist
	keys, err := infobot.GetKeys(kb, teamName)
	if err != nil {
		return fmt.Errorf("Error: Unable to fetch keys -- %v", err)
	}
	if utils.StringInSlice(key, keys) {
		return fmt.Errorf("Error: Key exists")
	}

	// create a new Info and write it to the team's kvstore
	info := infobot.NewInfo(key, value, kb.Username+" via cli")
	err = infobot.WriteInfo(kb, teamName, *info)
	if err != nil {
		return fmt.Errorf("Error: Unable to write new key to team -- %v", err)
	}

	fmt.Fprintf(c.App.Writer, "Successfully added key %s to %s\n", key, teamName)
	return nil
}
