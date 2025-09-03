package hikari

type Middleware func(HandlerFunc) HandlerFunc
