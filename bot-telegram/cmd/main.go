// ./bot-telegram/cmd/main.go

package main

import (
	"log"
	"strings"

	httpclient "github.com/Enziofael/nutrigo/backend/pkg/HTTPclient"
	client "github.com/Enziofael/nutrigo/bot-telegram/internal/client/v1"
	"github.com/Enziofael/nutrigo/bot-telegram/internal/factories"
	"github.com/Enziofael/nutrigo/bot-telegram/internal/handlers"
	"github.com/Enziofael/nutrigo/bot-telegram/internal/middlewares"
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
	srv.Context = srv.Context.WithValue("client", clt)

	{
		srv.Logger.Config.LogBehaviour = tg.LogAllDetailed & ^tg.LogDetailsOk
		srv.Apply(tg.DefaultLogging(srv), "Logger")

		srv.Apply(func(ctx tg.Context, update tgbotapi.Update, next tg.HandlerFunc) (status tg.HandleStatus, err *tg.BotError) {
			u, er := clt.GetUser(ctx.C, update)
			if er != nil {
				return tg.StatusError, tg.NewBotError(er.Error(), nil)
			}
			ctx = srv.Context.WithValue("user", u)
			return next(ctx, update)

		}, "UserGet middleware")
	}

	{
		callbackQuery := srv.Group(tg.IsCallbackQuery, "callbackQ")

		callbackQuery.Handle(tg.CallbackQuery("user_usage_request"), tg.HandlerFuncStub, "userUsage_request")
		callbackQuery.Handle(func(ctx tg.Context, update tgbotapi.Update) bool {
			if strings.HasPrefix(update.CallbackQuery.Data, "wdky") {
				return true
			}
			return false
		}, handlers.CallbackQuery_wdky, "WDKY")
		callbackQuery.Handle(func(ctx tg.Context, update tgbotapi.Update) bool {
			if strings.HasPrefix(update.CallbackQuery.Data, "nur") {
				return true
			}
			return false
		}, handlers.CallbackQuery_nur, "NUR")
	}
	{
		command := srv.Group(tg.IsCommand, "command")
		command.Apply(middlewares.VerifyUsagePermission, "usagePermissionVerify")

		command.Handle(tg.Command("start"), tg.HandlerFuncStub, "start")
		command.Handle(tg.Command("admin"), tg.HandlerFuncStub, "admin")

		command.Handle(tg.Command("ok"), handlers.OkHandler, "ok")
		command.Handle(tg.Command("warn"), handlers.WarnHandler, "warn")
		command.Handle(tg.Command("err"), handlers.ErrHandler, "err")
		command.Handle(tg.Command("fall"), handlers.FallHandler, "fall")
		command.Handle(tg.Command("cust"), handlers.CustomHandler, "custom")
		command.Handle(tg.Command("shut"), handlers.ShutdownHandler, "shutdown")
	}
	{
		food := srv.Group(tg.IsAny, "food")
		food.Apply(middlewares.VerifyUsagePermission, "usagePermissionVerify")
		food.Handle(tg.Text("🍏"), tg.HandlerFuncStub, "food", tg.NewMiddleware(middlewares.DeleteSourceMessage, ""))
	}
	{
		training := srv.Group(tg.IsAny, "training")
		training.Apply(middlewares.VerifyUsagePermission, "usagePermissionVerify")
		training.Handle(tg.Text("🏋️"), handlers.TrainingMenu, "training", tg.NewMiddleware(middlewares.DeleteSourceMessage, ""))
		training.Handle(tg.CallbackQuery("T"), handlers.TrainingMenuEdit, "training")
		training.Handle(func(ctx tg.Context, update tgbotapi.Update) bool {
			return strings.HasPrefix(update.CallbackData(), "TEP")
		}, handlers.TrainingMenu_TEP, "TEP")
		training.Handle(func(ctx tg.Context, update tgbotapi.Update) bool {
			return strings.HasPrefix(update.CallbackData(), "TPP")
		}, handlers.TrainingMenu_TPP, "TPP")
	}
	{
		calendar := srv.Group(tg.IsAny, "calendar")
		calendar.Apply(middlewares.VerifyUsagePermission, "usagePermissionVerify")
		calendar.Handle(tg.Text("📆"), tg.HandlerFuncStub, "calendar", tg.NewMiddleware(middlewares.DeleteSourceMessage, ""))
	}
	{
		profile := srv.Group(tg.IsAny, "profile")
		profile.Apply(middlewares.VerifyUsagePermission, "usagePermissionVerify")
		profile.Handle(tg.Text("👤"), tg.HandlerFuncStub, "profile", tg.NewMiddleware(middlewares.DeleteSourceMessage, ""))
	}
	{
		settings := srv.Group(tg.IsAny, "settings")
		settings.Apply(middlewares.VerifyUsagePermission, "usagePermissionVerify")
		settings.Handle(tg.Text("⚙️"), tg.HandlerFuncStub, "settings", tg.NewMiddleware(middlewares.DeleteSourceMessage, ""))
	}

	srv.Handle(tg.Any(), tg.HandlerFuncStub, "Any stub", tg.NewMiddleware(middlewares.VerifyUsagePermission, "usagePermissionVerify"))

	srv.Start()
}
