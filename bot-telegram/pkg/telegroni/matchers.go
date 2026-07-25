package telegroni

import (
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

func IsText(ctx Context, update tgbotapi.Update) bool {
	return update.Message != nil && update.Message.Text != "" && !update.Message.IsCommand()
}

func Text(s string) MatchFunc {
	return func(ctx Context, update tgbotapi.Update) bool {
		return IsText(ctx, update) && update.Message.Text == s
	}
}

func IsEditedMessage(ctx Context, update tgbotapi.Update) bool {
	return update.EditedMessage != nil
}

func EditedMessage() MatchFunc {
	return func(ctx Context, update tgbotapi.Update) bool {
		return IsEditedMessage(ctx, update)
	}
}

func IsInlineQuery(ctx Context, update tgbotapi.Update) bool {
	return update.InlineQuery != nil
}

func InlineQuery() MatchFunc {
	return func(ctx Context, update tgbotapi.Update) bool {
		return IsInlineQuery(ctx, update)
	}
}

func IsChosenInlineResult(ctx Context, update tgbotapi.Update) bool {
	return update.ChosenInlineResult != nil
}

func ChosenInlineResult() MatchFunc {
	return func(ctx Context, update tgbotapi.Update) bool {
		return IsChosenInlineResult(ctx, update)
	}
}

func IsPollAnswer(ctx Context, update tgbotapi.Update) bool {
	return update.PollAnswer != nil
}

func PollAnswer() MatchFunc {
	return func(ctx Context, update tgbotapi.Update) bool {
		return IsPollAnswer(ctx, update)
	}
}

func IsChannelPost(ctx Context, update tgbotapi.Update) bool {
	return update.ChannelPost != nil
}

func ChannelPost() MatchFunc {
	return func(ctx Context, update tgbotapi.Update) bool {
		return IsChannelPost(ctx, update)
	}
}

func IsShippingQuery(ctx Context, update tgbotapi.Update) bool {
	return update.ShippingQuery != nil
}

func ShippingQuery() MatchFunc {
	return func(ctx Context, update tgbotapi.Update) bool {
		return IsShippingQuery(ctx, update)
	}
}

func IsPoll(ctx Context, update tgbotapi.Update) bool {
	return update.Poll != nil
}

func Poll() MatchFunc {
	return func(ctx Context, update tgbotapi.Update) bool {
		return IsPoll(ctx, update)
	}
}

func IsMyChatMember(ctx Context, update tgbotapi.Update) bool {
	return update.MyChatMember != nil
}

func MyChatMember() MatchFunc {
	return func(ctx Context, update tgbotapi.Update) bool {
		return IsMyChatMember(ctx, update)
	}
}

func IsChatMember(ctx Context, update tgbotapi.Update) bool {
	return update.ChatMember != nil
}

func ChatMember() MatchFunc {
	return func(ctx Context, update tgbotapi.Update) bool {
		return IsChatMember(ctx, update)
	}
}

func IsChatJoinRequest(ctx Context, update tgbotapi.Update) bool {
	return update.ChatJoinRequest != nil
}

func ChatJoinRequest() MatchFunc {
	return func(ctx Context, update tgbotapi.Update) bool {
		return IsChatJoinRequest(ctx, update)
	}
}

func IsCommand(ctx Context, update tgbotapi.Update) bool {
	ok := update.Message != nil && update.Message.IsCommand()
	return ok
}

func Command(command string) MatchFunc {
	return func(ctx Context, update tgbotapi.Update) bool {
		return IsCommand(ctx, update) &&
			update.Message.Command() == command
	}
}

func IsCallbackQuery(ctx Context, update tgbotapi.Update) bool {
	return update.CallbackQuery != nil
}

func CallbackQuery(data string) MatchFunc {
	return func(ctx Context, update tgbotapi.Update) bool {
		return IsCallbackQuery(ctx, update) &&
			update.CallbackQuery.Data == data
	}
}

func IsAny(ctx Context, update tgbotapi.Update) bool {
	return true
}

func Any() MatchFunc {
	return func(ctx Context, update tgbotapi.Update) bool { return IsAny(ctx, update) }
}
