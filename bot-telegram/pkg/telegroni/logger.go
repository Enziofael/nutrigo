package types

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/Enziofael/nutrigo/bot-telegram/pkg/telegroni/consts"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// - Logger is a struct for containing output and error streams to log,
//
// - out is a io.Writer for default logs
//
// - err is a io.Writer for error logs
//
// - mu is a Mutex to prevent data race while writing
//
// - Colored defines if Logger can use ANSI colors for logging
//
// Use NewLogger() to create instanses
type Logger struct {
	out     io.Writer
	err     io.Writer
	mu      sync.Mutex
	Colored bool
	Config  LogConfig
}

// - NewLogger() creates a new Logger and returns a pointer to it
//
// - writers redefine output and error streams
// instead of os.Stdout and os.Stderr, respectively
//
// # If no writer is provided:
//
// out and err will remain standard streams
//
// # If exactly 1 writer is provided:
//
// it will be applied to both output and error streams it's not nil
//
// # If at least 2 writers are provided:
//
// the first one will be applied to output stream
// if it's not nil,
//
// the second one will be applied to error stream
// if it's not nil
//
// subsequent writers after the second writer will be ignored.
func NewLogger(writers ...*os.File) *Logger {

	out, err := os.Stdout, os.Stderr

	if len(writers) == 1 && writers[0] != nil {
		out, err = writers[0], writers[0]
	}
	if len(writers) >= 2 && writers[0] != nil {
		out = writers[0]
	}
	if len(writers) >= 2 && writers[1] != nil {
		err = writers[1]
	}

	return &Logger{
		out:     out,
		err:     err,
		Colored: isTerminal(),
	}
}

func isTerminal() bool {
	stat, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
}

func DefaultLogging(l *Logger) MiddlewareFunc {
	return l.logMiddleware
}

func (l *Logger) logMiddleware(ctx context.Context, update tgbotapi.Update, next HandlerFunc) (status string, err *BotError) {

	status, err = next(ctx, update)

	timestampHandled := time.Now()

	go func() {
		timestampRouted := ctx.Value(ContextKey_TimestampRouted).(time.Time)
		l.Log(NewLog(ctx, update, status, err, timestampRouted, timestampHandled))
	}()
	return status, err
}

func (l *Logger) Log(log Log) {
	s := log.String(l.Config, l.Colored)
	l.mu.Lock()
	l.out.Write([]byte(s))
	l.mu.Unlock()
}

func (l *Logger) Err(log Log) {
	s := log.String(l.Config, l.Colored)
	l.mu.Lock()
	l.err.Write([]byte(s))
	l.mu.Unlock()
}

func (l *Logger) Write(s string) {
	s = "[Bot]" + s
	l.mu.Lock()
	l.out.Write([]byte(s))
	l.mu.Unlock()
}

func (l *Logger) WriteErr(s string) {
	s = "[Bot] Error: " + s
	l.mu.Lock()
	l.err.Write([]byte(s))
	l.mu.Unlock()
}

type Log struct {
	Timestamp      time.Time
	Delay          *time.Duration
	Status         string
	ProcessingTime time.Duration
	Username       string
	RoutingPath    string
	Update         tgbotapi.Update
	Details        string
}

func NewLog(ctx context.Context, u tgbotapi.Update, status string, err *BotError, routed time.Time, handled time.Time) Log {
	return Log{
		Timestamp:      handled,
		Delay:          getDelay(u, handled),
		Status:         status,
		ProcessingTime: getProcessingTime(routed, handled),
		Username:       getUsername(u),
		RoutingPath:    getRoutingPath(ctx),
		Update:         u,
		Details:        getDetails(),
	}
}

func getDelay(u tgbotapi.Update, handled time.Time) *time.Duration {
	var sendAt time.Time
	switch {
	case u.Message != nil:
		sendAt = u.Message.Time().UTC()
	case u.EditedMessage != nil:
		sendAt = u.EditedMessage.Time().UTC()
	case u.ChannelPost != nil:
		sendAt = u.ChannelPost.Time().UTC()
	case u.EditedChannelPost != nil:
		sendAt = u.EditedChannelPost.Time().UTC()
	case u.CallbackQuery != nil && u.CallbackQuery.Message != nil:
		sendAt = u.CallbackQuery.Message.Time().UTC()
	}
	if sendAt.IsZero() {
		return nil
	}
	res := handled.UTC().Sub(sendAt)
	if res < 0 {
		res = 0
	}
	return &res
}

func getProcessingTime(routed time.Time, handled time.Time) time.Duration {
	res := handled.UTC().Sub(routed.UTC())
	if res < 0 {
		res = 0
	}
	return res
}

func getRoutingPath(ctx context.Context) string {
	return ctx.Value(ContextKey_Path).(string)
}

func getDetails() string {
	return "details stub"
}

func getUsername(u tgbotapi.Update) string {
	res := ""
	switch {
	case u.Message != nil:
		res = u.Message.From.UserName
	case u.EditedMessage != nil:
		res = u.EditedMessage.From.UserName
	case u.ChannelPost != nil:
		res = u.ChannelPost.Chat.UserName
	case u.EditedChannelPost != nil:
		res = u.EditedChannelPost.From.UserName
	case u.CallbackQuery != nil && u.CallbackQuery.Message != nil:
		res = u.CallbackQuery.Message.From.UserName
	default:
		res = fmt.Sprintf("getUsernameError %#v", u)
	}
	return res
}

type LogConfig struct {
	DelayWidth          int
	StatusWidth         int
	ProcessingTimeWidth int
	Usernamewidth       int
	UpdateWidth         int
	RoutingPathWidth    int
}

func (l Log) String(cfg LogConfig, colored bool) string {
	return "[BOT]" +
		formatTimestamp(l.Timestamp) + "|" +
		formatDelay(l.Delay, cfg.DelayWidth) + "|" +
		formatStatus(l.Status, cfg.StatusWidth) + "|" +
		formatProcessingTime(l.ProcessingTime, cfg.ProcessingTimeWidth) + "|" +
		formatUsername(l.Username, cfg.Usernamewidth) + "|" +
		formatUpdate(l.Update, colored, cfg.UpdateWidth) + "|" +
		formatRoutingPath(l.RoutingPath, cfg.RoutingPathWidth) +
		formatDetails(l.Details)
}

func formatTimestamp(timestamp time.Time) string {
	s := timestamp.Format(" 2006/01/02 - 15:04:05 ")

	return s
}

func formatDelay(delay *time.Duration, width int) string {
	var s string
	if delay == nil {
		s = padLeft(" --- ", width)
	} else if *delay < 1000*time.Millisecond {
		s = fmt.Sprintf(" %3.3fms ", float64(delay.Nanoseconds())/1000)
	} else if *delay < 10*time.Second {
		s = fmt.Sprintf(" %4.2fms ", float64(delay.Nanoseconds())/100)
	} else if *delay < 1000*time.Second {
		s = fmt.Sprintf(" %3.3fs ", float64(delay.Milliseconds())/1000)
	} else if *delay < 200*time.Minute {
		s = fmt.Sprintf(" %dm%ds ", int(delay.Minutes()), int(delay.Seconds())%60)
	} else if *delay < 23*time.Hour {
		s = fmt.Sprintf(" %dh%dm ", int(delay.Hours()), int(delay.Minutes())%60)
	} else {
		s = fmt.Sprintf(" %dh ", int(delay.Hours()))
	}
	return padLeft(s, width)
}

func formatStatus(status string, width int) string {
	return padLeft(status, width)
}

func formatProcessingTime(processingTime time.Duration, width int) string {
	var s string
	if processingTime < 1000*time.Millisecond {
		s = fmt.Sprintf(" %3.3fms ", float64(processingTime.Nanoseconds())/1000)
	} else if processingTime < 10*time.Second {
		s = fmt.Sprintf(" %4.2fms ", float64(processingTime.Nanoseconds())/100)
	} else if processingTime < 1000*time.Second {
		s = fmt.Sprintf(" %3.3fs ", float64(processingTime.Milliseconds())/1000)
	} else if processingTime < 200*time.Minute {
		s = fmt.Sprintf(" %dm%ds ", int(processingTime.Minutes()), int(processingTime.Seconds())%60)
	} else if processingTime < 23*time.Hour {
		s = fmt.Sprintf(" %dh%dm ", int(processingTime.Hours()), int(processingTime.Minutes())%60)
	} else {
		s = fmt.Sprintf(" %dh ", int(processingTime.Hours()))
	}
	return padLeft(s, width)
}

func formatUsername(username string, width int) string {
	if username != "" {
		username = "@" + username
	} else {
		username = "---"
	}
	return padLeft(username, width)
}

func formatRoutingPath(routingPath string, width int) string {
	routingPath += " " + routingPath
	return padRight(routingPath, width)
}

func formatUpdate(u tgbotapi.Update, colored bool, width int) string {
	switch {
	// ===== POLL ANSWER (голос) =====
	case u.PollAnswer != nil:
		return formatUpdateType(colored, width, consts.UT_PollAnswer, consts.ANSI_BG_BLUE, consts.ANSI_ITALIC)

	// ===== MESSAGE =====
	case u.Message != nil && u.Message.Poll != nil:
		return formatUpdateType(colored, width, consts.UT_Poll, consts.ANSI_BG_BLUE, consts.ANSI_BOLD)

	case u.Message != nil && u.Message.IsCommand():
		return formatUpdateType(colored, width, consts.UT_Command, consts.ANSI_BG_GREEN)

	case u.Message != nil && u.Message.Text != "":
		return formatUpdateType(colored, width, consts.UT_TextMessage, consts.ANSI_BG_CYAN)

	// ===== CALLBACK =====
	case u.CallbackQuery != nil:
		return formatUpdateType(colored, width, consts.UT_CallbackQuery, consts.ANSI_BG_BRIGHT_GREEN)

	// ===== EDITED =====
	case u.EditedMessage != nil:
		return formatUpdateType(colored, width, consts.UT_EditedMessage, consts.ANSI_BG_CYAN, consts.ANSI_ITALIC)

	// ===== CHANNEL =====
	case u.ChannelPost != nil:
		return formatUpdateType(colored, width, consts.UT_ChannelPost, consts.ANSI_BG_MAGENTA)

	case u.EditedChannelPost != nil:
		return formatUpdateType(colored, width, consts.UT_EditedChannelPost, consts.ANSI_BG_MAGENTA, consts.ANSI_ITALIC)

	// ===== INLINE =====
	case u.InlineQuery != nil:
		return formatUpdateType(colored, width, consts.UT_InlineQuery, consts.ANSI_BG_GREEN)

	case u.ChosenInlineResult != nil:
		return formatUpdateType(colored, width, consts.UT_ChosenInlineResult, consts.ANSI_BG_GREEN)

	// ===== PAYMENTS =====
	case u.ShippingQuery != nil:
		return formatUpdateType(colored, width, consts.UT_ShippingQuery, consts.ANSI_BG_RED, consts.ANSI_BOLD)

	case u.PreCheckoutQuery != nil:
		return formatUpdateType(colored, width, consts.UT_PreCheckoutQuery, consts.ANSI_BG_RED, consts.ANSI_BOLD)

	// ===== UNKNOWN =====
	default:
		return formatUpdateType(colored, width, consts.UT_Unknown, consts.ANSI_BG_BLACK, consts.ANSI_BRIGHT_RED, consts.ANSI_BOLD)
	}
}

func formatUpdateType(colored bool, width int, ut string, colors ...string) string {
	text := string(ut)
	text = padRight(text, width)
	if !colored {
		return fmt.Sprintf(" %s ", text)
	}
	prefix := ""
	for _, c := range colors {
		prefix += c
	}

	return prefix + text + consts.ANSI_RESET
}

func formatDetails(Details string) string {
	return " " + Details
}
