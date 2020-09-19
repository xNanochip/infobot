package run

import (
	"samhofi.us/x/keybase/v2"
	"samhofi.us/x/keybase/v2/types/chat1"
)

func (b *bot) advertiseCommands() func() {
	b.log_debug("Advertising commands")
	b.k.ClearCommands()
	opts := keybase.AdvertiseCommandsOptions{
		Alias: "InfoBot",
		Advertisements: []chat1.AdvertiseCommandAPIParam{
			{
				Typ: "public",
				Commands: []chat1.UserBotCommandInput{
					{
						Name:        "ping",
						Description: "Pings the bot",
					},
					{
						Name:        "ding",
						Description: "Pings the bot",
					},
				},
			},
		},
	}
	b.k.AdvertiseCommands(opts)

	return func() {
		b.log_debug("Clearing commands")
		b.k.ClearCommands()
	}
}
