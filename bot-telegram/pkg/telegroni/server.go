// Package telegroni provides a flexible routing system for Telegram bots.
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
//		srv := tgbot.New(tgbot.ServerConfig{
//			BotConfig: tgbot.BotConfig{APIToken: "your_bot_private_token"},
//		})
//
//		// Global middleware
//		srv.Apply(tgbot.NewMiddleware(LoggerMiddleware, "logger"))
//
//		// Routes registration
//		srv.Use(tgbot.Command("start"), handleStart, "start")
//
//		// Starting server's loop
//		srv.Start()
//	}
//
//	func handleStart(ctx context.Context, update tgbotapi.Update) {
//		// Do something
//	}
//
//	func LoggerMiddleware(ctx context.Context, update tgbotapi.Update, next tgbot.HandlerFunc) {
//		log.Println("Before")
//		next(ctx, update)
//		log.Println("After")
//	}
package telegroni

import (
	"context"
	"fmt"
	"log"
	"time"

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
	handle(ctx context.Context, u tgbotapi.Update) (matched bool, status string, err *BotError)
	// Name returns the name of the route for logging and debugging.
	Name() string
}

// MatchFunc determines whether a route or group should handle a given update.
type MatchFunc func(ctx context.Context, update tgbotapi.Update) bool

// HandlerFunc is the function that processes a matched update.
type HandlerFunc func(ctx context.Context, update tgbotapi.Update) (status string, err *BotError)

// MiddlewareFunc is a function that wraps a handler.
// It can execute code before and after the handler,
// and can choose to call the next handler or not.
type MiddlewareFunc func(ctx context.Context, update tgbotapi.Update, next HandlerFunc) (status string, err *BotError)

// ============================================
// Middleware
// ============================================

// Middleware represents a middleware with a name for logging.
type Middleware struct {
	middlewareFunc MiddlewareFunc
	name           string
}

// NewMiddleware creates a new middleware.
func NewMiddleware(mwf MiddlewareFunc, name string) Middleware {
	return Middleware{
		middlewareFunc: mwf,
		name:           name,
	}
}

// apply wraps a handler with the middleware.
func (m Middleware) apply(handler HandlerFunc) HandlerFunc {
	return func(ctx context.Context, update tgbotapi.Update) (status string, err *BotError) {
		return m.middlewareFunc(ctx, update, handler)
	}
}

// Name returns the middleware name.
func (m Middleware) Name() string {
	return m.name
}

// ============================================
// Handler
// ============================================

// Handler represents a single route with a match condition and a handler function.
// When the match condition returns true, the handler is executed asynchronously.
type Handler struct {
	handlerFunc HandlerFunc
	matchFunc   MatchFunc
	name        string
}

// NewHandler creates a new Handler with the given match and handler functions.
func NewHandler(hf HandlerFunc, mf MatchFunc, name string) *Handler {
	return &Handler{
		handlerFunc: hf,
		matchFunc:   mf,
		name:        name,
	}
}

// Name returns the handler name.
func (h *Handler) Name() string {
	return h.name
}

// handle checks the match condition and executes the handler if it matches.
// Returns true if the condition was satisfied and routing was successful (handler may still be running).
func (h *Handler) handle(ctx context.Context, u tgbotapi.Update) (matched bool, status string, err *BotError) {
	if !h.matchFunc(ctx, u) {
		return false, StatusWarn, NewBotError(fmt.Sprintf("[Handler] %s unmatched", h.name), nil)
	}
	status, err = h.handlerFunc(ctx, u)
	return true, status, err
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
	matchFunc   MatchFunc
	routes      []Route
	middlewares []Middleware
	name        string
}

// NewHandlerGroup creates a new HandlerGroup with the given match condition.
func NewHandlerGroup(mf MatchFunc, name string) *HandlerGroup {
	return &HandlerGroup{
		matchFunc:   mf,
		routes:      make([]Route, 0),
		middlewares: make([]Middleware, 0),
		name:        name,
	}
}

// Name returns the group name.
func (g *HandlerGroup) Name() string {
	return g.name
}

// Apply adds a middleware to the group.
func (g *HandlerGroup) Apply(mw Middleware) {
	g.middlewares = append(g.middlewares, mw)
}

// Use adds a new handler to the group's routing chain.
// Routes are processed in the order they are registered.
func (g *HandlerGroup) Use(mf MatchFunc, hf HandlerFunc, name string) {
	//log.Printf("[Group] Adding handler: %s -> %s", g.name, name)

	wrapped := hf
	for i := len(g.middlewares) - 1; i >= 0; i-- {
		wrapped = g.middlewares[i].apply(wrapped)
	}

	g.routes = append(g.routes, NewHandler(wrapped, mf, name))
}

// Group creates a new nested group and returns a pointer to it.
// The pointer allows modification of the nested group after creation.
// Routes are processed in the order they are registered.
func (g *HandlerGroup) Group(mf MatchFunc, name string) *HandlerGroup {
	//log.Printf("[Group] Adding nested group: %s -> %s", g.name, name)

	newGroup := NewHandlerGroup(mf, name)

	// Copy parent's middlewares to the new group
	newGroup.middlewares = make([]Middleware, len(g.middlewares))
	copy(newGroup.middlewares, g.middlewares)

	g.routes = append(g.routes, newGroup)
	return newGroup
}

// handle checks the group's condition and, if true, processes child routes
// sequentially. Stops at the first child route that handles the update.
// Returns true if any child route handled the update.
func (g *HandlerGroup) handle(ctx context.Context, u tgbotapi.Update) (matched bool, status string, err *BotError) {
	if !g.matchFunc(ctx, u) {
		return false, StatusWarn, NewBotError(fmt.Sprintf("[Group] %s: unmatched", g.name), nil)
	}

	for _, route := range g.routes {
		var innerErr *BotError
		matched, status, innerErr = route.handle(ctx, u)
		if innerErr != nil {
			err = NewBotError("", innerErr)
		}
		if matched {
			return true, status, err
		}
	}
	return false, StatusWarn, NewBotError(fmt.Sprintf("[Group] %s: no matched handler", g.name), err)
}

// ============================================
// Server
// ============================================

// Server manages the Telegram bot's main event loop and routing system.
// It holds the routing chain and the bot instance, orchestrating the
// entire update handling process.
type Server struct {
	Context     context.Context
	routes      []Route
	middlewares []Middleware
}

// ServerConfig configures the server.
type ServerConfig struct {
	BotConfig BotConfig
}

// BotConfig configures the bot.
type BotConfig struct {
	APIToken string
}

// New creates a new Server instance with the given config.
func New(cfg ServerConfig) (*Server, *BotError) {
	bot, err := tgbotapi.NewBotAPI(cfg.BotConfig.APIToken)
	if err != nil {
		return nil, NewBotError("Bot Creation failed", NewBotError(err.Error(), nil))
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "bot", bot)

	return &Server{
		Context:     ctx,
		routes:      make([]Route, 0),
		middlewares: make([]Middleware, 0),
	}, nil
}

// Apply adds a global middleware to the server.
func (s *Server) Apply(mw Middleware) {
	s.middlewares = append(s.middlewares, mw)
}

// Use adds a new handler to the server's routing chain.
// Routes are processed in the order they are registered.
func (s *Server) Use(mf MatchFunc, hf HandlerFunc, name string) *BotError {
	wrapped := hf
	for i := len(s.middlewares) - 1; i >= 0; i-- {
		wrapped = s.middlewares[i].apply(wrapped)
	}

	s.routes = append(s.routes, NewHandler(wrapped, mf, name))
	return nil
}

// Group creates a new group and returns a pointer to it.
// The pointer allows modification of the group after creation.
// Routes are processed in the order they are registered.
func (s *Server) Group(mf MatchFunc, name string) *HandlerGroup {
	newGroup := NewHandlerGroup(mf, name)

	newGroup.middlewares = make([]Middleware, len(s.middlewares))
	copy(newGroup.middlewares, s.middlewares)

	s.routes = append(s.routes, newGroup)
	return newGroup
}

// handle processes a single update through the routing chain.
func (s *Server) handle(ctx context.Context, u tgbotapi.Update) {
	var (
		matched bool
	)
	for _, route := range s.routes {
		if matched, _, _ = route.handle(ctx, u); matched {
			return
		}
	}
	if !matched {
		LogGlobalError(ctx, u)
	}
}

// Start begins the main event loop, polling Telegram for updates.
// It processes each update through the routing chain sequentially.
// The function blocks until the bot's update channel is closed.
func (s *Server) Start() *BotError {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := s.Context.Value("bot").(*tgbotapi.BotAPI).GetUpdatesChan(u)
	
	for update := range updates {
		ctx := context.WithValue(s.Context, "timestamp_recieved", time.Now())
		go s.handle(ctx, update)
	}

	return nil
}

// ============================================
// Stubs for testing
// ============================================

// HandlerFuncStub is a stub handler for testing.
func HandlerFuncStub(ctx context.Context, update tgbotapi.Update) (status string, err *BotError) {
	log.Printf("[Stub] Handler called for update: %v", update.UpdateID)
	return StatusOK, nil
}
