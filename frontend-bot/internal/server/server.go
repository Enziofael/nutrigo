// ./frontend-bot/internal/server/server.go

// Package server runs the Telegram bot's main event loop, polling for updates
// and delegating messages to specifig handlers.
//
// Usage:
//
//	bot, err := tgbotapi.NewBotAPI(token)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	server.Start(bot)
package server

import (
	"log"

	processUpdate "github.com/Enziofael/nutrigo/frontend-bot/internal/processUpdate"
	cfg "github.com/Enziofael/nutrigo/shared/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Starting bot cycle.
func Start(bot *tgbotapi.BotAPI) {
	if debugLogging := cfg.GetDebugLogging(); debugLogging {
		log.Printf("Debug logging enabled")
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := bot.GetUpdatesChan(u)
	log.Printf("Ready to receive updates\n\n")

	// В цикле проходимся по каналу апдейтов, при получении раскидываем их соответствующим хендлерам в горутины, сам цикл при этом идет дальше
	for update := range updates {
		//Callback update
		if update.CallbackQuery != nil {
			go processUpdate.CallbackQuery(bot, update)
			continue
		}

		//Command update
		if update.Message != nil && update.Message.IsCommand() {
			go processUpdate.Command(bot, update)
			continue
		}
	}
}
