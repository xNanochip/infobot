package team

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
	"samhofi.us/x/infobot/pkg/infobot"
	"samhofi.us/x/infobot/pkg/utils"
	"samhofi.us/x/keybase/v2"
)

var teamGetSettingsCmd = cli.Command{
	Name:      "get-settings",
	Usage:     "Fetch the team's settings",
	ArgsUsage: "<team name>",
	Action:    teamGetSettings,
}

func teamGetSettings(c *cli.Context) (err error) {
	var args = c.Args()
	if args.Len() == 0 {
		return fmt.Errorf("missing team name")
	}

	kb := keybase.New(keybase.SetHomePath(c.Path("home")))

	var (
		teamName = strings.ToLower(args.Get(0))
	)

	// fetch key info from the team's kvstore
	settings, err := infobot.FetchTeamSettings(kb, teamName)
	if err != nil {
		return fmt.Errorf("Error: Unable to fetch settings -- %v", err)
	}

	fmt.Fprintf(c.App.Writer, "%s\n", utils.ToJsonPretty(settings))
	return nil
}
