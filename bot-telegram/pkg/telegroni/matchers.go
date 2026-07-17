package telegroni

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func IsMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
	return update.Message != nil && !update.Message.IsCommand()
}

func Message(s string) MatchFunc {
	return func(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
		return IsMessage(bot, update) && update.Message.Text == s
	}
}

func IsEditedMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
	return update.EditedMessage != nil
}

func EditedMessage() MatchFunc {
	return func(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
		return IsEditedMessage(bot, update)
	}
}

func IsInlineQuery(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
	return update.InlineQuery != nil
}

func InlineQuery() MatchFunc {
	return func(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
		return IsInlineQuery(bot, update)
	}
}

func IsChosenInlineResult(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
	return update.ChosenInlineResult != nil
}

func ChosenInlineResult() MatchFunc {
	return func(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
		return IsChosenInlineResult(bot, update)
	}
}

func IsPollAnswer(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
	return update.PollAnswer != nil
}

func PollAnswer() MatchFunc {
	return func(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
		return IsPollAnswer(bot, update)
	}
}

func IsChannelPost(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
	return update.ChannelPost != nil
}

func ChannelPost() MatchFunc {
	return func(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
		return IsChannelPost(bot, update)
	}
}

func IsShippingQuery(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
	return update.ShippingQuery != nil
}

func ShippingQuery() MatchFunc {
	return func(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
		return IsShippingQuery(bot, update)
	}
}

func IsPoll(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
	return update.Poll != nil
}

func Poll() MatchFunc {
	return func(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
		return IsPoll(bot, update)
	}
}

func IsMyChatMember(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
	return update.MyChatMember != nil
}

func MyChatMember() MatchFunc {
	return func(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
		return IsMyChatMember(bot, update)
	}
}

func IsChatMember(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
	return update.ChatMember != nil
}

func ChatMember() MatchFunc {
	return func(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
		return IsChatMember(bot, update)
	}
}

func IsChatJoinRequest(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
	return update.ChatJoinRequest != nil
}

func ChatJoinRequest() MatchFunc {
	return func(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
		return IsChatJoinRequest(bot, update)
	}
}

func IsCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
	return update.Message != nil && update.Message.IsCommand()
}

func Command(command string) MatchFunc {
	return func(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
		return IsCommand(bot, update) &&
			update.Message.Command() == command
	}
}

func IsCallbackQuery(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
	return update.CallbackQuery != nil
}

func CallbackQuery(data string) MatchFunc {
	return func(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
		return IsCallbackQuery(bot, update) &&
			update.CallbackQuery.Data == data
	}
}

func IsAny(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
	return true
}

func Any() MatchFunc {
	return func(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool { return IsAny(bot, update) }
}
