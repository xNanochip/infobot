package run

import (
	"strings"

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
	switch m.Content.TypeName {
	case "text":
		if strings.HasPrefix(m.Content.Text.Body, "!ping") {
			_, err := b.k.ReactByConvID(m.ConvID, m.Id, "PONG!")
			if err != nil {
				b.log_error("Error sending reaction: %v", err)
			}
			return
		}
		if strings.HasPrefix(m.Content.Text.Body, "!ding") {
			_, err := b.k.ReactByConvID(m.ConvID, m.Id, "DONG!")
			if err != nil {
				b.log_error("Error sending reaction: %v", err)
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
