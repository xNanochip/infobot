package team

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
	"samhofi.us/x/infobot/pkg/infobot"
	"samhofi.us/x/infobot/pkg/utils"
	"samhofi.us/x/keybase/v2"
)

var teamSetCmd = cli.Command{
	Name:      "set",
	Usage:     "Adjust team settings",
	ArgsUsage: "<team name> <option> <value>",
	Action:    teamSet,
}

func teamSet(c *cli.Context) (err error) {
	var args = c.Args()
	if args.Len() == 0 {
		return fmt.Errorf("missing team name, option, and value")
	}
	if args.Len() == 1 {
		return fmt.Errorf("missing option and value")
	}
	if args.Len() == 2 {
		return fmt.Errorf("missing value")
	}

	kb := keybase.New(keybase.SetHomePath(c.Path("home")))

	var (
		teamName = strings.ToLower(args.Get(0))
		option   = args.Get(1)
		value    = strings.ToLower(args.Get(2))
	)

	settings, err := infobot.FetchTeamSettings(kb, teamName)
	if err != nil {
		return fmt.Errorf("Error: Unable to fetch team settings -- %v", err)
	}

	switch option {
	case "nonadmincreate":
		switch value {
		case "true":
			settings.NonAdminCreate = true
		case "false":
			settings.NonAdminCreate = false
		default:
			return fmt.Errorf("Error: Invalid value")
		}
	case "nonadminedit":
		switch value {
		case "true":
			settings.NonAdminEdit = true
		case "false":
			settings.NonAdminEdit = false
		default:
			return fmt.Errorf("Error: Invalid value")
		}
	case "nonadmindelete":
		switch value {
		case "true":
			settings.NonAdminDelete = true
		case "false":
			settings.NonAdminDelete = false
		default:
			return fmt.Errorf("Error: Invalid value")
		}
	default:
		return fmt.Errorf("Error: Invalid option")
	}

	err = infobot.WriteTeamSettings(kb, teamName, settings)
	if err != nil {
		return fmt.Errorf("Error: Unable to write team settings -- %v", err)
	}

	fmt.Fprintf(c.App.Writer, utils.ToJsonPretty(settings))

	return nil
}
