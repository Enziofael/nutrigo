package factories

import (
	"log"
	"strconv"

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

func Training_Send_Exercises(chatID int64, page int, pageCount int, exercises *[]models.Exercise) (a.MessageConfig, *tg.BotError) {
	// Rendering text
	text, err := tmplManager.RenderHTML("TrainingExercises", nil)
	if err != nil {
		return a.MessageConfig{}, tg.NewBotErrorf("Can't render html at Training_Send_Exercises: %w", err)
	}

	// Creating editMsg
	msg := a.NewMessage(chatID, text)

	// Creating keyboard
	kb, boterr := Training_ReplyMarkup_Exercises(page, pageCount, exercises)
	if boterr != nil {
		return a.MessageConfig{}, tg.NewBotErrorw("Can't create keyboard at Training_Send_Exercises", boterr)
	}
	msg.ReplyMarkup = kb

	// Parse mode
	msg.ParseMode = "HTML"

	return msg, nil
}

func Training_Edit_Exercises(chatID int64, messageID int, page int, pageCount int, exercises *[]models.Exercise) (a.EditMessageTextConfig, *tg.BotError) {
	// Rendering text
	text, err := tmplManager.RenderHTML("TrainingExercises", nil)
	if err != nil {
		return a.EditMessageTextConfig{}, tg.NewBotErrorf("Can't render html at Training_Edit_Exercises: %w", err)
	}

	// Creating editMsg
	editMsg := a.NewEditMessageText(chatID, messageID, text)

	// Creating keyboard
	kb, boterr := Training_ReplyMarkup_Exercises(page, pageCount, exercises)
	if boterr != nil {
		return a.EditMessageTextConfig{}, tg.NewBotErrorw("Can't create keyboard at Training_Edit_Exercises", boterr)
	}
	editMsg.ReplyMarkup = kb

	// Parse mode
	editMsg.ParseMode = "HTML"

	return editMsg, nil
}

func Training_ReplyMarkup_Exercises(page int, pageCount int, exercises *[]models.Exercise) (*a.InlineKeyboardMarkup, *tg.BotError) {

	// Rendering text
	textPrevious, err := tmplManager.RenderHTML("TrainingExercisesReplyMarkup_Previous", nil)
	if err != nil {
		return &a.InlineKeyboardMarkup{}, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercises: %w", err)
	}

	pages := struct {
		Page      int
		PageCount int
	}{Page: page + 1, PageCount: pageCount}
	textSearch, err := tmplManager.RenderHTML("TrainingExercisesReplyMarkup_Search", pages)
	if err != nil {
		return &a.InlineKeyboardMarkup{}, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercises: %w", err)
	}

	textNext, err := tmplManager.RenderHTML("TrainingExercisesReplyMarkup_Next", nil)
	if err != nil {
		return &a.InlineKeyboardMarkup{}, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercises: %w", err)
	}

	textBack, err := tmplManager.RenderHTML("TrainingExercisesReplyMarkup_Back", nil)
	if err != nil {
		return &a.InlineKeyboardMarkup{}, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercises: %w", err)
	}

	textCreate, err := tmplManager.RenderHTML("TrainingExercisesReplyMarkup_Create", nil)
	if err != nil {
		return &a.InlineKeyboardMarkup{}, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercises: %w", err)
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
	log.Printf("Page %d Total %d", page, pageCount)
	previousPage := page
	nextPage := page
	if page > 0 {
		previousPage -= 1
		log.Printf("Previous %d", previousPage)
	}
	if page < pageCount-1 {
		nextPage += 1
		log.Printf("Previous %d", nextPage)
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
		a.NewInlineKeyboardButtonData(textSearch, "TES"),
		a.NewInlineKeyboardButtonData(textNext, nextCBQData),
	))

	// Assebling navigation buttons for "Back" and "Create"
	kb.InlineKeyboard = append(kb.InlineKeyboard, a.NewInlineKeyboardRow(
		a.NewInlineKeyboardButtonData(textBack, "T"),
		a.NewInlineKeyboardButtonData(textCreate, "TEN"),
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
}

// IMPORTANT: check template for valid value
const FormValuesStartLine int = 3

func Training_Edit_Exercises_CreateForm(chatID int64, contextData models.ContextData) (a.EditMessageTextConfig, *tg.BotError) {

	// Rendering text
	text, err := tmplManager.RenderHTML("Training_Edit_Exercises_CreateForm", contextData.Values)
	if err != nil {
		return a.EditMessageTextConfig{}, tg.NewBotErrorf("Can't render html at Training_Exercises_CreateForm: %w", err)
	}

	// Highlighting focused field
	text = tmplManager.WrapLineHTML(text, contextData.Focus+FormValuesStartLine, "<b><u>", "</u></b>")

	// Creating editMsg
	editMsg := a.NewEditMessageText(chatID, contextData.MessageID, text)

	// Creating keyboard
	kb, boterr := Training_ReplyMarkup_Exercises_CreateForm()
	if boterr != nil {
		return a.EditMessageTextConfig{}, boterr
	}
	editMsg.ReplyMarkup = kb

	// Parse mode
	editMsg.ParseMode = "HTML"

	return editMsg, nil
}

func Training_ReplyMarkup_Exercises_CreateForm() (*a.InlineKeyboardMarkup, *tg.BotError) {

	// Rendering text
	textUp, err := tmplManager.RenderHTML("Training_ReplyMarkup_Exercises_CreateForm_Up", nil)
	if err != nil {
		return &a.InlineKeyboardMarkup{}, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercises: %w", err)
	}

	textDown, err := tmplManager.RenderHTML("Training_ReplyMarkup_Exercises_CreateForm_Down", nil)
	if err != nil {
		return &a.InlineKeyboardMarkup{}, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercises: %w", err)
	}

	textBack, err := tmplManager.RenderHTML("Training_ReplyMarkup_Exercises_CreateForm_Back", nil)
	if err != nil {
		return &a.InlineKeyboardMarkup{}, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercises: %w", err)
	}

	textSave, err := tmplManager.RenderHTML("Training_ReplyMarkup_Exercises_CreateForm_Save", nil)
	if err != nil {
		return &a.InlineKeyboardMarkup{}, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Exercises: %w", err)
	}

	// Assembling keyboard
	kb := a.NewInlineKeyboardMarkup(
		a.NewInlineKeyboardRow(
			a.NewInlineKeyboardButtonData(textUp, "TENU"),
			a.NewInlineKeyboardButtonData(textDown, "TEND"),
		),
		a.NewInlineKeyboardRow(
			a.NewInlineKeyboardButtonData(textBack, "TEP1"),
			a.NewInlineKeyboardButtonData(textSave, "TENS"),
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
		return &a.InlineKeyboardMarkup{}, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Programms: %w", err)
	}

	pages := struct {
		Page      int
		PageCount int
	}{Page: page, PageCount: pageCount}
	textSearch, err := tmplManager.RenderHTML("Training_ReplyMarkup_Programms_Search", pages)
	if err != nil {
		return &a.InlineKeyboardMarkup{}, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Programms: %w", err)
	}

	textNext, err := tmplManager.RenderHTML("Training_ReplyMarkup_Programms_Next", nil)
	if err != nil {
		return &a.InlineKeyboardMarkup{}, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Programms: %w", err)
	}

	textBack, err := tmplManager.RenderHTML("Training_ReplyMarkup_Programms_Back", nil)
	if err != nil {
		return &a.InlineKeyboardMarkup{}, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Programms: %w", err)
	}

	textCreate, err := tmplManager.RenderHTML("Training_ReplyMarkup_Programms_Create", nil)
	if err != nil {
		return &a.InlineKeyboardMarkup{}, tg.NewBotErrorf("Can't render html at Training_ReplyMarkup_Programms: %w", err)
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
	log.Printf("Page %d Total %d", page, pageCount)
	previousPage := page
	nextPage := page
	if page > 0 {
		previousPage -= 1
		log.Printf("Previous %d", previousPage)
	}
	if page < pageCount-1 {
		nextPage += 1
		log.Printf("Previous %d", nextPage)
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
