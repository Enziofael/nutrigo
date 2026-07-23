package telegroni

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// type HandlerFunc func(bot *tgbotapi.BotAPI, update tgbotapi.Update)
func Send(bot *tgbotapi.BotAPI, config tgbotapi.Chattable) {
	if _, err := bot.Send(config); err != nil {
		HandleSendError(bot, config)
	}
}

func Request(bot *tgbotapi.BotAPI, config tgbotapi.Chattable) {
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