package run

import (
	"fmt"
	"strings"

	"samhofi.us/x/infobot/pkg/infobot"
	"samhofi.us/x/infobot/pkg/utils"
	"samhofi.us/x/keybase/v2"
	"samhofi.us/x/keybase/v2/types/chat1"
	"samhofi.us/x/keybase/v2/types/stellar1"
)

func (b *bot) registerHandlers() {
	b.log_debug("Registering handlers")

	var (
		chat   = b.chatHandler
		conv   = b.convHandler
		wallet = b.walletHandler
		err    = b.errorHandler
	)
	b.handlers = keybase.Handlers{
		ChatHandler:         &chat,
		ConversationHandler: &conv,
		WalletHandler:       &wallet,
		ErrorHandler:        &err,
	}
}

func (b *bot) chatHandler(m chat1.MsgSummary) {
	var (
		teamName = m.Channel.Name
		userName = m.Sender.Username
		convID   = m.ConvID
	)

	if userName == b.k.Username {
		return
	}

	switch m.Content.TypeName {
	case "text":
		if strings.HasPrefix(m.Content.Text.Body, "!info add ") {
			// fetch team's settings for the bot
			settings, err := infobot.FetchTeamSettings(b.k, teamName)
			if err != nil {
				b.log_error("Unable to fetch team settings for team %s. -- %v", teamName, err)
				b.k.ReplyByConvID(convID, m.Id, "Unable to fetch team settings")
				return
			}
			if !settings.NonAdminCreate {
				if !utils.HasMinRole(b.k, "admin", userName, convID) {
					b.k.ReplyByConvID(convID, m.Id, "Only admins are allowed to add new keys")
					return
				}
			}

			// parse the key and value from the received message
			msg := strings.TrimSpace(strings.Replace(m.Content.Text.Body, "!info add ", "", 1))
			if msg == "" {
				b.k.ReplyByConvID(convID, m.Id, "Error: Missing key and value")
				return
			}

			key := strings.Fields(msg)[0]
			value := strings.TrimSpace(strings.Replace(msg, key, "", 1))
			if value == "" {
				b.k.ReplyByConvID(convID, m.Id, "Error: Missing value")
				return
			}

			// make sure key doesn't already exist
			keys, err := infobot.GetKeys(b.k, teamName)
			if err != nil {
				b.log_error("Unable to fetch keys for team %s. -- %v", teamName, err)
				b.k.ReplyByConvID(convID, m.Id, "Unable to fetch existing keys")
				return
			}
			if utils.StringInSlice(key, keys) {
				b.k.ReplyByConvID(convID, m.Id, "Error: Key already exists")
				return
			}

			// create a new Info and write it to the team's kvstore
			info := infobot.NewInfo(key, value, userName)
			err = infobot.WriteInfo(b.k, teamName, *info)
			if err != nil {
				b.log_error("Unable to write new key to team %s. -- %v", teamName, err)
				b.k.ReplyByConvID(convID, m.Id, "Failed to write new key to kvstore: %v", err)
				return
			}

			// react to the command message to let them know it was successful
			_, err = b.k.ReactByConvID(convID, m.Id, ":heavy_check_mark:")
			if err != nil {
				b.log_error("Error sending reaction: %v", err)
			}
			return
		}
		if strings.HasPrefix(m.Content.Text.Body, "!info edit ") {
			// fetch team's settings for the bot
			settings, err := infobot.FetchTeamSettings(b.k, teamName)
			if err != nil {
				b.log_error("Unable to fetch team settings for team %s. -- %v", teamName, err)
				b.k.ReplyByConvID(convID, m.Id, "Unable to fetch team settings")
				return
			}
			if !settings.NonAdminEdit {
				if !utils.HasMinRole(b.k, "admin", userName, convID) {
					b.k.ReplyByConvID(convID, m.Id, "Only admins are allowed to edit keys")
					return
				}
			}

			// parse the key and value from the received message
			msg := strings.TrimSpace(strings.Replace(m.Content.Text.Body, "!info edit ", "", 1))
			if msg == "" {
				b.k.ReplyByConvID(convID, m.Id, "Error: Missing key and value")
				return
			}

			key := strings.Fields(msg)[0]
			value := strings.TrimSpace(strings.Replace(msg, key, "", 1))
			if value == "" {
				b.k.ReplyByConvID(convID, m.Id, "Error: Missing value")
				return
			}

			// make sure key exists
			keys, err := infobot.GetKeys(b.k, teamName)
			if err != nil {
				b.log_error("Unable to fetch keys for team %s. -- %v", teamName, err)
				b.k.ReplyByConvID(convID, m.Id, "Unable to fetch existing keys")
				return
			}
			if !utils.StringInSlice(key, keys) {
				b.k.ReplyByConvID(convID, m.Id, "Error: Key not found")
				return
			}

			// fetch existing key info from the team's kvstore
			info, err := infobot.FetchKey(b.k, teamName, key)
			if err != nil {
				b.log_error("Unable to fetch key for team %s. -- %v", teamName, err)
				b.k.ReplyByConvID(convID, m.Id, "Unable to fetch key")
				return
			}

			// write a new edit item to the key
			edit := infobot.NewEdit(userName, value)
			info.Edits = append(info.Edits, *edit)
			info.Value = value

			// save edited key
			err = infobot.WriteInfo(b.k, teamName, info)
			if err != nil {
				b.log_error("Unable to write new key to team %s. -- %v", teamName, err)
				b.k.ReplyByConvID(convID, m.Id, "Failed to write new key to kvstore: %v", err)
				return
			}

			// react to the command message to let them know it was successful
			_, err = b.k.ReactByConvID(convID, m.Id, ":heavy_check_mark:")
			if err != nil {
				b.log_error("Error sending reaction: %v", err)
			}
			return
		}
		if strings.HasPrefix(m.Content.Text.Body, "!info delete ") {
			// fetch team's settings for the bot
			settings, err := infobot.FetchTeamSettings(b.k, teamName)
			if err != nil {
				b.log_error("Unable to fetch team settings for team %s. -- %v", teamName, err)
				b.k.ReplyByConvID(convID, m.Id, "Unable to fetch team settings")
				return
			}
			if !settings.NonAdminDelete {
				if !utils.HasMinRole(b.k, "admin", userName, convID) {
					b.k.ReplyByConvID(convID, m.Id, "Only admins are allowed to delete keys")
					return
				}
			}

			// parse the key and value from the received message
			msg := strings.TrimSpace(strings.Replace(m.Content.Text.Body, "!info delete ", "", 1))
			if msg == "" {
				b.k.ReplyByConvID(convID, m.Id, "Error: Missing key")
				return
			}

			key := strings.Fields(msg)[0]

			// make sure key exists
			keys, err := infobot.GetKeys(b.k, teamName)
			if err != nil {
				b.log_error("Unable to fetch keys for team %s. -- %v", teamName, err)
				b.k.ReplyByConvID(convID, m.Id, "Unable to fetch existing keys")
				return
			}
			if !utils.StringInSlice(key, keys) {
				b.k.ReplyByConvID(convID, m.Id, "Error: Key not found")
				return
			}

			// fetch existing key info from the team's kvstore
			info, err := infobot.FetchKey(b.k, teamName, key)
			if err != nil {
				b.log_error("Unable to fetch key for team %s. -- %v", teamName, err)
				b.k.ReplyByConvID(convID, m.Id, "Unable to fetch key")
				return
			}

			// delete key
			err = infobot.DeleteKey(b.k, teamName, info)
			if err != nil {
				b.log_error("Unable to delete key from team %s. -- %v", teamName, err)
				b.k.ReplyByConvID(convID, m.Id, "Failed to delete key from kvstore: %v", err)
				return
			}

			// react to the command message to let them know it was successful
			_, err = b.k.ReactByConvID(convID, m.Id, ":heavy_check_mark:")
			if err != nil {
				b.log_error("Error sending reaction: %v", err)
			}
			return
		}
		if strings.HasPrefix(m.Content.Text.Body, "!info read ") {
			// parse the key from the received message
			msg := strings.TrimSpace(strings.Replace(m.Content.Text.Body, "!info read ", "", 1))
			if msg == "" {
				b.k.ReplyByConvID(convID, m.Id, "Error: Missing key")
				return
			}

			key := strings.Fields(msg)[0]

			// fetch key info from the team's kvstore
			info, err := infobot.FetchKey(b.k, teamName, key)
			if err != nil {
				b.log_error("Unable to fetch key for team %s. -- %v", teamName, err)
				b.k.ReplyByConvID(convID, m.Id, "Unable to fetch key")
				return
			}
			_, err = b.k.ReplyByConvID(convID, m.Id, info.Value)
			if err != nil {
				b.log_error("Error sending reply: %v", err)
			}
			return
		}
		if strings.HasPrefix(m.Content.Text.Body, "!info audit ") {
			msg := strings.TrimSpace(strings.Replace(m.Content.Text.Body, "!info audit ", "", 1))
			if msg == "" {
				b.k.ReplyByConvID(convID, m.Id, "Error: Missing key")
				return
			}
			key := strings.Fields(msg)[0]

			// fetch key info from the team's kvstore
			info, err := infobot.FetchKey(b.k, teamName, key)
			if err != nil {
				b.log_error("Unable to fetch key for team %s. -- %v", teamName, err)
				b.k.ReplyByConvID(convID, m.Id, "Unable to fetch key")
				return
			}

			_, err = b.k.ReplyByConvID(convID, m.Id, utils.ToJsonPretty(info))
			if err != nil {
				b.log_error("Error sending reply: %v", err)
			}
			return
		}
		if strings.HasPrefix(m.Content.Text.Body, "!info set ") {
			availableOptions := []string{
				"nonadmincreate",
				"nonadminedit",
				"nonadmindelete",
			}
			if !utils.HasMinRole(b.k, "admin", userName, convID) {
				b.k.ReplyByConvID(convID, m.Id, "Only admins are allowed to edit team settings")
				return
			}

			// parse the option and value from the received message
			msg := strings.ToLower(strings.TrimSpace(strings.Replace(m.Content.Text.Body, "!info set ", "", 1)))
			if msg == "" {
				b.k.ReplyByConvID(convID, m.Id, "Error: Missing option and value")
				return
			}

			option := strings.Fields(msg)[0]

			if !utils.StringInSlice(option, availableOptions) {
				b.k.ReplyByConvID(convID, m.Id, "Error: Invalid option")
				return
			}

			value := strings.TrimSpace(strings.Replace(msg, option, "", 1))
			if value == "" {
				b.k.ReplyByConvID(convID, m.Id, "Error: Missing value")
				return
			}

			settings, err := infobot.FetchTeamSettings(b.k, teamName)
			if err != nil {
				b.log_error("Unable to fetch team settings for team %s. -- %v", teamName, err)
				b.k.ReplyByConvID(convID, m.Id, "Unable to fetch team settings")
				return
			}

			if value != "true" && value != "false" {
				b.k.ReplyByConvID(convID, m.Id, "Error: Invalid value %s", value)
				return
			}
			switch option {
			case "nonadmincreate":
				if value == "true" {
					settings.NonAdminCreate = true
				}
				if value == "false" {
					settings.NonAdminCreate = false
				}
			case "nonadminedit":
				if value == "true" {
					settings.NonAdminEdit = true
				}
				if value == "false" {
					settings.NonAdminEdit = false
				}
			case "nonadmindelete":
				if value == "true" {
					settings.NonAdminDelete = true
				}
				if value == "false" {
					settings.NonAdminDelete = false
				}
			}
			err = infobot.WriteTeamSettings(b.k, teamName, settings)
			if err != nil {
				b.log_error("Unable to write team settings for team %s. -- %v", teamName, err)
				b.k.ReplyByConvID(convID, m.Id, "Unable to write team settings")
				return
			}
			// react to the command message to let them know it was successful
			_, err = b.k.ReactByConvID(convID, m.Id, ":heavy_check_mark:")
			if err != nil {
				b.log_error("Error sending reaction: %v", err)
			}
			return
		}
		if strings.HasPrefix(m.Content.Text.Body, "!info settings") {
			// fetch key info from the team's kvstore
			settings, err := infobot.FetchTeamSettings(b.k, teamName)
			if err != nil {
				b.log_error("Unable to fetch settings for team %s. -- %v", teamName, err)
				b.k.ReplyByConvID(convID, m.Id, "Unable to fetch settings")
				return
			}

			_, err = b.k.ReplyByConvID(convID, m.Id, utils.ToJsonPretty(settings))
			if err != nil {
				b.log_error("Error sending reply: %v", err)
			}
			return
		}
		if utils.StringInSlice(b.k.Username, m.AtMentionUsernames) {
			key := strings.TrimSpace(strings.ReplaceAll(m.Content.Text.Body, fmt.Sprintf("@%s", b.k.Username), ""))
			if key == "" {
				_, err := b.k.ReplyByConvID(convID, m.Id, "Try @%s <key>", b.k.Username)
				if err != nil {
					b.log_error("Error sending reply: %v", err)
				}
				return
			}

			// fetch key info from the team's kvstore
			info, err := infobot.FetchKey(b.k, teamName, key)
			if err != nil {
				b.log_error("Unable to fetch key for team %s. -- %v", teamName, err)
				b.k.ReplyByConvID(convID, m.Id, "Unable to fetch key")
				return
			}
			_, err = b.k.ReplyByConvID(convID, m.Id, info.Value)
			if err != nil {
				b.log_error("Error sending reply: %v", err)
			}
			return
		}
	}
}

func (b *bot) convHandler(c chat1.ConvSummary) {
}

func (b *bot) walletHandler(p stellar1.PaymentDetailsLocal) {
}

func (b *bot) errorHandler(e error) {
}
