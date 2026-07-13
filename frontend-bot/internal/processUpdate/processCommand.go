// ./frontend-bot/internal/processUpdate/processCommands.go

package processUpdate

import (
	"github.com/Enziofael/nutrigo/frontend-bot/internal/actions"
	models "github.com/Enziofael/nutrigo/shared/domain-models/system/user"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Delegating commands to specific handlers
func Command(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	
	//go actions.StartStatus() или типо того. Статус бота (пишет, отправляет и тд)

	// Получаем пользователя. Обновляем state пользователя.
	// в handle надо передавать пользователя.

	go actions.DeleteMessage(bot, tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID))
	switch update.Message.Command() {
	case "start":
		handleStartCommand(bot, update)
	case "admin":
		handleAdminCommand(bot, update)
	default:
		handleUnknownCommand(bot, update)
	}
}

// Handler for /start command
//
// Checks the user's status and sends response depending on its value.
func handleStartCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	//Checking user's status
	// not found? -> Requesting form
	// requested? -> Info to wait
	// confirmed/admin? -> Main form
	// restricted? -> Appeal form
	// banned? -> Info message
	u := GetUserStub(update)

	switch u.Status {
	case models.StatusRequested:
		actions.SendRequested(bot, update)
	case models.StatusConfirmed, models.StatusAdmin:
		actions.SendMain(bot, update)
	case models.StatusRestricted:
		actions.SendRestricted(bot, update)
	case models.StatusBanned:
		actions.SendBanned(bot, update)
	default:
		actions.SendNew(bot, update)
	}
}

// Handler for /admin command
//
// Checks the user's status and send response depending on its value.
func handleAdminCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	//Checking user's status
	// admin? -> Send admin menu
	// else handleUnknownCommand()
	u := GetUserStub(update)

	if u.Status == models.StatusAdmin {
		actions.SendAdmin(bot, update)
	} else {
		handleUnknownCommand(bot, update)
	}
}

// Handler for non-existent commands
//
// Panics if sending failed
func handleUnknownCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	actions.SendUnknown(bot, update)
}

// DEV-ONLY stub function
// Return a user with status depending on command's argument
// Usage:
//
//	u := GetUserStub(update)
//
// In chat:
//
// /command status
func GetUserStub(update tgbotapi.Update) models.User {
	if update.Message != nil && update.Message.IsCommand() {
		return models.User{Status: models.Status(update.Message.CommandArguments())}
	} else if update.CallbackQuery != nil && update.CallbackQuery.Data == "usage_request" {
		return models.User{}
	} else {
		return models.User{}
	}
}
