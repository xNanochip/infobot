package team

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
	"samhofi.us/x/infobot/pkg/infobot"
	"samhofi.us/x/infobot/pkg/utils"
	"samhofi.us/x/keybase/v2"
)

var teamDeleteKeyCmd = cli.Command{
	Name:      "delete-key",
	Usage:     "Delete a key",
	ArgsUsage: "<team name> <key>",
	Action:    teamDeleteKey,
}

func teamDeleteKey(c *cli.Context) (err error) {
	var args = c.Args()
	if args.Len() == 0 {
		return fmt.Errorf("missing team name and key")
	}
	if args.Len() == 1 {
		return fmt.Errorf("missing key")
	}

	kb := keybase.New(keybase.SetHomePath(c.Path("home")))

	var (
		teamName = strings.ToLower(args.Get(0))
		key      = args.Get(1)
	)

	// make sure key exists
	keys, err := infobot.GetKeys(kb, teamName)
	if err != nil {
		return fmt.Errorf("Error: Unable to fetch keys -- %v", err)
	}
	if !utils.StringInSlice(key, keys) {
		return fmt.Errorf("Error: Key does not exist")
	}

	// fetch existing key info from the team's kvstore
	info, err := infobot.FetchKey(kb, teamName, key)
	if err != nil {
		return fmt.Errorf("Error: Unable to fetch key -- %v", err)
	}

	// delete key
	err = infobot.DeleteKey(kb, teamName, info)
	if err != nil {
		return fmt.Errorf("Error: Unable to delete key from team -- %v", err)
	}

	fmt.Fprintf(c.App.Writer, "Successfully deleted key %s from %s\n", key, teamName)
	return nil
}
