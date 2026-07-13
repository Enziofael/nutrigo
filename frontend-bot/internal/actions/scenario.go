// ./frontend-bot/internal/actions/scenario.go

package actions

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TODO: фабрики сообщений?
// TODO: обработка сбоя отправки. Локальная запись отложенной отправки? Чтобы кидать в очередь на отправку. Но это уже в планировщик надо.

// Sends message informing that confirmation request is under consideration.
//
// Panics when sending failed
func SendRequested(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// TODO: сделать норм соо
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ваша заявка на использование бота находится на рассмотрении\nМы обязательно оповестим Вас о результате!")

	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}
func EditRequested(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	//TODO: фабрики!!!
	edit := tgbotapi.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID, "Ваша заявка на использование бота находится на рассмотрении\nМы обязательно оповестим Вас о результате!")
	if _, err := bot.Request(edit); err != nil {
		panic(err)
	}
}

// Sends message with main menu
//
// Panics when sending failed
func SendMain(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// TODO: сформировать менюшку
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "{Тут основное меню выводим с кнопками}")

	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}

// Sends message informing that user's access has been restricted with appeal button.
//
// Panics when sending failed
func SendRestricted(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// TODO: сделать норм соо
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Доступ к боту ограничен. Возможно ваша заявка была отклонена, либо вы были заблокированы. {кнопка подать аппеляцию}")
	//Кнопка меняет сообщение и просит написать текст для аппеляции. Так же можно добавить коментарий к блокировке.
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}

// Sends message informing that user has been banned
//
// Panics when sending failed
func SendBanned(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// TODO: сделать норм соо
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ваш аккаунт находится в блокировке. :(")
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}

// Sends message informing that user must has a permission to use this bot with request button
//
// Panics when sending failed
func SendNew(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Для того, чтобы пользоваться ботом, необходимо разрешение администратора")

	btn := tgbotapi.NewInlineKeyboardButtonData("Подать заявку", "usage_request")
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(btn),
	)
	msg.ReplyMarkup = kb
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}

// Sends message informing that command is unknown
//
// Panics when sending failed
func SendUnknown(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command")
	
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}
func EditUnknown(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	//TODO: фабрики!!!
	edit := tgbotapi.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID, "Странно, я не знаю что делает эта кнопка :\\")
	if _, err := bot.Request(edit); err != nil {
		panic(err)
	}
}

// Sends message with admin menu
//
// Panics when sending failed
func SendAdmin(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "{тут админ панель. Управление доступами пользователей и тд, логи}")
	// TODO: сформировать панельку
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}
