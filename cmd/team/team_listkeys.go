package team

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
	"samhofi.us/x/infobot/pkg/infobot"
	"samhofi.us/x/infobot/pkg/utils"
	"samhofi.us/x/keybase/v2"
)

var teamListKeysCmd = cli.Command{
	Name:      "list-keys",
	Usage:     "List all keys",
	ArgsUsage: "<team name>",
	Action:    teamListKeys,
}

func teamListKeys(c *cli.Context) (err error) {
	var args = c.Args()
	if args.Len() != 1 {
		return fmt.Errorf("missing team name")
	}

	kb := keybase.New(keybase.SetHomePath(c.Path("home")))
	keys, err := infobot.GetKeys(kb, args.Get(0))
	if err != nil {
		return err
	}

	if c.Bool("json") {
		fmt.Fprintf(c.App.Writer, utils.ToJson(keys)+"\n")
		return nil
	}

	fmt.Fprintf(c.App.Writer, "%s\n", strings.Join(keys, "\n"))
	return nil
}
