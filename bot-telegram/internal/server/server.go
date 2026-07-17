// ./bot-telegram/internal/server/server.go

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

	processUpdate "github.com/Enziofael/nutrigo/bot-telegram/internal/routers"
	//cfg "github.com/Enziofael/nutrigo/shared/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Starting bot cycle.
func Start(bot *tgbotapi.BotAPI) {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := bot.GetUpdatesChan(u)
	log.Printf("Ready to receive updates\n\n")

	// В цикле проходимся по каналу апдейтов, при получении раскидываем их соответствующим роутерам в горутины, сам цикл при этом идет дальше.
	// TODO: Для пакета лучше переделать с регистрацией роутеров по типу с передачей функций
	for update := range updates {

		//Callback update
		if update.CallbackQuery != nil {
			go processUpdate.CallbackQueryRouter(bot, update)
			continue
		}

		//Command update
		if update.Message != nil && update.Message.IsCommand() {
			go processUpdate.CommandRouter(bot, update)
			continue
		}
	}
}
