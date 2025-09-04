package hikari

type Group struct {
	prefix      string
	middlewares []Middleware
	app         *App
}

func (g *Group) GET(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	g.handle("GET", pattern, handler, middlewares...)
}

func (g *Group) POST(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	g.handle("POST", pattern, handler, middlewares...)
}

func (g *Group) PUT(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	g.handle("PUT", pattern, handler, middlewares...)
}

func (g *Group) PATCH(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	g.handle("PATCH", pattern, handler, middlewares...)
}

func (g *Group) DELETE(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	g.handle("DELETE", pattern, handler, middlewares...)
}

func (g *Group) Use(middleware Middleware) {
	g.middlewares = append(g.middlewares, middleware)
}

func (g *Group) Group(prefix string, middlewares ...Middleware) *Group {
	newGroup := &Group{
		prefix:      g.prefix + prefix,
		middlewares: make([]Middleware, len(g.middlewares)+len(middlewares)),
		app:         g.app,
	}

	// Copy existing middlewares
	copy(newGroup.middlewares, g.middlewares)
	copy(newGroup.middlewares[len(g.middlewares):], middlewares)

	return newGroup
}

func (g *Group) handle(method, pattern string, handler HandlerFunc, middlewares ...Middleware) {
	allMiddlewares := make([]Middleware, 0, len(g.middlewares)+len(middlewares))
	copy(allMiddlewares, g.middlewares)
	copy(allMiddlewares[len(g.middlewares):], middlewares)

	fullPattern := g.prefix + pattern
	g.app.router.handle(method, fullPattern, handler, allMiddlewares...)
}
