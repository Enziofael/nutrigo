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
	"github.com/Enziofael/nutrigo/bot-telegram/pkg/telegroni/defaults"
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

	srv := tg.New(tg.NewConfig(cfg.GetBotToken()))

	clt := client.New(cfg.GetBackendFullURL(), httpclient.ClientConfig{
		APIToken: cfg.GetBackendAPIToken(),
	})
	srv.Context = context.WithValue(srv.Context, "client", clt)

	{
		//srv.Apply(tg.DefaultLogging())
		srv.Apply(tg.NewMiddleware(func(ctx context.Context, update tgbotapi.Update, next tg.HandlerFunc, "UserGet") (status string, err *tg.BotError) {
			u, er := clt.GetUser(ctx, update)
			if er != nil {
				return tg.StatusError, tg.NewBotError(er.Error(), nil)
			}
			ctx = context.WithValue(ctx, "user", u)
			return next(ctx, update)

		}, "UserGet middleware"))
	}

	{
		callbackQuery := srv.Group(defaults.IsCallbackQuery, "callbackquery group")

		callbackQuery.Handle(defaults.CallbackQuery("user_usage_request"), defaults.HandlerFuncStub, "userUsage_request")
	}
	{
		command := srv.Group(defaults.IsCommand, "command group")

		command.Handle(defaults.Command("start"), handlers.CommandStartHandler, "start command")
		command.Handle(defaults.Command("admin"), defaults.HandlerFuncStub, "admin command")
	}

	srv.Start()
}
