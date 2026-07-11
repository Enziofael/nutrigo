package server

import (
	"log"

	processMsg "github.com/Enziofael/nutrigo/frontend-bot/internal/usecases"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Starting bot cycle.
func Start(bot *tgbotapi.BotAPI) {
	log.Printf("Ready to receive updates")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		//Deleting recieved user's message
		go processMsg.Delete(bot, update)

		if update.Message.IsCommand() {
			go processMsg.Command(bot, update)
		}
	}
}
