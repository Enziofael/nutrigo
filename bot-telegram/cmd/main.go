// ./bot-telegram/cmd/main.go

package main

import (
	"log"

	"github.com/Enziofael/nutrigo/bot-telegram/internal/factories"
	tg "github.com/Enziofael/nutrigo/bot-telegram/pkg/telegroni"
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

	srv := tg.New(cfg.GetBotToken())

	callbackQuery := srv.Group(tg.IsCallbackQuery)
	{
		callbackQuery.Use(tg.CallbackQuery("user_usage_request"), tg.HandlerFuncStub)
	}
	command := srv.Group(tg.IsCommand)
	{
		command.Use(tg.Command("start"), tg.HandlerFuncStub)
		command.Use(tg.Command("admin"), tg.HandlerFuncStub)
	}

	srv.Start()
}
