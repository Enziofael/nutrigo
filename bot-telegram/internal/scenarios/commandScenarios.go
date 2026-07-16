// ./bot-telegram/internal/scenarios/scenarios.go

package scenarios

import (
	"github.com/Enziofael/nutrigo/bot-telegram/internal/actions"
	"github.com/Enziofael/nutrigo/bot-telegram/internal/factory"
	models "github.com/Enziofael/nutrigo/shared/domain-models/system/user"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// S etting status + defer to end
//           v
// G etting data/state
//           v
// P rocessing data
//           v
// U pdating/saving data/state
//           v
// C reating responce ( scenario -> factory -> templateManager -> factory -> scenario)
//           v
// S ending responce

func SendSuggestUsageRequest(bot *tgbotapi.BotAPI, update tgbotapi.Update, user models.User) {
	stop := actions.StartChatAction(bot, update.Message.Chat.ID, tgbotapi.ChatTyping)
	defer stop()

	msg := factory.SuggestUsageRequest(update)

	actions.SendText(bot, msg)
	actions.DeleteMessage(bot, tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID))
}

func SendWaitToConfirm(bot *tgbotapi.BotAPI, update tgbotapi.Update, user models.User) {
	stop := actions.StartChatAction(bot, update.Message.Chat.ID, tgbotapi.ChatTyping)
	defer stop()

	msg := factory.WaitToConfirm(update)

	actions.SendText(bot, msg)
}

func SendMainMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, user models.User) {
	stop := actions.StartChatAction(bot, update.Message.Chat.ID, tgbotapi.ChatTyping)
	defer stop()

}

func SendRestrictedAppealForm(bot *tgbotapi.BotAPI, update tgbotapi.Update, user models.User) {
	stop := actions.StartChatAction(bot, update.Message.Chat.ID, tgbotapi.ChatTyping)
	defer stop()

}

func SendYouWasBanned(bot *tgbotapi.BotAPI, update tgbotapi.Update, user models.User) {

}

func SendAdminMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, user models.User) {

}

func SendUnknownCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update, user models.User) {

}
