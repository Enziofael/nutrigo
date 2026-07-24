package defaults

import (
	"context"

	t "github.com/Enziofael/nutrigo/bot-telegram/pkg/telegroni"
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

// ============================================
// Stubs for testing
// ============================================

// HandlerFuncStub is a stub handler for testing.
func HandlerFuncStub(ctx context.Context, update tgbotapi.Update) (status string, err *t.BotError) {
	return t.StatusWarn, t.NewBotError("Stub handler matched", nil)
}
