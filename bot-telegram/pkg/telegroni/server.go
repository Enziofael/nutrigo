// - telegroni package provides a flexible routing system for Telegram bots
// in gin-gonic like style
//
// It allows registering routes with match conditions and handlers,
// grouping routes, and nesting groups for complex routing logic.
//
// # Key Components:
//
//   - **Server** — main entry point for routing.
//   - **HandlerGroup** — groups routes with shared match conditions.
//   - **Handler** — individual route handler.
//   - **Middleware** — wraps handlers for cross-cutting concerns.
//
// # Example:
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
//
// # See also:
//
//	package defaults // for pre-made matchers, functions, middlewares
package types

import (
	"context"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ======================= INDEX =========================
//          You can navigate by searching REGION
// =======================================================
//	Arguments are omitted for brevity
//
//	1. type ServerConfig struct
//	- Exported
//		- func NewConfig(...) *ServerConfig
//		- func (cfg *ServerConfig) SetLimit(...)
//		- func (cfg *ServerConfig) SetTimeout(...)
//		- func (cfg *ServerConfig) SetAllowedUpdates(...)
//
//	2. type Route interface
//
//	3. type Server struct
//	- Exported
//		- func New(...) *Server
//		- func (s *Server) Apply(...)
//		- func (s *Server) Handle(...)
//		- func (s *Server) Group(...) *HandlerGroup
//		- func (s *Server) Start() *BotError
//		- const ContextKey_Bot
//		- const ContextKey_TimestampRouted
//		- const ContextKey_Path
//	- Unexported
//		- func (s *Server) handle(...)
//		- func defaultRoutingFallback(...)
//
// =======================================================

// ======================== TYPE =========================
// REGION               ServerConfig
// =======================================================

// - ServerConfig is a struct for all values
// that are necessary for server configuration
//
// # It is used for Server creation
//
// # Related:
//
//	//Used in:
//	func New()
//	func defaults.Default()
//
// # See also:
//
//	func NewConfig()
//
// # Methods:
type ServerConfig struct {
	BotApiToken string
	ApiConfig   tgbotapi.UpdateConfig
}

// ====================== EXPORTED ======================

// - NewConfig() returns ServerConfig for creating a server.
//
// - BotApiToken is a token to acces your bot.
// You can find it in @BotFather.
// Revoke it and use a new one if you think it was сompromised
//
// It is HIGHLY RECOMMENDED NOT CREATE SERVER CONFIG MANUALLY
// using a strict literal because
//
// - empty BotApiToken will cause panic
// - invalid ApiConfig may impair server functionality
// which may cause losing updates.
//
// # ALWAYS USE NewConfig() unless you have a huge and specific reason not to
//
// You can manually set:
//
// - limit to control server load
// - timeout to control amount of queries to the Telegram API
// - allowed updates to control which types you want to get from Telegram API
//
// # Also you can set allowed updates
//
// # Related:
//
//	//Uses::
//	type ServerConfig struct
//
// # See also:
//
//	func (cfg *ServerConfig) SetLimit()
//	func (cfg *ServerConfig) SetTimeout()
//	func (cfg *ServerConfig) SetAllowedUpdates()
func NewConfig(BotApiToken string) ServerConfig {
	return ServerConfig{
		BotApiToken: BotApiToken,
		ApiConfig: tgbotapi.UpdateConfig{
			Offset:         0,
			Limit:          100,
			Timeout:        30,
			AllowedUpdates: []string{},
		},
	}
}

// - SetLimit() changes polling limit in ServerConfig
//
// - limit is a max amount of updates Telegram API will send at once.
//
// Default value equals to 100 (bot default buffer size)
//
// # Similar:
//
//	func (cfg *ServerConfig) SetTimeout()
//	func (cfg *ServerConfig) SetAllowedUpdates()
func (cfg ServerConfig) SetLimit(limit int) {
	cfg.ApiConfig.Limit = limit
}

// - SetTimeout() changes polling timeout in ServerConfig
//
// - timeout defines how long will Telegram API will wait for Updates
// if currently there're no new Updates before sending response
//
// # Default value equals to 30
//
// # Similar:
//
//	func (cfg *ServerConfig) SetLimit()
//	func (cfg *ServerConfig) SetAllowedUpdates()
func (cfg ServerConfig) SetTimeout(timeout int) {
	cfg.ApiConfig.Timeout = timeout
}

// - SetAllowedUpdates() changes polling filter in ServerConfig
//
// - allowed defines which Update's types Telegram API will send
//
// Default value equals is empty (all updates)
//
// Change carefully so that you don't miss the updates you want to process.
// Use constants from tgbotapi that starts with "UpdateType..."
//
// # Example:
//
//	cfg.SetAllowedUpdates([]string{
//		tgbotapi.UpdateTypeCallbackQuery,
//		tgbotapi.UpdateTypeMessage,
//	})
//
// # Similar:
//
//	func (cfg *ServerConfig) SetLimit()
//	func (cfg *ServerConfig) SetAllowedUpdates()
func (cfg ServerConfig) SetAllowedUpdates(allowed []string) {
	cfg.ApiConfig.AllowedUpdates = allowed
}

// ===================== INTERFACE =======================
// REGION                  Route
// =======================================================

// - Route defines the common interface for both Handlers and HandlerGroups
//
// It allows the Server to process routint uniformly without knowing if it's a Handler or a HandlerGroup,
// considering them as identical routing nodes in Server's own handle() function.
//
// Handler's and HandlerGroup's handle() implementations
// first process mathcing logic
//
// In the case of Handler,
// after matching it calls wrapped in MiddlewareFuncs (if any) HandlerFunc.
//
// In the case of HandlerGroup,
// after matching it starts routing inside nested Routes.
//
// # See also:
//
//	type Middleware struct
//	type HandlerGroup struct
//	type Server struct
//	type Route struct
type Route interface {
	handle(ctx context.Context, u tgbotapi.Update) (matched bool, status string, err *BotError)
}

// ======================== TYPE =========================
// REGION                  Server
// =======================================================

// - Server is a main struct for all Routes, global Middlewares
// and running Bot's server
//
// It manages main loop by getting Updates and passing them to
// Routes recursively in order until one of the Handlers will match.
// If none of them matched it will pass Update to the
// RoutingFallbackFunc.
//
// Server's logic is based on consecutive stages:
//
// 1. Routing tree registration.
//
// 2. Running server's loop.
//
// For each recieved Update server will:
//
// 2.1. Route it based on routing tree
//
// 2.2. Handle it by matched Handler asynchronously
//
// Routing tree consists of Server, HandlerGroups and Handlers,
// where Server is the root from where routing starts,
// HandlerGroups are brunches that groups similar handlers with common match logic
// and Handlers are leaves where Update starts its processing.
// It can be build by Group() and Handle() functions,
// wrapping whole tree, branche or leaf with necessary Middlewares.
//
// Finding matching Handler works by In-depth-first traversal
// by adding order based on MatchFuncs.
//
// Server can safely not have any matching Handler
// if specified RoutingFallbackFunc doesn't say the opposite,
// but it exactly means that user won't get any response and Update
// won't be processed.
//
// It is RECOMMENDED to add a Handler with Any() as the last one added
// or add a global Middleware that will handle unrouted Update.
// (for example logging unhandled Update's details).
//
// You can use defaults.Default() to create a Server with pre-applied default
// logger and panic recovery Middlewares.
//
// Or you can use New() to create a blank Server.
//
// You can find simple example in the package doc.
//
// # Related:
//
//	//Implements:
//	type Route interface
//
//	//Uses:
//	type Route interface
//	type Middleware struct
//
// # Similar:
//
//	type Handler struct
//	type Middleware struct
//	type Server struct
//	type Route interface
//
// # See also:
//
//	type BotError struct
//
//	package defaults // for  Any() function, which return MatchFunc always returning true,
//					 // Default() and other default functions
//
// # Methods:
type Server struct {
	Context         context.Context
	Routes          []Route
	Middlewares     []Middleware
	RoutingFallback RoutingFallbackFunc
	Config          ServerConfig
}

// ====================== EXPORTED ======================

// - New() creates new Server, creates a BotAPI, adds it to the Context
// and returns a pointer to it and BotError if something went with the BotAPI.
//
// - config contains all the values for the server configuration.
//
// You CAN create Server using a struct literal directly:
//
//	srv := Server{Context: context.Background(), Routes: make([]Route, 0), Middlewares: make([]Middleware, 0),}
//
// However, this is NOT RECOMMENDED
// because
// - if at least one of the Context, Routes, Config or Middlewares are nil n
// it will cause panic
// - if Config doesn't contain bot api token
// it will cause panic
//
// Use New() or defaults.Default() unless you have a specific reason not to
// (e.g., performance-critical code where you're certain name is unnecessary
// and internal calls are a concern).
//
// # Related:
//
//	//Uses:
//	type ServerConfig struct
//
// # Similar:
//
//	func NewHandler()
//	func NewMiddleware()
//	func NewHandlerGroup()
func New(config ServerConfig) *Server {
	return &Server{
		//Root context. Contains *BotAPI and Timestamp when Update was routed
		Context:         context.Background(),
		Routes:          make([]Route, 0),
		Middlewares:     make([]Middleware, 0),
		RoutingFallback: defaultRoutingFallback,
		Config:          config,
	}
}

// - Apply() appends a new global Middleware to Server
//
// - function contains all Middleware's logic.
//
// - name is primarily used for logging and debugging.
//
// All Middlewares will be applied recursively to ALL Handlers
// by apply(), that will wrap their own HandleFunctions.
//
// # Example:
//
//	s.Apply(myMiddlewareFunc, "myMiddleware")
//
// # Related:
//
//	//Uses:
//	type MiddlewareFunc func()
//	func NewMiddleware()
//
// # Similar:
//
//	func (s *Server) Apply()
//
// # See also:
//
//	func NewHandler() // for more details about name
func (s *Server) Apply(function MiddlewareFunc, name string) {
	s.Middlewares = append(s.Middlewares, NewMiddleware(function, name))
}

// - Handle() adds a new Handler to Server's Routes
// and wraps Handler with Server's Middlewares and provided
// local middlewares
//
// - matcher contains logic to define if this Handler should process Update.
//
// - function contains all Handler's logic.
//
// - name is primarily used for logging and debugging.
//
// - middlewares are intended to apply to this Handler only.
//
// # Example:
//
//	s.Handle(myMatcherFunc, myHandlerFunc, "myHandler",
//		someMiddleware1, someMiddleware2
//	)
//
// # Related:
//
//	//Uses:
//	type MatchFunc func()
//	type HandlerFunc func()
//	type Middleware struct
//	func NewHandler()
//
// # Similar:
//
//	func (s *Server) Handle()
//
// # See also:
//
//	func NewHandler() // for more details about name
func (s *Server) Handle(matcher MatchFunc, function HandlerFunc, name string, middlewares ...Middleware) {
	wrapped := function
	for i := len(middlewares) - 1; i >= 0; i-- {
		wrapped = middlewares[i].apply(wrapped)
	}
	for i := len(s.Middlewares) - 1; i >= 0; i-- {
		wrapped = s.Middlewares[i].apply(wrapped)
	}

	s.Routes = append(s.Routes, NewHandler(matcher, wrapped, name))
}

// - Group() adds a new HandlerGroup to Server's Routes
// and adds all Server's Middlewares to it,
// then returns a pointer to it
//
// - matcher contains logic to define if this HandlerGroup should process Update.
//
// - name is primarily used for logging and debugging.
//
// In server's main loop if thiw new HandlerGroup will match,
// it will start passing processing to each its
// Route in order they were added, until one of Handlers will match
// or none of them (then routing passes to the Route next to HandlerGroup)
//
// # Example:
//
//	grp := s.Group(myMatcherFunc, "myGroup")
//
// # Related:
//
//	//Uses:
//	type MatchFunc func()
//	func NewGroup()
//
// # See also:
//
//	func NewHandlerGroup() // for more details about name
//
// # Similar:
//
//	func (s *Server) Group()
func (s *Server) Group(matcher MatchFunc, name string) *HandlerGroup {
	newGroup := NewHandlerGroup(matcher, name)

	newGroup.Middlewares = make([]Middleware, len(s.Middlewares))
	copy(newGroup.Middlewares, s.Middlewares)

	s.Routes = append(s.Routes, newGroup)
	return newGroup
}

// - Start() begins the main loop, polling Telegram for updates.
// It processes each update through the routing chain sequentially.
// The function blocks until the bot's update channel is closed.
func (s *Server) Start() *BotError {
	bot, err := tgbotapi.NewBotAPI(s.Config.BotApiToken)
	if err != nil {
		return NewBotError("Bot Creation failed", NewBotError(err.Error(), nil))
	}
	s.Context = context.WithValue(s.Context, ContextKey_Bot, bot)

	u := s.Config.ApiConfig

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		ctx := context.WithValue(s.Context, ContextKey_TimestampRouted, time.Now())
		go s.handle(ctx, update)
	}

	return nil
}

// - ContextKey_Bot for getting bot from context
const ContextKey_Bot = "bot"

// - ContextKey_TimestampRouted for getting timestamp when update was routed from context
const ContextKey_TimestampRouted = "timestamp_routed"

// - ContextKey_Bot for getting routed update's route path from context
const ContextKey_Path = "route_path"

// ===================== UNEXPORTED ======================

// - handle() starts passing processing to each Server's
// Route recursively in turn they were added, until one of Handlers will match
// or none of them.
//
// If no Handler matched update passing to the Server's RoutingFallbackFunc.
//
// It is used internally in the server's main loop.
//
// # Related:
//
//	//Used in:
//	func (s *Server) Start()
func (s *Server) handle(ctx context.Context, u tgbotapi.Update) {
	var matched bool
	var status string
	var err *BotError

	ctx = context.WithValue(ctx, ContextKey_Path, "")

	for _, route := range s.Routes {
		matched, status, err = route.handle(ctx, u)
		if matched {
			break
		}
	}
	if !matched {
		s.RoutingFallback(u, status, err)
	}
}

// - defaultRoutingFallback() is a default RoutingFallbackFunc
//
// # It logs update details, error and status
//
// # See also:
//
//	func (s *Server) Start()
//	type Server struct
func defaultRoutingFallback(u tgbotapi.Update, status string, err *BotError) {
	log.Printf("[BOT] FALLBACK Unhandled update:\n\nUpdate:\n%#v\n\nError:\n%#v\n\nStatus:\n%s", u, err, status)
}
