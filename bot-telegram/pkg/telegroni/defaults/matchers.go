package defaults

import (
	"context"

	t "github.com/Enziofael/nutrigo/bot-telegram/pkg/telegroni"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func IsMessage(ctx context.Context, update tgbotapi.Update) bool {
	return update.Message != nil && !update.Message.IsCommand()
}

func Message(s string) t.MatchFunc {
	return func(ctx context.Context, update tgbotapi.Update) bool {
		return IsMessage(ctx, update) && update.Message.Text == s
	}
}

func IsEditedMessage(ctx context.Context, update tgbotapi.Update) bool {
	return update.EditedMessage != nil
}

func EditedMessage() t.MatchFunc {
	return func(ctx context.Context, update tgbotapi.Update) bool {
		return IsEditedMessage(ctx, update)
	}
}

func IsInlineQuery(ctx context.Context, update tgbotapi.Update) bool {
	return update.InlineQuery != nil
}

func InlineQuery() t.MatchFunc {
	return func(ctx context.Context, update tgbotapi.Update) bool {
		return IsInlineQuery(ctx, update)
	}
}

func IsChosenInlineResult(ctx context.Context, update tgbotapi.Update) bool {
	return update.ChosenInlineResult != nil
}

func ChosenInlineResult() t.MatchFunc {
	return func(ctx context.Context, update tgbotapi.Update) bool {
		return IsChosenInlineResult(ctx, update)
	}
}

func IsPollAnswer(ctx context.Context, update tgbotapi.Update) bool {
	return update.PollAnswer != nil
}

func PollAnswer() t.MatchFunc {
	return func(ctx context.Context, update tgbotapi.Update) bool {
		return IsPollAnswer(ctx, update)
	}
}

func IsChannelPost(ctx context.Context, update tgbotapi.Update) bool {
	return update.ChannelPost != nil
}

func ChannelPost() t.MatchFunc {
	return func(ctx context.Context, update tgbotapi.Update) bool {
		return IsChannelPost(ctx, update)
	}
}

func IsShippingQuery(ctx context.Context, update tgbotapi.Update) bool {
	return update.ShippingQuery != nil
}

func ShippingQuery() t.MatchFunc {
	return func(ctx context.Context, update tgbotapi.Update) bool {
		return IsShippingQuery(ctx, update)
	}
}

func IsPoll(ctx context.Context, update tgbotapi.Update) bool {
	return update.Poll != nil
}

func Poll() t.MatchFunc {
	return func(ctx context.Context, update tgbotapi.Update) bool {
		return IsPoll(ctx, update)
	}
}

func IsMyChatMember(ctx context.Context, update tgbotapi.Update) bool {
	return update.MyChatMember != nil
}

func MyChatMember() t.MatchFunc {
	return func(ctx context.Context, update tgbotapi.Update) bool {
		return IsMyChatMember(ctx, update)
	}
}

func IsChatMember(ctx context.Context, update tgbotapi.Update) bool {
	return update.ChatMember != nil
}

func ChatMember() t.MatchFunc {
	return func(ctx context.Context, update tgbotapi.Update) bool {
		return IsChatMember(ctx, update)
	}
}

func IsChatJoinRequest(ctx context.Context, update tgbotapi.Update) bool {
	return update.ChatJoinRequest != nil
}

func ChatJoinRequest() t.MatchFunc {
	return func(ctx context.Context, update tgbotapi.Update) bool {
		return IsChatJoinRequest(ctx, update)
	}
}

func IsCommand(ctx context.Context, update tgbotapi.Update) bool {
	ok := update.Message != nil && update.Message.IsCommand()
	return ok
}

func Command(command string) t.MatchFunc {
	return func(ctx context.Context, update tgbotapi.Update) bool {
		return IsCommand(ctx, update) &&
			update.Message.Command() == command
	}
}

func IsCallbackQuery(ctx context.Context, update tgbotapi.Update) bool {
	return update.CallbackQuery != nil
}

func CallbackQuery(data string) t.MatchFunc {
	return func(ctx context.Context, update tgbotapi.Update) bool {
		return IsCallbackQuery(ctx, update) &&
			update.CallbackQuery.Data == data
	}
}

func IsAny(ctx context.Context, update tgbotapi.Update) bool {
	return true
}

func Any() t.MatchFunc {
	return func(ctx context.Context, update tgbotapi.Update) bool { return IsAny(ctx, update) }
}
