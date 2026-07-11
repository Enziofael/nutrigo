// ./frontend-bot/internal/usecases/processCallback.go

package processUpdate

import (
	"github.com/Enziofael/nutrigo/frontend-bot/internal/actions"
	models "github.com/Enziofael/nutrigo/shared/domain-models/system/user"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Delegating callback queries to specific handlers
func CallbackQuery(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	switch update.CallbackQuery.Data {
	case "usage_request":
		handleUsageRequestCallbackQuery(bot, update)
	default:
		handleUnknownCallbackQuery(bot, update)
	}
}

// Handler for usage_request callback query
func handleUsageRequestCallbackQuery(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	//Checking if user already exists
	// if so -> /start
	// if doesn't creates user in db, creates use_request in db, edits message
	u := GetUserStub(update)

	switch u.Status {
	case models.StatusAdmin, models.StatusBanned, models.StatusConfirmed, models.StatusRequested, models.StatusRestricted:
		handleStartCommand(bot, update)
	default:
		if _, err := CreateUserStub(update); err != nil {
			panic("CreateUserStub() failed")
		}
		actions.EditRequested(bot, update)
	}
}

func handleUnknownCallbackQuery(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	actions.EditUnknown(bot, update)
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
func CreateUserStub(update tgbotapi.Update) (models.User, error) {
	//add user to db
	// add request to db
	return models.User{Status: models.StatusRestricted}, nil
}
