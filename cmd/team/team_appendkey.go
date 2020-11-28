package team

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
	"samhofi.us/x/infobot/pkg/infobot"
	"samhofi.us/x/infobot/pkg/utils"
	"samhofi.us/x/keybase/v2"
)

var teamAppendKeyCmd = cli.Command{
	Name:      "append-key",
	Usage:     "Append a new line to a key",
	ArgsUsage: "<team name> <key> <new line>",
	Action:    teamAppendKey,
}

func teamAppendKey(c *cli.Context) (err error) {
	var args = c.Args()
	if args.Len() == 0 {
		return fmt.Errorf("missing team name, key, and new line")
	}
	if args.Len() == 1 {
		return fmt.Errorf("missing key and new line")
	}
	if args.Len() == 2 {
		return fmt.Errorf("missing new line")
	}

	kb := keybase.New(keybase.SetHomePath(c.Path("home")))

	var (
		teamName = strings.ToLower(args.Get(0))
		key      = args.Get(1)
		value    = strings.Join(args.Slice()[2:args.Len()], " ")
	)

	// make sure key exists
	keys, err := infobot.GetKeys(kb, teamName)
	if err != nil {
		return fmt.Errorf("Error: Unable to fetch keys -- %v", err)
	}
	if !utils.StringInSlice(key, keys) {
		return fmt.Errorf("Error: Key not found")
	}

	// append to the key
	err = infobot.AppendKey(kb, teamName, key, kb.Username+" (via cli)", value)
	if err != nil {
		return fmt.Errorf("Error: Unable to write new key to team -- %v", err)
	}

	fmt.Fprintf(c.App.Writer, "Successfully appended to key %s in %s\n", key, teamName)
	return nil
}
