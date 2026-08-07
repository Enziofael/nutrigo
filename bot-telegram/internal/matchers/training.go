package matchers

import (
	"github.com/Enziofael/nutrigo/bot-telegram/internal/client/v1"
	"github.com/Enziofael/nutrigo/bot-telegram/pkg/telegroni"
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

func UserContextEquals(contextExpected string) telegroni.MatchFunc {
	return func(ctx telegroni.Context, update tgbotapi.Update) bool {
		if update.Message == nil {
			return false
		}
		clt := ctx.Value("client").(*client.Client)
		u, err := clt.GetUserStrict(ctx.C, update.SentFrom().ID)
		if err != nil {
			return false
		}
		if u.Context == contextExpected {
			return true
		}
		return false
	}
}
