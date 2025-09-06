package hikari

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type App struct {
	addr        string
	router      *router
	middlewares []Middleware
	server      *http.Server
	logger      *zap.Logger

	requestTimeout time.Duration

	wsManager *WebSocketManager
}

func New(addr string) *App {
	// Create a development logger with pretty colors
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05")
	config.EncoderConfig.EncodeCaller = nil // Remove caller info for cleaner logs

	logger, _ := config.Build()

	return &App{
		addr:           addr,
		router:         newRouter(logger),
		middlewares:    []Middleware{},
		logger:         logger,
		requestTimeout: 30 * time.Second, // Default request timeout
		server: &http.Server{
			Addr:         addr,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		},
	}
}

func (a *App) buildHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		isWebSocket := a.wsManager != nil && req.Header.Get("Upgrade") == "websocket" && req.Header.Get("Connection") == "Upgrade"

		var reqCtx context.Context
		var cancel context.CancelFunc

		if isWebSocket {
			// For WebSocket requests, use the request's context without timeout
			reqCtx = req.Context()
			cancel = func() {}
		} else {
			// For regular HTTP requests, apply the request timeout
			reqCtx, cancel = context.WithTimeout(req.Context(), a.requestTimeout)
		}
		defer cancel()

		routerHandler := func(c *Context) {
			a.router.serveContext(c)
		}

		// Apply middlewares in reverse order to achieve correct execution order:
		// Execution flow: Recovery -> Logger -> User Middlewares -> Handler
		handler := routerHandler

		// Apply user middlewares first (in reverse order)
		for i := len(a.middlewares) - 1; i >= 0; i-- {
			handler = a.middlewares[i](handler)
		}

		// Apply built-in middlewares (logger wraps user middlewares)
		handler = a.loggerMiddleware(handler)
		// Recovery wraps everything (outermost layer)
		handler = a.recoveryMiddleware(handler)

		// Create context and call the handler
		ctx := &Context{
			Writer:  newResponseWriter(w),
			Request: req,
			Params:  make(map[string]string),
			Logger:  a.logger,

			Context: reqCtx,
			storage: make(map[string]interface{}),
		}

		handler(ctx)
	})
}

func (a *App) ListenAndServe() {
	// Set the handler with middlewares applied
	a.server.Handler = a.buildHandler()

	// Log server startup
	a.logger.Info("Starting HTTP server",
		zap.String("address", a.addr),
	)

	for _, route := range a.router.routes {
		a.logger.Info("HTTP route configured",
			zap.String("method", route.method),
			zap.String("pattern", route.pattern),
		)
	}

	// Listing all configured WebSocket and HTTP routes
	if a.wsManager != nil {
		for hubName := range a.wsManager.hubs {
			a.logger.Info("WebSocket Hub configured",
				zap.String("hub", hubName),
			)
		}
	}

	// Channel to receive server errors
	serverErr := make(chan error, 1)

	// Start the server in a goroutine
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// Channel to listen for interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Wait for either an error or interrupt signal
	select {
	case err := <-serverErr:
		a.logger.Error("Server error", zap.Error(err))
		panic(err)
	case <-stop:
		a.logger.Info("Shutdown signal received, gracefully stopping server...")
	}

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.Shutdown(ctx); err != nil {
		a.logger.Error("Server shutdown error", zap.Error(err))
		panic(err)
	}

	a.logger.Info("Server stopped gracefully")
}

func (a *App) Shutdown(ctx context.Context) error {
	defer a.logger.Sync() // Flush any remaining logs
	return a.server.Shutdown(ctx)
}

func (a *App) SetRequestTimeout(d time.Duration) {
	a.requestTimeout = d
}

func (a *App) WithWebSocket(config *WebSocketConfig) {
	a.wsManager = NewWebSocketManager(config, a.logger)
}

func (a *App) GET(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	a.router.handle(http.MethodGet, pattern, handler, middlewares...)
}

func (a *App) POST(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	a.router.handle(http.MethodPost, pattern, handler, middlewares...)
}

func (a *App) PUT(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	a.router.handle(http.MethodPut, pattern, handler, middlewares...)
}

func (a *App) PATCH(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	a.router.handle(http.MethodPatch, pattern, handler, middlewares...)
}

func (a *App) DELETE(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	a.router.handle(http.MethodDelete, pattern, handler, middlewares...)
}

func (a *App) WebSocket(path, hubName string, handler WebSocketHandler, middlewares ...Middleware) {
	if a.wsManager != nil {
		a.wsManager.RegisterHub(hubName)
	}

	wsHandler := func(c *Context) {
		if a.wsManager == nil {
			a.logger.Error("WebSocket manager not initialized. Call WithWebSocket() before using WebSocket routes.")
			c.JSON(http.StatusInternalServerError, map[string]string{"error": "WebSocket not configured"})
			return
		}

		err := a.wsManager.Upgrade(c, hubName, handler)
		if err != nil {
			a.logger.Error("WebSocket upgrade failed", zap.Error(err))
			c.JSON(http.StatusInternalServerError, map[string]string{"error": "WebSocket upgrade failed"})
			return
		}
	}

	a.GET(path, wsHandler, middlewares...)
}

func (a *App) GetWebSocketHub(name string) (*WebSocketHub, bool) {
	if a.wsManager == nil {
		a.logger.Error("WebSocket manager not initialized. Call WithWebSocket() before using WebSocket hubs.")
		return nil, false
	}
	return a.wsManager.GetHub(name)
}

func (a *App) Group(prefix string, middlewares ...Middleware) *Group {
	return &Group{
		prefix:      prefix,
		middlewares: middlewares,
		app:         a,
	}
}

func (a *App) Use(middleware Middleware) {
	a.middlewares = append(a.middlewares, middleware)
}
