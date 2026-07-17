// Package server provides a flexible routing system for Telegram bots.
// It allows registering routes with match conditions and handlers,
// grouping routes, and nesting groups for complex routing logic.
//
// Example:
//
//	import (
//		tgbot "github.com/Enziofael/nutrigo/bot-telegram/pkg/telegroni"
//		tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
//	)
//
//	func main() {
//		// Creating a server
//		srv := tgbot.New("your_bot_private_token")
//
//		// Routes registration
//		srv.Use(tgbot.Command("start"), handleStart)
//
//		// Starting server's loop
//		srv.Start()
//	}
//
//	func handleStart(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
//		// Do something
//	}
package telegroni

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ============================================
// Route
// ============================================

// Route defines the common interface for both single handlers and handler groups.
// It allows the server to process routes uniformly without knowing their internal structure.
type Route interface {
	// handle processes the update and returns true if the route matched and handled it.
	// If false, the server continues to the next route in the chain.
	handle(b *tgbotapi.BotAPI, u tgbotapi.Update) (matched bool)
}

// MatchFunc determines whether a route or group should handle a given update.
// It receives the bot instance and the update to allow complex matching logic.
type MatchFunc func(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool

// HandlerFunc is the function that processes a matched update.
// It runs asynchronously in a goroutine to avoid blocking the main event loop.
type HandlerFunc func(bot *tgbotapi.BotAPI, update tgbotapi.Update)

// ============================================
// Handler
// ============================================

// Handler represents a single route with a match condition and a handler function.
// When the match condition returns true, the handler is executed asynchronously.
type Handler struct {
	handlerFunc HandlerFunc
	matchFunc   MatchFunc
}

// makeHandler creates a new Handler with the given match and handler functions.
func makeHandler(hf HandlerFunc, mf MatchFunc) Handler {
	return Handler{handlerFunc: hf, matchFunc: mf}
}

// handle checks the match condition and executes the handler if it matches.
// Returns true if the condition was satisfied and rouing was successful (handler may still be running).
func (r Handler) handle(b *tgbotapi.BotAPI, u tgbotapi.Update) (matched bool) {
	if !r.matchFunc(b, u) {
		return false
	}
	go r.handlerFunc(b, u)
	return true
}

// ============================================
// HandlerGroup
// ============================================

// HandlerGroup is a collection of routes that share a common match condition.
// If the group's condition matches, it delegates routing to its child routes
// in the order they were registered. The first route that handles the update
// stops further processing within the group.
//
// HandlerGroups can be nested to create hierarchical routing structures.
type HandlerGroup struct {
	matchFunc MatchFunc
	routes    []Route
}

// makeHandlerGroup creates a new HandlerGroup with the given match condition.
func makeHandlerGroup(mf MatchFunc) HandlerGroup {
	return HandlerGroup{matchFunc: mf, routes: make([]Route, 0)}
}

// handle checks the group's condition and, if true, processes child routes
// sequentially. Stops at the first child route that handles the update.
// Returns true if any child route handled the update.
func (g HandlerGroup) handle(b *tgbotapi.BotAPI, u tgbotapi.Update) (matched bool) {
	if !g.matchFunc(b, u) {
		return false
	}
	for _, route := range g.routes {
		if route.handle(b, u) {
			return true
		}
	}
	return false
}

// Use adds a new handler to the group's routing chain.
//
// Routes are processed in the order they are registered.
func (g *HandlerGroup) Use(mf MatchFunc, hf HandlerFunc) {
	g.routes = append(g.routes, makeHandler(hf, mf))
}

// Group creates a new nested group and returns a pointer to it.
// The pointer allows modification of the nested group after creation.
//
// Routes are processed in the order they are registered.
func (g *HandlerGroup) Group(mf MatchFunc) *HandlerGroup {
	newGroup := makeHandlerGroup(mf)
	g.routes = append(g.routes, newGroup)
	return &newGroup
}

// ============================================
// Server
// ============================================

// Server manages the Telegram bot's main event loop and routing system.
// It holds the routing chain and the bot instance, orchestrating the
// entire update handling process.
type Server struct {
	routes []Route
	bot    *tgbotapi.BotAPI
}

// New creates a new Server instance with the given bot.
// The server starts with an empty routing chain.
func New(bot_token string) *Server {
	bot, err := tgbotapi.NewBotAPI(bot_token)
	if err != nil {
		log.Fatalf("Panic: bot creation failed. Error: \"%s\"", err)
	}

	//------
	bot.Debug = true
	//------
	log.Printf(`Bot init successful.
	Authorized on account @%s
	ID: %d
	CanJoinGroups: %t
	CanReadAllGroupMessages: %t
	SupportsInlineQueries: %t`,
		bot.Self.UserName, bot.Self.ID, bot.Self.CanJoinGroups, bot.Self.CanReadAllGroupMessages, bot.Self.SupportsInlineQueries)

	return &Server{bot: bot, routes: make([]Route, 0)}
}

// Use adds a new handler to the server's routing chain.
//
// Routes are processed in the order they are registered.
func (s *Server) Use(mf MatchFunc, hf HandlerFunc) error {
	s.routes = append(s.routes, makeHandler(hf, mf))
	return nil
}

// Group creates a new group and returns a pointer to it.
// The pointer allows modification of the group after creation.
//
// Routes are processed in the order they are registered.
func (s *Server) Group(mf MatchFunc) *HandlerGroup {
	newGroup := makeHandlerGroup(mf)
	s.routes = append(s.routes, newGroup)
	return &newGroup
}

// Start begins the main event loop, polling Telegram for updates.
// It processes each update through the routing chain sequentially.
// The function blocks until the bot's update channel is closed.
//
// Each update is processed asynchronously after it was routed
// and handeled to the FIRST matched handler,
// which runs in its own goroutine to prevent blocking the main loop.
func (s *Server) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := s.bot.GetUpdatesChan(u)
	log.Printf("Bot is running. Ready to receive updates\n\n")

	// В цикле проходимся по каналу апдейтов, при получении раскидываем их соответствующим роутерам в горутины, сам цикл при этом идет дальше.
	// TODO: Для пакета лучше переделать с регистрацией роутеров по типу с передачей функций
	for update := range updates {
		matched := false
		for _, r := range s.routes {
			if r.handle(s.bot, update) {
				matched = true
				break // Stop processing after first match
			}
		}

		if !matched {
			log.Printf("Warning: unhandled update (type: %T)", update)
		}
	}
	return nil
}

// PREMADE

func IsCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
	return update.Message != nil && update.Message.IsCommand()
}

func Command(command string) MatchFunc {
	return func(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
		return IsCommand(bot, update) &&
			update.Message.Command() == command
	}
}

func IsCallbackQuery(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
	return update.CallbackQuery != nil
}

func CallbackQuery(data string) MatchFunc {
	return func(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
		return IsCallbackQuery(bot, update) &&
			update.CallbackQuery.Data == data
	}
}

func IsAny(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
	return true
}

func Any() MatchFunc {
	return func(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool { return IsAny(bot, update) }
}

func HandlerFuncStub(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
}
