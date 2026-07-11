package processMsg

import (
	"log"

	cfg "github.com/Enziofael/nutrigo/shared/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Delete(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if cfg.GetDebugLogging() {
		log.Printf("Deleting message(text: \"%s\", id: %d)", update.Message.Text, update.Message.MessageID)
	}

	deleteMsgConfig := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
	if _, err := bot.Request(deleteMsgConfig); err != nil {
		panic(err)
	}

	if cfg.GetDebugLogging() {
		log.Print("Deleted!")
	}
}
