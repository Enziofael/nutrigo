package server

import (
	"log"

	processMsg "github.com/Enziofael/nutrigo/frontend-bot/internal/usecases"
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
		if update.Message == nil {
			continue
		}

		if cfg.GetDebugLogging() {
			log.Printf("Message from @%s: \"%s\"", update.Message.From.UserName, update.Message.Text)
		}

		go processMsg.Delete(bot, update)

		if update.Message.IsCommand() {
			go processMsg.Command(bot, update)
		} else {

		}
	}
}
