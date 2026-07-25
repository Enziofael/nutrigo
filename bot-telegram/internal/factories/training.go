package factories

import (
	"strconv"

	tg "github.com/Enziofael/nutrigo/bot-telegram/pkg/telegroni"
	"github.com/Enziofael/nutrigo/shared/models/v1"
	a "github.com/OvyFlash/telegram-bot-api"
)

func TrainingMenu(tgID int64) (a.MessageConfig, *tg.BotError) {
	text, err := tmplManager.RenderHTML("TrainingMenu", nil)
	if err != nil {
		return a.MessageConfig{}, tg.NewBotError(err.Error(), nil)
	}

	msg := a.NewMessage(tgID, text)

	kb, boterr := TrainingMenuReplyMarkup()
	if boterr != nil {
		return a.MessageConfig{}, boterr
	}
	msg.ReplyMarkup = kb

	msg.ParseMode = "HTML"

	return msg, nil
}

func TrainingMenuReplyMarkup() (a.InlineKeyboardMarkup, *tg.BotError) {
	textExercises, err := tmplManager.RenderHTML("TrainingMenuReplyMarkup_textExercises", nil)
	if err != nil {
		return a.InlineKeyboardMarkup{}, tg.NewBotError(err.Error(), nil)
	}
	textTrainingProgramms, err := tmplManager.RenderHTML("TrainingMenuReplyMarkup_textTrainingProgramms", nil)
	if err != nil {
		return a.InlineKeyboardMarkup{}, tg.NewBotError(err.Error(), nil)
	}

	kb := a.NewInlineKeyboardMarkup(
		a.NewInlineKeyboardRow(
			a.NewInlineKeyboardButtonData(textExercises, "TEP1"),
		),
		a.NewInlineKeyboardRow(
			a.NewInlineKeyboardButtonData(textTrainingProgramms, "TPP1"),
		),
	)

	return kb, nil
}

func TrainingExercises(tgID int64, messageID int, page int, pageCount int, exercises *[]models.Exercise) (a.EditMessageTextConfig, *tg.BotError) {
	text, err := tmplManager.RenderHTML("TrainingExercises", nil)
	if err != nil {
		return a.EditMessageTextConfig{}, tg.NewBotError(err.Error(), nil)
	}

	editMsg := a.NewEditMessageText(tgID, messageID, text)

	kb, boterr := TrainingExercisesReplyMarkup(page, pageCount, exercises)
	if boterr != nil {
		return a.EditMessageTextConfig{}, boterr
	}
	editMsg.ReplyMarkup = kb

	editMsg.ParseMode = "HTML"

	return editMsg, nil
}

func TrainingExercisesReplyMarkup(page int, pageCount int, exercises *[]models.Exercise) (*a.InlineKeyboardMarkup, *tg.BotError) {
	kb := a.NewInlineKeyboardMarkup()
	for _, exercise := range *exercises {
		kb.InlineKeyboard = append(kb.InlineKeyboard,
			a.NewInlineKeyboardRow(
				a.NewInlineKeyboardButtonData(exercise.Name, "TE"+strconv.FormatInt(exercise.ID, 10))),
		)

	}

	textPrevious, err := tmplManager.RenderHTML("TrainingExercisesReplyMarkup_textPrevious", nil)
	if err != nil {
		return &a.InlineKeyboardMarkup{}, tg.NewBotError(err.Error(), nil)
	}

	pages := struct {
		Page      int
		PageCount int
	}{Page: page, PageCount: pageCount}
	textSearch, err := tmplManager.RenderHTML("TrainingExercisesReplyMarkup_textSearch", pages)
	if err != nil {
		return &a.InlineKeyboardMarkup{}, tg.NewBotError(err.Error(), nil)
	}

	textNext, err := tmplManager.RenderHTML("TrainingExercisesReplyMarkup_textNext", nil)
	if err != nil {
		return &a.InlineKeyboardMarkup{}, tg.NewBotError(err.Error(), nil)
	}

	previousPage := page
	nextPage := page
	if page > 1 {
		previousPage -= 1
	}
	if page < pageCount-1 {
		nextPage += 1
	}

	kb.InlineKeyboard = append(kb.InlineKeyboard, a.NewInlineKeyboardRow(
		a.NewInlineKeyboardButtonData(textPrevious, "TE"+strconv.Itoa(previousPage)),
		a.NewInlineKeyboardButtonData(textSearch, "TES"),
		a.NewInlineKeyboardButtonData(textNext, "TE"+strconv.Itoa(nextPage)),
	))

	textBack, err := tmplManager.RenderHTML("TrainingExercisesReplyMarkup_textBack", nil)
	if err != nil {
		return &a.InlineKeyboardMarkup{}, tg.NewBotError(err.Error(), nil)
	}
	textAdd, err := tmplManager.RenderHTML("TrainingExercisesReplyMarkup_textAdd", nil)
	kb.InlineKeyboard = append(kb.InlineKeyboard, a.NewInlineKeyboardRow(
		a.NewInlineKeyboardButtonData(textBack, "T"),
		a.NewInlineKeyboardButtonData(textAdd, "TEA"),
	))

	return &kb, nil
}

func TrainingProgramms(tgID int64, messageID int, page int, pageCount int, programms *[]models.Exercise) (a.EditMessageTextConfig, *tg.BotError) {
	text, err := tmplManager.RenderHTML("TrainingProgramms", nil)
	if err != nil {
		return a.EditMessageTextConfig{}, tg.NewBotError(err.Error(), nil)
	}

	editMsg := a.NewEditMessageText(tgID, messageID, text)

	kb, boterr := TrainingProgrammsReplyMarkup(page, pageCount, programms)
	if boterr != nil {
		return a.EditMessageTextConfig{}, boterr
	}
	editMsg.ReplyMarkup = kb

	editMsg.ParseMode = "HTML"

	return editMsg, nil
}

func TrainingProgrammsReplyMarkup(page int, pageCount int, programms *[]models.Exercise) (*a.InlineKeyboardMarkup, *tg.BotError) {
	kb := a.NewInlineKeyboardMarkup()
	for _, programm := range *programms {
		kb.InlineKeyboard = append(kb.InlineKeyboard,
			a.NewInlineKeyboardRow(
				a.NewInlineKeyboardButtonData(programm.Name, "TP"+strconv.FormatInt(programm.ID, 10))),
		)
	}

	textPrevious, err := tmplManager.RenderHTML("TrainingProgrammsReplyMarkup_textPrevious", nil)
	if err != nil {
		return &a.InlineKeyboardMarkup{}, tg.NewBotError(err.Error(), nil)
	}

	pages := struct {
		Page      int
		PageCount int
	}{Page: page, PageCount: pageCount}
	textSearch, err := tmplManager.RenderHTML("TrainingProgrammsReplyMarkup_textSearch", pages)
	if err != nil {
		return &a.InlineKeyboardMarkup{}, tg.NewBotError(err.Error(), nil)
	}

	textNext, err := tmplManager.RenderHTML("TrainingProgrammsReplyMarkup_textNext", nil)
	if err != nil {
		return &a.InlineKeyboardMarkup{}, tg.NewBotError(err.Error(), nil)
	}

	previousPage := page
	nextPage := page
	if page > 1 {
		previousPage -= 1
	}
	if page < pageCount-1 {
		nextPage += 1
	}

	kb.InlineKeyboard = append(kb.InlineKeyboard, a.NewInlineKeyboardRow(
		a.NewInlineKeyboardButtonData(textPrevious, "TP"+strconv.Itoa(previousPage)),
		a.NewInlineKeyboardButtonData(textSearch, "TPS"),
		a.NewInlineKeyboardButtonData(textNext, "TP"+strconv.Itoa(nextPage)),
	))

	textBack, err := tmplManager.RenderHTML("TrainingExercisesReplyMarkup_textBack", nil)
	if err != nil {
		return &a.InlineKeyboardMarkup{}, tg.NewBotError(err.Error(), nil)
	}
	textAdd, err := tmplManager.RenderHTML("TrainingExercisesReplyMarkup_textAdd", nil)
	kb.InlineKeyboard = append(kb.InlineKeyboard, a.NewInlineKeyboardRow(
		a.NewInlineKeyboardButtonData(textBack, "T"),
		a.NewInlineKeyboardButtonData(textAdd, "TPA"),
	))

	return &kb, nil
}

func TrainingMenuEdit(tgID int64, messageID int) (a.EditMessageTextConfig, *tg.BotError) {
	text, err := tmplManager.RenderHTML("TrainingMenu", nil)
	if err != nil {
		return a.EditMessageTextConfig{}, tg.NewBotError(err.Error(), nil)
	}

	editMsg := a.NewEditMessageText(tgID, messageID, text)

	kb, boterr := TrainingMenuReplyMarkup()
	if boterr != nil {
		return a.EditMessageTextConfig{}, boterr
	}
	editMsg.ReplyMarkup = &kb

	editMsg.ParseMode = "HTML"

	return editMsg, nil
}
