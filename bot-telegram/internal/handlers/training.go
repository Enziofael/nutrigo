package handlers

import (
	//"time"

	"strconv"
	"strings"

	"github.com/Enziofael/nutrigo/bot-telegram/internal/scenarios"
	tg "github.com/Enziofael/nutrigo/bot-telegram/pkg/telegroni"
	"github.com/Enziofael/nutrigo/shared/models/v1"

	//models "github.com/Enziofael/nutrigo/shared/models/v1"
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

func Training_Send_MainMenu(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	if update.SentFrom() == nil {
		return tg.StatusError, tg.NewBotError("Invalid routing to Training_Send_MainMenu. Can't get user's tgID")
	}

	return scenarios.Training_Send_MainMenu(ctx, update.SentFrom().ID)
}

func Training_Edit_MainMenu(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	if update.CallbackQuery == nil {
		return tg.StatusError, tg.NewBotError("Invalid routing to Training_Edit_MainMenu. Expected CallbackQuery")
	}

	return scenarios.Training_Edit_MainMenu(ctx, update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID)
}

func Training_Edit_Exercises(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	if update.CallbackQuery == nil {
		return tg.StatusError, tg.NewBotError("Invalid routing to Training_Edit_Exercises. Expected CallbackQuery")
	}

	suff, found := strings.CutPrefix(update.CallbackData(), "TEP")
	if !found {
		return tg.StatusError, tg.NewBotError("Hasn't found prefix TEP")
	}
	page, err := strconv.Atoi(suff)
	if err != nil {
		return tg.StatusError, tg.NewBotErrorf("Can't parse page's number: %w", err)
	}

	return scenarios.Training_Edit_Exercises(ctx, update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, page)
}

func Training_Edit_Exercises_CreateForm(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	if update.CallbackQuery == nil {
		return tg.StatusError, tg.NewBotError("Invalid routing to Training_Edit_Exercises_CreateForm. Expected CallbackQuery")
	}

	contextData := ctx.Value("user").(*models.User).ContextData
	if contextData.MessageID == 0 || len(contextData.Values) == 0 {
		contextData = models.ContextData{MessageID: update.CallbackQuery.Message.MessageID, Values: make(map[string]string)}
	}

	return scenarios.Training_Edit_Exercises_CreateForm(ctx, update.CallbackQuery.From.ID, contextData)
}

func Training_Edit_Exercises_CreateForm_UP(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	if update.CallbackQuery == nil {
		return tg.StatusError, tg.NewBotError("Invalid routing to Training_Edit_Exercises_CreateForm. Expected CallbackQuery")
	}

	contextData := ctx.Value("user").(*models.User).ContextData
	contextData.Focus = contextData.Focus - 1
	if contextData.Focus < 0 {
		contextData.Focus = 2
	}
	if contextData.Focus > 2 {
		contextData.Focus = 0
	}

	return scenarios.Training_Edit_Exercises_CreateForm(ctx, update.CallbackQuery.From.ID, contextData)
}

func Training_Edit_Exercises_CreateForm_DOWN(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	if update.CallbackQuery == nil {
		return tg.StatusError, tg.NewBotError("Invalid routing to Training_Edit_Exercises_CreateForm. Expected CallbackQuery")
	}

	contextData := ctx.Value("user").(*models.User).ContextData
	contextData.Focus = contextData.Focus + 1
	if contextData.Focus < 0 {
		contextData.Focus = 2
	}
	if contextData.Focus > 2 {
		contextData.Focus = 0
	}

	return scenarios.Training_Edit_Exercises_CreateForm(ctx, update.CallbackQuery.From.ID, contextData)
}

func Training_Input_Exercises_CreateForm(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	if update.Message == nil || update.Message.Text == "" {
		return tg.StatusError, tg.NewBotError("Invalid routing to Training_Input_Exercises_CreateForm. Expected text message")
	}

	contextData := ctx.Value("user").(*models.User).ContextData

	return scenarios.Training_Input_Exercises_CreateForm(ctx, update.SentFrom().ID, contextData, update.Message.Text)
}

func Training_Edit_Exercises_CreateForm_SAVE(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	if update.CallbackQuery == nil {
		return tg.StatusError, tg.NewBotError("Invalid routing to Training_Edit_Exercises_CreateForm_SAVE. Expected CallbackQuery")
	}

	contextData := ctx.Value("user").(*models.User).ContextData
	if contextData.Values["Name"] == "" {
		return tg.StatusWarn, tg.NewBotError("Can't save exercise with empty name")
	}

	return scenarios.Training_Input_Exercises_CreateForm_Save(ctx, update.SentFrom().ID, contextData)
}

func Training_Edit_Exercise(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	if update.CallbackQuery == nil {
		return tg.StatusError, tg.NewBotError("Invalid routing to Training_Edit_Exercise. Expected CallbackQuery")
	}

	suff, found := strings.CutPrefix(update.CallbackData(), "TEI")
	if !found {
		return tg.StatusError, tg.NewBotError("Hasn't found prefix TEI")
	}
	exerciseID, err := strconv.ParseInt(suff, 10, 64)
	if err != nil {
		return tg.StatusError, tg.NewBotErrorf("Can't parse exercise's ID: %w", err)
	}

	return scenarios.Training_Edit_Exercise(ctx, update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, exerciseID)
}

func Training_Edit_Programms(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	if update.CallbackQuery == nil {
		return tg.StatusError, tg.NewBotError("Invalid routing to Training_Edit_Programms. Expected CallbackQuery")
	}

	suff, found := strings.CutPrefix(update.CallbackData(), "TPP")
	if !found {
		return tg.StatusError, tg.NewBotError("Hasn't found prefix TPP")
	}
	page, err := strconv.Atoi(suff)
	if err != nil {
		return tg.StatusError, tg.NewBotErrorf("Can't parse page's number: %w", err)
	}
	return scenarios.Training_Edit_Programms(ctx, update, page)
}
