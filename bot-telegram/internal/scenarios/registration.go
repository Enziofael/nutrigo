package scenarios

import (
	client "github.com/Enziofael/nutrigo/bot-telegram/internal/client/v1"
	f "github.com/Enziofael/nutrigo/bot-telegram/internal/factories"
	tg "github.com/Enziofael/nutrigo/bot-telegram/pkg/telegroni"
	models "github.com/Enziofael/nutrigo/shared/models/v1"
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

func WeDontKnowYou(ctx tg.Context, u tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	msg, boterr := f.WeDontKnowYou(ctx, u)
	if boterr != nil {
		return tg.StatusError, boterr
	}
	if _, err := ctx.Bot.Request(msg); err != nil {
		return tg.StatusError, tg.NewBotError(err.Error(), nil)
	}

	return tg.StatusOK, nil
}

func NewUsageRequest(ctx tg.Context, u tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	if _, err := ctx.Value("client").(*client.Client).CreateUser(ctx.C, u); err != nil {
		return tg.StatusError, tg.NewBotError(err.Error(), nil)
	}

	adminMsg, boterr := f.NewUsageRequest(ctx, u)
	if boterr != nil {
		return tg.StatusError, boterr
	}
	if _, err := ctx.Bot.Request(adminMsg); err != nil {
		return tg.StatusError, tg.NewBotError(err.Error(), nil)
	}

	return YourUsageRequestSend(ctx, u)
}

func YourUsageRequestSend(ctx tg.Context, u tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	msg, boterr := f.YourUsageRequestSend(ctx, u)
	if boterr != nil {
		return tg.StatusError, boterr
	}
	if _, err := ctx.Bot.Request(msg); err != nil {
		return tg.StatusError, tg.NewBotError(err.Error(), nil)
	}

	return tg.StatusOK, nil
}

func ConfirmUsageRequest(ctx tg.Context, u tgbotapi.Update, tgID int64) (tg.HandleStatus, *tg.BotError) {
	if _, err := ctx.Value("client").(*client.Client).PatchUserStatus(models.StatusConfirmed, tgID, ctx.C); err != nil {
		return tg.StatusError, tg.NewBotError(err.Error(), nil)
	}
	editMsg, boterr := f.ConfirmUsageRequest(ctx, u)
	if boterr != nil {
		return tg.StatusError, boterr
	}
	if _, err := ctx.Bot.Request(editMsg); err != nil {
		return tg.StatusError, tg.NewBotError(err.Error(), nil)
	}

	return NotifyUsageRequestConfirmed(ctx, tgID)
}

func RejectUsageRequest(ctx tg.Context, u tgbotapi.Update, tgID int64) (tg.HandleStatus, *tg.BotError) {
	if err := ctx.Value("client").(*client.Client).Delete(tgID, ctx.C); err != nil {
		return tg.StatusError, tg.NewBotError(err.Error(), nil)
	}
	editMsg, boterr := f.RejectUsageRequest(ctx, u)
	if boterr != nil {
		return tg.StatusError, boterr
	}
	if _, err := ctx.Bot.Request(editMsg); err != nil {
		return tg.StatusError, tg.NewBotError(err.Error(), nil)
	}

	return NotifyUsageRequestRejected(ctx, tgID)
}

func BlockUsageRequest(ctx tg.Context, u tgbotapi.Update, tgID int64) (tg.HandleStatus, *tg.BotError) {
	if _, err := ctx.Value("client").(*client.Client).PatchUserStatus(models.StatusBanned, tgID, ctx.C); err != nil {
		return tg.StatusError, tg.NewBotError(err.Error(), nil)
	}
	editMsg, boterr := f.BlockUsageRequest(ctx, u)
	if boterr != nil {
		return tg.StatusError, boterr
	}
	if _, err := ctx.Bot.Request(editMsg); err != nil {
		return tg.StatusError, tg.NewBotError(err.Error(), nil)
	}

	return NotifyUsageRequestBlocked(ctx, tgID)
}

func NotifyUsageRequestConfirmed(ctx tg.Context, tgID int64) (tg.HandleStatus, *tg.BotError) {
	notificationMsg, boterr := f.NotifyUsageRequestConfirmed(ctx, tgID)
	if boterr != nil {
		return tg.StatusError, boterr
	}
	if _, err := ctx.Bot.Request(notificationMsg); err != nil {
		return tg.StatusError, tg.NewBotError(err.Error(), nil)
	}

	return tg.StatusOK, nil
}

func NotifyUsageRequestRejected(ctx tg.Context, tgID int64) (tg.HandleStatus, *tg.BotError) {

	notificationMsg, boterr := f.NotifyUsageRequestRejected(ctx, tgID)
	if boterr != nil {
		return tg.StatusError, boterr
	}
	if _, err := ctx.Bot.Request(notificationMsg); err != nil {
		return tg.StatusError, tg.NewBotError(err.Error(), nil)
	}

	return tg.StatusOK, nil
}

func NotifyUsageRequestBlocked(ctx tg.Context, tgID int64) (tg.HandleStatus, *tg.BotError) {

	notificationMsg, boterr := f.NotifyUsageRequestBlocked(ctx, tgID)
	if boterr != nil {
		return tg.StatusError, boterr
	}
	if _, err := ctx.Bot.Request(notificationMsg); err != nil {
		return tg.StatusError, tg.NewBotError(err.Error(), nil)
	}

	return tg.StatusOK, nil
}
