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
