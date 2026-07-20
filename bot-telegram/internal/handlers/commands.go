package handlers

import (
	"context"
	"log"

	client "github.com/Enziofael/nutrigo/bot-telegram/internal/client/v1"
	"github.com/Enziofael/nutrigo/bot-telegram/internal/factories"
	v1 "github.com/Enziofael/nutrigo/shared/models/v1"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CommandStartHandler(ctx context.Context, update tgbotapi.Update) {
	user, err := ctx.Value("client").(*client.Client).GetUser(ctx, update)
	if err != nil {
		log.Printf("ERROR:%v", err)
	}
	if user == nil {
		user = &v1.User{
			Status: "unknown",
		}
	}
	log.Print(user)
	msg := factories.GetMe(update, user.Status)
	ctx.Value("bot").(*tgbotapi.BotAPI).Send(msg)
}

func CommandAdminHandler(ctx context.Context, update tgbotapi.Update) {

}
