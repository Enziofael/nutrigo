package processMsg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Delete(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	deleteMsgConfig := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
	if _, err := bot.Request(deleteMsgConfig); err != nil {
		panic(err)
	}
}
