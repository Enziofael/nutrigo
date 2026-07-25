package factories

import (
	"log"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

func SuggestUsageRequestKeyboard(update tgbotapi.Update) tgbotapi.InlineKeyboardMarkup {
	text, err := tmplManager.RenderHTML("SuggestUsageRequest_btn_request", nil)
	if err != nil {
		log.Panicf("Render error: %v", err)
		text = err.Error()
	}

	btn := tgbotapi.NewInlineKeyboardButtonData(text, "usage_request")
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(btn),
	)

	return kb
}
