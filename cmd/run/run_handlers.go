package run

import (
	"strings"

	"samhofi.us/x/infobot/pkg/utils"
	"samhofi.us/x/keybase/v2"
	"samhofi.us/x/keybase/v2/types/chat1"
	"samhofi.us/x/keybase/v2/types/stellar1"
)

func (b *bot) registerHandlers() {
	b.logDebug("Registering handlers")

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
		userName = m.Sender.Username
		convID   = m.ConvID
	)

	if userName == b.k.Username {
		return
	}

	switch m.Content.TypeName {
	case "text":
		if strings.HasPrefix(m.Content.Text.Body, "!info add ") {
			if err := b.cmdInfoAdd(m); err != nil {
				b.k.ReplyByConvID(convID, m.Id, "Error: %v", err)
			}
			return
		}
		if strings.HasPrefix(m.Content.Text.Body, "!info edit ") {
			if err := b.cmdInfoEdit(m); err != nil {
				b.k.ReplyByConvID(convID, m.Id, "Error: %v", err)
			}
			return
		}
		if strings.HasPrefix(m.Content.Text.Body, "!info delete ") {
			if err := b.cmdInfoDelete(m); err != nil {
				b.k.ReplyByConvID(convID, m.Id, "Error: %v", err)
			}
			return
		}
		if strings.HasPrefix(m.Content.Text.Body, "!info read ") {
			if err := b.cmdInfoRead(m); err != nil {
				b.k.ReplyByConvID(convID, m.Id, "Error: %v", err)
			}
			return
		}
		if strings.HasPrefix(m.Content.Text.Body, "!info audit ") {
			if err := b.cmdInfoAudit(m); err != nil {
				b.k.ReplyByConvID(convID, m.Id, "Error: %v", err)
			}
			return
		}
		if strings.HasPrefix(m.Content.Text.Body, "!info set ") {
			if err := b.cmdInfoSet(m); err != nil {
				b.k.ReplyByConvID(convID, m.Id, "Error: %v", err)
			}
			return
		}
		if strings.HasPrefix(m.Content.Text.Body, "!info settings") {
			if err := b.cmdInfoSettings(m); err != nil {
				b.k.ReplyByConvID(convID, m.Id, "Error: %v", err)
			}
			return
		}
		if strings.HasPrefix(m.Content.Text.Body, "!info keys") {
			if err := b.cmdInfoKeys(m); err != nil {
				b.k.ReplyByConvID(convID, m.Id, "Error: %v", err)
			}
			return
		}
		if utils.StringInSlice(b.k.Username, m.AtMentionUsernames) {
			if err := b.cmdAtMention(m); err != nil {
				b.k.ReplyByConvID(convID, m.Id, "Error: %v", err)
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
