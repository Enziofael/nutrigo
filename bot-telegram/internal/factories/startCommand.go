package factories

import (
	"log"

	v1 "github.com/Enziofael/nutrigo/shared/models/v1"
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

func MainMenu(update tgbotapi.Update, user *v1.User) tgbotapi.MessageConfig {
	text, err := tmplManager.RenderHTML("MainMenu", nil)
	if err != nil {
		log.Panicf("Render error: %v", err)
		text = err.Error()
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ParseMode = "HTML"

	return msg
}

func YouWasBanned(update tgbotapi.Update, user *v1.User) tgbotapi.MessageConfig {
	text, err := tmplManager.RenderHTML("YouWasBanned", nil)
	if err != nil {
		log.Panicf("Render error: %v", err)
		text = err.Error()
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ParseMode = "HTML"

	return msg
}

func PermissionRequested(update tgbotapi.Update, user *v1.User) tgbotapi.MessageConfig {
	text, err := tmplManager.RenderHTML("PermissionRequested", nil)
	if err != nil {
		log.Panicf("Render error: %v", err)
		text = err.Error()
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ParseMode = "HTML"

	return msg
}

func YouWasRestricted(update tgbotapi.Update, user *v1.User) tgbotapi.MessageConfig {
	text, err := tmplManager.RenderHTML("YouWasRestricted", nil)
	if err != nil {
		log.Panicf("Render error: %v", err)
		text = err.Error()
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ParseMode = "HTML"

	return msg
}

func SuggestPermissionRequest(update tgbotapi.Update, user *v1.User) tgbotapi.MessageConfig {
	text, err := tmplManager.RenderHTML("SuggestPermissionRequest", nil)
	if err != nil {
		log.Panicf("Render error: %v", err)
		text = err.Error()
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ParseMode = "HTML"

	return msg
}
