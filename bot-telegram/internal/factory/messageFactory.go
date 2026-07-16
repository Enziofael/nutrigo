package factory

import (
	"log"

	"github.com/Enziofael/nutrigo/bot-telegram/internal/templates"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var tmplManager *templates.Manager

func SetTemplateManager(m *templates.Manager) {
	tmplManager = m
}

func SuggestUsageRequest(update tgbotapi.Update) tgbotapi.MessageConfig {
	text, err := tmplManager.RenderHTML("SuggestUsageRequest", nil)
	if err != nil {
		log.Panicf("Render error: %v", err)
		text = err.Error()
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ParseMode = "HTML"

	kb := SuggestUsageRequestKeyboard(update)
	msg.ReplyMarkup = kb

	return msg
}

func WaitToConfirm(update tgbotapi.Update) tgbotapi.MessageConfig {
	text, err := tmplManager.RenderHTML("WaitToConfirm", nil)
	if err != nil {
		log.Panicf("Render error: %v", err)
		text = err.Error()
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ParseMode = "HTML"

	return msg
}
