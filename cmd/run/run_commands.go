package run

import (
	"fmt"
	"strings"

	"samhofi.us/x/infobot/pkg/infobot"
	"samhofi.us/x/infobot/pkg/utils"
	"samhofi.us/x/keybase/v2"
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
	errAdminRequiredToEditLocked   = "this key is locked, and can only be edited by an admin"
	errAdminRequiredToLock         = "Only admins are allowed to lock keys"
	errAdminRequiredToUnlock       = "Only admins are allowed to unlock keys"
	errAdminRequiredToDelete       = "only admins are allowed to delete keys"
	errAdminRequiredToDeleteLocked = "this key is locked, and can only be deleted by an admin"
	errAdminRequiredToEditSettings = "Only admins are allowed to edit team settings"
)

func (b *bot) cmdInfoAdd(m chat1.MsgSummary) error {
	var (
		teamName = m.Channel.Name
		userName = m.Sender.Username
		convID   = m.ConvID
	)

	// this is needed in case a user adds the bot to their own conversation with themself
	if m.Channel.MembersType == keybase.USER && !strings.Contains(teamName, ",") {
		teamName = teamName + "," + teamName
	}

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

	// this is needed in case a user adds the bot to their own conversation with themself
	if m.Channel.MembersType == keybase.USER && !strings.Contains(teamName, ",") {
		teamName = teamName + "," + teamName
	}

	// fetch team's settings for the bot
	settings, err := infobot.FetchTeamSettings(b.k, teamName)
	if err != nil {
		b.logError("Unable to fetch team settings for team %s. -- %v", teamName, err)
		return fmt.Errorf(errFetchingTeamSettings)
	}

	// i don't want to check if they're an admin twice, and this makes more sense to check earlier rather than later,
	// so i'm creating these vars so i know later if i've already checked this
	var (
		adminChecked = false
		isAdmin      = false
	)
	if !settings.NonAdminEdit {
		adminChecked = true
		if !utils.HasMinRole(b.k, "admin", userName, convID) {
			return fmt.Errorf(errAdminRequiredToEdit)
		}
		isAdmin = true
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

	// fetch key info from the team's kvstore
	info, err := infobot.FetchKey(b.k, teamName, key)
	if err != nil {
		b.logError("Unable to fetch key for team %s. -- %v", teamName, err)
		return fmt.Errorf(errFetchingKey)
	}

	// handle locked keys
	if info.Locked && !isAdmin {
		if adminChecked {
			return fmt.Errorf(errAdminRequiredToEditLocked)
		}
		if !utils.HasMinRole(b.k, "admin", userName, convID) {
			return fmt.Errorf(errAdminRequiredToEditLocked)
		}
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

func (b *bot) cmdInfoLock(m chat1.MsgSummary) error {
	var (
		teamName = m.Channel.Name
		userName = m.Sender.Username
		convID   = m.ConvID
	)

	// this is needed in case a user adds the bot to their own conversation with themself
	if m.Channel.MembersType == keybase.USER && !strings.Contains(teamName, ",") {
		teamName = teamName + "," + teamName
	}

	// only admins can lock keys
	if !utils.HasMinRole(b.k, "admin", userName, convID) {
		return fmt.Errorf(errAdminRequiredToLock)
	}

	// parse the key and value from the received message
	msg := strings.TrimSpace(strings.Replace(m.Content.Text.Body, "!info lock ", "", 1))
	if msg == "" {
		return fmt.Errorf(errMissingKeyValue)
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

	// lock the key
	err = infobot.LockKey(b.k, teamName, key, userName)
	if err != nil {
		if err.Error() == "key is already locked" {
			b.k.ReactByConvID(convID, m.Id, "no change")
			return nil
		}

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

func (b *bot) cmdInfoUnlock(m chat1.MsgSummary) error {
	var (
		teamName = m.Channel.Name
		userName = m.Sender.Username
		convID   = m.ConvID
	)

	// this is needed in case a user adds the bot to their own conversation with themself
	if m.Channel.MembersType == keybase.USER && !strings.Contains(teamName, ",") {
		teamName = teamName + "," + teamName
	}

	// only admins can unlock keys
	if !utils.HasMinRole(b.k, "admin", userName, convID) {
		return fmt.Errorf(errAdminRequiredToUnlock)
	}

	// parse the key and value from the received message
	msg := strings.TrimSpace(strings.Replace(m.Content.Text.Body, "!info unlock ", "", 1))
	if msg == "" {
		return fmt.Errorf(errMissingKeyValue)
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

	// unlock the key
	err = infobot.UnlockKey(b.k, teamName, key, userName)
	if err != nil {
		if err.Error() == "key is already unlocked" {
			b.k.ReactByConvID(convID, m.Id, "no change")
			return nil
		}

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

	// this is needed in case a user adds the bot to their own conversation with themself
	if m.Channel.MembersType == keybase.USER && !strings.Contains(teamName, ",") {
		teamName = teamName + "," + teamName
	}

	// fetch team's settings for the bot
	settings, err := infobot.FetchTeamSettings(b.k, teamName)
	if err != nil {
		b.logError("Unable to fetch team settings for team %s. -- %v", teamName, err)
		return fmt.Errorf(errFetchingTeamSettings)
	}

	// i don't want to check if they're an admin twice, and this makes more sense to check earlier rather than later,
	// so i'm creating these vars so i know later if i've already checked this
	var (
		adminChecked = false
		isAdmin      = false
	)
	if !settings.NonAdminDelete {
		adminChecked = true
		if !utils.HasMinRole(b.k, "admin", userName, convID) {
			return fmt.Errorf(errAdminRequiredToDelete)
		}
		isAdmin = true
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

	// handle locked keys
	if info.Locked && !isAdmin {
		if adminChecked {
			return fmt.Errorf(errAdminRequiredToDeleteLocked)
		}
		if !utils.HasMinRole(b.k, "admin", userName, convID) {
			return fmt.Errorf(errAdminRequiredToDeleteLocked)
		}
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

	// this is needed in case a user adds the bot to their own conversation with themself
	if m.Channel.MembersType == keybase.USER && !strings.Contains(teamName, ",") {
		teamName = teamName + "," + teamName
	}

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

	// this is needed in case a user adds the bot to their own conversation with themself
	if m.Channel.MembersType == keybase.USER && !strings.Contains(teamName, ",") {
		teamName = teamName + "," + teamName
	}

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

	// this is needed in case a user adds the bot to their own conversation with themself
	if m.Channel.MembersType == keybase.USER && !strings.Contains(teamName, ",") {
		teamName = teamName + "," + teamName
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

	switch option {
	case "nonadmincreate":
		switch value {
		case "true":
			settings.NonAdminCreate = true
		case "false":
			settings.NonAdminCreate = false
		default:
			return fmt.Errorf(errInvalidValue)
		}
	case "nonadminedit":
		switch value {
		case "true":
			settings.NonAdminEdit = true
		case "false":
			settings.NonAdminEdit = false
		default:
			return fmt.Errorf(errInvalidValue)
		}
	case "nonadmindelete":
		switch value {
		case "true":
			settings.NonAdminDelete = true
		case "false":
			settings.NonAdminDelete = false
		default:
			return fmt.Errorf(errInvalidValue)
		}
	default:
		return fmt.Errorf(errInvalidOption)
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

	// this is needed in case a user adds the bot to their own conversation with themself
	if m.Channel.MembersType == keybase.USER && !strings.Contains(teamName, ",") {
		teamName = teamName + "," + teamName
	}

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

	// this is needed in case a user adds the bot to their own conversation with themself
	if m.Channel.MembersType == keybase.USER && !strings.Contains(teamName, ",") {
		teamName = teamName + "," + teamName
	}

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

	// this is needed in case a user adds the bot to their own conversation with themself
	if m.Channel.MembersType == keybase.USER && !strings.Contains(teamName, ",") {
		teamName = teamName + "," + teamName
	}

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
		// log error, but do not send error message to channel
		b.logError("Unable to fetch key for team %s. -- %v", teamName, err)
		return nil
	}
	_, err = b.k.SendMessageByConvID(convID, info.Value)
	if err != nil {
		b.logError("Error sending reply: %v", err)
	}

	return nil
}
