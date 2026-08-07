// ./bot-telegram/cmd/main.go

package main

import (
	"fmt"
	"log"

	httpclient "github.com/Enziofael/nutrigo/backend/pkg/HTTPclient"
	client "github.com/Enziofael/nutrigo/bot-telegram/internal/client/v1"
	"github.com/Enziofael/nutrigo/bot-telegram/internal/factories"
	"github.com/Enziofael/nutrigo/bot-telegram/internal/handlers"
	"github.com/Enziofael/nutrigo/bot-telegram/internal/matchers"
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
		srv.Logger.Config.LogBehaviour = tg.LogAll
		srv.Apply(tg.DefaultLogging(srv), "Logger")

		srv.Apply(func(ctx tg.Context, update tgbotapi.Update, next tg.HandlerFunc) (status tg.HandleStatus, err *tg.BotError) {
			u, er := clt.GetUser(ctx.C, update.SentFrom().ID, update.SentFrom().UserName)
			srv.Logger.Write(fmt.Sprintf("%v", u))
			if er != nil {
				return tg.StatusError, tg.NewBotErrorf("Can't get user at UserGet middleware: %w", er)
			}
			ctx = srv.Context.WithValue("user", u)
			return next(ctx, update)

		}, "UserGet middleware")
	}

	callbackQuery := srv.Group(tg.IsCallbackQuery, "callbackQ")
	{

		callbackQuery.Handle(tg.CallbackQuery("user_usage_request"), tg.HandlerFuncStub, "userUsage_request")
		callbackQuery.Handle(tg.CallbackQuery("^wdky"), handlers.CallbackQuery_wdky, "WDKY")
		callbackQuery.Handle(tg.CallbackQuery("^nur"), handlers.CallbackQuery_nur, "NUR")
	}

	command := srv.Group(tg.IsCommand, "command")
	{
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

	food := srv.Group(tg.IsAny, "food")
	{
		food.Apply(middlewares.VerifyUsagePermission, "usagePermissionVerify")
		food.Handle(tg.Text("🍏"), tg.HandlerFuncStub, "food", tg.NewMiddleware(middlewares.DeleteSourceMessage, ""))
	}

	training := srv.Group(tg.IsAny, "training")
	{
		training.Apply(middlewares.VerifyUsagePermission, "usagePermissionVerify")

		training.Handle(tg.Text("🏋️"), handlers.Training_Send_MainMenu, "menu", tg.NewMiddleware(middlewares.DeleteSourceMessage, ""))
		training.Handle(tg.CallbackQuery("^T$"), handlers.Training_Edit_MainMenu, "training")

		training.Handle(tg.CallbackQuery(`^TEP\d*$`), handlers.Training_Edit_Exercises, "TEP")
		training.Handle(tg.CallbackQuery(`^TEN$`), handlers.Training_Edit_Exercises_CreateForm, "TEN")
		training.Handle(tg.CallbackQuery(`^TENU$`), handlers.Training_Edit_Exercises_CreateForm_UP, "TENU")
		training.Handle(tg.CallbackQuery(`^TEND$`), handlers.Training_Edit_Exercises_CreateForm_DOWN, "TEND")
		training.Handle(tg.CallbackQuery(`^TENS$`), handlers.Training_Edit_Exercises_CreateForm_SAVE, "TENS")
		training.Handle(tg.CallbackQuery(`^TEI\d*$`), handlers.Training_Edit_Exercise, "TEI")

		training.Handle(matchers.UserContextEquals(`Training_Exercises_CreateForm`), handlers.Training_Input_Exercises_CreateForm, "TEN INPUT", tg.NewMiddleware(middlewares.DeleteSourceMessage, ""))

		training.Handle(tg.CallbackQuery(`^TPP\d*$`), handlers.Training_Edit_Programms, "TPP")
	}

	calendar := srv.Group(tg.IsAny, "calendar")
	{
		calendar.Apply(middlewares.VerifyUsagePermission, "usagePermissionVerify")
		calendar.Handle(tg.Text("📆"), tg.HandlerFuncStub, "calendar", tg.NewMiddleware(middlewares.DeleteSourceMessage, ""))
	}

	profile := srv.Group(tg.IsAny, "profile")
	{
		profile.Apply(middlewares.VerifyUsagePermission, "usagePermissionVerify")
		profile.Handle(tg.Text("👤"), tg.HandlerFuncStub, "profile", tg.NewMiddleware(middlewares.DeleteSourceMessage, ""))
	}

	settings := srv.Group(tg.IsAny, "settings")
	{
		settings.Apply(middlewares.VerifyUsagePermission, "usagePermissionVerify")
		settings.Handle(tg.Text("⚙️"), tg.HandlerFuncStub, "settings", tg.NewMiddleware(middlewares.DeleteSourceMessage, ""))
	}

	{
		srv.Handle(tg.CallbackQuery("-"), handlers.IgnoreHandler, "skipping callback")
		srv.Handle(tg.Any(), tg.HandlerFuncStub, "Any stub", tg.NewMiddleware(middlewares.VerifyUsagePermission, "usagePermissionVerify"))
	}

	srv.Start()
}
