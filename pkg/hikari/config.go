package hikari

import (
	"net/http"
	"time"
)

type WebSocketConfig struct {
	ReadBufferSize    int
	WriteBufferSize   int
	HandshakeTimeout  time.Duration
	CheckOrigin       func(r *http.Request) bool
	EnableCompression bool
	PingInterval      time.Duration
	PongTimeout       time.Duration
	RegisterTimeout   time.Duration
}

func DefaultWebSocketConfig() *WebSocketConfig {
	return &WebSocketConfig{
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
		HandshakeTimeout:  10 * time.Second,
		CheckOrigin:       func(r *http.Request) bool { return true },
		EnableCompression: true,
		PingInterval:      30 * time.Second,
		PongTimeout:       60 * time.Second,
		RegisterTimeout:   30 * time.Second,
	}
}
