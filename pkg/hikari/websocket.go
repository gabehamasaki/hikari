package hikari

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type WebSocketConnection struct {
	conn    *websocket.Conn
	send    chan []byte
	hub     *WebSocketHub
	id      string
	mu      sync.RWMutex
	closed  bool
	logger  *zap.Logger
	ctx     context.Context
	cancel  context.CancelFunc
	request *http.Request
}

type WebSocketHub struct {
	name        string
	connections map[string]*WebSocketConnection
	register    chan *WebSocketConnection
	unregister  chan *WebSocketConnection
	broadcast   chan []byte
	mu          sync.RWMutex
	logger      *zap.Logger
	ctx         context.Context
	cancel      context.CancelFunc
}

type WebSocketManager struct {
	config   *WebSocketConfig
	upgrader websocket.Upgrader
	hubs     map[string]*WebSocketHub
	mu       sync.RWMutex
	logger   *zap.Logger
}

func NewWebSocketManager(config *WebSocketConfig, logger *zap.Logger) *WebSocketManager {
	if config == nil {
		config = DefaultWebSocketConfig()
	}

	return &WebSocketManager{
		config: config,
		upgrader: websocket.Upgrader{
			ReadBufferSize:    config.ReadBufferSize,
			WriteBufferSize:   config.WriteBufferSize,
			HandshakeTimeout:  config.HandshakeTimeout,
			CheckOrigin:       config.CheckOrigin,
			EnableCompression: config.EnableCompression,
		},
		hubs:   make(map[string]*WebSocketHub),
		logger: logger,
	}
}

func (wm *WebSocketManager) RegisterHub(name string) *WebSocketHub {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	hub := &WebSocketHub{
		name:        name,
		connections: make(map[string]*WebSocketConnection),
		register:    make(chan *WebSocketConnection),
		unregister:  make(chan *WebSocketConnection),
		broadcast:   make(chan []byte),
		logger:      wm.logger.With(zap.String("hub", name)),
		ctx:         ctx,
		cancel:      cancel,
	}

	wm.hubs[name] = hub
	go hub.run()

	return hub
}

func (wm *WebSocketManager) GetHub(name string) (*WebSocketHub, bool) {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	hub, ok := wm.hubs[name]
	return hub, ok
}

func (wm *WebSocketManager) RemoveHub(name string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	if hub, ok := wm.hubs[name]; ok {
		hub.cancel()
		delete(wm.hubs, name)
		wm.logger.Info("WebSocket hub removed", zap.String("hub", name))
	}
}

func (wm *WebSocketManager) Upgrade(c *Context, hubName string, handler WebSocketHandler) error {
	hub, ok := wm.GetHub(hubName)
	if !ok {
		hub = wm.RegisterHub(hubName)
	}

	conn, err := wm.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		wm.logger.Error("WebSocket upgrade failed", zap.Error(err))
		return err
	}

	connId := generateConnectionID()
	wsConn := &WebSocketConnection{
		conn:    conn,
		send:    make(chan []byte, 256),
		hub:     hub,
		id:      connId,
		logger:  wm.logger.With(zap.String("conn_id", connId)),
		request: c.Request,
		ctx:     c,
	}
	select {
	case hub.register <- wsConn:
		wm.logger.Info("WebSocket connection registered", zap.String("conn_id", connId))
	case <-time.After(wm.config.RegisterTimeout):
		conn.Close()
		return fmt.Errorf("failed to register connection: timeout")
	}
	wm.logger.Info("WebSocket connection established", zap.String("conn_id", connId))

	go wsConn.writePump(wm.config)

	wsConn.readPump(wm.config, handler, c)

	wm.logger.Info("WebSocket connection ended", zap.String("conn_id", connId))
	return nil
}

func (h *WebSocketHub) run() {
	defer func() {
		h.logger.Info("WebSocket hub shutting down")
		h.cancel()
		close(h.register)
		close(h.unregister)
		close(h.broadcast)
	}()

	for {
		select {
		case <-h.ctx.Done():
			return

		case conn := <-h.register:
			h.mu.Lock()
			h.connections[conn.id] = conn
			h.mu.Unlock()
			h.logger.Info("WebSocket connection registered", zap.String("conn_id", conn.id))

		case conn := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.connections[conn.id]; ok {
				delete(h.connections, conn.id)
				close(conn.send)
				h.logger.Info("WebSocket connection unregistered", zap.String("conn_id", conn.id))
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for _, conn := range h.connections {
				select {
				case conn.send <- message:
				default:
					close(conn.send)
					delete(h.connections, conn.id)
					h.logger.Warn("WebSocket connection send buffer full, connection closed")
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (c *WebSocketConnection) readPump(config *WebSocketConfig, handler WebSocketHandler, originalContext *Context) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
		c.logger.Info("WebSocket connection closed")
	}()

	c.conn.SetReadDeadline(time.Now().Add(config.PongTimeout))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(config.PongTimeout))
		return nil
	})

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			messageType, message, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					c.logger.Error("WebSocket read error", zap.Error(err))
				}
				return
			}

			if handler != nil {
				wsContext := &WSContext{
					Context:     originalContext,
					connection:  c,
					messageType: messageType,
					data:        message,
				}

				go func() {
					defer func() {
						if r := recover(); r != nil {
							c.logger.Error("WebSocket handler panic",
								zap.Any("panic", r),
								zap.String("conn_id", c.id),
							)
						}
					}()

					handler(wsContext)
				}()
			}
		}
	}
}

func (c *WebSocketConnection) writePump(config *WebSocketConfig) {
	ticker := time.NewTicker(config.PingInterval)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case <-c.ctx.Done():
			return
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				c.logger.Error("WebSocket write error", zap.Error(err))
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.logger.Error("WebSocket ping error", zap.Error(err))
				return
			}
		}
	}
}

func (c *WebSocketConnection) Send(message []byte) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.closed {
		return nil
	}
	select {
	case c.send <- message:
		return nil
	default:
		close(c.send)
		c.closed = true
		c.logger.Warn("WebSocket send buffer full, message dropped", zap.String("hub", c.hub.name))
		return nil
	}
}

func (c *WebSocketConnection) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.closed {
		c.closed = true
		c.cancel()
		close(c.send)
		c.conn.Close()
		c.logger.Info("WebSocket connection closed manually", zap.String("hub", c.hub.name))
	}
}

func (h *WebSocketHub) Broadcast(message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	select {
	case h.broadcast <- message:
	default:
		h.logger.Warn("WebSocket hub broadcast buffer full, message dropped", zap.String("hub", h.name))
	}
}

func (h *WebSocketHub) SendToConnection(connId string, message []byte) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if conn, ok := h.connections[connId]; ok {
		conn.Send(message)
		return true
	}
	return false
}

func (h *WebSocketHub) GetConnectionCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.connections)
}

func generateConnectionID() string {
	return fmt.Sprintf("conn_%d_%d", time.Now().UnixNano(), rand.Int63())
}
