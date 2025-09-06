package hikari

import (
	"time"

	"go.uber.org/zap"
)

// Built-in Logger Middleware (applied after recovery)
func (a *App) loggerMiddleware(next HandlerFunc) HandlerFunc {
	return func(c *Context) {
		start := time.Now()

		// Create a contextual logger with request information
		reqLogger := a.logger.With(
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.String("remote_addr", c.Request.RemoteAddr),
			//zap.String("user_agent", c.Request.Header.Get("User-Agent")),
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
