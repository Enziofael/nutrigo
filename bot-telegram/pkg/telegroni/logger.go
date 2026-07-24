package telegroni

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Enziofael/nutrigo/bot-telegram/pkg/telegroni/internal/consts"
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

// ======================= INDEX =========================
//          You can navigate by searching REGION
// =======================================================
//	Arguments are omitted for brevity
//
//	1. type Handler struct
//	- Exported
//		- func NewHandler(...) *Handler
//  - Unexported
//  	- func (h *Handler) handle(...) (bool, string, *BotError)
//
//	2. type HandlerGroup struct
//	- Exported
//		- func NewHandlerGroup(...) *HandlerGroup
//		- func (g *HandlerGroup) Apply(...)
//		- func (g *HandlerGroup) Handle(...)
//		- func (g *HandlerGroup) Group(...) *HandlerGroup
//	- Unexported
//		- func (g *HandlerGroup) handle(...) (bool, string, *BotError)
//
//	3. type Middleware struct
//	- Exported
//		- func NewMiddleware(...) Middleware
//	- Unexported
//		- func (m Middleware) apply(...) HandlerFunc
//
//	4. Func() types
//	- Exported
//		- type MatchFunc func(...) bool
//		- type HandlerFunc func(...) (string, *BotError)
//		- type MiddlewareFunc func(...) (string, *BotError)
//		- type RoutingFallbackFunc func(...)
//
// =======================================================

// ======================== TYPE =========================
// REGION                  Logger
// =======================================================

// - Logger is a struct for containing output and error streams to log,
//
// - out is a io.Writer for default logs
//
// - err is a io.Writer for error logs
//
// - file is a *os.File for parallel logging to a file
//
// - mu is a Mutex to prevent data race while writing
//
// - Colored defines if Logger can use ANSI colors for logging
//
// - Config allows to set up logger behaviour and logs appearence
//
// # You can define where and what to log by setting a proper Config
//
// Use NewLogger() to create instanses
type Logger struct {
	out    io.Writer
	err    io.Writer
	file   *os.File
	mu     sync.Mutex
	Config LogConfig
}

// - NewLogger() creates a new Logger and returns a pointer to it
//
// - writers redefine output, error and file streams
// instead of os.Stdout, os.Stderr, nil respectively
//
// # If no writer is provided:
//
// out and err will remain standard streams,
// file logging is disabled.
//
// # If exactly 1 writer is provided:
//
// it will be applied to both output and error streams if it's not nil
//
// # If at least 2 writers are provided:
//
// the first one will be applied to output stream
// if it's not nil,
//
// the second one will be applied to error stream
// if it's not nil
//
// # If at least 3 writers are provided:
//
// the first three writers will be applied to output,
// error and file streamps respectively
//
// subsequent writers after the second writer will be ignored.
func NewLogger(writers ...*os.File) *Logger {

	out, err := os.Stdout, os.Stderr
	var file *os.File

	if len(writers) == 1 && writers[0] != nil {
		out, err = writers[0], writers[0]
	}
	if len(writers) >= 2 && writers[0] != nil {
		out = writers[0]
	}
	if len(writers) >= 2 && writers[1] != nil {
		err = writers[1]
	}
	if len(writers) >= 3 && writers[2] != nil {
		file = writers[3]
	}

	l := &Logger{
		out:    out,
		err:    err,
		file:   file,
		Config: defaultLogConfig,
	}

	l.Config.Colored = isTerminal()

	return l
}

func (l *Logger) Log(log Log) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.checkAndWrite(l.out, log, l.Config.LogBehaviour)
}

func (l *Logger) Err(log Log) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.checkAndWrite(l.err, log, l.Config.ErrBehaviour)
}

func (l *Logger) FileWrite(log Log) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.checkAndWrite(l.file, log, l.Config.FileBehaviour)
}

func (l *Logger) checkAndWrite(str io.Writer, log Log, behaviour byte) {
	var s string
	var m byte

	if behaviour == 0 || log.Status.code == 0 {
		return
	}

	if log.Status.code > 0b1_0000 {
		log.Status.code = log.Status.code / 0b1_0000
	}

	m = behaviour >> 4
	for code := byte(0b1000); code >= 0b0001; code = code >> 1 {
		if (m&code == code) && (log.Status.code&code == code) {
			s = l.String(l.Config, log)
			str.Write([]byte(s))
			break
		}
	}

	m = behaviour & 0b0000_1111
	for code := byte(0b1000); code >= 0b0001; code = code >> 1 {
		if (m&code == code) && (log.Status.code&code == code) {
			s = l.formatUpdateDetails(reflect.ValueOf(log.Update), 0)
			str.Write([]byte(s))
			break
		}
	}
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
	Status         HandleStatus
	ProcessingTime time.Duration
	Username       string
	RoutingPath    string
	Update         tgbotapi.Update
	Details        string
}

func NewLog(ctx Context, u tgbotapi.Update, status HandleStatus, err *BotError, recieved time.Time, handled time.Time) Log {
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

func getRoutingPath(ctx Context) string {
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

func (l *Logger) String(cfg LogConfig, log Log) string {
	return "[BOT]" +
		l.formatTimestamp(log.Timestamp, cfg.TimestampWidth) + "|" +
		l.formatDelay(log.Delay, cfg.DelayWidth) + "|" +
		l.formatStatus(log.Status, cfg.StatusWidth) + "|" +
		l.formatProcessingTime(log.ProcessingTime, cfg.ProcessingTimeWidth) + "|" +
		l.formatUpdate(log.Update, cfg.UpdateWidth) + "|" +
		l.formatUsername(log.Username, cfg.Usernamewidth) +
		l.formatRoutingPath(log.RoutingPath, cfg.RoutingPathWidth) +
		l.formatDetails(log.Details) + "\n"
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

	if !l.Config.Colored {
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

func (l *Logger) formatStatus(status HandleStatus, width int) string {
	s := " " + status.Name + " "
	s = padRight(s, width)
	if !l.Config.Colored {
		return s
	}
	switch status.code {
	case StatusOK.code:
		return consts.ANSI_BG_GREEN + s + consts.ANSI_RESET
	case StatusWarn.code:
		return consts.ANSI_BG_YELLOW + s + consts.ANSI_RESET
	case StatusError.code:
		return consts.ANSI_BG_RED + consts.ANSI_BOLD + s + consts.ANSI_RESET
	case StatusFallback.code:
		return consts.ANSI_BG_BLACK + consts.ANSI_BOLD + consts.ANSI_RED + s + consts.ANSI_RESET
	default:
		return consts.ANSI_BG_BLUE + consts.ANSI_BOLD + consts.ANSI_WHITE + s + consts.ANSI_RESET
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
	if !l.Config.Colored {
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
	if !l.Config.Colored {
		return username
	}
	return consts.ANSI_BLUE + username + consts.ANSI_RESET
}

func (l *Logger) formatRoutingPath(routingPath string, width int) string {
	routingPath = " " + routingPath + " "
	routingPath = padRight(routingPath, width)
	if !l.Config.Colored {
		return routingPath
	}
	return consts.ANSI_BLUE + routingPath + consts.ANSI_RESET
}

func (l *Logger) formatUpdate(u tgbotapi.Update, width int) string {
	switch {
	// ===== POLL ANSWER (голос) =====
	case u.PollAnswer != nil:
		return l.formatUpdateType(width, UT_PollAnswer, consts.ANSI_BG_BLUE, consts.ANSI_ITALIC)

	// ===== MESSAGE =====
	case u.Message != nil && u.Message.Poll != nil:
		return l.formatUpdateType(width, UT_Poll, consts.ANSI_BG_BLUE, consts.ANSI_WHITE)

	case u.Message != nil && u.Message.IsCommand():
		return l.formatUpdateType(width, UT_Command, consts.ANSI_BG_YELLOW, consts.ANSI_BLACK)

	case u.Message != nil && u.Message.Text != "":
		return l.formatUpdateType(width, UT_TextMessage, consts.ANSI_BG_BLUE, consts.ANSI_WHITE)

	case u.Message != nil:
		return l.formatUpdateType(width, UT_Message, consts.ANSI_BG_BLUE, consts.ANSI_WHITE)

	// ===== CALLBACK =====
	case u.CallbackQuery != nil:
		return l.formatUpdateType(width, UT_CallbackQuery, consts.ANSI_BG_GREEN)

	// ===== EDITED =====
	case u.EditedMessage != nil:
		return l.formatUpdateType(width, UT_EditedMessage, consts.ANSI_BG_CYAN, consts.ANSI_ITALIC)

	// ===== CHANNEL =====
	case u.ChannelPost != nil:
		return l.formatUpdateType(width, UT_ChannelPost, consts.ANSI_BG_MAGENTA)

	case u.EditedChannelPost != nil:
		return l.formatUpdateType(width, UT_EditedChannelPost, consts.ANSI_BG_MAGENTA, consts.ANSI_ITALIC)

	// ===== INLINE =====
	case u.InlineQuery != nil:
		return l.formatUpdateType(width, UT_InlineQuery, consts.ANSI_BG_GREEN)

	case u.ChosenInlineResult != nil:
		return l.formatUpdateType(width, UT_ChosenInlineResult, consts.ANSI_BG_GREEN)

	// ===== PAYMENTS =====
	case u.ShippingQuery != nil:
		return l.formatUpdateType(width, UT_ShippingQuery, consts.ANSI_BG_RED, consts.ANSI_BOLD)

	case u.PreCheckoutQuery != nil:
		return l.formatUpdateType(width, UT_PreCheckoutQuery, consts.ANSI_BG_RED, consts.ANSI_BOLD)

	// ===== UNKNOWN =====
	default:
		return l.formatUpdateType(width, UT_Unknown, consts.ANSI_BG_BLACK, consts.ANSI_BRIGHT_RED, consts.ANSI_BOLD)
	}
}

func (l *Logger) formatUpdateType(width int, ut string, colors ...string) string {
	text := string(ut)
	text = " " + text + " "

	if !l.Config.Colored {
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

func (l *Logger) formatUpdateDetails(v reflect.Value, level int) string {
	var sb strings.Builder
	l.formatUpdateDetailsToBuilder(&sb, v, level)
	return sb.String()
}

func (l *Logger) formatUpdateDetailsToBuilder(sb *strings.Builder, v reflect.Value, level int) {
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
				l.formatUpdateDetailsToBuilder(sb, fv.Index(j), level+2)
			}
			sb.WriteString(tabInner)
			sb.WriteString("]\n")
			continue
		}

		if fv.Kind() == reflect.Struct || fv.Kind() == reflect.Ptr || fv.Kind() == reflect.Interface {

			l.formatUpdateDetailsToBuilder(sb, fv, level+1)
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

type LogConfig struct {
	TimestampWidth      int
	DelayWidth          int
	StatusWidth         int
	ProcessingTimeWidth int
	UpdateWidth         int
	Usernamewidth       int
	RoutingPathWidth    int

	LogBehaviour  byte
	ErrBehaviour  byte
	FileBehaviour byte

	Colored bool
}

const (
	LogStatusOK   byte = 0b1000_0000
	LogStatusWarn byte = 0b0100_0000
	LogStatusErr  byte = 0b0010_0000
	LogStatusFall byte = 0b0001_0000

	LogDetailsOk   byte = 0b0000_1000
	LogDetailsWarn byte = 0b0000_0100
	LogDetailsErr  byte = 0b0000_0010
	LogDetailsFall byte = 0b0000_0001

	LogAllDetailed byte = 0b1111_1111
	LogAll         byte = 0b1111_0000
	LogDetails     byte = 0b0000_1111
	LogNothing     byte = 0b0000_0000
)

var defaultLogConfig LogConfig = LogConfig{
	TimestampWidth:      0,
	DelayWidth:          10, //10
	StatusWidth:         6,  //6
	ProcessingTimeWidth: 11, //11
	UpdateWidth:         17, //17
	Usernamewidth:       0,
	RoutingPathWidth:    0,

	LogBehaviour:  LogAll | LogDetailsErr | LogDetailsFall,
	ErrBehaviour:  LogNothing,
	FileBehaviour: LogNothing,
}

const (
	UT_TextMessage        = "TEXT"
	UT_Message            = "MESSAGE"
	UT_Command            = "COMMAND"
	UT_CallbackQuery      = "CALLBACKQ"
	UT_EditedMessage      = "EDIT"
	UT_ChannelPost        = "CHANNELPOST"
	UT_EditedChannelPost  = "EDITCHANPOST"
	UT_InlineQuery        = "INLINEQ"
	UT_ChosenInlineResult = "CHOSENINLINERES"
	UT_ShippingQuery      = "SHIPPINGQ"
	UT_PreCheckoutQuery   = "PRECHECKOUT_Q"
	UT_Poll               = "POLL"
	UT_PollAnswer         = "POLLANSWER"
	UT_Unknown            = "UNKNOWN"
)

var hints sync.Once

func (l *Logger) printHint() {
	s := "Hint: behaviour bits: Uppercase - default log, downcase - details;\n            O_k, W_arn, E_rr, F_all\tOWEF owef\n"
	s += fmt.Sprintf("            Log output behaviour:\t%04b %04b\n", l.Config.LogBehaviour>>4, l.Config.LogBehaviour&0x0F)
	s += fmt.Sprintf("            Err output behaviour:\t%04b %04b\n", l.Config.ErrBehaviour>>4, l.Config.ErrBehaviour&0x0F)
	s += fmt.Sprintf("            File output behaviour:\t%04b %04b\n", l.Config.FileBehaviour>>4, l.Config.FileBehaviour&0x0F)
	l.Write(s)

	l.Write("Default logging columns:\n      Recieved timestamp; responce delay; status; handling time; update type; username; routing path; details\n")
}

func DefaultLogging(s *Server) MiddlewareFunc {
	return s.Logger.defaultLogMiddleware
}

func (l *Logger) defaultLogMiddleware(ctx Context, update tgbotapi.Update, next HandlerFunc) (status HandleStatus, err *BotError) {

	status, err = next(ctx, update)

	timestampHandled := time.Now()

	hints.Do(l.printHint)

	go func() {
		timestampRecieved := ctx.Value(ContextKey_TimestampRecieved).(time.Time)
		log := NewLog(ctx, update, status, err, timestampRecieved, timestampHandled)

		l.Log(log)
		l.Err(log)
		l.FileWrite(log)
	}()
	return status, err
}
