package actions

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartChatAction(bot *tgbotapi.BotAPI, chatID int64, action string) func() {
	ctx, stopChatAction := context.WithCancel(context.Background())

	config := tgbotapi.NewChatAction(chatID, action)
	if _, err := bot.Request(config); err != nil {
		HandleRequestError(bot, config)
	}

	go func() {

		ticker := time.NewTicker(4 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if _, err := bot.Request(config); err != nil {
					HandleRequestError(bot, config)
				}
			}
		}
	}()

	return stopChatAction
}
