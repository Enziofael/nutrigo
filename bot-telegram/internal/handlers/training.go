package handlers

import (
	//"time"

	"strconv"
	"strings"

	"github.com/Enziofael/nutrigo/bot-telegram/internal/scenarios"
	tg "github.com/Enziofael/nutrigo/bot-telegram/pkg/telegroni"

	//models "github.com/Enziofael/nutrigo/shared/models/v1"
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

func TrainingMenu(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	return scenarios.TrainingMenu(ctx, update)
}

func TrainingMenu_TEP(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	suff, found := strings.CutPrefix(update.CallbackData(), "TEP")
	if !found {
		return tg.StatusError, tg.NewBotError("Hasn't found prefix TEP", nil)
	}
	page, err := strconv.Atoi(suff)
	if err != nil {
		return tg.StatusError, tg.NewBotError(err.Error(), nil)
	}
	return scenarios.TrainingExercises(ctx, update, page)
}

func TrainingMenu_TPP(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	suff, found := strings.CutPrefix(update.CallbackData(), "TPP")
	if !found {
		return tg.StatusError, tg.NewBotError("Hasn't found prefix TPP", nil)
	}
	page, err := strconv.Atoi(suff)
	if err != nil {
		return tg.StatusError, tg.NewBotError(err.Error(), nil)
	}
	return scenarios.TrainingProgramms(ctx, update, page)
}

func TrainingMenuEdit(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	return scenarios.TrainingMenuEdit(ctx, update)
}
