package factories

import (
	"strconv"

	t "github.com/Enziofael/nutrigo/bot-telegram/internal/templates"
	tg "github.com/Enziofael/nutrigo/bot-telegram/pkg/telegroni"
	models "github.com/Enziofael/nutrigo/shared/models/v1"
	a "github.com/OvyFlash/telegram-bot-api"
)

var tmplManager *t.Manager

func SetTemplateManager(m *t.Manager) {
	tmplManager = m
}

func WeDontKnowYou(ctx tg.Context, u a.Update) (a.MessageConfig, *tg.BotError) {
	text, err := tmplManager.RenderHTML("WeDontKnowYou", nil)
	if err != nil {
		return a.MessageConfig{}, tg.NewBotError(err.Error(), nil)
	}

	msg := a.NewMessage(u.FromChat().ID, text)

	kb, boterr := WeDontKnowYouReplyMarkup(ctx.Value("user").(*models.User).TgID)
	if boterr != nil {
		return a.MessageConfig{}, boterr
	}
	msg.ReplyMarkup = kb

	msg.ParseMode = "HTML"

	return msg, nil
}

func WeDontKnowYouReplyMarkup(tgId int64) (a.InlineKeyboardMarkup, *tg.BotError) {
	text, err := tmplManager.RenderHTML("WeDontKnowYouReplyMarkup", nil)
	if err != nil {
		return a.InlineKeyboardMarkup{}, tg.NewBotError(err.Error(), nil)
	}

	kb := a.NewInlineKeyboardMarkup(
		a.NewInlineKeyboardRow(
			a.NewInlineKeyboardButtonData(text, "wdky"),
		),
	)

	return kb, nil
}

func NewUsageRequest(ctx tg.Context, u a.Update) (a.MessageConfig, *tg.BotError) {
	text, err := tmplManager.RenderHTML("NewUsageRequest", ctx.Value("user"))
	if err != nil {
		return a.MessageConfig{}, tg.NewBotError(err.Error(), nil)
	}

	msg := a.NewMessage(490590745, text)

	kb, boterr := NewUsageRequestReplyMarkup(ctx.Value("user").(*models.User).TgID)
	if boterr != nil {
		return a.MessageConfig{}, boterr
	}
	msg.ReplyMarkup = kb

	msg.ParseMode = "HTML"

	return msg, nil
}

func NewUsageRequestReplyMarkup(tgId int64) (a.InlineKeyboardMarkup, *tg.BotError) {
	textConfirm, err := tmplManager.RenderHTML("NewUsageRequestReplyMarkup_textConfirm", nil)
	if err != nil {
		return a.InlineKeyboardMarkup{}, tg.NewBotError(err.Error(), nil)
	}
	textReject, err := tmplManager.RenderHTML("NewUsageRequestReplyMarkup_textReject", nil)
	if err != nil {
		return a.InlineKeyboardMarkup{}, tg.NewBotError(err.Error(), nil)
	}
	textBlock, err := tmplManager.RenderHTML("NewUsageRequestReplyMarkup_textBlock", nil)
	if err != nil {
		return a.InlineKeyboardMarkup{}, tg.NewBotError(err.Error(), nil)
	}

	kb := a.NewInlineKeyboardMarkup(
		a.NewInlineKeyboardRow(
			a.NewInlineKeyboardButtonData(textConfirm, "nur_c"+strconv.Itoa(int(tgId))),
			a.NewInlineKeyboardButtonData(textReject, "nur_r"+strconv.Itoa(int(tgId))),
			a.NewInlineKeyboardButtonData(textBlock, "nur_b"+strconv.Itoa(int(tgId))),
		),
	)

	return kb, nil
}

func YourUsageRequestSend(ctx tg.Context, u a.Update) (a.MessageConfig, *tg.BotError) {
	text, err := tmplManager.RenderHTML("YourUsageRequestSend", nil)
	if err != nil {
		return a.MessageConfig{}, tg.NewBotError(err.Error(), nil)
	}

	msg := a.NewMessage(u.FromChat().ID, text)

	msg.ParseMode = "HTML"

	return msg, nil
}

func ConfirmUsageRequest(ctx tg.Context, u a.Update) (a.EditMessageTextConfig, *tg.BotError) {
	text, err := tmplManager.RenderHTML("ConfirmUsageRequest", u.CallbackQuery.Message)
	if err != nil {
		return a.EditMessageTextConfig{}, tg.NewBotError(err.Error(), nil)
	}

	msg := a.NewEditMessageText(u.FromChat().ID, u.CallbackQuery.Message.MessageID, text)

	msg.ReplyMarkup = nil

	msg.ParseMode = "HTML"

	return msg, nil
}

func RejectUsageRequest(ctx tg.Context, u a.Update) (a.EditMessageTextConfig, *tg.BotError) {
	text, err := tmplManager.RenderHTML("RejectUsageRequest", u.CallbackQuery.Message)
	if err != nil {
		return a.EditMessageTextConfig{}, tg.NewBotError(err.Error(), nil)
	}

	msg := a.NewEditMessageText(u.FromChat().ID, u.CallbackQuery.Message.MessageID, text)

	msg.ReplyMarkup = nil

	msg.ParseMode = "HTML"

	return msg, nil
}

func BlockUsageRequest(ctx tg.Context, u a.Update) (a.EditMessageTextConfig, *tg.BotError) {
	text, err := tmplManager.RenderHTML("BlockUsageRequest", u.CallbackQuery.Message)
	if err != nil {
		return a.EditMessageTextConfig{}, tg.NewBotError(err.Error(), nil)
	}

	msg := a.NewEditMessageText(u.FromChat().ID, u.CallbackQuery.Message.MessageID, text)

	msg.ReplyMarkup = nil

	msg.ParseMode = "HTML"

	return msg, nil
}

func NotifyUsageRequestConfirmed(ctx tg.Context, tgID int64) (a.MessageConfig, *tg.BotError) {
	text, err := tmplManager.RenderHTML("NotifyUsageRequestConfirmed", nil)
	if err != nil {
		return a.MessageConfig{}, tg.NewBotError(err.Error(), nil)
	}

	msg := a.NewMessage(tgID, text)

	kb := MainMenuReplyMarkup()
	msg.ReplyMarkup = kb

	msg.ParseMode = "HTML"

	return msg, nil
}

func NotifyUsageRequestRejected(ctx tg.Context, tgID int64) (a.MessageConfig, *tg.BotError) {
	text, err := tmplManager.RenderHTML("NotifyUsageRequestRejected", nil)
	if err != nil {
		return a.MessageConfig{}, tg.NewBotError(err.Error(), nil)
	}

	msg := a.NewMessage(tgID, text)

	msg.ParseMode = "HTML"

	return msg, nil
}

func NotifyUsageRequestBlocked(ctx tg.Context, tgID int64) (a.MessageConfig, *tg.BotError) {
	text, err := tmplManager.RenderHTML("NotifyUsageRequestBlocked", nil)
	if err != nil {
		return a.MessageConfig{}, tg.NewBotError(err.Error(), nil)
	}

	msg := a.NewMessage(tgID, text)

	msg.ParseMode = "HTML"

	return msg, nil
}
