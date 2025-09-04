package hikari

import (
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type route struct {
	method      string
	pattern     string
	handler     HandlerFunc
	middlewares []Middleware
}

type router struct {
	routes []route
	logger *zap.Logger
}

func newRouter(logger *zap.Logger) *router {
	return &router{
		routes: []route{},
		logger: logger,
	}
}

func (r *router) handle(method, pattern string, handler HandlerFunc, middlewares ...Middleware) {
	normalizedPattern := buildPattern("", pattern, r.logger)

	r.routes = append(r.routes, route{
		method:      method,
		pattern:     normalizedPattern,
		handler:     handler,
		middlewares: middlewares,
	})
}

func (r *router) handleNormalized(method, pattern string, handler HandlerFunc, middlewares ...Middleware) {
	r.routes = append(r.routes, route{
		method:      method,
		pattern:     pattern,
		handler:     handler,
		middlewares: middlewares,
	})
}

func splitPath(p string) []string {
	p = strings.Trim(p, "/")
	if p == "" {
		return []string{}
	}

	return strings.Split(p, "/")
}

func (r *router) serveContext(ctx *Context) {
	for _, rt := range r.routes {
		if rt.method != ctx.Request.Method {
			continue
		}

		requestPath := normalizedPattern(ctx.Request.URL.Path)

		pParts := splitPath(rt.pattern)
		rParts := splitPath(requestPath)

		// Check for wildcard pattern (*)
		hasWildcard := len(pParts) > 0 && pParts[len(pParts)-1] == "*"

		// For wildcard routes, we need at least as many parts as the pattern (minus the wildcard)
		if hasWildcard {
			if len(rParts) < len(pParts)-1 {
				continue
			}
		} else {
			// For non-wildcard routes, parts must match exactly
			if len(pParts) != len(rParts) {
				continue
			}
		}

		params := map[string]string{}
		matched := true

		// Check parts up to wildcard (if any)
		partsToCheck := len(pParts)
		if hasWildcard {
			partsToCheck = len(pParts) - 1
		}

		for i := 0; i < partsToCheck; i++ {
			if strings.HasPrefix(pParts[i], ":") {
				params[strings.TrimPrefix(pParts[i], ":")] = rParts[i]
			} else if pParts[i] != rParts[i] {
				matched = false
				break
			}
		}

		// If we have a wildcard, capture the remaining path
		if matched && hasWildcard && len(rParts) > partsToCheck {
			remainingParts := rParts[partsToCheck:]
			params["*"] = strings.Join(remainingParts, "/")
		}

		if matched {
			// Update the existing context with route parameters
			ctx.Params = params

			handler := rt.handler
			// Apply user middlewares first (in reverse order)
			for i := len(rt.middlewares) - 1; i >= 0; i-- {
				handler = rt.middlewares[i](handler)
			}

			handler(ctx)

			return
		}
	}
	http.NotFound(ctx.Writer, ctx.Request)
}
