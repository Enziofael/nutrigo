package handlers

import (
	"context"

	"github.com/Enziofael/nutrigo/bot-telegram/internal/factories"
	tg "github.com/Enziofael/nutrigo/bot-telegram/pkg/telegroni"
	v1 "github.com/Enziofael/nutrigo/shared/models/v1"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CommandStartHandler(ctx context.Context, update tgbotapi.Update) (string, *tg.BotError) {
	user, ok := ctx.Value("user").(*v1.User)
	if !ok || user == nil {
		user = &v1.User{
			TgID:   update.Message.From.ID,
			TgTag:  update.Message.From.UserName,
			Status: "",
		}
	}

	var msg tgbotapi.MessageConfig
	switch user.Status {
	case v1.StatusAdmin, v1.StatusConfirmed:
		msg = factories.MainMenu(update, user)
	case v1.StatusBanned:
		msg = factories.YouWasBanned(update, user)
	case v1.StatusRequested:
		msg = factories.PermissionRequested(update, user)
	case v1.StatusRestricted:
		msg = factories.YouWasRestricted(update, user)
	default:
		msg = factories.SuggestPermissionRequest(update, user)
	}

	ctx.Value("bot").(*tgbotapi.BotAPI).Send(msg)

	return tg.StatusOK, nil
}
