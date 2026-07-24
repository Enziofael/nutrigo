package types

import (
	"context"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Enziofael/nutrigo/bot-telegram/pkg/telegroni/consts"
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
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
		Config:  defaultLogConfig,
	}
}

func isTerminal() bool {
	stat, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
}

func DefaultLogging(s *Server) MiddlewareFunc {
	return s.Logger.defaultLogMiddleware
}

var once sync.Once

func (l *Logger) defaultLogMiddleware(ctx context.Context, update tgbotapi.Update, next HandlerFunc) (status string, err *BotError) {

	once.Do(func() {
		l.Write("Default logging columns:\n      Recieved timestamp; responce delay; status; handling time; update type; username; routing path; details\n")
	})

	status, err = next(ctx, update)

	timestampHandled := time.Now()

	go func() {
		timestampRecieved := ctx.Value(ContextKey_TimestampRecieved).(time.Time)
		log := NewLog(ctx, update, status, err, timestampRecieved, timestampHandled)
		l.Log(log)
		if l.Config.LogUpdateDetails {
			l.Write(formatUpdateDetails(reflect.ValueOf(update), 0))
		} else if l.Config.LogUpdateDetailsOnWarn && status == StatusWarn {
			l.Write("Warn: " + err.Error() + "\n" + formatUpdateDetails(reflect.ValueOf(update), 0))
		} else if l.Config.LogUpdateDetailsOnError && status == StatusError {
			l.Write("Error: " + err.Error() + "\n" + formatUpdateDetails(reflect.ValueOf(update), 0))
		} else if status != StatusOK && status != StatusWarn && status != StatusError {
			l.Write("Error: " + err.Error() + "\n" + formatUpdateDetails(reflect.ValueOf(update), 0))
		}

		if err != nil || status != StatusOK {
			if status != StatusWarn || l.Config.ErrWarn {
				l.WriteErr(formatUpdateDetails(reflect.ValueOf(update), 0))
			}
		}
	}()
	return status, err
}

func (l *Logger) Log(log Log) {
	s := l.String(l.Config, log)
	l.mu.Lock()
	l.out.Write([]byte(s + "\n"))
	l.mu.Unlock()
}

func (l *Logger) Err(log Log) {
	s := l.String(l.Config, log)
	l.mu.Lock()
	l.err.Write([]byte(s + "\n"))
	l.mu.Unlock()
}

func (l *Logger) Write(s string) {
	s = "[BOT] " + s
	l.mu.Lock()
	l.out.Write([]byte(s + "\n"))
	l.mu.Unlock()
}

func (l *Logger) WriteErr(s string) {
	s = "[BOT] " + s
	l.mu.Lock()
	l.err.Write([]byte(s + "\n"))
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

func NewLog(ctx context.Context, u tgbotapi.Update, status string, err *BotError, recieved time.Time, handled time.Time) Log {
	return Log{
		Timestamp:      handled,
		Delay:          getDelay(u, handled),
		Status:         status,
		ProcessingTime: getProcessingTime(recieved, handled),
		Username:       getUsername(u),
		RoutingPath:    getRoutingPath(ctx),
		Update:         u,
		Details:        getDetails(u),
	}
}

func getDelay(u tgbotapi.Update, handled time.Time) *time.Duration {
	handled = handled.UTC()
	var sendAt time.Time

	switch {
	case u.Message != nil:
		sendAt = u.Message.Time()
	case u.EditedMessage != nil:
		sendAt = u.EditedMessage.Time()
	case u.ChannelPost != nil:
		sendAt = u.ChannelPost.Time()
	case u.EditedChannelPost != nil:
		sendAt = u.EditedChannelPost.Time()
	case u.CallbackQuery != nil && u.CallbackQuery.Message != nil:
		sendAt = u.CallbackQuery.Message.Time()
	}

	if sendAt.IsZero() {
		return nil
	}

	sendAt = sendAt.UTC().Truncate(time.Millisecond)
	res := handled.Sub(sendAt)
	if res <= 0 {
		res = 0
	}
	return &res
}

func getProcessingTime(recieved time.Time, handled time.Time) time.Duration {
	res := handled.UTC().Sub(recieved.UTC())
	if res < 0 {
		res = 0
	}
	return res
}

func getRoutingPath(ctx context.Context) string {
	return ctx.Value(ContextKey_Path).(string)
}

func getDetails(u tgbotapi.Update) string {
	s := "id:" + strconv.Itoa(u.UpdateID) + " "

	if u.Message != nil && u.Message.Text != "" {
		s = s + "text:\"" + u.Message.Text + "\" "
	}
	if u.Message != nil && len(u.Message.Entities) != 0 {
		for _, e := range u.Message.Entities {
			s = s + "{"
			s = s + "type:\"" + e.Type + "\" "
			s = s + "length:\"" + strconv.Itoa(e.Length) + "\" "
			if e.Offset != 0 {
				s = s + "offset:\"" + strconv.Itoa(e.Offset) + "\" "
			}
			if e.URL != "" {
				s = s + "textLink:\"" + e.URL + "\" "
			}
			if e.User != nil {
				s = s + "mention:\"@" + e.User.UserName + "\" "
			}
			if e.Language != "" {
				s = s + "codeLang:\"" + e.Language + "\" "
			}
			if e.CustomEmojiID != "" {
				s = s + "customEmoji:\"" + e.CustomEmojiID + "\" "
			}
			if e.UnixTime != 0 {
				s = s + "dateTime:\"" + strconv.Itoa(e.Length) + "\" "
			}
			// DateTimeFormat not implemented
			s = s + "} "
		}
	}
	if u.Message != nil && u.Message.Photo != nil && u.Message.Caption != "" {
		s = s + "caption:\"" + u.Message.Caption + "\" "
	}
	if u.Message != nil && u.Message.Photo != nil && u.Message.MediaGroupID != "" {
		s = s + "mediaGroup:\"" + u.Message.MediaGroupID + "\" "
	}
	if u.Message != nil && u.Message.Document != nil {
		s = s + "fileName:\"" + u.Message.Document.FileName + "\" mimeType:\"" + u.Message.Document.MimeType + "\" "
	}
	if u.Message != nil && u.Message.Sticker != nil {
		s = s + "setName:\"" + u.Message.Sticker.SetName + "\" emoji:\"" + u.Message.Sticker.Emoji + "\" "
	}
	s = strings.ReplaceAll(s, "\n", " ")
	return s
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
	TimestampWidth      int
	DelayWidth          int
	StatusWidth         int
	ProcessingTimeWidth int
	UpdateWidth         int
	Usernamewidth       int
	RoutingPathWidth    int

	LogUpdateDetails        bool
	LogUpdateDetailsOnWarn  bool
	LogUpdateDetailsOnError bool

	ErrWarn bool
}

var defaultLogConfig LogConfig = LogConfig{
	TimestampWidth:      0,
	DelayWidth:          9,  //9
	StatusWidth:         6,  //6
	ProcessingTimeWidth: 11, //11
	UpdateWidth:         17, //17
	Usernamewidth:       0,
	RoutingPathWidth:    0,

	LogUpdateDetails:        false,
	LogUpdateDetailsOnWarn:  false,
	LogUpdateDetailsOnError: true,

	ErrWarn: false,
}

func (l *Logger) String(cfg LogConfig, log Log) string {
	return "[BOT]" +
		l.formatTimestamp(log.Timestamp, cfg.TimestampWidth) + "|" +
		l.formatDelay(log.Delay, cfg.DelayWidth) + "|" +
		l.formatStatus(log.Status, cfg.StatusWidth) + "|" +
		l.formatProcessingTime(log.ProcessingTime, cfg.ProcessingTimeWidth) + "|" +
		l.formatUpdate(log.Update, cfg.UpdateWidth) + "|" +
		l.formatUsername(log.Username, cfg.Usernamewidth) +
		l.formatRoutingPath(log.RoutingPath, cfg.RoutingPathWidth) +
		l.formatDetails(log.Details)
}

func (l *Logger) formatTimestamp(timestamp time.Time, width int) string {
	s := timestamp.Format(" 2006/01/02 - 15:04:05 ")

	return padRight(s, width)
}

func (l *Logger) formatDelay(delay *time.Duration, width int) string {
	var s string
	if delay == nil {
		s = padLeft(" --- ", width)
	} else if *delay == 0 {
		s = " 0s "
	} else if *delay < 1000*time.Second {
		s = fmt.Sprintf(" %3.3fs ", float64(delay.Milliseconds())/1000)
	} else if *delay < 200*time.Minute {
		s = fmt.Sprintf(" %dm%ds ", int(delay.Minutes()), int(delay.Seconds())%60)
	} else if *delay < 23*time.Hour {
		s = fmt.Sprintf(" %dh%dm ", int(delay.Hours()), int(delay.Minutes())%60)
	} else {
		s = fmt.Sprintf(" %dh ", int(delay.Hours()))
	}

	s = padLeft(s, width)

	if !l.Colored {
		return s
	}

	if *delay < 2*time.Second {
		return consts.ANSI_BG_WHITE + consts.ANSI_BLACK + s + consts.ANSI_RESET
	}
	if *delay < 5*time.Second {
		return consts.ANSI_BG_YELLOW + consts.ANSI_BOLD + consts.ANSI_WHITE + s + consts.ANSI_RESET
	}
	if *delay < 10*time.Second {
		return consts.ANSI_BG_RED + consts.ANSI_BOLD + consts.ANSI_WHITE + s + consts.ANSI_RESET
	}
	return consts.ANSI_BG_BLACK + consts.ANSI_BOLD + consts.ANSI_RED + s + consts.ANSI_RESET
}

func (l *Logger) formatStatus(status string, width int) string {
	s := " " + status + " "
	s = padRight(s, width)
	if !l.Colored {
		return s
	}
	switch status {
	case StatusOK:
		return consts.ANSI_BG_GREEN + s + consts.ANSI_RESET
	case StatusWarn:
		return consts.ANSI_BG_YELLOW + s + consts.ANSI_RESET
	case StatusError:
		return consts.ANSI_BG_RED + consts.ANSI_BOLD + s + consts.ANSI_RESET
	default:
		return consts.ANSI_BG_BLACK + consts.ANSI_BOLD + consts.ANSI_RED + s + consts.ANSI_RESET
	}
}

func (l *Logger) formatProcessingTime(processingTime time.Duration, width int) string {
	var s string
	if processingTime < 1000*time.Millisecond {
		s = fmt.Sprintf(" %3.3fms ", float64(processingTime.Microseconds())/1000)
	} else if processingTime < 10*time.Second {
		s = fmt.Sprintf(" %4.2fms ", float64(processingTime.Microseconds())/100)
	} else if processingTime < 1000*time.Second {
		s = fmt.Sprintf(" %3.3fs ", float64(processingTime.Milliseconds())/1000)
	} else if processingTime < 200*time.Minute {
		s = fmt.Sprintf(" %dm%ds ", int(processingTime.Minutes()), int(processingTime.Seconds())%60)
	} else if processingTime < 23*time.Hour {
		s = fmt.Sprintf(" %dh%dm ", int(processingTime.Hours()), int(processingTime.Minutes())%60)
	} else {
		s = fmt.Sprintf(" %dh ", int(processingTime.Hours()))
	}
	s = padLeft(s, width)
	if !l.Colored {
		return s
	}
	if processingTime < 300*time.Millisecond {
		return consts.ANSI_BG_WHITE + consts.ANSI_BLACK + s + consts.ANSI_RESET
	}
	if processingTime < time.Second {
		return consts.ANSI_BG_GREEN + consts.ANSI_WHITE + s + consts.ANSI_RESET
	}
	if processingTime < 5*time.Second {
		return consts.ANSI_BG_YELLOW + consts.ANSI_BLACK + consts.ANSI_BOLD + s + consts.ANSI_RESET
	}
	if processingTime < 10*time.Second {
		return consts.ANSI_BG_RED + consts.ANSI_BLACK + consts.ANSI_BOLD + s + consts.ANSI_RESET
	}
	return consts.ANSI_BG_BLACK + consts.ANSI_RED + consts.ANSI_BOLD + s + consts.ANSI_RESET
}

func (l *Logger) formatUsername(username string, width int) string {
	if username != "" {
		username = " @" + username + ":"
	} else {
		username = " --- "
	}
	username = padRight(username, width)
	if !l.Colored {
		return username
	}
	return consts.ANSI_BLUE + username + consts.ANSI_RESET
}

func (l *Logger) formatRoutingPath(routingPath string, width int) string {
	routingPath = " " + routingPath + " "
	routingPath = padRight(routingPath, width)
	if !l.Colored {
		return routingPath
	}
	return consts.ANSI_BLUE + routingPath + consts.ANSI_RESET
}

func (l *Logger) formatUpdate(u tgbotapi.Update, width int) string {
	switch {
	// ===== POLL ANSWER (голос) =====
	case u.PollAnswer != nil:
		return l.formatUpdateType(width, consts.UT_PollAnswer, consts.ANSI_BG_BLUE, consts.ANSI_ITALIC)

	// ===== MESSAGE =====
	case u.Message != nil && u.Message.Poll != nil:
		return l.formatUpdateType(width, consts.UT_Poll, consts.ANSI_BG_BLUE, consts.ANSI_WHITE)

	case u.Message != nil && u.Message.IsCommand():
		return l.formatUpdateType(width, consts.UT_Command, consts.ANSI_BG_YELLOW, consts.ANSI_BLACK)

	case u.Message != nil && u.Message.Text != "":
		return l.formatUpdateType(width, consts.UT_TextMessage, consts.ANSI_BG_BLUE, consts.ANSI_WHITE)

	case u.Message != nil:
		return l.formatUpdateType(width, consts.UT_Message, consts.ANSI_BG_BLUE, consts.ANSI_WHITE)

	// ===== CALLBACK =====
	case u.CallbackQuery != nil:
		return l.formatUpdateType(width, consts.UT_CallbackQuery, consts.ANSI_BG_GREEN)

	// ===== EDITED =====
	case u.EditedMessage != nil:
		return l.formatUpdateType(width, consts.UT_EditedMessage, consts.ANSI_BG_CYAN, consts.ANSI_ITALIC)

	// ===== CHANNEL =====
	case u.ChannelPost != nil:
		return l.formatUpdateType(width, consts.UT_ChannelPost, consts.ANSI_BG_MAGENTA)

	case u.EditedChannelPost != nil:
		return l.formatUpdateType(width, consts.UT_EditedChannelPost, consts.ANSI_BG_MAGENTA, consts.ANSI_ITALIC)

	// ===== INLINE =====
	case u.InlineQuery != nil:
		return l.formatUpdateType(width, consts.UT_InlineQuery, consts.ANSI_BG_GREEN)

	case u.ChosenInlineResult != nil:
		return l.formatUpdateType(width, consts.UT_ChosenInlineResult, consts.ANSI_BG_GREEN)

	// ===== PAYMENTS =====
	case u.ShippingQuery != nil:
		return l.formatUpdateType(width, consts.UT_ShippingQuery, consts.ANSI_BG_RED, consts.ANSI_BOLD)

	case u.PreCheckoutQuery != nil:
		return l.formatUpdateType(width, consts.UT_PreCheckoutQuery, consts.ANSI_BG_RED, consts.ANSI_BOLD)

	// ===== UNKNOWN =====
	default:
		return l.formatUpdateType(width, consts.UT_Unknown, consts.ANSI_BG_BLACK, consts.ANSI_BRIGHT_RED, consts.ANSI_BOLD)
	}
}

func (l *Logger) formatUpdateType(width int, ut string, colors ...string) string {
	text := string(ut)
	text = " " + text + " "

	if !l.Colored {
		return padRight(text, width)
	}

	prefix := ""
	for _, c := range colors {
		prefix += c
	}

	text = padRight(text, width)
	text = prefix + text + consts.ANSI_RESET
	return text
}

func (l *Logger) formatDetails(Details string) string {
	return " " + Details
}

func formatUpdateDetails(v reflect.Value, level int) string {
	var sb strings.Builder
	formatUpdateDetailsToBuilder(&sb, v, level)
	return sb.String()
}

func formatUpdateDetailsToBuilder(sb *strings.Builder, v reflect.Value, level int) {
	if !v.IsValid() || v.IsZero() {
		return
	}

	tab := strings.Repeat(" ", level*4)

	if (v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface) && !v.IsNil() {
		v = v.Elem()
	}
	if !v.IsValid() || v.IsZero() {
		return
	}

	t := v.Type()

	sb.WriteString(tab)
	sb.WriteString(t.Name())
	sb.WriteString(" {\n")

	tabInner := tab + "    "

	for i := 0; i < v.NumField(); i++ {
		fv := v.Field(i)
		if fv.IsZero() {
			continue
		}
		ft := t.Field(i).Name

		if fv.IsValid() && fv.Kind() == reflect.Slice {
			sb.WriteString(tabInner)
			sb.WriteString(ft)
			sb.WriteString(" [\n")
			for j := 0; j < fv.Len(); j++ {
				formatUpdateDetailsToBuilder(sb, fv.Index(j), level+2)
			}
			sb.WriteString(tabInner)
			sb.WriteString("]\n")
			continue
		}

		if fv.Kind() == reflect.Struct || fv.Kind() == reflect.Ptr || fv.Kind() == reflect.Interface {

			formatUpdateDetailsToBuilder(sb, fv, level+1)
			continue
		}

		sb.WriteString(tabInner)
		sb.WriteString(ft)
		sb.WriteString(" : ")
		if fv.Kind() == reflect.String {
			sb.WriteString("\"")
			sb.WriteString(fv.String())
			sb.WriteString("\"\n")
		} else {
			fmt.Fprint(sb, fv.Interface())
			sb.WriteString("\n")
		}

		//formatUpdateDetailsToBuilder(sb, fv.Elem(), level+1)
	}

	sb.WriteString(tab)
	sb.WriteString("}\n")
}
