package factories

import (
	a "github.com/OvyFlash/telegram-bot-api"
)

func MainMenuReplyMarkup() a.ReplyKeyboardMarkup {

	kb := a.NewReplyKeyboard(
		a.NewKeyboardButtonRow(
			a.NewKeyboardButton("🍏"),
			a.NewKeyboardButton("🏋️"),
			a.NewKeyboardButton("👤"),
			a.NewKeyboardButton("📆"),
			a.NewKeyboardButton("⚙️"),
		),
	)

	return kb
}
