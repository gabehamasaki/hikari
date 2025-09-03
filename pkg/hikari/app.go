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
	router      *Router
	middlewares []Middleware
	server      *http.Server
	logger      *zap.Logger
}

func New(addr string) *App {
	// Create a development logger with pretty colors
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05")
	config.EncoderConfig.EncodeCaller = nil // Remove caller info for cleaner logs

	logger, _ := config.Build()

	return &App{
		addr:        addr,
		router:      NewRouter(),
		middlewares: []Middleware{},
		logger:      logger,
		server: &http.Server{
			Addr:         addr,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		},
	}
}

func (a *App) GET(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	a.router.Handle(http.MethodGet, pattern, handler, middlewares...)
}

func (a *App) POST(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	a.router.Handle(http.MethodPost, pattern, handler, middlewares...)
}

func (a *App) PUT(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	a.router.Handle(http.MethodPut, pattern, handler, middlewares...)
}

func (a *App) PATCH(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	a.router.Handle(http.MethodPatch, pattern, handler, middlewares...)
}

func (a *App) DELETE(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	a.router.Handle(http.MethodDelete, pattern, handler, middlewares...)
}

func (a *App) Use(middleware Middleware) {
	a.middlewares = append(a.middlewares, middleware)
}

func (a *App) buildHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		routerHandler := func(c *Context) {
			a.router.ServeContext(c)
		}

		// Built-in Recovery Middleware (always applied first)
		recoveryMiddleware := func(next HandlerFunc) HandlerFunc {
			return func(c *Context) {
				defer func() {
					if r := recover(); r != nil {
						a.logger.Error("Request panic recovered",
							zap.Any("panic", r),
							zap.String("method", c.Method()),
							zap.String("path", c.Path()),
						)
						http.Error(c.Writer, "Internal Server Error", http.StatusInternalServerError)
					}
				}()
				next(c)
			}
		}

		// Built-in Logger Middleware (applied after recovery)
		loggerMiddleware := func(next HandlerFunc) HandlerFunc {
			return func(c *Context) {
				start := time.Now()

				// Create a contextual logger with request information
				reqLogger := a.logger.With(
					zap.String("method", c.Method()),
					zap.String("path", c.Path()),
					zap.String("remote_addr", c.Request.RemoteAddr),
					zap.String("user_agent", c.Request.Header.Get("User-Agent")),
				)

				// Replace the context logger with the enriched one
				c.Logger = reqLogger

				reqLogger.Info("Request started")
				next(c)

				duration := time.Since(start)
				status := c.GetStatus()

				// Choose log level based on status code
				switch {
				case status >= 500:
					reqLogger.Error("Request completed",
						zap.Int("status", status),
						zap.Duration("duration", duration),
					)
				case status >= 400:
					reqLogger.Warn("Request completed",
						zap.Int("status", status),
						zap.Duration("duration", duration),
					)
				default:
					reqLogger.Info("Request completed",
						zap.Int("status", status),
						zap.Duration("duration", duration),
					)
				}
			}
		}

		// Apply middlewares in reverse order to achieve correct execution order:
		// Execution flow: Recovery -> Logger -> User Middlewares -> Handler
		handler := routerHandler

		// Apply user middlewares first (in reverse order)
		for i := len(a.middlewares) - 1; i >= 0; i-- {
			handler = a.middlewares[i](handler)
		}

		// Apply built-in middlewares (logger wraps user middlewares)
		handler = loggerMiddleware(handler)
		// Recovery wraps everything (outermost layer)
		handler = recoveryMiddleware(handler)

		// Create context and call the handler
		ctx := &Context{
			Writer:  newResponseWriter(w),
			Request: req,
			Params:  make(map[string]string),
			Logger:  a.logger,
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
