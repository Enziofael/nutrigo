package factories

import (
	"fmt"
	"strconv"
	"strings"

	tg "github.com/Enziofael/nutrigo/bot-telegram/pkg/telegroni"
	"github.com/Enziofael/nutrigo/shared/models/v1"
	a "github.com/OvyFlash/telegram-bot-api"
)

// ===============================
// MainMenu
// ===============================

func Training_Send_MainMenu(chatID int64) (a.MessageConfig, *tg.BotError) {
	// Rendering text
	text, err := tmplManager.RenderHTML("TrainingMenu", nil)
	if err != nil {
		return a.MessageConfig{}, tg.NewBotErrorf("Can't render html at Training_Send_MainMenu: %w", err)
	}

	// Creating msg
	msg := a.NewMessage(chatID, text)

	// Creating keyboard
	kb, boterr := Training_ReplyMarkup_MainMenu()
	if boterr != nil {
		return a.MessageConfig{}, tg.NewBotErrorw("Can't create keyboard at Training_Send_MainMenu", boterr)
	}
	msg.ReplyMarkup = kb

	// Parse mode
	msg.ParseMode = "HTML"

	return msg, nil
}

func Training_Edit_MainMenu(chatID int64, messageID int) (a.EditMessageTextConfig, *tg.BotError) {
	// Rendering text
	text, err := tmplManager.RenderHTML("TrainingMenu", nil)
	if err != nil {
		return a.EditMessageTextConfig{}, tg.NewBotErrorf("Can't render html at Training_Edit_MainMenu: %w", err)
	}

	// Creating editMsg
	editMsg := a.NewEditMessageText(chatID, messageID, text)

	// Creating keyboard
	kb, boterr := Training_ReplyMarkup_MainMenu()
	if boterr != nil {
		return a.EditMessageTextConfig{}, tg.NewBotErrorw("Can't create keyboard at Training_Edit_MainMenu", boterr)
	}
	editMsg.ReplyMarkup = &kb

	editMsg.ParseMode = "HTML"

	return editMsg, nil
}

func Training_ReplyMarkup_MainMenu() (a.InlineKeyboardMarkup, *tg.BotError) {
	// Rendering text
	textExercises, err := tmplManager.RenderHTML("TrainingMenuReplyMarkup_textExercises", nil)
	if err != nil {
		return a.InlineKeyboardMarkup{}, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_MainMenu: %w", err)
	}

	textTrainingProgramms, err := tmplManager.RenderHTML("TrainingMenuReplyMarkup_textTrainingProgramms", nil)
	if err != nil {
		return a.InlineKeyboardMarkup{}, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_MainMenu: %w", err)
	}

	// Assembling keyboard
	kb := a.NewInlineKeyboardMarkup(
		a.NewInlineKeyboardRow(
			a.NewInlineKeyboardButtonData(textExercises, "TEP0"),
		),
		a.NewInlineKeyboardRow(
			a.NewInlineKeyboardButtonData(textTrainingProgramms, "TPP0"),
		),
	)

	return kb, nil
}

// ===============================
// Exercises
// ===============================

func Training_Send_Exercises(chatID int64, page int, pageCount int, exercises *[]models.Exercise) (a.SendRichMessageConfig, *tg.BotError) {
	// Rendering text
	text, err := tmplManager.RenderHTML("TrainingExercises", nil)
	if err != nil {
		return a.SendRichMessageConfig{}, tg.NewBotErrorf("Can't render html at Training_Send_Exercises: %w", err)
	}
	exerciseTables := "\n"
	for _, exercise := range *exercises {
		if exercise.Description != nil {
			nbsp := strings.ReplaceAll(*exercise.Description, " ", "\u00A0")
			exercise.Description = &nbsp
		} else {
			empty := ""
			exercise.Description = &empty
		}
		text, err := tmplManager.RenderHTML("TrainingExercises_ExerciseTable", exercise)
		if err != nil {
			return a.SendRichMessageConfig{}, tg.NewBotErrorf("Can't render html at Training_Send_Exercises: %w", err)
		}
		exerciseTables += "\n" + text
	}
	text += exerciseTables

	// Creating msg
	rich := a.NewInputRichMessageHTML(text)
	msg := a.NewSendRichMessage(chatID, rich)

	// Creating keyboard
	kb, boterr := Training_ReplyMarkup_Exercises(page, pageCount, exercises)
	if boterr != nil {
		return a.SendRichMessageConfig{}, tg.NewBotErrorw("Can't create keyboard at Training_Send_Exercises", boterr)
	}
	msg.ReplyMarkup = kb

	return msg, nil
}

func Training_Edit_Exercises(chatID int64, messageID int, page int, pageCount int, exercises *[]models.Exercise) (a.EditMessageTextConfig, *tg.BotError) {
	// Rendering text
	text, err := tmplManager.RenderHTML("TrainingExercises", nil)
	if err != nil {
		return a.EditMessageTextConfig{}, tg.NewBotErrorf("Can't render html at Training_Edit_Exercises: %w", err)
	}
	exerciseTables, err := tmplManager.RenderHTML("TrainingExercises_ExerciseTable_Head", nil)
	if err != nil {
		return a.EditMessageTextConfig{}, tg.NewBotErrorf("Can't render html at Training_Edit_Exercises: %w", err)
	}
	for _, exercise := range *exercises {
		if exercise.Description != nil {
			nbsp := strings.ReplaceAll(*exercise.Description, " ", "\u00A0")
			exercise.Description = &nbsp
		} else {
			empty := ""
			exercise.Description = &empty
		}
		text, err := tmplManager.RenderHTML("TrainingExercises_ExerciseTable_Body", exercise)
		if err != nil {
			return a.EditMessageTextConfig{}, tg.NewBotErrorf("Can't render html at Training_Edit_Exercises: %w", err)
		}
		exerciseTables += text
	}
	tail, err := tmplManager.RenderHTML("TrainingExercises_ExerciseTable_Tail", nil)
	if err != nil {
		return a.EditMessageTextConfig{}, tg.NewBotErrorf("Can't render html at Training_Edit_Exercises: %w", err)
	}
	text += exerciseTables + tail

	// Creating editMsg
	rich := a.NewInputRichMessageHTML(text)
	editMsg := a.NewEditMessageText(chatID, messageID, "")
	editMsg.RichMessage = rich

	// Creating keyboard
	kb, boterr := Training_ReplyMarkup_Exercises(page, pageCount, exercises)
	if boterr != nil {
		return a.EditMessageTextConfig{}, tg.NewBotErrorw("Can't create keyboard at Training_Edit_Exercises", boterr)
	}
	editMsg.ReplyMarkup = kb

	return editMsg, nil
}

func Training_ReplyMarkup_Exercises(page int, pageCount int, exercises *[]models.Exercise) (*a.InlineKeyboardMarkup, *tg.BotError) {

	// Rendering text
	textPrevious, err := tmplManager.RenderHTML("TrainingExercisesReplyMarkup_Previous", nil)
	if err != nil {
		return nil, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercises: %w", err)
	}

	pages := struct {
		Page      int
		PageCount int
	}{Page: page + 1, PageCount: pageCount}
	if pageCount == 0 {
		pages.Page = 0
	}
	textSearch, err := tmplManager.RenderHTML("TrainingExercisesReplyMarkup_Search", pages)
	if err != nil {
		return nil, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercises: %w", err)
	}

	textNext, err := tmplManager.RenderHTML("TrainingExercisesReplyMarkup_Next", nil)
	if err != nil {
		return nil, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercises: %w", err)
	}

	textBack, err := tmplManager.RenderHTML("TrainingExercisesReplyMarkup_Back", nil)
	if err != nil {
		return nil, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercises: %w", err)
	}

	textCreate, err := tmplManager.RenderHTML("TrainingExercisesReplyMarkup_Create", nil)
	if err != nil {
		return nil, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercises: %w", err)
	}

	// Assembling exercises' buttons
	kb := a.NewInlineKeyboardMarkup()
	for _, exercise := range *exercises {
		kb.InlineKeyboard = append(kb.InlineKeyboard,
			a.NewInlineKeyboardRow(
				a.NewInlineKeyboardButtonData(exercise.Name, "TEI"+strconv.FormatInt(exercise.ID, 10))),
		)
	}

	// Next & previous page calculating
	previousPage := page
	nextPage := page
	if page > 0 {
		previousPage -= 1
	}
	if page < pageCount-1 {
		nextPage += 1
	}

	var previousCBQData string = "-"
	var nextCBQData string = "-"
	if previousPage != page {
		previousCBQData = "TEP" + strconv.Itoa(previousPage)
	}
	if nextPage != page {
		nextCBQData = "TEP" + strconv.Itoa(nextPage)
	}

	// Assembling navigation buttons for pages
	kb.InlineKeyboard = append(kb.InlineKeyboard, a.NewInlineKeyboardRow(
		a.NewInlineKeyboardButtonData(textPrevious, previousCBQData),
		a.NewInlineKeyboardButtonData(textSearch, "-"),
		a.NewInlineKeyboardButtonData(textNext, nextCBQData),
	))

	// Assebling navigation buttons for "Back" and "Create"
	kb.InlineKeyboard = append(kb.InlineKeyboard, a.NewInlineKeyboardRow(
		a.NewInlineKeyboardButtonData(textBack, "T"),
		a.NewInlineKeyboardButtonData(textCreate, "TEN"),
	))

	return &kb, nil
}

func Training_Edit_Exercises_Search(chatID int64, messageID int, page int, pageCount int, searchRes *models.ExerciseSearchResponse) (a.EditMessageTextConfig, *tg.BotError) {
	// Rendering text
	text, err := tmplManager.RenderHTML("TrainingExercises_Search", searchRes)
	if err != nil {
		return a.EditMessageTextConfig{}, tg.NewBotErrorf("Can't render html at Training_Edit_Exercises_Search: %w", err)
	}

	exerciseTables, err := tmplManager.RenderHTML("TrainingExercises_Search_ExerciseTable_Head", nil)
	if err != nil {
		return a.EditMessageTextConfig{}, tg.NewBotErrorf("Can't render html at Training_Edit_Exercises_Search: %w", err)
	}
	for _, exercise := range searchRes.Exercises {
		if exercise.Description != nil {
			nbsp := strings.ReplaceAll(*exercise.Description, " ", "\u00A0")
			exercise.Description = &nbsp
		} else {
			empty := ""
			exercise.Description = &empty
		}
		text, err := tmplManager.RenderHTML("TrainingExercises_Search_ExerciseTable_Body", exercise)
		if err != nil {
			return a.EditMessageTextConfig{}, tg.NewBotErrorf("Can't render html at Training_Edit_Exercises_Search: %w", err)
		}
		exerciseTables += text
	}
	tail, err := tmplManager.RenderHTML("TrainingExercises_Search_ExerciseTable_Tail", nil)
	if err != nil {
		return a.EditMessageTextConfig{}, tg.NewBotErrorf("Can't render html at Training_Edit_Exercises_Search: %w", err)
	}
	text += exerciseTables + tail

	// Creating editMsg
	rich := a.NewInputRichMessageHTML(text)
	editMsg := a.NewEditMessageText(chatID, messageID, "")
	editMsg.RichMessage = rich

	// Creating keyboard
	kb, boterr := Training_ReplyMarkup_Exercises_Search(page, pageCount, &searchRes.Exercises)
	if boterr != nil {
		return a.EditMessageTextConfig{}, tg.NewBotErrorw("Can't create keyboard at Training_Edit_Exercises_Search", boterr)
	}
	editMsg.ReplyMarkup = kb

	return editMsg, nil
}

func Training_ReplyMarkup_Exercises_Search(page int, pageCount int, exercises *[]models.Exercise) (*a.InlineKeyboardMarkup, *tg.BotError) {
	pages := struct {
		Page      int
		PageCount int
	}{Page: page + 1, PageCount: pageCount}
	if pageCount == 0 {
		pages.Page = 0
	}

	// Rendering text
	textPrevious, err := tmplManager.RenderHTML("TrainingExercises_SearchReplyMarkup_Previous", nil)
	if err != nil {
		return nil, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercises_Search: %w", err)
	}

	textPages, err := tmplManager.RenderHTML("TrainingExercises_SearchReplyMarkup_Pages", pages)
	if err != nil {
		return nil, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercises_Search: %w", err)
	}

	textNext, err := tmplManager.RenderHTML("TrainingExercises_SearchReplyMarkup_Next", nil)
	if err != nil {
		return nil, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercises_Search: %w", err)
	}

	textBack, err := tmplManager.RenderHTML("TrainingExercises_SearchReplyMarkup_Back", nil)
	if err != nil {
		return nil, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercises_Search: %w", err)
	}

	// Assembling exercises' buttons
	kb := a.NewInlineKeyboardMarkup()
	for _, exercise := range *exercises {
		kb.InlineKeyboard = append(kb.InlineKeyboard,
			a.NewInlineKeyboardRow(
				a.NewInlineKeyboardButtonData(exercise.Name, "TEI"+strconv.FormatInt(exercise.ID, 10))),
		)
	}

	// Next & previous page calculating
	previousPage := page
	nextPage := page
	if page > 0 {
		previousPage -= 1
	}
	if page < pageCount-1 {
		nextPage += 1
	}

	var previousCBQData string = "-"
	var nextCBQData string = "-"
	if previousPage != page {
		previousCBQData = "TES" + strconv.Itoa(previousPage)
	}
	if nextPage != page {
		nextCBQData = "TES" + strconv.Itoa(nextPage)
	}

	// Assembling navigation buttons for pages
	kb.InlineKeyboard = append(kb.InlineKeyboard, a.NewInlineKeyboardRow(
		a.NewInlineKeyboardButtonData(textPrevious, previousCBQData),
		a.NewInlineKeyboardButtonData(textPages, "-"),
		a.NewInlineKeyboardButtonData(textNext, nextCBQData),
	))

	// Assebling navigation buttons for "Back" and "Create"
	kb.InlineKeyboard = append(kb.InlineKeyboard, a.NewInlineKeyboardRow(
		a.NewInlineKeyboardButtonData(textBack, "TEP0"),
	))

	return &kb, nil
}

// ==============================================
// Exercises create form
// ==============================================

type CreateExerciseFormValues struct {
	Name        string
	Description string
	Technique   string
	Unit        string
}

func Training_Edit_Exercises_CreateForm(chatID int64, contextData models.ContextData) (a.EditMessageTextConfig, *tg.BotError) {

	// Rendering text
	text, err := tmplManager.RenderHTML("Training_Edit_Exercises_CreateForm", contextData.Values)
	if err != nil {
		return a.EditMessageTextConfig{}, tg.NewBotErrorf("Can't render html at Training_Exercises_CreateForm: %w", err)
	}

	// Highlighting focused field
	text = tmplManager.WrapLineHTML(text, contextData.Focus+1, "➡️</tg-emoji>", "<tg-emoji emoji-id=\"5215416746453776052\">🔴</tg-emoji> <u>", "</u>")

	// Creating editMsg
	editMsg := a.NewEditMessageText(chatID, contextData.MessageID, text)

	// Creating keyboard
	unit := models.WeightUnit(contextData.Values["Unit"])
	kb, boterr := Training_ReplyMarkup_Exercises_CreateForm(unit)
	if boterr != nil {
		return a.EditMessageTextConfig{}, boterr
	}
	editMsg.ReplyMarkup = kb

	// Parse mode
	editMsg.ParseMode = "HTML"
	editMsg.LinkPreviewOptions.IsDisabled = true

	return editMsg, nil
}

func Training_ReplyMarkup_Exercises_CreateForm(unit models.WeightUnit) (*a.InlineKeyboardMarkup, *tg.BotError) {

	// Rendering text
	textUp, err := tmplManager.RenderHTML("Training_ReplyMarkup_Exercises_CreateForm_Up", nil)
	if err != nil {
		return nil, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercises_CreateForm: %w", err)
	}

	textDown, err := tmplManager.RenderHTML("Training_ReplyMarkup_Exercises_CreateForm_Down", nil)
	if err != nil {
		return nil, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercises_CreateForm: %w", err)
	}

	textBack, err := tmplManager.RenderHTML("Training_ReplyMarkup_Exercises_CreateForm_Back", nil)
	if err != nil {
		return nil, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercises_CreateForm: %w", err)
	}

	textSave, err := tmplManager.RenderHTML("Training_ReplyMarkup_Exercises_CreateForm_Save", nil)
	if err != nil {
		return nil, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercises_CreateForm: %w", err)
	}

	textKg, err := tmplManager.RenderHTML("Training_ReplyMarkup_Exercises_CreateForm_Kg", nil)
	if err != nil {
		return nil, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercises_CreateForm: %w", err)
	}
	if unit == models.WeightUnit_Kg {
		textKg = "✅" + textKg
	}

	textMixed, err := tmplManager.RenderHTML("Training_ReplyMarkup_Exercises_CreateForm_Mixed", nil)
	if err != nil {
		return nil, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercises_CreateForm: %w", err)
	}
	if unit == models.WeightUnit_Mixed {
		textMixed = "✅" + textMixed
	}

	textLbs, err := tmplManager.RenderHTML("Training_ReplyMarkup_Exercises_CreateForm_Lbs", nil)
	if err != nil {
		return nil, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercises_CreateForm: %w", err)
	}
	if unit == models.WeightUnit_Lbs {
		textLbs = "✅" + textLbs
	}

	// Assembling keyboard
	kb := a.NewInlineKeyboardMarkup(
		a.NewInlineKeyboardRow(
			a.NewInlineKeyboardButtonData(textUp, "TENU"),
			a.NewInlineKeyboardButtonData(textDown, "TEND"),
		),
		a.NewInlineKeyboardRow(
			a.NewInlineKeyboardButtonData(textKg, "TEN_KG"),
			a.NewInlineKeyboardButtonData(textMixed, "TEN_MIX"),
			a.NewInlineKeyboardButtonData(textLbs, "TEN_LBS"),
		),
		a.NewInlineKeyboardRow(
			a.NewInlineKeyboardButtonData(textBack, "TEP1"),
			a.NewInlineKeyboardButtonData(textSave, "TENS"),
		),
	)

	return &kb, nil
}

type EditExerciseFormValues struct {
	ID          int64
	Name        string
	Description string
	Technique   string
	Unit        string
}

func Training_Edit_Exercises_EditForm(chatID int64, contextData models.ContextData) (a.EditMessageTextConfig, *tg.BotError) {
	// Rendering text
	text, err := tmplManager.RenderHTML("Training_Edit_Exercises_EditForm", contextData.Values)
	if err != nil {
		return a.EditMessageTextConfig{}, tg.NewBotErrorf("Can't render html at Training_Edit_Exercises_EditForm: %w", err)
	}

	// Highlighting focused field
	text = tmplManager.WrapLineHTML(text, contextData.Focus+1, "➡️</tg-emoji>", "<tg-emoji emoji-id=\"5215416746453776052\">🔴</tg-emoji> <u>", "</u>")

	// Creating editMsg
	editMsg := a.NewEditMessageText(chatID, contextData.MessageID, text)

	// Creating keyboard
	unit := models.WeightUnit(contextData.Values["Unit"])
	exerciseID, err := strconv.ParseInt(contextData.Values["ExerciseID"], 10, 64)
	if err != nil || exerciseID == 0 {
		return a.EditMessageTextConfig{}, tg.NewBotErrorf("Can't parse exercise id at Training_Edit_Exercises_EditForm: %w", err)
	}
	kb, boterr := Training_ReplyMarkup_Exercises_EditForm(exerciseID, unit)
	if boterr != nil {
		return a.EditMessageTextConfig{}, boterr
	}
	editMsg.ReplyMarkup = kb

	// Parse mode
	editMsg.ParseMode = "HTML"
	editMsg.LinkPreviewOptions.IsDisabled = true

	return editMsg, nil
}

func Training_ReplyMarkup_Exercises_EditForm(exerciseID int64, unit models.WeightUnit) (*a.InlineKeyboardMarkup, *tg.BotError) {

	// Rendering text
	textUp, err := tmplManager.RenderHTML("Training_ReplyMarkup_Exercises_EditForm_Up", nil)
	if err != nil {
		return nil, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercises_EditForm: %w", err)
	}

	textDown, err := tmplManager.RenderHTML("Training_ReplyMarkup_Exercises_EditForm_Down", nil)
	if err != nil {
		return nil, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercises_EditForm: %w", err)
	}

	textBack, err := tmplManager.RenderHTML("Training_ReplyMarkup_Exercises_EditForm_Back", nil)
	if err != nil {
		return nil, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercises_EditForm: %w", err)
	}

	textSave, err := tmplManager.RenderHTML("Training_ReplyMarkup_Exercises_EditForm_Save", nil)
	if err != nil {
		return nil, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercises_EditForm: %w", err)
	}

	textKg, err := tmplManager.RenderHTML("Training_ReplyMarkup_Exercises_EditForm_Kg", nil)
	if err != nil {
		return nil, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercises_EditForm: %w", err)
	}
	if unit == models.WeightUnit_Kg {
		textKg = "✅" + textKg
	}

	textMixed, err := tmplManager.RenderHTML("Training_ReplyMarkup_Exercises_EditForm_Mixed", nil)
	if err != nil {
		return nil, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercises_EditForm: %w", err)
	}
	if unit == models.WeightUnit_Mixed {
		textMixed = "✅" + textMixed
	}

	textLbs, err := tmplManager.RenderHTML("Training_ReplyMarkup_Exercises_EditForm_Lbs", nil)
	if err != nil {
		return nil, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercises_EditForm: %w", err)
	}
	if unit == models.WeightUnit_Lbs {
		textLbs = "✅" + textLbs
	}

	// Assembling keyboard
	kb := a.NewInlineKeyboardMarkup(
		a.NewInlineKeyboardRow(
			a.NewInlineKeyboardButtonData(textUp, fmt.Sprintf("TEEU%d", exerciseID)),
			a.NewInlineKeyboardButtonData(textDown, fmt.Sprintf("TEED%d", exerciseID)),
		),
		a.NewInlineKeyboardRow(
			a.NewInlineKeyboardButtonData(textKg, "TEE_KG"),
			a.NewInlineKeyboardButtonData(textMixed, "TEE_MIX"),
			a.NewInlineKeyboardButtonData(textLbs, "TEE_LBS"),
		),
		a.NewInlineKeyboardRow(
			a.NewInlineKeyboardButtonData(textBack, fmt.Sprintf("TEI%d", exerciseID)),
			a.NewInlineKeyboardButtonData(textSave, fmt.Sprintf("TEES%d", exerciseID)),
		),
	)

	return &kb, nil
}

func Training_Edit_Exercise(chatID int64, messageID int, exercise *models.Exercise) (a.EditMessageTextConfig, *tg.BotError) {
	exercise.WeightUnit = models.WeightUnit(strings.ToUpper(string(exercise.WeightUnit)))
	if exercise.Description == nil {
		empty := ""
		exercise.Description = &empty
	}
	// Rendering text
	text, err := tmplManager.RenderHTML("TrainingExercise", exercise)
	if err != nil {
		return a.EditMessageTextConfig{}, tg.NewBotErrorf("Can't render html at Training_Edit_Exercise: %w", err)
	}

	// Creating editMsg
	editMsg := a.NewEditMessageText(chatID, messageID, text)

	// Creating keyboard
	kb, boterr := Training_ReplyMarkup_Exercise(exercise.ID)
	if boterr != nil {
		return a.EditMessageTextConfig{}, tg.NewBotErrorw("Can't create keyboard at Training_Edit_Exercise", boterr)
	}
	editMsg.ReplyMarkup = kb

	// Parse mode
	editMsg.ParseMode = "HTML"

	return editMsg, nil
}

func Training_ReplyMarkup_Exercise(exerciseID int64) (*a.InlineKeyboardMarkup, *tg.BotError) {
	// Rendering text
	textPerform, err := tmplManager.RenderHTML("TrainingExerciseReplyMarkup_Perform", nil)
	if err != nil {
		return nil, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercise: %w", err)
	}
	textEdit, err := tmplManager.RenderHTML("TrainingExerciseReplyMarkup_Edit", nil)
	if err != nil {
		return nil, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercise: %w", err)
	}
	textDelete, err := tmplManager.RenderHTML("TrainingExerciseReplyMarkup_Delete", nil)
	if err != nil {
		return nil, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercise: %w", err)
	}
	textBack, err := tmplManager.RenderHTML("TrainingExerciseReplyMarkup_Back", nil)
	if err != nil {
		return nil, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercise: %w", err)
	}

	// Assembling exercise's buttons
	kb := a.NewInlineKeyboardMarkup(
		a.NewInlineKeyboardRow(
			a.NewInlineKeyboardButtonData(textPerform, fmt.Sprintf("TEEnN%d", exerciseID)),
		),
		a.NewInlineKeyboardRow(
			a.NewInlineKeyboardButtonData(textBack, "TEP0"),
			a.NewInlineKeyboardButtonData(textEdit, fmt.Sprintf("TEE%d", exerciseID)),
			a.NewInlineKeyboardButtonData(textDelete, fmt.Sprintf("TED%d", exerciseID)),
		),
	)

	return &kb, nil
}

func Training_Edit_Exercise_DeleteForm(chatID int64, messageID int, exercise *models.Exercise) (a.EditMessageTextConfig, *tg.BotError) {
	// Rendering text
	text, err := tmplManager.RenderHTML("Training_Edit_Exercises_DeleteForm", exercise)
	if err != nil {
		return a.EditMessageTextConfig{}, tg.NewBotErrorf("Can't render html at Training_Edit_Exercise_DeleteForm: %w", err)
	}

	// Creating editMsg
	editMsg := a.NewEditMessageText(chatID, messageID, text)

	// Creating keyboard
	kb, boterr := Training_ReplyMarkup_Exercise_DeleteForm(exercise.ID)
	if boterr != nil {
		return a.EditMessageTextConfig{}, tg.NewBotErrorw("Can't create keyboard at Training_Edit_Programms", boterr)
	}
	editMsg.ReplyMarkup = kb

	// Parse mode
	editMsg.ParseMode = "HTML"

	return editMsg, nil
}

func Training_ReplyMarkup_Exercise_DeleteForm(exerciseID int64) (*a.InlineKeyboardMarkup, *tg.BotError) {
	// Rendering text
	textBack, err := tmplManager.RenderHTML("Training_Edit_Exercises_DeleteForm_Back", nil)
	if err != nil {
		return nil, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercise_DeleteForm: %w", err)
	}
	textDelete, err := tmplManager.RenderHTML("Training_Edit_Exercises_DeleteForm_Delete", nil)
	if err != nil {
		return nil, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercise_DeleteForm: %w", err)
	}

	// Assembling exercise's buttons
	kb := a.NewInlineKeyboardMarkup(
		a.NewInlineKeyboardRow(
			a.NewInlineKeyboardButtonData(textBack, fmt.Sprintf("TEI%d", exerciseID)),
			a.NewInlineKeyboardButtonData(textDelete, fmt.Sprintf("TEDC%d", exerciseID)),
		),
	)

	return &kb, nil
}

// ================================================
// Programms
// ================================================

func Training_Edit_Programms(chatID int64, messageID int, page int, pageCount int, programms *[]models.Exercise) (a.EditMessageTextConfig, *tg.BotError) {
	// Rendering text
	text, err := tmplManager.RenderHTML("TrainingProgramms", nil)
	if err != nil {
		return a.EditMessageTextConfig{}, tg.NewBotErrorf("Can't render html at Training_Edit_Programms: %w", err)
	}

	// Creating editMsg
	editMsg := a.NewEditMessageText(chatID, messageID, text)

	// Creating keyboard
	kb, boterr := Training_ReplyMarkup_Programms(page, pageCount, programms)
	if boterr != nil {
		return a.EditMessageTextConfig{}, tg.NewBotErrorw("Can't create keyboard at Training_Edit_Programms", boterr)
	}
	editMsg.ReplyMarkup = kb

	// Parse mode
	editMsg.ParseMode = "HTML"

	return editMsg, nil
}

func Training_ReplyMarkup_Programms(page int, pageCount int, programms *[]models.Exercise) (*a.InlineKeyboardMarkup, *tg.BotError) {

	// Rendering text
	textPrevious, err := tmplManager.RenderHTML("Training_ReplyMarkup_Programms_Previous", nil)
	if err != nil {
		return nil, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Programms: %w", err)
	}

	pages := struct {
		Page      int
		PageCount int
	}{Page: page, PageCount: pageCount}
	textSearch, err := tmplManager.RenderHTML("Training_ReplyMarkup_Programms_Search", pages)
	if err != nil {
		return nil, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Programms: %w", err)
	}

	textNext, err := tmplManager.RenderHTML("Training_ReplyMarkup_Programms_Next", nil)
	if err != nil {
		return nil, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Programms: %w", err)
	}

	textBack, err := tmplManager.RenderHTML("Training_ReplyMarkup_Programms_Back", nil)
	if err != nil {
		return nil, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Programms: %w", err)
	}

	textCreate, err := tmplManager.RenderHTML("Training_ReplyMarkup_Programms_Create", nil)
	if err != nil {
		return nil, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Programms: %w", err)
	}

	// Assembling programms' buttons
	kb := a.NewInlineKeyboardMarkup()
	for _, programm := range *programms {
		kb.InlineKeyboard = append(kb.InlineKeyboard,
			a.NewInlineKeyboardRow(
				a.NewInlineKeyboardButtonData(programm.Name, "TPI"+strconv.FormatInt(programm.ID, 10))),
		)
	}

	// Next & previous page calculating
	previousPage := page
	nextPage := page
	if page > 0 {
		previousPage -= 1
	}
	if page < pageCount-1 {
		nextPage += 1
	}

	var previousCBQData string = "-"
	var nextCBQData string = "-"
	if previousPage != page {
		previousCBQData = "TPP" + strconv.Itoa(previousPage)
	}
	if nextPage != page {
		nextCBQData = "TPP" + strconv.Itoa(nextPage)
	}

	// Assembling navigation buttons for pages
	kb.InlineKeyboard = append(kb.InlineKeyboard, a.NewInlineKeyboardRow(
		a.NewInlineKeyboardButtonData(textPrevious, previousCBQData),
		a.NewInlineKeyboardButtonData(textSearch, "TPS"),
		a.NewInlineKeyboardButtonData(textNext, nextCBQData),
	))

	// Assebling navigation buttons for "Back" and "Create"
	kb.InlineKeyboard = append(kb.InlineKeyboard, a.NewInlineKeyboardRow(
		a.NewInlineKeyboardButtonData(textBack, "T"),
		a.NewInlineKeyboardButtonData(textCreate, "TPA"),
	))

	return &kb, nil
}
