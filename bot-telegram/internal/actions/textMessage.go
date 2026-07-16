// ./bot-telegram/internal/actions/testMessage.go

package actions

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Send

func SendText(bot *tgbotapi.BotAPI, config tgbotapi.MessageConfig) {
	if _, err := bot.Send(config); err != nil {
		HandleSendError(bot, config)
	}
}

// Edit

func EditText(bot *tgbotapi.BotAPI, config tgbotapi.EditMessageTextConfig) {
	if _, err := bot.Request(config); err != nil {
		HandleRequestError(bot, config)
	}
}
