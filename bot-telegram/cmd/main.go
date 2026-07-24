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
		srv.Apply(tg.DefaultLogging(srv), "Logger")

		srv.Apply(func(ctx context.Context, update tgbotapi.Update, next tg.HandlerFunc) (status tg.HandleStatus, err *tg.BotError) {
			u, er := clt.GetUser(ctx, update)
			if er != nil {
				return tg.StatusError, tg.NewBotError(er.Error(), nil)
			}
			ctx = context.WithValue(ctx, "user", u)
			return next(ctx, update)

		}, "UserGet middleware")
	}

	{
		callbackQuery := srv.Group(tg.IsCallbackQuery, "callbackQ")

		callbackQuery.Handle(tg.CallbackQuery("user_usage_request"), tg.HandlerFuncStub, "userUsage_request")
	}
	{
		command := srv.Group(tg.IsCommand, "command")

		command.Handle(tg.Command("start"), handlers.CommandStartHandler, "start")
		command.Handle(tg.Command("admin"), tg.HandlerFuncStub, "admin")

		command.Handle(tg.Command("ok"), handlers.OkHandler, "ok")
		command.Handle(tg.Command("warn"), handlers.WarnHandler, "warn")
		command.Handle(tg.Command("err"), handlers.ErrHandler, "err")
		command.Handle(tg.Command("fall"), handlers.FallHandler, "fall")
		command.Handle(tg.Command("cust"), handlers.CustomHandler, "custom")
		command.Handle(tg.Command("shut"), handlers.ShutdownHandler, "shutdown")
	}

	srv.Handle(tg.Any(), tg.HandlerFuncStub, "any")

	srv.Start()
}
