package handlers

import (
	//"time"

	"strconv"
	"strings"

	"github.com/Enziofael/nutrigo/bot-telegram/internal/client/v1"
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

func Training_Input_Exercises(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	if update.Message == nil || update.Message.Text == "" {
		return tg.StatusError, tg.NewBotError("Invalid routing to Training_Input_Exercises_CreateForm. Expected text message")
	}

	contextData := ctx.Value("user").(*models.User).ContextData

	if status, err := scenarios.Training_Input_Exercises(ctx, update.Message.From.ID, contextData, update.Message.Text); status == tg.StatusError || err != nil {
		return tg.StatusError, err
	}

	return scenarios.Training_Edit_Exercises_Search(ctx, update.Message.From.ID, 0)
}

func Training_Edit_Exercises_Search(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	if update.CallbackQuery == nil {
		return tg.StatusError, tg.NewBotError("Invalid routing to Training_Edit_Exercises_Search. Expected CallbackQuery")
	}

	suff, found := strings.CutPrefix(update.CallbackData(), "TES")
	if !found {
		return tg.StatusError, tg.NewBotError("Hasn't found prefix TES")
	}
	page, err := strconv.Atoi(suff)
	if err != nil {
		return tg.StatusError, tg.NewBotErrorf("Can't parse page's number: %w", err)
	}

	return scenarios.Training_Edit_Exercises_Search(ctx, update.SentFrom().ID, page)
}

func Training_Edit_Exercises_CreateForm(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	if update.CallbackQuery == nil {
		return tg.StatusError, tg.NewBotError("Invalid routing to Training_Edit_Exercises_CreateForm. Expected CallbackQuery")
	}

	context := ctx.Value("user").(*models.User).Context
	contextData := ctx.Value("user").(*models.User).ContextData
	if contextData.MessageID == 0 || len(contextData.Values) == 0 || context != "Training_Exercises_CreateForm" {
		contextData = models.ContextData{MessageID: update.CallbackQuery.Message.MessageID, Values: make(map[string]string)}
		contextData.Values["Unit"] = string(models.WeightUnit_Kg)
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

func Training_Edit_Exercises_CreateForm_KG(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	if update.CallbackQuery == nil {
		return tg.StatusError, tg.NewBotError("Invalid routing to Training_Edit_Exercises_CreateForm_KG. Expected CallbackQuery")
	}

	contextData := ctx.Value("user").(*models.User).ContextData
	contextData.Values["Unit"] = string(models.WeightUnit_Kg)

	return scenarios.Training_Edit_Exercises_CreateForm(ctx, update.SentFrom().ID, contextData)
}

func Training_Edit_Exercises_CreateForm_LBS(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	if update.CallbackQuery == nil {
		return tg.StatusError, tg.NewBotError("Invalid routing to Training_Edit_Exercises_CreateForm_LBS. Expected CallbackQuery")
	}

	contextData := ctx.Value("user").(*models.User).ContextData
	contextData.Values["Unit"] = string(models.WeightUnit_Lbs)

	return scenarios.Training_Edit_Exercises_CreateForm(ctx, update.SentFrom().ID, contextData)
}

func Training_Edit_Exercises_CreateForm_MIXED(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	if update.CallbackQuery == nil {
		return tg.StatusError, tg.NewBotError("Invalid routing to Training_Edit_Exercises_CreateForm_MIXED. Expected CallbackQuery")
	}

	contextData := ctx.Value("user").(*models.User).ContextData
	contextData.Values["Unit"] = string(models.WeightUnit_Mixed)

	return scenarios.Training_Edit_Exercises_CreateForm(ctx, update.SentFrom().ID, contextData)
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
	if _, err := models.SanitizeUrl(contextData.Values["Technique"]); err != nil {
		return tg.StatusWarn, tg.NewBotErrorf("Can't save exercise with invalid url: %s", contextData.Values["Technique"])
	}

	return scenarios.Training_Input_Exercises_CreateForm_Save(ctx, update.SentFrom().ID, contextData)
}

func Training_Edit_Exercises_EditForm(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	if update.CallbackQuery == nil {
		return tg.StatusError, tg.NewBotError("Invalid routing to Training_Edit_Exercises_CreateForm. Expected CallbackQuery")
	}

	suff, found := strings.CutPrefix(update.CallbackData(), "TEE")
	if !found {
		return tg.StatusError, tg.NewBotError("Hasn't found prefix TEE")
	}
	exerciseID, err := strconv.ParseInt(suff, 10, 64)
	if err != nil {
		return tg.StatusError, tg.NewBotErrorf("Can't parse exerciseID: %w", err)
	}

	client := ctx.Value("client").(*client.Client)
	exercise, err := client.GetExercise(ctx.C, exerciseID)
	if err != nil {
		return tg.StatusError, tg.NewBotErrorf("Can't patch context to Training_Exercises_EditForm: %w", err)
	}
	s := ""
	if exercise.Description == nil {
		exercise.Description = &s
	}
	if exercise.Technique == nil {
		exercise.Technique = &s
	}

	context := ctx.Value("user").(*models.User).Context
	contextData := ctx.Value("user").(*models.User).ContextData
	if contextData.MessageID == 0 || len(contextData.Values) == 0 || context != "Training_Exercises_EditForm" {
		contextData = models.ContextData{MessageID: update.CallbackQuery.Message.MessageID, Values: make(map[string]string)}
		contextData.Values["Unit"] = string(exercise.WeightUnit)
		contextData.Values["ExerciseID"] = suff
		contextData.Values["Name"] = exercise.Name
		contextData.Values["Description"] = *exercise.Description
		contextData.Values["Technique"] = *exercise.Technique
	}

	return scenarios.Training_Edit_Exercises_EditForm(ctx, update.CallbackQuery.From.ID, contextData)
}

func Training_Edit_Exercises_EditForm_UP(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	if update.CallbackQuery == nil {
		return tg.StatusError, tg.NewBotError("Invalid routing to Training_Edit_Exercises_EditForm_UP. Expected CallbackQuery")
	}

	contextData := ctx.Value("user").(*models.User).ContextData
	contextData.Focus = contextData.Focus - 1
	if contextData.Focus < 0 {
		contextData.Focus = 2
	}
	if contextData.Focus > 2 {
		contextData.Focus = 0
	}

	return scenarios.Training_Edit_Exercises_EditForm(ctx, update.CallbackQuery.From.ID, contextData)
}

func Training_Edit_Exercises_EditForm_DOWN(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	if update.CallbackQuery == nil {
		return tg.StatusError, tg.NewBotError("Invalid routing to Training_Edit_Exercises_EditForm_DOWN. Expected CallbackQuery")
	}

	contextData := ctx.Value("user").(*models.User).ContextData
	contextData.Focus = contextData.Focus + 1
	if contextData.Focus < 0 {
		contextData.Focus = 2
	}
	if contextData.Focus > 2 {
		contextData.Focus = 0
	}

	return scenarios.Training_Edit_Exercises_EditForm(ctx, update.CallbackQuery.From.ID, contextData)
}

func Training_Edit_Exercises_EditForm_KG(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	if update.CallbackQuery == nil {
		return tg.StatusError, tg.NewBotError("Invalid routing to Training_Edit_Exercises_EditForm_KG. Expected CallbackQuery")
	}

	contextData := ctx.Value("user").(*models.User).ContextData
	contextData.Values["Unit"] = string(models.WeightUnit_Kg)

	return scenarios.Training_Edit_Exercises_EditForm(ctx, update.SentFrom().ID, contextData)
}

func Training_Edit_Exercises_EditForm_LBS(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	if update.CallbackQuery == nil {
		return tg.StatusError, tg.NewBotError("Invalid routing to Training_Edit_Exercises_EditForm_LBS. Expected CallbackQuery")
	}

	contextData := ctx.Value("user").(*models.User).ContextData
	contextData.Values["Unit"] = string(models.WeightUnit_Lbs)

	return scenarios.Training_Edit_Exercises_EditForm(ctx, update.SentFrom().ID, contextData)
}

func Training_Edit_Exercises_EditForm_MIXED(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	if update.CallbackQuery == nil {
		return tg.StatusError, tg.NewBotError("Invalid routing to Training_Edit_Exercises_EditForm_MIXED. Expected CallbackQuery")
	}

	contextData := ctx.Value("user").(*models.User).ContextData
	contextData.Values["Unit"] = string(models.WeightUnit_Mixed)

	return scenarios.Training_Edit_Exercises_EditForm(ctx, update.SentFrom().ID, contextData)
}

func Training_Input_Exercises_EditForm(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	if update.Message == nil || update.Message.Text == "" {
		return tg.StatusError, tg.NewBotError("Invalid routing to Training_Input_Exercises_EditForm. Expected text message")
	}

	contextData := ctx.Value("user").(*models.User).ContextData

	return scenarios.Training_Input_Exercises_EditForm(ctx, update.SentFrom().ID, contextData, update.Message.Text)
}

func Training_Edit_Exercises_EditForm_SAVE(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	if update.CallbackQuery == nil {
		return tg.StatusError, tg.NewBotError("Invalid routing to Training_Edit_Exercise_EditForm_SAVE. Expected CallbackQuery")
	}

	contextData := ctx.Value("user").(*models.User).ContextData
	if contextData.Values["Name"] == "" {
		return tg.StatusWarn, tg.NewBotError("Can't save exercise with empty name")
	}
	if _, err := models.SanitizeUrl(contextData.Values["Technique"]); err != nil {
		return tg.StatusWarn, tg.NewBotErrorf("Can't save exercise with invalid url: %s", contextData.Values["Technique"])
	}

	return scenarios.Training_Input_Exercises_EditForm_Save(ctx, update.SentFrom().ID, contextData)
}

func Training_Edit_Exercise_DeleteForm(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	if update.CallbackQuery == nil {
		return tg.StatusError, tg.NewBotError("Invalid routing to Training_Edit_Exercise_DeleteForm. Expected CallbackQuery")
	}

	suff, found := strings.CutPrefix(update.CallbackData(), "TED")
	if !found {
		return tg.StatusError, tg.NewBotError("Hasn't found prefix TED")
	}
	exerciseID, err := strconv.ParseInt(suff, 10, 64)
	if err != nil {
		return tg.StatusError, tg.NewBotErrorf("Can't parse exerciseID: %w", err)
	}

	client := ctx.Value("client").(*client.Client)
	exercise, err := client.GetExercise(ctx.C, exerciseID)
	if err != nil {
		return tg.StatusError, tg.NewBotErrorf("Can't patch context to Training_Exercises_EditForm: %w", err)
	}

	return scenarios.Training_Edit_Exercise_DeleteForm(ctx, update.SentFrom().ID, update.CallbackQuery.Message.MessageID, exercise)
}

func Training_Edit_Exercise_DeleteForm_CONFIRM(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	if update.CallbackQuery == nil {
		return tg.StatusError, tg.NewBotError("Invalid routing to Training_Edit_Exercise_DeleteForm_CONFIRM. Expected CallbackQuery")
	}

	suff, found := strings.CutPrefix(update.CallbackData(), "TEDC")
	if !found {
		return tg.StatusError, tg.NewBotError("Hasn't found prefix TEDC")
	}
	exerciseID, err := strconv.ParseInt(suff, 10, 64)
	if err != nil {
		return tg.StatusError, tg.NewBotErrorf("Can't parse exerciseID: %w", err)
	}

	return scenarios.Training_Edit_Exercise_DeleteForm_CONFIRM(ctx, update.SentFrom().ID, update.CallbackQuery.Message.MessageID, exerciseID)
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
