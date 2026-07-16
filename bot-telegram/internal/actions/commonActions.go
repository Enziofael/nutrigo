// ./bot-telegram/internal/actions/deleteMessage.go

package actions

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	Send   = "action_send"
	Edit   = "action_edit"
	Delete = "action_delete"
)

func DeleteMessage(bot *tgbotapi.BotAPI, config tgbotapi.DeleteMessageConfig) {
	if _, err := bot.Request(config); err != nil {
		HandleRequestError(bot, config)
	}
}

func HandleSendError(bot *tgbotapi.BotAPI, config tgbotapi.Chattable) {
	panic(config)
}

func HandleRequestError(bot *tgbotapi.BotAPI, config tgbotapi.Chattable) {
	panic(config)
}
