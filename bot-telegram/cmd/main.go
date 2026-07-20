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
)

func init() {
	m, err := templates.NewManager()
	if err != nil {
		log.Fatal(err)
	}
	factories.SetTemplateManager(m)
}

func main() {

	srv := tg.New(tg.ServerConfig{
		BotConfig: tg.BotConfig{APIToken: cfg.GetBotToken()},
	})

	clt := client.New(cfg.GetBackendFullURL(), httpclient.ClientConfig{
		APIToken: cfg.GetBackendAPIToken(),
	})
	srv.Context = context.WithValue(srv.Context, "client", clt)

	{
		callbackQuery := srv.Group(tg.IsCallbackQuery, "callbackquery group")

		callbackQuery.Use(tg.CallbackQuery("user_usage_request"), tg.HandlerFuncStub, "userUsage_request")
	}
	{
		command := srv.Group(tg.IsCommand, "command group")

		command.Use(tg.Command("start"), handlers.CommandStartHandler, "start command")
		command.Use(tg.Command("admin"), handlers.CommandAdminHandler, "admin command")
	}

	srv.Start()
}
