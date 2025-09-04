package hikari

import (
	"net/http"

	"go.uber.org/zap"
)

// Built-in Recovery Middleware (always applied first)
func (a *App) recoveryMiddleware(next HandlerFunc) HandlerFunc {
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
