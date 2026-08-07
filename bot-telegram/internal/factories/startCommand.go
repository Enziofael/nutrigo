package factories

import (
	"log"

	tg "github.com/Enziofael/nutrigo/bot-telegram/pkg/telegroni"
	v1 "github.com/Enziofael/nutrigo/shared/models/v1"
	a "github.com/OvyFlash/telegram-bot-api"
)

func YouWasBanned(update a.Update, user *v1.User) a.MessageConfig {
	text, err := tmplManager.RenderHTML("YouWasBanned", nil)
	if err != nil {
		log.Panicf("Render error: %v", err)
		text = err.Error()
	}

	msg := a.NewMessage(update.Message.Chat.ID, text)
	msg.ParseMode = "HTML"

	return msg
}

func PermissionRequested(update a.Update, user *v1.User) a.MessageConfig {
	text, err := tmplManager.RenderHTML("PermissionRequested", nil)
	if err != nil {
		log.Panicf("Render error: %v", err)
		text = err.Error()
	}

	msg := a.NewMessage(update.Message.Chat.ID, text)
	msg.ParseMode = "HTML"

	return msg
}

func YouWasRestricted(update a.Update, user *v1.User) a.MessageConfig {
	text, err := tmplManager.RenderHTML("YouWasRestricted", nil)
	if err != nil {
		log.Panicf("Render error: %v", err)
		text = err.Error()
	}

	msg := a.NewMessage(update.Message.Chat.ID, text)
	msg.ParseMode = "HTML"

	return msg
}

func SuggestPermissionRequest(update a.Update, user *v1.User) a.MessageConfig {
	text, err := tmplManager.RenderHTML("SuggestPermissionRequest", nil)
	if err != nil {
		log.Panicf("Render error: %v", err)
		text = err.Error()
	}

	msg := a.NewMessage(update.Message.Chat.ID, text)
	msg.ParseMode = "HTML"

	return msg
}

func Common_Send_ReplyMenu(tgID int64) (a.MessageConfig, *tg.BotError) {
	text, err := tmplManager.RenderHTML("Hello", nil)
	if err != nil {
		return a.MessageConfig{}, tg.NewBotErrorf("Can't render html at Training_Send_MainMenu: %w", err)
	}

	msg := a.NewMessage(tgID, text)

	kb := MainMenuReplyMarkup()

	msg.ReplyMarkup = kb

	return msg, nil
}
