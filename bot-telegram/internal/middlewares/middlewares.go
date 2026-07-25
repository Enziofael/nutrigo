package middlewares

import (
	"github.com/Enziofael/nutrigo/bot-telegram/internal/scenarios"
	tg "github.com/Enziofael/nutrigo/bot-telegram/pkg/telegroni"
	models "github.com/Enziofael/nutrigo/shared/models/v1"
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

func VerifyUsagePermission(ctx tg.Context, update tgbotapi.Update, next tg.HandlerFunc) (status tg.HandleStatus, err *tg.BotError) {
	us := ctx.Value("user").(*models.User).Status
	switch us {
	case "":
		return scenarios.WeDontKnowYou(ctx, update)
	case models.StatusRequested:
		return scenarios.YourUsageRequestSend(ctx, update)
	case models.StatusRestricted:
		return scenarios.NotifyUsageRequestRejected(ctx, update.SentFrom().ID)
	case models.StatusBanned:
		return scenarios.NotifyUsageRequestBlocked(ctx, update.SentFrom().ID)
	}
	return next(ctx, update)
}

func DeleteSourceMessage(ctx tg.Context, update tgbotapi.Update, next tg.HandlerFunc) (status tg.HandleStatus, err *tg.BotError) {
	var messageID int
	switch {
	case update.Message != nil:
		messageID = update.Message.MessageID
	case update.CallbackQuery != nil:
		messageID = update.CallbackQuery.Message.MessageID
	default:
		return next(ctx, update)
	}
	del := tgbotapi.NewDeleteMessage(update.SentFrom().ID, messageID)
	if _, err := ctx.Bot.Request(del); err != nil {
		return tg.StatusError, tg.NewBotError(err.Error(), nil)
	}
	return next(ctx, update)
}
