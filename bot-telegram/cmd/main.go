// ./bot-telegram/cmd/main.go

package main

import (
	"context"
	"log"

	httpclient "github.com/Enziofael/nutrigo/backend/pkg/HTTPclient"
	client "github.com/Enziofael/nutrigo/bot-telegram/internal/client/v1"
	"github.com/Enziofael/nutrigo/bot-telegram/internal/factories"
	"github.com/Enziofael/nutrigo/bot-telegram/internal/handlers"
	"github.com/Enziofael/nutrigo/bot-telegram/internal/templates"
	tg "github.com/Enziofael/nutrigo/bot-telegram/pkg/telegroni"
	cfg "github.com/Enziofael/nutrigo/shared/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func init() {
	m, err := templates.NewManager()
	if err != nil {
		log.Fatal(err)
	}
	factories.SetTemplateManager(m)
}

func main() {

	srv, err := tg.New(tg.ServerConfig{
		BotConfig: tg.BotConfig{APIToken: cfg.GetBotToken()},
	})

	if err != nil {
		log.Fatal(err)
	}

	clt := client.New(cfg.GetBackendFullURL(), httpclient.ClientConfig{
		APIToken: cfg.GetBackendAPIToken(),
	})
	srv.Context = context.WithValue(srv.Context, "client", clt)

	{
		srv.Apply(tg.DefaultLogMiddleware)
		srv.Apply(tg.NewMiddleware(func(ctx context.Context, update tgbotapi.Update, next tg.HandlerFunc) (status string, err *tg.BotError) {
			u, er := clt.GetUser(ctx, update)
			if er != nil {
				return tg.StatusError, tg.NewBotError(er.Error(), nil)
			}
			ctx = context.WithValue(ctx, "user", u)
			return next(ctx, update)

		}, "UserGet middleware"))
	}

	{
		callbackQuery := srv.Group(tg.IsCallbackQuery, "callbackquery group")

		callbackQuery.Use(tg.CallbackQuery("user_usage_request"), tg.HandlerFuncStub, "userUsage_request")
	}
	{
		command := srv.Group(tg.IsCommand, "command group")

		command.Use(tg.Command("start"), handlers.CommandStartHandler, "start command")
		command.Use(tg.Command("admin"), tg.HandlerFuncStub, "admin command")
	}

	srv.Start()
}
