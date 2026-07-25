package handlers

import (
	"strconv"
	"strings"

	"github.com/Enziofael/nutrigo/bot-telegram/internal/scenarios"
	tg "github.com/Enziofael/nutrigo/bot-telegram/pkg/telegroni"
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

func CallbackQuery_wdky(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	return scenarios.NewUsageRequest(ctx, update)
}

func CallbackQuery_nur(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	if strings.HasPrefix(update.CallbackData(), "nur_c") {
		tgIDstr, _ := strings.CutPrefix(update.CallbackData(), "nur_c")
		tgID, err := strconv.ParseInt(tgIDstr, 10, 64)
		if err != nil {
			boterr := tg.NewBotError(err.Error(), nil)
			return tg.StatusError, tg.NewBotError("Invalid TgID in callback", boterr)
		}
		return scenarios.ConfirmUsageRequest(ctx, update, tgID)
	}

	if strings.HasPrefix(update.CallbackData(), "nur_r") {
		tgIDstr, _ := strings.CutPrefix(update.CallbackData(), "nur_r")
		tgID, err := strconv.ParseInt(tgIDstr, 10, 64)
		if err != nil {
			boterr := tg.NewBotError(err.Error(), nil)
			return tg.StatusError, tg.NewBotError("Invalid TgID in callback", boterr)
		}
		return scenarios.RejectUsageRequest(ctx, update, tgID)
	}

	if strings.HasPrefix(update.CallbackData(), "nur_b") {
		tgIDstr, _ := strings.CutPrefix(update.CallbackData(), "nur_b")
		tgID, err := strconv.ParseInt(tgIDstr, 10, 64)
		if err != nil {
			boterr := tg.NewBotError(err.Error(), nil)
			return tg.StatusError, tg.NewBotError("Invalid TgID in callback", boterr)
		}
		return scenarios.BlockUsageRequest(ctx, update, tgID)
	}

	return tg.StatusError, tg.NewBotError("Invalid data postfix", nil)
}
