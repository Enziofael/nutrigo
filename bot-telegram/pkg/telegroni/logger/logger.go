package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Enziofael/nutrigo/bot-telegram/pkg/telegroni"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	reset     = "\033[0m"
	bold      = "\033[1m"
	dim       = "\033[2m"
	italic    = "\033[3m"
	underline = "\033[4m"

	black   = "\033[30m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	white   = "\033[37m"

	bgBlack   = "\033[40m"
	bgRed     = "\033[41m"
	bgGreen   = "\033[42m"
	bgYellow  = "\033[43m"
	bgBlue    = "\033[44m"
	bgMagenta = "\033[45m"
	bgCyan    = "\033[46m"
	bgWhite   = "\033[47m"
)

func isTerminal() bool {
	stat, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
}

type Logger struct {
	out      io.Writer
	mu       sync.Mutex
	Colored  bool
	statusCh chan string
}

func New() *Logger {
	return &Logger{
		out:      os.Stdout,
		Colored:  isTerminal(),
		statusCh: make(chan string, 100),
	}
}

func DefaultLogger() telegroni.Middleware {
	return telegroni.NewMiddleware(mwLogger, "Default logger")
}

func (l *Logger) StatusCh() chan<- string {
	return l.statusCh
}

var defaultLogger = New()

func LogGlobalError(update tgbotapi.Update) {
	startTimestamp := time.Unix(int64(update.Message.Date), 0)
	statusS := formatStatus(StatusError)
	duration := time.Since(startTimestamp)
	durationS := formatDuration(duration, defaultLogger.Colored)
	timestamp := startTimestamp.Format("2006/01/02 15:04:05")
	updateS, usernameS, detailsS := defaultLogger.getUpdateData(update)
	line := fmt.Sprintf("[BOT] %s |%s|%s|%s|%s\"%s\"\n",
		timestamp, statusS, durationS, usernameS, updateS, detailsS)

	defaultLogger.mu.Lock()
	defaultLogger.out.Write([]byte(line))
	defaultLogger.mu.Unlock()
}

func mwLogger(ctx context.Context, update tgbotapi.Update, next telegroni.HandlerFunc) {
	startTimestamp := time.Now()
	timestamp := startTimestamp.Format("2006/01/02 15:04:05")

	next(ctx, update)

	duration := time.Since(startTimestamp)
	durationS := formatDuration(duration, defaultLogger.Colored)

	var status string
	select {
	case status = <-defaultLogger.statusCh:
	default:
		status = "OK"
	}

	statusS := status
	if defaultLogger.Colored {
		switch status {
		case StatusOK:
			statusS = bgGreen + formatStatus(status) + reset
		case StatusError:
			statusS = bgRed + formatStatus(status) + reset
		case StatusWarn:
			statusS = bgYellow + formatStatus(status) + reset
		default:
			statusS = bgBlack + red + bold + formatStatus(status) + reset
		}
	}

	updateS, usernameS, detailsS := defaultLogger.getUpdateData(update)
	line := fmt.Sprintf("[BOT] %s |%s|%s|%s|%s\"%s\"\n",
		timestamp, statusS, durationS, usernameS, updateS, detailsS)

	defaultLogger.mu.Lock()
	defaultLogger.out.Write([]byte(line))
	defaultLogger.mu.Unlock()
}

type UpdateType string

const (
	TextMessage        UpdateType = " TEXT"
	Command                       = " COMMAND"
	CallbackQuery                 = " CALLBACKQ"
	EditedMessage                 = " EDIT"
	ChannelPost                   = " CHANNELPOST"
	EditedChannelPost             = " EDITCHANPOST"
	InlineQuery                   = " INLINEQ"
	ChosenInlineResult            = " CHOSENINLINERES"
	ShippingQuery                 = " SHIPPINGQ"
	PreCheckoutQuery              = " PRECHECKOUTQ"
	Poll                          = " POLL"
	PollAnswer                    = " POLLANSWER"
	Unknown                       = " UNKNOWN"
)

const (
	StatusError = "ERR"
	StatusOK    = "OK"
	StatusWarn  = "WARN"
)

func (l *Logger) getUpdateData(u tgbotapi.Update) (ut string, username string, details string) {
	switch {
	// ===== POLL ANSWER (голос) =====
	case u.PollAnswer != nil:
		return l.updateTypeFormated(PollAnswer, bgBlue, italic),
			formatUsername(u.PollAnswer.User.UserName),
			fmt.Sprintf("poll_id=%s, options=%v", u.PollAnswer.PollID, u.PollAnswer.OptionIDs)

	// ===== MESSAGE =====
	case u.Message != nil && u.Message.Poll != nil:
		return l.updateTypeFormated(Poll, bgBlue, bold),
			formatUsername(u.Message.From.UserName),
			fmt.Sprintf("poll_id=%q, q=%q, total=%d, options=%v", u.Message.Poll.ID, u.Message.Poll.Question, len(u.Message.Poll.Options))

	case u.Message != nil && u.Message.IsCommand():
		return l.updateTypeFormated(Command, bgGreen),
			formatUsername(u.Message.From.UserName),
			fmt.Sprintf("%q", u.Message.Text)

	case u.Message != nil:
		return l.updateTypeFormated(TextMessage, bgCyan),
			formatUsername(u.Message.From.UserName),
			fmt.Sprintf("%q", u.Message.Text)

	// ===== CALLBACK =====
	case u.CallbackQuery != nil:
		return l.updateTypeFormated(CallbackQuery, bgGreen),
			formatUsername(u.CallbackQuery.From.UserName),
			fmt.Sprintf("%q", u.CallbackQuery.Data)

	// ===== EDITED =====
	case u.EditedMessage != nil:
		return l.updateTypeFormated(EditedMessage, bgCyan, italic),
			formatUsername(u.EditedMessage.From.UserName),
			fmt.Sprintf("%q", u.EditedMessage.Text)

	// ===== CHANNEL =====
	case u.ChannelPost != nil:
		return l.updateTypeFormated(ChannelPost, bgMagenta),
			formatUsername(u.ChannelPost.From.UserName),
			fmt.Sprintf("%q", u.ChannelPost.Text)

	case u.EditedChannelPost != nil:
		return l.updateTypeFormated(EditedChannelPost, bgMagenta, italic),
			formatUsername(u.EditedChannelPost.From.UserName),
			fmt.Sprintf("%q", u.EditedChannelPost.Text)

	// ===== INLINE =====
	case u.InlineQuery != nil:
		return l.updateTypeFormated(InlineQuery, bgGreen),
			formatUsername(u.InlineQuery.From.UserName),
			fmt.Sprintf("query=%q, offset=%q", u.InlineQuery.Query, u.InlineQuery.Offset)

	case u.ChosenInlineResult != nil:
		return l.updateTypeFormated(ChosenInlineResult, bgGreen),
			formatUsername(u.ChosenInlineResult.From.UserName),
			fmt.Sprintf("result_id=%s, query=%q", u.ChosenInlineResult.ResultID, u.ChosenInlineResult.Query)

	// ===== PAYMENTS =====
	case u.ShippingQuery != nil:
		return l.updateTypeFormated(ShippingQuery, bgRed, bold),
			formatUsername(u.ShippingQuery.From.UserName),
			fmt.Sprintf("shipping_id=%s", u.ShippingQuery.ID)

	case u.PreCheckoutQuery != nil:
		return l.updateTypeFormated(PreCheckoutQuery, bgRed, bold),
			formatUsername(u.PreCheckoutQuery.From.UserName),
			fmt.Sprintf("invoice_payload=%s", u.PreCheckoutQuery.InvoicePayload)

	// ===== UNKNOWN =====
	default:
		return l.updateTypeFormated(Unknown, bgBlack, red, bold),
			"",
			fmt.Sprintf("UpdateID=%d", u.UpdateID)
	}
}

const (
	UPDATE_TYPE_WIDTH = 17
	USERNAME_WIDTH    = 20
)

func (l *Logger) updateTypeFormated(ut UpdateType, colors ...string) string {
	prefix := ""
	for _, c := range colors {
		prefix += c
	}

	text := string(ut)
	if len(text) < UPDATE_TYPE_WIDTH {
		text = text + strings.Repeat(" ", UPDATE_TYPE_WIDTH-len(text))
	}
	if l.Colored {
		text = fmt.Sprintf("%s%s%s", prefix, text, reset)
	}
	return text
}

func pad(leftPad int, s string, rightPad int) string {
	left := strings.Repeat(" ", leftPad-len(s))
	right := strings.Repeat(" ", rightPad-len(s))
	return fmt.Sprintf("%s%s%s ", left, s, right)
}

func formatUsername(username string) string {
	us := "@" + username + " "
	left := USERNAME_WIDTH
	if left < len(us) {
		left = len(us) + 1
	}
	return pad(left, username, len(us))
}

func formatStatus(status string) string {
	s := " " + status + " "
	return pad(len(s), s, 6)
}

func formatDuration(d time.Duration, colored bool) string {
	var durationS string
	switch {
	case d < time.Microsecond*500:
		durationS = fmt.Sprintf("\t%.3gμs ", float64(d.Nanoseconds())/1000)
	case d < time.Millisecond*500:
		durationS = fmt.Sprintf("\t%.3gms ", float64(d.Microseconds())/1000)
	default:
		durationS = fmt.Sprintf("\t%.3gs ", float64(d.Milliseconds())/1000)
	}

	if colored {
		switch {
		case d < time.Millisecond*10:
			durationS = fmt.Sprintf("%s%s%s", bgWhite, durationS, reset)
		case d < time.Millisecond*50:
			durationS = fmt.Sprintf("%s%s%s", bgGreen, durationS, reset)
		case d < time.Second:
			durationS = fmt.Sprintf("%s%s%s", bgYellow, durationS, reset)
		default:
			durationS = fmt.Sprintf("%s%s%s", bgRed, durationS, reset)
		}
	}
	return durationS
}
