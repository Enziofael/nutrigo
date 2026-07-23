package defaults

import (
	"context"
	"log"

	t "github.com/Enziofael/nutrigo/bot-telegram/pkg/telegroni"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ============================================
// Stubs for testing
// ============================================

// HandlerFuncStub is a stub handler for testing.
func HandlerFuncStub(ctx context.Context, update tgbotapi.Update) (status string, err *t.BotError) {
	log.Printf("[Stub] Handler called for update: %v", update.UpdateID)
	return t.StatusOK, nil
}
