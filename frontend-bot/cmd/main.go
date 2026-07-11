package main

import (
	"log"

	env "github.com/Enziofael/nutrigo/shared/enviroment"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	env.Load()

	bot, err := tgbotapi.NewBotAPI(env.Get("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	//bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s\t %d", update.Message.From.UserName, update.Message.Text, update.Message.From.ID)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
}
