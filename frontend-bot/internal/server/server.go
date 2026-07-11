package server

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Start(bot *tgbotapi.BotAPI) {
	log.Printf("Ready to receive updates")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		go deleteMsg(bot, update)

		if update.Message.IsCommand() {
			go processCommand(bot, update)
		}
	}
}

func deleteMsg(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	deleteMsgConfig := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
	if _, err := bot.Request(deleteMsgConfig); err != nil {
		panic(err)
	}
}

func processCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

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
