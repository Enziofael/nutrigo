// ./bot-telegram/internal/routers/CommmandRouter.go

// ТАК КАК handle фунции маршрутизируют обработку ПЕРСОНАЛИЗИРОВАННО
// То их задача это:
// Получить от handler'а бота, апдейт, пользователя.
// Посмотреть данные пользователя и решить ЧТО нужно сделать.
// в зависимости от состояния
// И на основе этого вызвать СЛОЙ СЦЕНАРИЯ
// Сценарный слой вызовет сервисные слои с бизнес логикой, вызовет action'ы, установит новое состояние для пользователя.
// СЕРВИСНЫЕ СЛОИ уже будут работать с бд

package processUpdate

import (
	"github.com/Enziofael/nutrigo/bot-telegram/internal/scenarios"
	models "github.com/Enziofael/nutrigo/shared/domain-models/system/user"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// # Router function
//
// Delegating commands to specific handlers
// by defining which exactly command bot recieved
func CommandRouter(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	user := GetUserStub(update)

	switch update.Message.Command() {
	case "start":
		handleStartCommand(bot, update, user)
	case "admin":
		handleAdminCommand(bot, update, user)
	default:
		handleUnknownCommand(bot, update, user)
	}
}

// Handler for /start command
//
// Checks the user's status and sends response depending on its value.
func handleStartCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update, user models.User) {
	//Checking user's status
	// not found? -> Requesting form
	// requested? -> Info to wait
	// confirmed/admin? -> Main form
	// restricted? -> Appeal form
	// banned? -> Info message

	switch user.Status {
	case models.StatusRequested:
		scenarios.SendWaitToConfirm(bot, update, user)
	case models.StatusConfirmed, models.StatusAdmin:
		scenarios.SendMainMenu(bot, update, user)
	case models.StatusRestricted:
		scenarios.SendRestrictedAppealForm(bot, update, user)
	case models.StatusBanned:
		scenarios.SendYouWasBanned(bot, update, user)
	default:
		scenarios.SendSuggestUsageRequest(bot, update, user)
	}
}

// Handler for /admin command
//
// Checks the user's status and send response depending on its value.
func handleAdminCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update, user models.User) {
	//Checking user's status
	// admin? -> Send admin menu
	// else handleUnknownCommand()

	if user.Status == models.StatusAdmin {
		scenarios.SendAdminMenu(bot, update, user)
	} else {
		handleUnknownCommand(bot, update, user)
	}
}

// Handler for non-existent commands
//
// Panics if sending failed
func handleUnknownCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update, user models.User) {
	scenarios.SendUnknownCommand(bot, update, user)
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
