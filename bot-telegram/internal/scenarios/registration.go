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
		return tg.StatusError, tg.NewBotErrorf("Can't request sending msg at WeDontKnowYou: %w", err)
	}

	return tg.StatusOK, nil
}

func NewUsageRequest(ctx tg.Context, u tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	if _, err := ctx.Value("client").(*client.Client).CreateUser(ctx.C, u.SentFrom().ID, u.SentFrom().UserName); err != nil {
		return tg.StatusError, tg.NewBotErrorf("Can't get client at NewUsageRequest: %w", err)
	}

	adminMsg, boterr := f.NewUsageRequest(ctx, u)
	if boterr != nil {
		return tg.StatusError, tg.NewBotErrorw("Can't fabricate adminMsg at NewUsageRequest", boterr)
	}
	if _, err := ctx.Bot.Request(adminMsg); err != nil {
		return tg.StatusError, tg.NewBotErrorf("Can't request sending adminMsg at NewUsageRequest: %w", err)
	}

	return YourUsageRequestSend(ctx, u)
}

func YourUsageRequestSend(ctx tg.Context, u tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	msg, boterr := f.YourUsageRequestSend(ctx, u)
	if boterr != nil {
		return tg.StatusError, boterr
	}
	if _, err := ctx.Bot.Send(msg); err != nil {
		return tg.StatusError, tg.NewBotErrorf("Can't send msg at YourUsageRequestSend: %w", err)
	}

	return tg.StatusOK, nil
}

func ConfirmUsageRequest(ctx tg.Context, u tgbotapi.Update, tgID int64) (tg.HandleStatus, *tg.BotError) {
	if _, err := ctx.Value("client").(*client.Client).PatchUserStatus(ctx.C, tgID, models.StatusConfirmed); err != nil {
		return tg.StatusError, tg.NewBotErrorf("Can't get client at ConfirmUsageRequest: %w", err)
	}

	editMsg, boterr := f.ConfirmUsageRequest(ctx, u)
	if boterr != nil {
		return tg.StatusError, tg.NewBotErrorw("Can't fabricate editMsg at ConfirmUsageRequest", boterr)
	}

	if _, err := ctx.Bot.Request(editMsg); err != nil {
		return tg.StatusError, tg.NewBotErrorf("Can't request editMsg at ConfirmUsageRequest: %w", err)
	}

	return NotifyUsageRequestConfirmed(ctx, tgID)
}

func RejectUsageRequest(ctx tg.Context, u tgbotapi.Update, tgID int64) (tg.HandleStatus, *tg.BotError) {
	if err := ctx.Value("client").(*client.Client).DeleteUser(ctx.C, tgID); err != nil {
		return tg.StatusError, tg.NewBotErrorf("Can't get client at RejectUsageRequest: %w", err)
	}

	editMsg, boterr := f.RejectUsageRequest(ctx, u)
	if boterr != nil {
		return tg.StatusError, tg.NewBotErrorw("Can't fabricate editMsg at RejectUsageRequest", boterr)
	}

	if _, err := ctx.Bot.Request(editMsg); err != nil {
		return tg.StatusError, tg.NewBotErrorf("Can't request editMsg at RejectUsageRequest: %w", err)
	}

	return NotifyUsageRequestRejected(ctx, tgID)
}

func BlockUsageRequest(ctx tg.Context, u tgbotapi.Update, tgID int64) (tg.HandleStatus, *tg.BotError) {
	if _, err := ctx.Value("client").(*client.Client).PatchUserStatus(ctx.C, tgID, models.StatusBanned); err != nil {
		return tg.StatusError, tg.NewBotErrorf("Can't get client at BlockUsageRequest: %w", err)
	}

	editMsg, boterr := f.BlockUsageRequest(ctx, u)
	if boterr != nil {
		return tg.StatusError, tg.NewBotErrorw("Can't fabricate editMsg at BlockUsageRequest", boterr)
	}
	if _, err := ctx.Bot.Request(editMsg); err != nil {
		return tg.StatusError, tg.NewBotErrorf("Can't request editMsg at BlockUsageRequest: %w", err)
	}

	return NotifyUsageRequestBlocked(ctx, tgID)
}

func NotifyUsageRequestConfirmed(ctx tg.Context, tgID int64) (tg.HandleStatus, *tg.BotError) {

	notificationMsg, boterr := f.NotifyUsageRequestConfirmed(ctx, tgID)
	if boterr != nil {
		return tg.StatusError, tg.NewBotErrorw("Can't fabricate notificationMsg at NotifyUsageRequestConfirmed", boterr)
	}
	if _, err := ctx.Bot.Send(notificationMsg); err != nil {
		return tg.StatusError, tg.NewBotErrorf("Can't send notificationMsg at NotifyUsageRequestConfirmed: %w", err)
	}

	return tg.StatusOK, nil
}

func NotifyUsageRequestRejected(ctx tg.Context, tgID int64) (tg.HandleStatus, *tg.BotError) {

	notificationMsg, boterr := f.NotifyUsageRequestRejected(ctx, tgID)
	if boterr != nil {
		return tg.StatusError, tg.NewBotErrorw("Can't fabricate notificationMsg at NotifyUsageRequestRejected", boterr)
	}
	if _, err := ctx.Bot.Send(notificationMsg); err != nil {
		return tg.StatusError, tg.NewBotErrorf("Can't send notificationMsg at NotifyUsageRequestRejected: %w", err)
	}

	return tg.StatusOK, nil
}

func NotifyUsageRequestBlocked(ctx tg.Context, tgID int64) (tg.HandleStatus, *tg.BotError) {

	notificationMsg, boterr := f.NotifyUsageRequestBlocked(ctx, tgID)
	if boterr != nil {
		return tg.StatusError, tg.NewBotErrorw("Can't fabricate notificationMsg at NotifyUsageRequestBlocked", boterr)
	}
	if _, err := ctx.Bot.Send(notificationMsg); err != nil {
		return tg.StatusError, tg.NewBotErrorf("Can't send notificationMsg at NotifyUsageRequestBlocked: %w", err)
	}

	return tg.StatusOK, nil
}
