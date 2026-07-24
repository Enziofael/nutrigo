package telegroni

import (
	"context"
	"fmt"
	"time"

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
// REGION                  Handler
// =======================================================

// - Handler is a struct for HandlerFunc that processes Update
// if Update matches conditions determined by MatchFunc
//
// Handler is intended as a part of routing -
// to define which scenario to start by
// processing data received from sources such as database,
// external APIs, etc.
//
// Handler is NOT intended and NOT recommended to do some business-logic,
// sending bot responses, change any data in databases, etc. directly in HandlerFunc.
//
// "Routing is for defining WHAT was received"
// "Handler is for defining WHAT should be done with it"
// "Scenario is for DOING something"
//
// But you CAN do this if HandlerFunc is pretty small, doesn't do anything hard or
// you have a specific reason to.
//
// Handler can safely return
// before calling any scenario and sending any response to the user,
// but it exactly means that user won't get any response
// (for example ignoring unauthorized actions if you want to)
//
// # Example:
//
//	handler := NewHandler(func(ctx context.Context, update tgbotapi.Update) (status HandleStatus, err *BotError) {
//		// Processing update, ctx, etc.
//		// Defining scenario example:
//		switch value{
//		case 1:
//			status, err = RunScenario1(ctx, update)
//		case 2:
//			status, err = RunScenario1(ctx, update)
//		case 3:
//			return status, err
//		default:
//			status, err = RunScenario1(ctx, update)
//		}
//
//		// Post processing scenario results, status, errors
//
//		// returning status and error to the Middlewares and Server
//		return status, err
//	})
//
// # Related:
//
//	//Uses:
//	type HandlerFunc func()
//	type MatchFunc func()
//
//	//Used in:
//	func (s *Server) Handle()
//	func (g *HandlerGroup) Handle()
//
// # Similar:
//
//	type Middleware struct
//	type HandlerGroup struct
//	type Server struct
//	type Route struct
//
// # Methods:
type Handler struct {
	HandlerFunc HandlerFunc
	MatchFunc   MatchFunc
	Name        string
}

// ====================== EXPORTED ======================

// - NewHandler() creates new Handler and returns a pointer to it
//
// - matcher contains logic to define if this Handler should process Update
//
// - function contains all Handler's logic.
//
// - name is primarily used for logging and debugging.
//
// If name is empty, it will be generated automatically by autoname().
// autoname is formatted like "H0"
// Where H - Handler, 0 - autoname's number
//
// You CAN create Handler using a struct literal directly:
//
//	h := Handler{HandlerFunc: myFunc, MatchFunc: myMatcher, Name: "middlewareName"}
//
// However, this is NOT RECOMMENDED because
// if name is empty it won't appear in logs, making debugging harder
//
// Use NewHandler() unless you have a specific reason not to
// (e.g., performance-critical code where you're certain name is unnecessary
// and internal calls are a concern).
//
// # Related:
//
//	//Uses:
//	type HandlerFunc func()
//	type MatchFunc func()
//
// # Similar:
//
//	func NewMiddleware()
//	func NewHandlerGroup()
//	func NewServer()
func NewHandler(matcher MatchFunc, function HandlerFunc, name string) *Handler {
	if name == "" {
		name = autoname(typeHandler)
	}
	return &Handler{
		HandlerFunc: function,
		MatchFunc:   matcher,
		Name:        name,
	}
}

// ===================== UNEXPORTED ======================

// - handle() checks if MatchFunc returns true and executes HandlerFunc if so
//
// returns true Handlers matched, false otherwise
//
// It is used internally in the server's main loop.
//
// # Related:
//
//	//Implements:
//	type Route interface
//
//	//Used in:
//	func (s *Server) handle()
//	func (g *HandlerGroup) handle()
func (h *Handler) handle(ctx context.Context, u tgbotapi.Update) (matched bool, status HandleStatus, err *BotError) {
	if !h.MatchFunc(ctx, u) {
		return false, StatusWarn, NewBotError(fmt.Sprintf("[Handler] %s unmatched", h.Name), nil)
	}

	ctx = appendPath(ctx, h.Name)
	ctx = context.WithValue(ctx, ContextKey_TimestampRouted, time.Now)
	status, err = h.HandlerFunc(ctx, u)

	return true, status, err
}

// ======================== TYPE =========================
// REGION               HandlerGroup
// =======================================================

// - HandlerGroup is a struct for collection of Routes that pass processing Update
// to Routes if Update matches conditions determined by MatchFunc
//
// HandlerGroup is intended as a part of routing
// for grouping Handlers that share common match conditions.
//
// "Routing is for defining WHAT was received"
// "Handler is for defining WHAT should be done with it"
// "Scenario is for DOING something"
//
// Matched HandlerGroup may safely contain no matching Handler,
// but it exactly means that user WON'T get any response
// if there're no other matched HandlerGroups with matched Handler.
//
// It is RECOMMENDED to add a Handler with Any() as the last one.
//
// # Example:
//
//	g.Handle(defaults.Any(), myAnyHandlerFunc, "AnyHandler")
//
// # Related:
//
//	//Uses:
//	type MatchFunc func()
//	type Route interface
//	type Middleware struct
//
// # Similar:
//
//	type Handler struct
//	type Middleware struct
//	type Server struct
//	type Route struct
//
// # See also:
//
//	package defaults // for Any() function, which return MatchFunc always returning true
//					 // and other default functions
//
// # Methods:
type HandlerGroup struct {
	MatchFunc   MatchFunc
	Routes      []Route
	Middlewares []Middleware
	Name        string
}

// ====================== EXPORTED ======================

// - NewHandlerGroup() creates new Handler and returns a pointer to it
//
// - matcher contains logic to define if this HandlerGroup should process Update
//
// - name is primarily used for logging and debugging.
//
// If name is empty, it will be generated automatically by autoname().
// autoname is formatted like "G0"
// Where G - HandlerGroup, 0 - autoname's number
//
// You CAN create HandlerGroup using a struct literal directly:
//
//	hg := HandlerGroup{MatchFunc: myMatcher, Name: "middlewareName", Routes: make([]Route, 0), Middlewares: make([]Middleware, 0),}
//
// However, this is NOT RECOMMENDED
// because
//
// - if name is empty it won't appear in logs, making debugging harder
//
// - if at least one of the Routes or Middlewares is nil
// it will cause panic
//
// Use NewHandlerGroup() unless you have a specific reason not to
// (e.g., performance-critical code where you're certain name is unnecessary
// and internal calls are a concern).
//
// # Related:
//
//	//Uses:
//	type MatchFunc func()
//
// # Similar:
//
//	func NewMiddleware()
//	func NewHandlerGroup()
//	func NewServer()
func NewHandlerGroup(matcher MatchFunc, name string) *HandlerGroup {
	if name == "" {
		name = autoname(typeHandlerGroup)
	}
	return &HandlerGroup{
		MatchFunc:   matcher,
		Routes:      make([]Route, 0),
		Middlewares: make([]Middleware, 0),
		Name:        name,
	}
}

// - Apply() appends a new Middleware to HandlerGroup
//
// - function contains all Middleware's logic.
//
// - name is primarily used for logging and debugging.
//
// All Middlewares will be applied recursively to all Handlers
// in this HandlerGroup by apply(), that will wrap their own
// HandleFunctions.
//
// # Example:
//
//	g.Apply(myMiddlewareFunc, "myMiddleware")
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
func (g *HandlerGroup) Apply(function MiddlewareFunc, name string) {
	g.Middlewares = append(g.Middlewares, NewMiddleware(function, name))
}

// - Handle() adds a new Handler to HandlerGroup's Routes
// and wraps Handler with HandlerGroup's Middlewares and provided
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
// In server's main loop if HandlerGroup will match,
// it will start passing processing to each
// Route in order they were added, until one of Handlers will match
// or none of them (then routing passes to the Route next to HandlerGroup)
//
// # Example:
//
//	g.Handle(myMatcherFunc, myHandlerFunc, "myHandler",
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
func (g *HandlerGroup) Handle(matcher MatchFunc, function HandlerFunc, name string, middlewares ...Middleware) {
	wrapped := function
	for i := len(middlewares) - 1; i >= 0; i-- {
		wrapped = middlewares[i].apply(wrapped)
	}
	for i := len(g.Middlewares) - 1; i >= 0; i-- {
		wrapped = g.Middlewares[i].apply(wrapped)
	}

	g.Routes = append(g.Routes, NewHandler(matcher, wrapped, name))
}

// - Group() adds a new nested HandlerGroup to HandlerGroup's Routes
// and adds all HandlerGroup's Middlewares to the nested one,
// then returns a pointer to it
//
// - matcher contains logic to define if this nested HandlerGroup should process Update.
//
// - name is primarily used for logging and debugging.
//
// In server's main loop if this new HandlerGroup will match,
// it will start passing processing to each its
// Route in order they were added, until one of Handlers will match
// or none of them (then routing passes to the Route next to HandlerGroup)
//
// # Example:
//
//	grp := g.Group(myMatcherFunc, "myGroup")
//
// # Related:
//
//	//Uses:
//	type MatchFunc func()
//	func NewHandlerGroup()
//
// # See also:
//
//	func NewHandlerGroup() // for more details about name
//
// # Similar:
//
//	func (s *Server) Group()
func (g *HandlerGroup) Group(matcher MatchFunc, name string) *HandlerGroup {
	newGroup := NewHandlerGroup(matcher, name)

	newGroup.Middlewares = make([]Middleware, len(g.Middlewares))
	copy(newGroup.Middlewares, g.Middlewares)

	g.Routes = append(g.Routes, newGroup)
	return newGroup
}

// ===================== UNEXPORTED ======================

// - handle() checks if MatchFunc returns true and starts passing processing to each
// Route in order they were added, until one of Handlers will match
// or none of them (then routing passes to the Route next to HandlerGroup).
//
// returns true if one of the nested Handlers matched, false otherwise
//
// It is used internally in the server's main loop.
//
// # Related:
//
//	//Implements:
//	type Route interface
//
//	//Used in:
//	func (s *Server) handle()
func (g *HandlerGroup) handle(ctx context.Context, u tgbotapi.Update) (matched bool, status HandleStatus, err *BotError) {
	if !g.MatchFunc(ctx, u) {
		return false, StatusWarn, NewBotError(fmt.Sprintf("[Group] %s: unmatched", g.Name), nil)
	}

	ctx = appendPath(ctx, g.Name)

	for _, route := range g.Routes {
		var innerErr *BotError
		matched, status, innerErr = route.handle(ctx, u)
		if innerErr != nil {
			err = NewBotError("", innerErr)
		}
		if matched {
			return true, status, err
		}
	}
	return false, StatusWarn, NewBotError(fmt.Sprintf("[Group] %s: no matched handler", g.Name), err)
}

// ======================== TYPE =========================
// REGION                Middleware
// =======================================================

// - Middleware struct for MiddlewareFunc that runs before and after the handler it wraps.
//
// It is applied recursively to all nested handlers in Server/HandlerGroup by Apply() func,
// or can be applied directly to Handler by Handle() func.
//
// Server -> HandlerGroup (and the nested ones) -> Handler
//
// Middleware is commonly used for code-reuse to reduce boilerplate code
// for logic that should be done for each Update on different levels,
// such as logging, auth, adding common values to context etc.
//
// Middleware works like a wrapper function for a Handler,
// that runs "before-code", then calls next() to pass control
// to the next middleware
// (or the handler if it's the last middleware in a chain).
// When the handler returns, middlewares run "after-code",
// passing control to the previous middleware in the chain,
// returning Handler results (modified by middlewares if needed):
//
// MW1 -> next() -> MW2 -> next() -> Handler -> return -> MW2 -> return -> MW1
//
// Middleware can abort the chain by returning before
// the next() call, preventing Handler execution.
//
// MW1 -> next() -> MW2 -> return -> MW1
//
// # Example:
//
//	mw := NewMiddleware(func(ctx context.Context, update tgbotapi.Update, next HandlerFunc) (status HandleStatus, err *BotError){
//			//Doing something before handling
//
//			//Calling next() to pass control to the next middleware (or Handler)
//			status, err = next(ctx, update)
//
//			//Doing something after handling
//
//			//Returning values we got after handling
//			return status, err // You can use just "return next(ctx, update)"
//		}, "MiddlewareName")
//
// # Related:
//
//	//Uses:
//	type MiddlewareFunc func()
//
//	func (s *Server) Handle()
//	func (g *HandlerGroup) Handle()
//
// # See also:
//
//	package defaults // for pre-made middlewares
//
// # Similar:
//
//	type Handler struct
//	type HandlerGroup struct
//	type Server struct
//	type Route struct
//
// # Methods:
type Middleware struct {
	MiddlewareFunc MiddlewareFunc
	Name           string
}

// ====================== EXPORTED ======================

// - NewMiddleware() returns Middleware instance
//
// - function contains all Middleware's logic.
//
// - name is primarily used for logging and debugging.
//
// If name is empty, it will be generated automatically by autoname().
// autoname is formatted like "M0"
// Where M - Middleware, 0 - autoname's number
//
// You CAN create Middleware using a struct literal directly:
//
//	m := Middleware{MiddlewareFunc: myFunc, Name: "middlewareName"}
//
// However, this is NOT RECOMMENDED because
// if name is empty it won't appear in logs, making debugging harder
//
// Use NewMiddleware() unless you have a specific reason not to
// (e.g., performance-critical code where you're certain name is unnecessary
// and internal calls are a concern).
//
// # Related:
//
//	//Uses:
//	type MiddlewareFunc func()
//
// # Similar:
//
//	func NewHandler()
//	func NewHandlerGroup()
//	func NewServer()
func NewMiddleware(function MiddlewareFunc, name string) Middleware {
	if name == "" {
		name = autoname(typeMiddleware)
	}
	return Middleware{
		MiddlewareFunc: function,
		Name:           name,
	}
}

// ===================== UNEXPORTED ======================

// - apply() wraps a HandlerFunc with the middleware and returns
// a wrapped HandlerFunc
//
// - function is the HandlerFunc to wrap.
//
// It is used internally for Handle() functions when registering a Handler,
// wrapping it with all recursively applied Middlewares
//
// # Related:
//
//	//Uses:
//	type HandlerFunc func()
//	type MiddlewareFunc func()
//
//	//Used in:
//	func (s *Server) Handle()
//	func (g *HandlerGroup) Handle()
//
// # See also:
//
//	type Middleware struct
//	func (s *Server) Apply()
//	func (g *HandlerGroup) Apply()
//	func (h *Handler) Apply()
func (m Middleware) apply(function HandlerFunc) HandlerFunc {
	return func(ctx context.Context, update tgbotapi.Update) (status HandleStatus, err *BotError) {
		return m.MiddlewareFunc(ctx, update, function)
	}
}

// ======================== TYPE =========================
// REGION                  Func()
// =======================================================

// ====================== EXPORTED ======================

// - MatchFunc determines whether a route or group should handle a given update.
type MatchFunc func(ctx context.Context, update tgbotapi.Update) bool

// - HandlerFunc is the function that processes a matched update.
type HandlerFunc func(ctx context.Context, update tgbotapi.Update) (status HandleStatus, err *BotError)

// - HandlerFuncStub is a stub handler for testing.
func HandlerFuncStub(ctx context.Context, update tgbotapi.Update) (status HandleStatus, err *BotError) {
	return StatusWarn, NewBotError("Stub handler matched", nil)
}

// - MiddlewareFunc is a function that wraps a handler.
// It can execute code before and after the handler,
// and can choose to call the next handler or not.
type MiddlewareFunc func(ctx context.Context, update tgbotapi.Update, next HandlerFunc) (status HandleStatus, err *BotError)

// - RoutingFallbackFunc process unhandled Update.
//
// Server calls it when no Handler matched.
type RoutingFallbackFunc func(ctx context.Context, l *Logger, update tgbotapi.Update, status HandleStatus, err *BotError)

type HandleStatus struct {
	code byte
	Name string
}

func NewStatus(code byte, name string) HandleStatus {
	return HandleStatus{code: code, Name: name}
}

var (
	StatusOK       = HandleStatus{0b1000, "OK"} //8
	StatusWarn     = HandleStatus{0b0100, "WARN"} //4
	StatusError    = HandleStatus{0b0010, "ERR"} //2
	StatusFallback = HandleStatus{0b0001, "FALL"} //1
)
