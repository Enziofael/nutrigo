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
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
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
		srv.Logger.Config.LogUpdateDetails = false
		srv.Logger.Config.LogUpdateDetailsOnError = false
		srv.Apply(tg.DefaultLogging(srv), "Logger")

		srv.Apply(func(ctx context.Context, update tgbotapi.Update, next tg.HandlerFunc) (status string, err *tg.BotError) {
			u, er := clt.GetUser(ctx, update)
			if er != nil {
				return tg.StatusError, tg.NewBotError(er.Error(), nil)
			}
			ctx = context.WithValue(ctx, "user", u)
			return next(ctx, update)

		}, "UserGet middleware")
	}

	{
		callbackQuery := srv.Group(defaults.IsCallbackQuery, "callbackQ")

		callbackQuery.Handle(defaults.CallbackQuery("user_usage_request"), defaults.HandlerFuncStub, "userUsage_request")
	}
	{
		command := srv.Group(defaults.IsCommand, "command")

		command.Handle(defaults.Command("start"), handlers.CommandStartHandler, "start")
		command.Handle(defaults.Command("admin"), defaults.HandlerFuncStub, "admin")
	}

	srv.Start()
}
