package handlers

import (
	"time"

	"github.com/Enziofael/nutrigo/bot-telegram/internal/scenarios"
	tg "github.com/Enziofael/nutrigo/bot-telegram/pkg/telegroni"
	v1 "github.com/Enziofael/nutrigo/shared/models/v1"
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

func CommandStartHandler(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	if update.Message != nil && update.Message.IsCommand() {
		return scenarios.Common_Send_ReplyMenu(ctx, update.SentFrom().ID)
	}
	return tg.StatusError, tg.NewBotError("Invalid routing to CommandStartHandler. Expected Command")
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
