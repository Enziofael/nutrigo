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

	processUpdate "github.com/Enziofael/nutrigo/frontend-bot/internal/actions"
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

	for update := range updates {

		if update.CallbackQuery != nil {
			if cfg.GetDebugLogging() {
				log.Printf("CallbackQuery from @%s: \"%s\"", update.CallbackQuery.From, update.CallbackData())
			}
			go processUpdate.CallbackQuery(bot, update)
		}

		if update.Message == nil {
			continue
		}

		if cfg.GetDebugLogging() {
			log.Printf("Message from @%s: \"%s\"", update.Message.From.UserName, update.Message.Text)
		}

		go processUpdate.Delete(bot, update)

		if update.Message.IsCommand() {
			go processUpdate.Command(bot, update)
		}

	}
}
