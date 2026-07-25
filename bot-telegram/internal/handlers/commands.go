package handlers

import (
	"time"

	"github.com/Enziofael/nutrigo/bot-telegram/internal/factories"
	tg "github.com/Enziofael/nutrigo/bot-telegram/pkg/telegroni"
	v1 "github.com/Enziofael/nutrigo/shared/models/v1"
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

func CommandStartHandler(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
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

	ctx.Bot.Request(msg)

	return tg.StatusOK, nil
}

func OkHandler(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	return tg.StatusOK, nil
}
func WarnHandler(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	return tg.StatusWarn, nil
}
func ErrHandler(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	return tg.StatusError, nil
}
func FallHandler(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	return tg.StatusFallback, nil
}
func CustomHandler(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	return tg.NewStatus(0b11000000, "CUST"), nil
}

func ShutdownHandler(ctx tg.Context, udpate tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	isAdmin := ctx.Value("user").(*v1.User).Status == "admin"
	if isAdmin {
		go func() {
			time.Sleep(5 * time.Second)
			ctx.Value(tg.ContextKey_Bot).(*tgbotapi.BotAPI).StopReceivingUpdates()
		}()
		return tg.StatusOK, nil
	}
	return tg.StatusWarn, nil
}

