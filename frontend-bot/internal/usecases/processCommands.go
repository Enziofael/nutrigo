package processMsg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Command(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	switch update.Message.Command() {
	case "start":
		handleStartCommand(bot, update)
	case "admin":
		handleAdminCommand(bot, update)
	default:
		handleUnknownCommand(bot, update)
	}
}

func handleStartCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Starting bot")

	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}

func handleAdminCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Admin panel")

	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}

func handleUnknownCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command")

	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}
