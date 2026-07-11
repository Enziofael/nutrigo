package main

import (
	"log"

	"github.com/Enziofael/nutrigo/frontend-bot/internal/server"
	env "github.com/Enziofael/nutrigo/shared/enviroment"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(env.Get("TELEGRAM_BOT_API_TOKEN"))
	if err != nil {
		log.Fatalf("Panic: bot creation failed. Invalid token? Error: \"%s\"", err)
	}

	//------
	// bot.Debug = true
	//------

	log.Printf(`Bot init successful.
	Authorized on account @%s
	ID: %d
	CanJoinGroups: %t
	CanReadAllGroupMessages: %t
	SupportsInlineQueries: %t`,
		bot.Self.UserName, bot.Self.ID, bot.Self.CanJoinGroups, bot.Self.CanReadAllGroupMessages, bot.Self.SupportsInlineQueries)

	server.Start(bot)
}
