// ./bot-telegram/cmd/main.go

package main

import (
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

	// Adding client
	clt := client.New(cfg.GetBackendFullURL(), httpclient.ClientConfig{
		APIToken: cfg.GetBackendAPIToken(),
	})
	srv.Context = srv.Context.WithValue("client", clt)

	// Logger setup & UserGet middleware
	{
		srv.Logger.Config.LogBehaviour = tg.LogAll
		srv.Apply(tg.DefaultLogging(srv), "Logger")

		srv.Apply(func(ctx tg.Context, update tgbotapi.Update, next tg.HandlerFunc) (status tg.HandleStatus, err *tg.BotError) {
			u, er := clt.GetUser(ctx.C, update.SentFrom().ID, update.SentFrom().UserName)
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

		command.Handle(tg.Command("start"), handlers.CommandStartHandler, "start")
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

		cbq := training.Group(tg.IsCallbackQuery, "CBQ")
		{
			// Main training menu
			cbq.Handle(tg.CallbackQuery("^T$"), handlers.Training_Edit_MainMenu, "training")

			// Exercises menu page \d
			cbq.Handle(tg.CallbackQuery(`^TEP\d*$`), handlers.Training_Edit_Exercises, "TEP")
			// Exercises search page \d
			cbq.Handle(tg.CallbackQuery(`^TES\d*$`), handlers.Training_Edit_Exercises_Search, "TES")
			{
				// Exercise create form
				cbq.Handle(tg.CallbackQuery(`^TEN$`), handlers.Training_Edit_Exercises_CreateForm, "TEN")
				// Exercise create form previous field (to up)
				cbq.Handle(tg.CallbackQuery(`^TENU$`), handlers.Training_Edit_Exercises_CreateForm_UP, "TENU")
				// Exercise create form next field (to down)
				cbq.Handle(tg.CallbackQuery(`^TEND$`), handlers.Training_Edit_Exercises_CreateForm_DOWN, "TEND")
				// Exercise create form weight unit preferences
				cbq.Handle(tg.CallbackQuery(`^TEN_KG$`), handlers.Training_Edit_Exercises_CreateForm_KG, "TEN_KG")
				cbq.Handle(tg.CallbackQuery(`^TEN_LBS$`), handlers.Training_Edit_Exercises_CreateForm_LBS, "TEN_LBS")
				cbq.Handle(tg.CallbackQuery(`^TEN_MIX$`), handlers.Training_Edit_Exercises_CreateForm_MIXED, "TEN_MIX")
				// Exercise create form save
				cbq.Handle(tg.CallbackQuery(`^TENS$`), handlers.Training_Edit_Exercises_CreateForm_SAVE, "TENS")
			}
			{
				// Exercise item menu id \d
				cbq.Handle(tg.CallbackQuery(`^TEI\d*$`), handlers.Training_Edit_Exercise, "TEI")
				{
					// Exercise item edit form id \d
					cbq.Handle(tg.CallbackQuery(`^TEE\d*$`), handlers.Training_Edit_Exercises_EditForm, "TEE")
					// Exercise item edit form previous field (to up)
					cbq.Handle(tg.CallbackQuery(`^TEEU\d*$`), handlers.Training_Edit_Exercises_EditForm_UP, "TEEU")
					// Exercise item edit form next field (to down)
					cbq.Handle(tg.CallbackQuery(`^TEED\d*$`), handlers.Training_Edit_Exercises_EditForm_DOWN, "TEED")
					// Exercise create form weight unit preferences
					cbq.Handle(tg.CallbackQuery(`^TEE_KG$`), handlers.Training_Edit_Exercises_EditForm_KG, "TEE_KG")
					cbq.Handle(tg.CallbackQuery(`^TEE_LBS$`), handlers.Training_Edit_Exercises_EditForm_LBS, "TEE_LBS")
					cbq.Handle(tg.CallbackQuery(`^TEE_MIX$`), handlers.Training_Edit_Exercises_EditForm_MIXED, "TEE_MIX")
					// Exercise item edit form save
					cbq.Handle(tg.CallbackQuery(`^TEES\d*$`), handlers.Training_Edit_Exercises_EditForm_SAVE, "TEEU")
				}
				{
					// Exercise item delete form id \d
					cbq.Handle(tg.CallbackQuery(`^TED\d*$`), handlers.Training_Edit_Exercise_DeleteForm, "TED")
					// Exercise item delete form confirm
					cbq.Handle(tg.CallbackQuery(`^TEDC\d*$`), handlers.Training_Edit_Exercise_DeleteForm_CONFIRM, "TEDC")
				}
				{
					// ExerciseEntry create form (performing solo exercise) exercise id \d
					cbq.Handle(tg.CallbackQuery(`^TEEnN\d*$`), tg.HandlerFuncStub /*handlers.Training_Edit_ExerciseEntry_CreateForm*/, "TEEnN")

				}
			}

			// Programms menu page \d
			cbq.Handle(tg.CallbackQuery(`^TPP\d*$`), handlers.Training_Edit_Programms, "TPP")
			{
				// Programm create form
				cbq.Handle(tg.CallbackQuery(`^TPN$`), tg.HandlerFuncStub /*handlers.Training_Edit_Programms_CreateForm*/, "TPN")
				// Programm create form previous field (to up)
				cbq.Handle(tg.CallbackQuery(`^TPNU$`), tg.HandlerFuncStub /*handlers.Training_Edit_Programms_CreateForm_UP*/, "TPNU")
				// Programm create form next field (to down)
				cbq.Handle(tg.CallbackQuery(`^TPND$`), tg.HandlerFuncStub /*handlers.Training_Edit_Programms_CreateForm_DOWN*/, "TPND")
				// Programm create form save
				cbq.Handle(tg.CallbackQuery(`^TPNS$`), tg.HandlerFuncStub /*handlers.Training_Edit_Programms_CreateForm_SAVE*/, "TPNS")
			}
			{
				// Programm item menu id \d
				cbq.Handle(tg.CallbackQuery(`^TPI\d*$`), handlers.Training_Edit_Exercise, "TPI")
			}
		}

		text := training.Group(tg.IsText, "TEXT")
		{
			text.Handle(tg.Text("🏋️"), handlers.Training_Send_MainMenu, "menu", tg.NewMiddleware(middlewares.DeleteSourceMessage, ""))
			text.Handle(matchers.UserContextEquals(`Training_Exercises`), handlers.Training_Input_Exercises, "TEP INPUT", tg.NewMiddleware(middlewares.DeleteSourceMessage, ""))
			text.Handle(matchers.UserContextEquals(`Training_Exercises_Search`), handlers.Training_Input_Exercises, "TEP INPUT", tg.NewMiddleware(middlewares.DeleteSourceMessage, ""))
			text.Handle(matchers.UserContextEquals(`Training_Exercises_CreateForm`), handlers.Training_Input_Exercises_CreateForm, "TEN INPUT", tg.NewMiddleware(middlewares.DeleteSourceMessage, ""))
			text.Handle(matchers.UserContextEquals(`Training_Exercises_EditForm`), handlers.Training_Input_Exercises_EditForm, "TEE INPUT", tg.NewMiddleware(middlewares.DeleteSourceMessage, ""))
		}

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
		// Отображение предпросмотра ссылок
		// Количество вывода упражнение на страницу
	}

	{
		srv.Handle(tg.CallbackQuery("-"), handlers.IgnoreHandler, "skipping callback")
		srv.Handle(tg.IsText, tg.HandlerFuncStub, "Any text stub",
			tg.NewMiddleware(middlewares.VerifyUsagePermission, "usagePermissionVerify"),
			tg.NewMiddleware(middlewares.DeleteSourceMessage, ""))
		srv.Handle(tg.Any(), tg.HandlerFuncStub, "Any stub", tg.NewMiddleware(middlewares.VerifyUsagePermission, "usagePermissionVerify"))
	}

	srv.Start()
}
