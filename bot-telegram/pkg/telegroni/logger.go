package telegroni

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

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

const (
	UpdateTypeTextMessage        = " TEXT"
	UpdateTypeCommand            = " COMMAND"
	UpdateTypeCallbackQuery      = " CALLBACKQ"
	UpdateTypeEditedMessage      = " EDIT"
	UpdateTypeChannelPost        = " CHANNELPOST"
	UpdateTypeEditedChannelPost  = " EDITCHANPOST"
	UpdateTypeInlineQuery        = " INLINEQ"
	UpdateTypeChosenInlineResult = " CHOSENINLINERES"
	UpdateTypeShippingQuery      = " SHIPPINGQ"
	UpdateTypePreCheckoutQuery   = " PRECHECKOUTQ"
	UpdateTypePoll               = " POLL"
	UpdateTypePollAnswer         = " POLLANSWER"
	UpdateTypeUnknown            = " UNKNOWN"
)

const (
	StatusError = "ERR"
	StatusOK    = "OK"
	StatusWarn  = "WARN"
)

const (
	UPDATE_TYPE_WIDTH = 17
	USERNAME_WIDTH    = 15
)

func isTerminal() bool {
	stat, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
}

type Logger struct {
	out     io.Writer
	mu      sync.Mutex
	Colored bool
}

func NewLogger() *Logger {
	return &Logger{
		out:     os.Stdout,
		Colored: isTerminal(),
	}
}

var DefaultLogMiddleware = NewMiddleware(mwLogger, "Default logger")
var defaultLogger = NewLogger()

func LogGlobalError(ctx context.Context, update tgbotapi.Update) {
	startTimestamp := ctx.Value("timestamp_recieved").(time.Time)

	statusS := formatStatus(StatusError)
	end := time.Now()
	duration := end.Sub(startTimestamp)
	durationS := formatDuration(duration, defaultLogger.Colored)
	timestamp := startTimestamp.Format("2006/01/02 - 15:04:05")
	dateS, updateS, usernameS, detailsS := defaultLogger.getUpdateData(update, end)
	err := NewBotError("error: unhandled update", nil)

	line := fmt.Sprintf("[BOT] %s |%s|%s|%s|%s|%s %s %s\n",
		timestamp, dateS, statusS, durationS, usernameS, updateS, detailsS, err.Message)

	defaultLogger.mu.Lock()
	defaultLogger.out.Write([]byte(line))
	defaultLogger.mu.Unlock()
}

func mwLogger(ctx context.Context, update tgbotapi.Update, next HandlerFunc) (status string, err *BotError) {
	startTimestamp := ctx.Value("timestamp_recieved").(time.Time)
	timestamp := startTimestamp.Format("2006/01/02 - 15:04:05")

	status, err = next(ctx, update)
	if err == nil {
		err = NewBotError("", nil)
	} else {
		err = NewBotError(fmt.Sprintf("error: \"%v\"", err), nil)
	}

	end := time.Now()
	duration := end.Sub(startTimestamp)

	durationS := formatDuration(duration, defaultLogger.Colored)

	statusS := formatStatus(status)
	dateS,
		updateS,
		usernameS,
		detailsS := defaultLogger.getUpdateData(update, end)

	line := fmt.Sprintf("[BOT] %s |%s|%s|%s|%s|%s %s %s\n",
		timestamp, dateS, statusS, durationS, usernameS, updateS, detailsS, err.Message)

	defaultLogger.mu.Lock()
	defaultLogger.out.Write([]byte(line))
	defaultLogger.mu.Unlock()

	return status, err
}

func (l *Logger) getUpdateData(u tgbotapi.Update, end time.Time) (date string, ut string, username string, details string) {
	switch {
	// ===== POLL ANSWER (голос) =====
	case u.PollAnswer != nil:
		return formatDate(time.Unix(0, 0), end),
			l.updateTypeFormated(UpdateTypePollAnswer, bgBlue, italic),
			formatUsername(u.PollAnswer.User.UserName),
			fmt.Sprintf("poll_id=%s, options=%v", u.PollAnswer.PollID, u.PollAnswer.OptionIDs)

	// ===== MESSAGE =====
	case u.Message != nil && u.Message.Poll != nil:
		return formatDate(u.Message.Time(), end),
			l.updateTypeFormated(UpdateTypePoll, bgBlue, bold),
			formatUsername(u.Message.From.UserName),
			fmt.Sprintf("poll_id=%q, q=%q, total=%d, options=%v", u.Message.Poll.ID, u.Message.Poll.Question, len(u.Message.Poll.Options), u.Message.Poll.Options)

	case u.Message != nil && u.Message.IsCommand():
		return formatDate(u.Message.Time(), end),
			l.updateTypeFormated(UpdateTypeCommand, bgGreen),
			formatUsername(u.Message.From.UserName),
			fmt.Sprintf("command:%q", u.Message.Text)

	case u.Message != nil && u.Message.Text != "":
		return formatDate(u.Message.Time(), end),
			l.updateTypeFormated(UpdateTypeTextMessage, bgCyan),
			formatUsername(u.Message.From.UserName),
			fmt.Sprintf("text:%q", u.Message.Text)

	// ===== CALLBACK =====
	case u.CallbackQuery != nil:
		return formatDate(time.Unix(0, 0), end),
			l.updateTypeFormated(UpdateTypeCallbackQuery, bgGreen),
			formatUsername(u.CallbackQuery.From.UserName),
			fmt.Sprintf("%q", u.CallbackQuery.Data)

	// ===== EDITED =====
	case u.EditedMessage != nil:
		return formatDate(u.EditedMessage.Time(), end),
			l.updateTypeFormated(UpdateTypeEditedMessage, bgCyan, italic),
			formatUsername(u.EditedMessage.From.UserName),
			fmt.Sprintf("%q", u.EditedMessage.Text)

	// ===== CHANNEL =====
	case u.ChannelPost != nil:
		return formatDate(u.ChannelPost.Time(), end),
			l.updateTypeFormated(UpdateTypeChannelPost, bgMagenta),
			formatUsername(u.ChannelPost.From.UserName),
			fmt.Sprintf("%q", u.ChannelPost.Text)

	case u.EditedChannelPost != nil:
		return formatDate(u.EditedChannelPost.Time(), end),
			l.updateTypeFormated(UpdateTypeEditedChannelPost, bgMagenta, italic),
			formatUsername(u.EditedChannelPost.From.UserName),
			fmt.Sprintf("%q", u.EditedChannelPost.Text)

	// ===== INLINE =====
	case u.InlineQuery != nil:
		return formatDate(time.Unix(0, 0), end),
			l.updateTypeFormated(UpdateTypeInlineQuery, bgGreen),
			formatUsername(u.InlineQuery.From.UserName),
			fmt.Sprintf("query=%q, offset=%q", u.InlineQuery.Query, u.InlineQuery.Offset)

	case u.ChosenInlineResult != nil:
		return formatDate(time.Unix(0, 0), end),
			l.updateTypeFormated(UpdateTypeChosenInlineResult, bgGreen),
			formatUsername(u.ChosenInlineResult.From.UserName),
			fmt.Sprintf("result_id=%s, query=%q", u.ChosenInlineResult.ResultID, u.ChosenInlineResult.Query)

	// ===== PAYMENTS =====
	case u.ShippingQuery != nil:
		return formatDate(time.Unix(0, 0), end),
			l.updateTypeFormated(UpdateTypeShippingQuery, bgRed, bold),
			formatUsername(u.ShippingQuery.From.UserName),
			fmt.Sprintf("shipping_id=%s", u.ShippingQuery.ID)

	case u.PreCheckoutQuery != nil:
		return formatDate(time.Unix(0, 0), end),
			l.updateTypeFormated(UpdateTypePreCheckoutQuery, bgRed, bold),
			formatUsername(u.PreCheckoutQuery.From.UserName),
			fmt.Sprintf("invoice_payload=%s", u.PreCheckoutQuery.InvoicePayload)

	// ===== UNKNOWN =====
	default:
		return formatDate(time.Unix(0, 0), end),
			l.updateTypeFormated(UpdateTypeUnknown, bgBlack, red, bold),
			formatUsername(""),
			fmt.Sprintf("UpdateID=%v", u.UpdateID)
	}
}

func formatDate(t time.Time, end time.Time) string {

	if t.Equal(time.Unix(0, 0)) {
		return fmt.Sprintf("%s%s%s%s", bgBlack, red, padLeft("---", 11), reset)
	}
	t = t.UTC()
	end = end.UTC()

	t = t.Truncate(time.Second)
	dur := end.Sub(t)
	if dur < 0 {
		dur = 0
	}
	dateS := ""
	switch {
	case dur < time.Millisecond*500:
		dateS = fmt.Sprintf("%3.3fms ", float64(dur.Microseconds())/1000)
	default:
		dateS = fmt.Sprintf("%3.3fs ", float64((dur-dur.Truncate(time.Hour)).Milliseconds())/1000)
	}

	dateS = padLeft(dateS, 11)

	if defaultLogger.Colored {
		switch {
		case dur < time.Second*2:
			dateS = fmt.Sprintf("%s%s%s", bgGreen, dateS, reset)
		case dur < time.Second*5:
			dateS = fmt.Sprintf("%s%s%s", bgYellow, dateS, reset)
		default:
			dateS = fmt.Sprintf("%s%s%s", bgRed, dateS, reset)
		}
	}

	return dateS
}

func (l *Logger) updateTypeFormated(ut string, colors ...string) string {
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

func padLeft(s string, length int) string {
	slen := utf8.RuneCountInString(s)
	if slen > length {
		length = slen
	}
	pad := strings.Repeat(" ", length-slen)
	return fmt.Sprintf("%s%s", pad, s)
}
func padRight(s string, length int) string {
	slen := utf8.RuneCountInString(s)
	if slen > length {
		length = slen
	}
	pad := strings.Repeat(" ", length-slen)
	return fmt.Sprintf("%s%s", s, pad)
}

func formatUsername(username string) string {
	us := fmt.Sprintf(" @%s", padRight(username, USERNAME_WIDTH))

	return us
}

func formatStatus(status string) string {
	statusS := padRight(status, 4)

	if defaultLogger.Colored {
		switch status {
		case StatusOK:
			statusS = fmt.Sprintf("%s %s %s", bgGreen, statusS, reset)
		case StatusError:
			statusS = fmt.Sprintf("%s %s %s", bgRed, statusS, reset)
		case StatusWarn:
			statusS = fmt.Sprintf("%s %s %s", bgYellow, statusS, reset)
		default:
			statusS = fmt.Sprintf("%s%s%s %s %s", bgBlack, red, bold, statusS, reset)
		}
	}
	return statusS
}

func formatDuration(d time.Duration, colored bool) string {
	var durationS string
	switch {
	case d < time.Microsecond*500:
		durationS = fmt.Sprintf("%3.3fμs ", float64(d.Nanoseconds())/1000)
	case d < time.Millisecond*500:
		durationS = fmt.Sprintf("%3.3fms ", float64(d.Microseconds())/1000)
	default:
		durationS = fmt.Sprintf("%3.3fs ", float64(d.Milliseconds())/1000)
	}

	durationS = padLeft(durationS, 11)

	if colored {
		switch {
		case d < time.Millisecond*100:
			durationS = fmt.Sprintf("%s%s%s", bgWhite, durationS, reset)
		case d < time.Millisecond*200:
			durationS = fmt.Sprintf("%s%s%s", bgGreen, durationS, reset)
		case d < time.Second:
			durationS = fmt.Sprintf("%s%s%s", bgYellow, durationS, reset)
		default:
			durationS = fmt.Sprintf("%s%s%s", bgRed, durationS, reset)
		}
	}

	return durationS
}
