package team

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
	"samhofi.us/x/infobot/pkg/infobot"
	"samhofi.us/x/infobot/pkg/utils"
	"samhofi.us/x/keybase/v2"
)

var teamUnlockKeyCmd = cli.Command{
	Name:      "unlock-key",
	Usage:     "Unlock a key",
	ArgsUsage: "<team name> <key>",
	Action:    teamUnlockKey,
}

func teamUnlockKey(c *cli.Context) (err error) {
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
		return fmt.Errorf("Error: Key not found")
	}

	// unlock the key
	err = infobot.UnlockKey(kb, teamName, key, "[cli]")
	if err != nil {
		if err.Error() == "key is already unlocked" {
			fmt.Fprintf(c.App.Writer, "Key is already unlocked, no changes made\n")
			return nil
		}

		return fmt.Errorf("Error: Unable to unlock key -- %v", err)
	}

	fmt.Fprintf(c.App.Writer, "Successfully unlocked key %s in %s\n", key, teamName)
	return nil
}
