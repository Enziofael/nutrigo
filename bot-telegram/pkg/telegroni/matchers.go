package telegroni

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func IsMessage(ctx context.Context, update tgbotapi.Update) bool {
	return update.Message != nil && !update.Message.IsCommand()
}

func Message(s string) MatchFunc {
	return func(ctx context.Context, update tgbotapi.Update) bool {
		return IsMessage(ctx, update) && update.Message.Text == s
	}
}

func IsEditedMessage(ctx context.Context, update tgbotapi.Update) bool {
	return update.EditedMessage != nil
}

func EditedMessage() MatchFunc {
	return func(ctx context.Context, update tgbotapi.Update) bool {
		return IsEditedMessage(ctx, update)
	}
}

func IsInlineQuery(ctx context.Context, update tgbotapi.Update) bool {
	return update.InlineQuery != nil
}

func InlineQuery() MatchFunc {
	return func(ctx context.Context, update tgbotapi.Update) bool {
		return IsInlineQuery(ctx, update)
	}
}

func IsChosenInlineResult(ctx context.Context, update tgbotapi.Update) bool {
	return update.ChosenInlineResult != nil
}

func ChosenInlineResult() MatchFunc {
	return func(ctx context.Context, update tgbotapi.Update) bool {
		return IsChosenInlineResult(ctx, update)
	}
}

func IsPollAnswer(ctx context.Context, update tgbotapi.Update) bool {
	return update.PollAnswer != nil
}

func PollAnswer() MatchFunc {
	return func(ctx context.Context, update tgbotapi.Update) bool {
		return IsPollAnswer(ctx, update)
	}
}

func IsChannelPost(ctx context.Context, update tgbotapi.Update) bool {
	return update.ChannelPost != nil
}

func ChannelPost() MatchFunc {
	return func(ctx context.Context, update tgbotapi.Update) bool {
		return IsChannelPost(ctx, update)
	}
}

func IsShippingQuery(ctx context.Context, update tgbotapi.Update) bool {
	return update.ShippingQuery != nil
}

func ShippingQuery() MatchFunc {
	return func(ctx context.Context, update tgbotapi.Update) bool {
		return IsShippingQuery(ctx, update)
	}
}

func IsPoll(ctx context.Context, update tgbotapi.Update) bool {
	return update.Poll != nil
}

func Poll() MatchFunc {
	return func(ctx context.Context, update tgbotapi.Update) bool {
		return IsPoll(ctx, update)
	}
}

func IsMyChatMember(ctx context.Context, update tgbotapi.Update) bool {
	return update.MyChatMember != nil
}

func MyChatMember() MatchFunc {
	return func(ctx context.Context, update tgbotapi.Update) bool {
		return IsMyChatMember(ctx, update)
	}
}

func IsChatMember(ctx context.Context, update tgbotapi.Update) bool {
	return update.ChatMember != nil
}

func ChatMember() MatchFunc {
	return func(ctx context.Context, update tgbotapi.Update) bool {
		return IsChatMember(ctx, update)
	}
}

func IsChatJoinRequest(ctx context.Context, update tgbotapi.Update) bool {
	return update.ChatJoinRequest != nil
}

func ChatJoinRequest() MatchFunc {
	return func(ctx context.Context, update tgbotapi.Update) bool {
		return IsChatJoinRequest(ctx, update)
	}
}

func IsCommand(ctx context.Context, update tgbotapi.Update) bool {
	ok := update.Message != nil && update.Message.IsCommand()
	log.Printf("IsCommand: message=%v, ok=%v", update.Message != nil, ok)
	return ok
}

func Command(command string) MatchFunc {
	return func(ctx context.Context, update tgbotapi.Update) bool {
		log.Printf("Command: expected=%v, actual=%v", command, update.Message.Command())
		return IsCommand(ctx, update) &&
			update.Message.Command() == command
	}
}

func IsCallbackQuery(ctx context.Context, update tgbotapi.Update) bool {
	return update.CallbackQuery != nil
}

func CallbackQuery(data string) MatchFunc {
	return func(ctx context.Context, update tgbotapi.Update) bool {
		return IsCallbackQuery(ctx, update) &&
			update.CallbackQuery.Data == data
	}
}

func IsAny(ctx context.Context, update tgbotapi.Update) bool {
	return true
}

func Any() MatchFunc {
	return func(ctx context.Context, update tgbotapi.Update) bool { return IsAny(ctx, update) }
}
