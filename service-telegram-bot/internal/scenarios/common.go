package scenarios

import (
	f "github.com/Enziofael/nutrigo/bot-telegram/internal/factories"
	tg "github.com/Enziofael/nutrigo/bot-telegram/pkg/telegroni"
)

func Common_Send_ReplyMenu(ctx tg.Context, tgID int64) (tg.HandleStatus, *tg.BotError) {
	msg, boterr := f.Common_Send_ReplyMenu(tgID)
	if boterr != nil {
		return tg.StatusError, tg.NewBotErrorw("Can't fabricate ReplyMenu", boterr)
	}

	// Sending msg
	_, err := ctx.Bot.Send(msg)
	if err != nil {
		return tg.StatusError, tg.NewBotErrorf("Can't send msg in Common_Send_ReplyMenu: %w", err)
	}

	return tg.StatusOK, nil
}
