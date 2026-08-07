package scenarios

import (
	//client "github.com/Enziofael/nutrigo/bot-telegram/internal/client/v1"

	"github.com/Enziofael/nutrigo/bot-telegram/internal/client/v1"
	f "github.com/Enziofael/nutrigo/bot-telegram/internal/factories"
	tg "github.com/Enziofael/nutrigo/bot-telegram/pkg/telegroni"
	"github.com/Enziofael/nutrigo/shared/models/v1"

	//models "github.com/Enziofael/nutrigo/shared/models/v1"
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

// Scenarios are systemised by its names:
// Module_Action_Name(_Type)
//
// This file is for Training module scenarios
//
// Actions:
// Send - send new message
// Edit - edit existing message with new content
// Input - process user's input message (text, media, etc.)
// 		   ! Input should be validated in handler, not in a scenario
//		   ! Scenario should not work with updates, all the values are passed by handler

// ===================================================================
// TrainingMenu
// ===================================================================

// Navgraph:
// Training module (prefix training)
// └── Menu
//	   ├── Exercises
//     │   ├── Exercise
//     │   ├── Search exercise
//     │   └── New exercise form
//     └── Programms

func Training_Send_MainMenu(ctx tg.Context, chatID int64) (tg.HandleStatus, *tg.BotError) {

	// Creating msg
	msg, boterr := f.Training_Send_MainMenu(chatID)
	if boterr != nil {
		return tg.StatusError, tg.NewBotErrorw("Can't fabricate training menu", boterr)
	}

	// Sending msg
	r, err := ctx.Bot.Send(msg)
	if err != nil {
		return tg.StatusError, tg.NewBotErrorf("Can't send training menu: %w", err)
	}

	// Patching context to new message
	client := ctx.Value("client").(*client.Client)
	_, clierr := client.PatchUserContext(ctx.C, chatID, "Training_MainMenu", models.ContextData{MessageID: r.MessageID})
	if clierr != nil {
		return tg.StatusError, tg.NewBotErrorf("Can't patch context to Training_MainMenu for user %d: %w", chatID, clierr)
	}

	return tg.StatusOK, nil
}

func Training_Edit_MainMenu(ctx tg.Context, chatID int64, messageID int) (tg.HandleStatus, *tg.BotError) {

	// Patching context
	client := ctx.Value("client").(*client.Client)
	_, clierr := client.PatchUserContext(ctx.C, chatID, "Training_MainMenu", models.ContextData{MessageID: messageID})
	if clierr != nil {
		return tg.StatusError, tg.NewBotErrorf("Can't patch context to Training_MainMenu for user %d: %w", chatID, clierr)
	}

	// Creating editMsg
	editMsg, boterr := f.Training_Edit_MainMenu(chatID, messageID)
	if boterr != nil {
		return tg.StatusError, boterr
	}

	// Requesting editing
	if _, err := ctx.Bot.Request(editMsg); err != nil {
		return tg.StatusError, tg.NewBotError(err.Error())
	}

	return tg.StatusOK, nil
}

// ===================================================================
// Exercises
// ===================================================================

const exercisesPerPage int = 4

func Training_Edit_Exercises(ctx tg.Context, chatID int64, messageID int, page int) (tg.HandleStatus, *tg.BotError) {

	// Patching context
	client := ctx.Value("client").(*client.Client)
	if _, err := client.PatchUserContext(ctx.C, chatID, "Training_Exercises", models.ContextData{MessageID: messageID}); err != nil {
		return tg.StatusError, tg.NewBotErrorf("Can't patch context to Training_Exercises: %w", err)
	}

	// Getting page count
	count, err := client.CountExercises(ctx.C, chatID)
	if err != nil {
		return tg.StatusError, tg.NewBotErrorf("Can't get exercises count at Training_Edit_Exercises: %w", err)
	}
	pageCount := 0
	if count > 0 {
		pageCount = (count + exercisesPerPage - 1) / exercisesPerPage
	}
	offset := page * exercisesPerPage

	// Getting page
	exercises, err := client.ListExercises(ctx.C, chatID, offset, exercisesPerPage, "rating", "desc")
	if err != nil {
		return tg.StatusError, tg.NewBotErrorf("Can't get exercises at Training_Edit_Exercises: %w", err)
	}

	// Creating editMsg
	editMsg, boterr := f.Training_Edit_Exercises(chatID,
		messageID, page, pageCount, exercises)
	if boterr != nil {
		return tg.StatusError, tg.NewBotErrorw("Can't fabricate editMsg at Training_Edit_Exercises", boterr)
	}

	// Requesting editing
	if _, err := ctx.Bot.Request(editMsg); err != nil {
		return tg.StatusError, tg.NewBotErrorf("Can't request editMsg at Training_Edit_Exercises: %w", err)
	}

	return tg.StatusOK, nil
}

func Training_Edit_Exercises_CreateForm(ctx tg.Context, chatID int64, contextData models.ContextData) (tg.HandleStatus, *tg.BotError) {
	if contextData.Focus < 0 {
		contextData.Focus = 2
	}
	if contextData.Focus > 2 {
		contextData.Focus = 0
	}

	// Patching context
	client := ctx.Value("client").(*client.Client)
	if _, err := client.PatchUserContext(ctx.C, chatID, "Training_Exercises_CreateForm", contextData); err != nil {
		return tg.StatusError, tg.NewBotErrorf("Can't patch context to Training_Exercises_CreateForm: %w", err)
	}

	// Creating editMsg
	editMsg, boterr := f.Training_Edit_Exercises_CreateForm(chatID, contextData)
	if boterr != nil {
		return tg.StatusError, tg.NewBotErrorw("Can't fabricate editMsg at Training_Edit_Exercises_CreateForm", boterr)
	}

	// Requesting editing
	if _, err := ctx.Bot.Request(editMsg); err != nil {
		return tg.StatusError, tg.NewBotErrorf("Can't request editMsg at Training_Edit_Exercises_CreateForm: %w", err)
	}

	return tg.StatusOK, nil
}

func Training_Input_Exercises_CreateForm(ctx tg.Context, chatID int64, contextData models.ContextData, input string) (tg.HandleStatus, *tg.BotError) {

	// Apply input
	switch contextData.Focus {
	case 0:
		contextData.Values["Name"] = input
	case 1:
		contextData.Values["Description"] = input
	case 2:
		contextData.Values["Technique"] = input
	}
	contextData.Focus = contextData.Focus + 1
	if contextData.Focus < 0 {
		contextData.Focus = 2
	}
	if contextData.Focus > 2 {
		contextData.Focus = 0
	}

	// Patching context with applyed input
	client := ctx.Value("client").(*client.Client)
	if _, err := client.PatchUserContext(ctx.C, chatID, "Training_Exercises_CreateForm", contextData); err != nil {
		return tg.StatusError, tg.NewBotErrorf("Can't patch context to Training_Exercises_CreateForm: %w", err)
	}

	// Editing form message
	return Training_Edit_Exercises_CreateForm(ctx, chatID, contextData)
}

func Training_Input_Exercises_CreateForm_Save(ctx tg.Context, chatID int64, contextData models.ContextData) (tg.HandleStatus, *tg.BotError) {
	form := f.CreateExerciseFormValues{
		Name:        contextData.Values["Name"],
		Description: contextData.Values["Description"],
		Technique:   contextData.Values["Technique"],
		Unit:        contextData.Values["Unit"],
	}

	client := ctx.Value("client").(*client.Client)

	e, err := client.CreateExercise(ctx.C, chatID, form.Name)
	if err != nil {
		return tg.StatusError, tg.NewBotErrorf("Can't create exercise at Training_Input_Exercises_CreateForm_Save: %w", err)
	}

	patchReq := models.ExercisePatchRequest{
		ID:          e.ID,
		Description: &form.Description,
		Technique:   &form.Technique,
		WeightUnit:  &form.Unit,
	}

	if _, err := client.PatchExercise(ctx.C, patchReq); err != nil {
		_ = client.DeleteExercise(ctx.C, patchReq.ID)
		return tg.StatusError, tg.NewBotErrorf("Can't patch exercise at Training_Input_Exercises_CreateForm_Save: %w", err)
	}

	return Training_Edit_Exercises(ctx, chatID, contextData.MessageID, 0)
}

func Training_Edit_Exercise(ctx tg.Context, chatID int64, messageID int, exerciseID int64) (tg.HandleStatus, *tg.BotError) {
	cd := models.ContextData{
		MessageID: messageID,
		Focus:     int(exerciseID),
	}

	// Patching context
	client := ctx.Value("client").(*client.Client)
	if _, err := client.PatchUserContext(ctx.C, chatID, "Training_Exercise", cd); err != nil {
		return tg.StatusError, tg.NewBotErrorf("Can't patch context to Training_Exercise: %w", err)
	}

	/*
		exercise, err := client.GetExercise(ctx.C, exerciseID)
		if err != nil {
			return tg.StatusError, tg.NewBotErrorf("Can't get exercise at Training_Edit_Exercise: %w", err)
		}

		// Creating editMsg
		editMsg, boterr := f.Training_Edit_Exercise(exercise)
		if boterr != nil {
			return tg.StatusError, tg.NewBotErrorw("Can't fabricate editMsg at Training_Edit_Exercises_CreateForm", boterr)
		}

		// Requesting editing
		if _, err := ctx.Bot.Request(editMsg); err != nil {
			return tg.StatusError, tg.NewBotErrorf("Can't request editMsg at Training_Edit_Exercises_CreateForm: %w", err)
		}
	*/
	return tg.StatusOK, nil
}

// ===================================================================
// Programms
// ===================================================================

func Training_Edit_Programms(ctx tg.Context, update tgbotapi.Update, page int) (tg.HandleStatus, *tg.BotError) {
	client := ctx.Value("client").(*client.Client)
	client.PatchUserContext(ctx.C, update.SentFrom().ID, "TrainingProgramms", models.ContextData{MessageID: update.CallbackQuery.Message.MessageID})

	pageCount := 1
	programms := &[]models.Exercise{
		{Name: "Программа 1", ID: 1},
		{Name: "Программа 2", ID: 2},
		{Name: "Программа 3", ID: 3},
		{Name: "Программа 4", ID: 4},
		{Name: "Программа 5", ID: 5},
	}
	editMsg, boterr := f.Training_Edit_Programms(update.SentFrom().ID,
		update.CallbackQuery.Message.MessageID, page, pageCount, programms)
	if boterr != nil {
		return tg.StatusError, boterr
	}
	if _, err := ctx.Bot.Request(editMsg); err != nil {
		return tg.StatusError, tg.NewBotErrorf("Can't request editMsg at Training_Edit_Programms: %w", err)
	}

	return tg.StatusOK, nil
}
