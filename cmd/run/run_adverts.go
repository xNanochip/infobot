package run

import (
	"samhofi.us/x/keybase/v2"
	"samhofi.us/x/keybase/v2/types/chat1"
)

func (b *bot) advertiseCommands() {
	b.logDebug("Advertising commands")
	b.k.ClearCommands()
	opts := keybase.AdvertiseCommandsOptions{
		Alias: "InfoBot",
		Advertisements: []chat1.AdvertiseCommandAPIParam{
			{
				Typ: "public",
				Commands: []chat1.UserBotCommandInput{
					{
						Name:        "info add",
						Usage:       "<key> <value>",
						Description: "Add a new key",
					},
					{
						Name:        "info edit",
						Usage:       "<key> <new value>",
						Description: "Edit a key",
					},
					{
						Name:        "info append",
						Usage:       "<key> <new line>",
						Description: "Append a new line to a key",
					},
					{
						Name:        "info lock",
						Usage:       "<key>",
						Description: "Lock a key",
					},
					{
						Name:        "info unlock",
						Usage:       "<key>",
						Description: "Unlock a key",
					},
					{
						Name:        "info delete",
						Usage:       "<key>",
						Description: "Delete a key",
					},
					{
						Name:        "info read",
						Usage:       "<key>",
						Description: "Read a key's value",
					},
					{
						Name:        "info audit",
						Usage:       "<key>",
						Description: "Print all info about a key, including its edit history",
					},
					{
						Name:        "info set",
						Usage:       "<option> <value>",
						Description: "Modify team settings",
					},
					{
						Name:        "info settings",
						Description: "Read team settings",
					},
					{
						Name:        "info keys",
						Description: "List all available keys",
					},
				},
			},
		},
	}
	b.k.AdvertiseCommands(opts)
}

func (b *bot) clearCommands() {
	b.logDebug("Clearing commands")
	b.k.ClearCommands()
}
