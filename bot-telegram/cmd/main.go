// ./bot-telegram/cmd/main.go

package main

import (
	"log"

	"github.com/Enziofael/nutrigo/bot-telegram/internal/factories"
	app "github.com/Enziofael/nutrigo/bot-telegram/internal/server"
	"github.com/Enziofael/nutrigo/bot-telegram/internal/templates"
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

	srv := app.New(cfg.GetBotToken())

	callbackQuery := srv.Group(app.IsCallbackQuery)
	{
		callbackQuery.Use(app.CallbackQuery("user_usage_request"), app.HandlerFuncStub)
	}
	command := srv.Group(app.IsCommand)
	{
		command.Use(app.Command("start"), app.HandlerFuncStub)
		command.Use(app.Command("admin"), app.HandlerFuncStub)
	}

	srv.Start()
}
