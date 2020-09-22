package run

import (
	"fmt"
	"strings"

	"samhofi.us/x/infobot/pkg/infobot"
	"samhofi.us/x/infobot/pkg/utils"
	"samhofi.us/x/keybase/v2/types/chat1"
)

const (
	errFetchingTeamSettings        = "unable to fetch team settings"
	errFetchingKeys                = "unable to fetch existing keys"
	errFetchingKey                 = "unable to fetch key"
	errKeyNotFound                 = "key not found"
	errKeyExists                   = "key already exists"
	errMissingKeyValue             = "missing key and value"
	errMissingKey                  = "missing key"
	errMissingValue                = "missing value"
	errMissingOptionValue          = "missing option and value"
	errMissingOption               = "missing option"
	errInvalidOption               = "invalid option"
	errInvalidValue                = "invalid value"
	errFailedWritingKey            = "failed to write key to kvstore"
	errFailedDeletingKey           = "failed to delete key from kvstore"
	errFailedWritingSettings       = "failed to write team settings to kvstore"
	errAdminRequiredToAdd          = "only admins are allowed to add new keys"
	errAdminRequiredToEdit         = "only admins are allowed to edit keys"
	errAdminRequiredToDelete       = "only admins are allowed to delete keys"
	errAdminRequiredToEditSettings = "Only admins are allowed to edit team settings"
)

func (b *bot) cmdInfoAdd(m chat1.MsgSummary) error {
	var (
		teamName = m.Channel.Name
		userName = m.Sender.Username
		convID   = m.ConvID
	)

	// fetch team's settings for the bot
	settings, err := infobot.FetchTeamSettings(b.k, teamName)
	if err != nil {
		b.logError("Unable to fetch team settings for team %s. -- %v", teamName, err)
		return fmt.Errorf(errFetchingTeamSettings)
	}
	if !settings.NonAdminCreate {
		if !utils.HasMinRole(b.k, "admin", userName, convID) {
			return fmt.Errorf(errAdminRequiredToAdd)
		}
	}

	// parse the key and value from the received message
	msg := strings.TrimSpace(strings.Replace(m.Content.Text.Body, "!info add ", "", 1))
	if msg == "" {
		return fmt.Errorf(errMissingKeyValue)
	}

	key := strings.Fields(msg)[0]
	value := strings.TrimSpace(strings.Replace(msg, key, "", 1))
	if value == "" {
		return fmt.Errorf(errMissingValue)
	}

	// make sure key doesn't already exist
	keys, err := infobot.GetKeys(b.k, teamName)
	if err != nil {
		b.logError("Unable to fetch keys for team %s. -- %v", teamName, err)
		return fmt.Errorf(errFetchingKeys)
	}
	if utils.StringInSlice(key, keys) {
		return fmt.Errorf(errKeyExists)
	}

	// create a new Info and write it to the team's kvstore
	info := infobot.NewInfo(key, value, userName)
	err = infobot.WriteInfo(b.k, teamName, *info)
	if err != nil {
		b.logError("Unable to write new key to team %s. -- %v", teamName, err)
		return fmt.Errorf(errFailedWritingKey)
	}

	// react to the command message to let them know it was successful
	_, err = b.k.ReactByConvID(convID, m.Id, ":heavy_check_mark:")
	if err != nil {
		b.logError("Error sending reaction: %v", err)
	}
	return nil
}

func (b *bot) cmdInfoEdit(m chat1.MsgSummary) error {
	var (
		teamName = m.Channel.Name
		userName = m.Sender.Username
		convID   = m.ConvID
	)

	// fetch team's settings for the bot
	settings, err := infobot.FetchTeamSettings(b.k, teamName)
	if err != nil {
		b.logError("Unable to fetch team settings for team %s. -- %v", teamName, err)
		return fmt.Errorf(errFetchingTeamSettings)
	}
	if !settings.NonAdminEdit {
		if !utils.HasMinRole(b.k, "admin", userName, convID) {
			return fmt.Errorf(errAdminRequiredToEdit)
		}
	}

	// parse the key and value from the received message
	msg := strings.TrimSpace(strings.Replace(m.Content.Text.Body, "!info edit ", "", 1))
	if msg == "" {
		return fmt.Errorf(errMissingKeyValue)
	}

	key := strings.Fields(msg)[0]
	value := strings.TrimSpace(strings.Replace(msg, key, "", 1))
	if value == "" {
		return fmt.Errorf(errMissingValue)
	}

	// make sure key exists
	keys, err := infobot.GetKeys(b.k, teamName)
	if err != nil {
		b.logError("Unable to fetch keys for team %s. -- %v", teamName, err)
		return fmt.Errorf(errFetchingKeys)
	}
	if !utils.StringInSlice(key, keys) {
		return fmt.Errorf(errKeyNotFound)
	}

	// edit the key
	err = infobot.EditKey(b.k, teamName, key, userName, value)
	if err != nil {
		b.logError("Unable to write new key to team %s. -- %v", teamName, err)
		return fmt.Errorf(errFailedWritingKey)
	}

	// react to the command message to let them know it was successful
	_, err = b.k.ReactByConvID(convID, m.Id, ":heavy_check_mark:")
	if err != nil {
		b.logError("Error sending reaction: %v", err)
	}
	return nil
}

func (b *bot) cmdInfoDelete(m chat1.MsgSummary) error {
	var (
		teamName = m.Channel.Name
		userName = m.Sender.Username
		convID   = m.ConvID
	)

	// fetch team's settings for the bot
	settings, err := infobot.FetchTeamSettings(b.k, teamName)
	if err != nil {
		b.logError("Unable to fetch team settings for team %s. -- %v", teamName, err)
		return fmt.Errorf(errFetchingTeamSettings)
	}
	if !settings.NonAdminDelete {
		if !utils.HasMinRole(b.k, "admin", userName, convID) {
			return fmt.Errorf(errAdminRequiredToDelete)
		}
	}

	// parse the key and value from the received message
	msg := strings.TrimSpace(strings.Replace(m.Content.Text.Body, "!info delete ", "", 1))
	if msg == "" {
		return fmt.Errorf(errMissingKey)
	}

	key := msg

	// make sure key exists
	keys, err := infobot.GetKeys(b.k, teamName)
	if err != nil {
		b.logError("Unable to fetch keys for team %s. -- %v", teamName, err)
		return fmt.Errorf(errFetchingKeys)
	}
	if !utils.StringInSlice(key, keys) {
		return fmt.Errorf(errKeyNotFound)
	}

	// fetch existing key info from the team's kvstore
	info, err := infobot.FetchKey(b.k, teamName, key)
	if err != nil {
		b.logError("Unable to fetch key for team %s. -- %v", teamName, err)
		return fmt.Errorf(errFetchingKey)
	}

	// delete key
	err = infobot.DeleteKey(b.k, teamName, info)
	if err != nil {
		b.logError("Unable to delete key from team %s. -- %v", teamName, err)
		return fmt.Errorf(errFailedDeletingKey)
	}

	// react to the command message to let them know it was successful
	_, err = b.k.ReactByConvID(convID, m.Id, ":heavy_check_mark:")
	if err != nil {
		b.logError("Error sending reaction: %v", err)
	}

	return nil
}

func (b *bot) cmdInfoRead(m chat1.MsgSummary) error {
	var (
		teamName = m.Channel.Name
		convID   = m.ConvID
	)

	// parse the key from the received message
	msg := strings.TrimSpace(strings.Replace(m.Content.Text.Body, "!info read ", "", 1))
	if msg == "" {
		return fmt.Errorf(errMissingKey)
	}

	key := msg

	// fetch key info from the team's kvstore
	info, err := infobot.FetchKey(b.k, teamName, key)
	if err != nil {
		b.logError("Unable to fetch key for team %s. -- %v", teamName, err)
		return fmt.Errorf(errFetchingKey)
	}
	_, err = b.k.SendMessageByConvID(convID, info.Value)
	if err != nil {
		b.logError("Error sending reply: %v", err)
	}

	return nil
}

func (b *bot) cmdInfoAudit(m chat1.MsgSummary) error {
	var (
		teamName = m.Channel.Name
		convID   = m.ConvID
	)

	msg := strings.TrimSpace(strings.Replace(m.Content.Text.Body, "!info audit ", "", 1))
	if msg == "" {
		return fmt.Errorf(errMissingKey)
	}
	key := msg

	// fetch key info from the team's kvstore
	info, err := infobot.FetchKey(b.k, teamName, key)
	if err != nil {
		b.logError("Unable to fetch key for team %s. -- %v", teamName, err)
		return fmt.Errorf(errFetchingKey)
	}

	_, err = b.k.SendMessageByConvID(convID, utils.ToJsonPretty(info))
	if err != nil {
		b.logError("Error sending reply: %v", err)
	}

	return nil
}

func (b *bot) cmdInfoSet(m chat1.MsgSummary) error {
	var (
		teamName = m.Channel.Name
		userName = m.Sender.Username
		convID   = m.ConvID
	)

	availableOptions := []string{
		"nonadmincreate",
		"nonadminedit",
		"nonadmindelete",
	}
	if !utils.HasMinRole(b.k, "admin", userName, convID) {
		return fmt.Errorf(errAdminRequiredToEditSettings)
	}

	// parse the option and value from the received message
	msg := strings.ToLower(strings.TrimSpace(strings.Replace(m.Content.Text.Body, "!info set ", "", 1)))
	if msg == "" {
		return fmt.Errorf(errMissingOptionValue)
	}

	option := strings.Fields(msg)[0]

	if !utils.StringInSlice(option, availableOptions) {
		return fmt.Errorf(errInvalidOption)
	}

	value := strings.TrimSpace(strings.Replace(msg, option, "", 1))
	if value == "" {
		b.k.ReplyByConvID(convID, m.Id, "Error: Missing value")
		return fmt.Errorf(errMissingValue)
	}

	settings, err := infobot.FetchTeamSettings(b.k, teamName)
	if err != nil {
		b.logError("Unable to fetch team settings for team %s. -- %v", teamName, err)
		return fmt.Errorf(errFetchingTeamSettings)
	}

	if value != "true" && value != "false" {
		return fmt.Errorf(errInvalidValue)
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
		b.logError("Unable to write team settings for team %s. -- %v", teamName, err)
		return fmt.Errorf(errFailedWritingSettings)
	}
	// react to the command message to let them know it was successful
	_, err = b.k.ReactByConvID(convID, m.Id, ":heavy_check_mark:")
	if err != nil {
		b.logError("Error sending reaction: %v", err)
	}

	return nil
}

func (b *bot) cmdInfoSettings(m chat1.MsgSummary) error {
	var (
		teamName = m.Channel.Name
		convID   = m.ConvID
	)
	// fetch key info from the team's kvstore
	settings, err := infobot.FetchTeamSettings(b.k, teamName)
	if err != nil {
		b.logError("Unable to fetch settings for team %s. -- %v", teamName, err)
		return fmt.Errorf(errFetchingTeamSettings)
	}

	_, err = b.k.SendMessageByConvID(convID, utils.ToJsonPretty(settings))
	if err != nil {
		b.logError("Error sending reply: %v", err)
	}

	return nil
}

func (b *bot) cmdInfoKeys(m chat1.MsgSummary) error {
	var (
		teamName = m.Channel.Name
		convID   = m.ConvID
	)

	// fetch key info from the team's kvstore
	keys, err := infobot.GetKeys(b.k, teamName)
	if err != nil {
		b.logError("Unable to fetch keys for team %s. -- %v", teamName, err)
		return fmt.Errorf(errFetchingKeys)
	}

	if len(keys) == 0 {
		_, err = b.k.ReplyByConvID(convID, m.Id, "There aren't any keys yet")
		if err != nil {
			b.logError("Error sending reply: %v", err)
		}
		return nil
	}

	_, err = b.k.ReplyByConvID(convID, m.Id, strings.Join(keys, ", "))
	if err != nil {
		b.logError("Error sending reply: %v", err)
	}

	return nil
}

func (b *bot) cmdAtMention(m chat1.MsgSummary) error {
	var (
		teamName = m.Channel.Name
		convID   = m.ConvID
	)

	key := strings.TrimSpace(strings.ReplaceAll(m.Content.Text.Body, fmt.Sprintf("@%s", b.k.Username), ""))
	if key == "" {
		_, err := b.k.ReplyByConvID(convID, m.Id, "Try @%s <key>", b.k.Username)
		if err != nil {
			b.logError("Error sending reply: %v", err)
		}
		return nil
	}

	// fetch key info from the team's kvstore
	info, err := infobot.FetchKey(b.k, teamName, key)
	if err != nil {
		b.logError("Unable to fetch key for team %s. -- %v", teamName, err)
		return fmt.Errorf(errFetchingKey)
	}
	_, err = b.k.SendMessageByConvID(convID, info.Value)
	if err != nil {
		b.logError("Error sending reply: %v", err)
	}

	return nil
}
