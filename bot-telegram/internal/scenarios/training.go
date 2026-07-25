package scenarios

import (
	//client "github.com/Enziofael/nutrigo/bot-telegram/internal/client/v1"
	f "github.com/Enziofael/nutrigo/bot-telegram/internal/factories"
	tg "github.com/Enziofael/nutrigo/bot-telegram/pkg/telegroni"
	"github.com/Enziofael/nutrigo/shared/models/v1"

	//models "github.com/Enziofael/nutrigo/shared/models/v1"
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

func TrainingMenu(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	msg, boterr := f.TrainingMenu(update.SentFrom().ID)
	if boterr != nil {
		return tg.StatusError, boterr
	}
	if _, err := ctx.Bot.Request(msg); err != nil {
		return tg.StatusError, tg.NewBotError(err.Error(), nil)
	}

	return tg.StatusOK, nil
}

func TrainingExercises(ctx tg.Context, update tgbotapi.Update, page int) (tg.HandleStatus, *tg.BotError) {
	pageCount := 1
	exercises := &[]models.Exercise{
		{Name: "Упражнение 1", ID: 1},
		{Name: "Упражнение 2", ID: 2},
		{Name: "Упражнение 3", ID: 3},
		{Name: "Упражнение 4", ID: 4},
		{Name: "Упражнение 5", ID: 5},
	}
	editMsg, boterr := f.TrainingExercises(update.SentFrom().ID,
		update.CallbackQuery.Message.MessageID, page, pageCount, exercises)
	if boterr != nil {
		return tg.StatusError, boterr
	}
	if _, err := ctx.Bot.Request(editMsg); err != nil {
		return tg.StatusError, tg.NewBotError(err.Error(), nil)
	}

	return tg.StatusOK, nil
}

func TrainingProgramms(ctx tg.Context, update tgbotapi.Update, page int) (tg.HandleStatus, *tg.BotError) {
	pageCount := 1
	programms := &[]models.Exercise{
		{Name: "Программа 1", ID: 1},
		{Name: "Программа 2", ID: 2},
		{Name: "Программа 3", ID: 3},
		{Name: "Программа 4", ID: 4},
		{Name: "Программа 5", ID: 5},
	}
	editMsg, boterr := f.TrainingProgramms(update.SentFrom().ID,
		update.CallbackQuery.Message.MessageID, page, pageCount, programms)
	if boterr != nil {
		return tg.StatusError, boterr
	}
	if _, err := ctx.Bot.Request(editMsg); err != nil {
		return tg.StatusError, tg.NewBotError(err.Error(), nil)
	}

	return tg.StatusOK, nil
}

func TrainingMenuEdit(ctx tg.Context, update tgbotapi.Update) (tg.HandleStatus, *tg.BotError) {
	editMsg, boterr := f.TrainingMenuEdit(update.SentFrom().ID, update.CallbackQuery.Message.MessageID)
	if boterr != nil {
		return tg.StatusError, boterr
	}
	if _, err := ctx.Bot.Request(editMsg); err != nil {
		return tg.StatusError, tg.NewBotError(err.Error(), nil)
	}

	return tg.StatusOK, nil
}
