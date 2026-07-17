// ./bot-telegram/internal/routers/CallbackQueryRouter.go

package processUpdate

import (
	"github.com/Enziofael/nutrigo/bot-telegram/internal/actions"
	models "github.com/Enziofael/nutrigo/shared/domain-models/system/user"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// # Router function
//
// Delegating commands to specific handlers
// by defining which exactly 'callback query' bot recieved
func CallbackQueryRouter(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	user := GetUserStub(update)

	switch update.CallbackQuery.Data {
	case "usage_request":
		handleUsageRequestCallbackQuery(bot, update, user)
	default:
		handleUnknownCallbackQuery(bot, update, user)
	}
}

// Handler for usage_request callback query
func handleUsageRequestCallbackQuery(bot *tgbotapi.BotAPI, update tgbotapi.Update, user models.User) {
	//Checking if user already exists
	// if so -> /start
	// if doesn't creates user in db, creates use_request in db, edits message

	switch user.Status {
	case models.StatusAdmin, models.StatusBanned, models.StatusConfirmed, models.StatusRequested, models.StatusRestricted:
		handleStartCommand(bot, update, user)
	default:
		if _, err := CreateUserStub(update); err != nil {
			panic("CreateUserStub() failed")
		}
		actions.EditRequested(bot, update)
	}
}

func handleUnknownCallbackQuery(bot *tgbotapi.BotAPI, update tgbotapi.Update, user models.User) {

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
