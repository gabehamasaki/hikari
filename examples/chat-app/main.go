package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gabehamasaki/hikari-go/pkg/hikari"
	"go.uber.org/zap"
)

// ChatMessage representa uma mensagem de chat
type ChatMessage struct {
	Type      string    `json:"type"`
	Username  string    `json:"username"`
	Message   string    `json:"message"`
	Room      string    `json:"room"`
	Timestamp time.Time `json:"timestamp"`
	ConnID    string    `json:"conn_id,omitempty"`
}

// UserJoinLeave representa eventos de entrada/saída
type UserJoinLeave struct {
	Type      string    `json:"type"`
	Username  string    `json:"username"`
	Room      string    `json:"room"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// RoomStats representa estatísticas de uma sala
type RoomStats struct {
	Room        string `json:"room"`
	UserCount   int    `json:"user_count"`
	MessageSent int    `json:"messages_sent"`
}

func main() {
	app := hikari.New(":8080")

	// Configurar WebSocket com configurações customizadas
	wsConfig := &hikari.WebSocketConfig{
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
		HandshakeTimeout:  10 * time.Second,
		CheckOrigin:       func(r *http.Request) bool { return true },
		EnableCompression: true,
		PingInterval:      30 * time.Second,
		PongTimeout:       60 * time.Second,
	}
	app.WithWebSocket(wsConfig)

	// WebSocket Routes
	// Chat geral
	app.WebSocket("/ws/general", "general", generalChatHandler)

	// Sala de tecnologia
	app.WebSocket("/ws/tech", "tech", techChatHandler)

	// Sala aleatória
	app.WebSocket("/ws/random", "random", randomChatHandler)

	// Sala VIP (com middleware de autenticação simples)
	app.WebSocket("/ws/vip", "vip", vipChatHandler, authMiddleware)

	// Servir arquivos estáticos
	app.GET("/", func(c *hikari.Context) {
		c.File("./static/index.html")
	})

	app.GET("/static/*", func(c *hikari.Context) {
		filePath := c.Wildcard()
		var contentType string
		switch {
		case filePath[len(filePath)-4:] == ".css":
			contentType = "text/css"
		case filePath[len(filePath)-3:] == ".js":
			contentType = "application/javascript"
		case filePath[len(filePath)-5:] == ".html":
			contentType = "text/html"
		default:
			contentType = "application/octet-stream"
		}
		c.SetHeader("Content-Type", contentType)
		c.File("./static/" + filePath)
	})

	// API Routes
	apiGroup := app.Group("/api/v1")
	{
		// Endpoint para listar salas disponíveis
		apiGroup.GET("/rooms", func(c *hikari.Context) {
			rooms := []string{"general", "tech", "random", "vip"}
			c.JSON(http.StatusOK, hikari.H{"rooms": rooms})
		})

		// Endpoint para estatísticas de uma sala
		apiGroup.GET("/rooms/:room/stats", func(c *hikari.Context) {
			roomName := c.Param("room")

			if hub, exists := app.GetWebSocketHub(roomName); exists {
				stats := RoomStats{
					Room:      roomName,
					UserCount: hub.GetConnectionCount(),
				}
				c.JSON(http.StatusOK, stats)
			} else {
				c.JSON(http.StatusNotFound, hikari.H{"error": "Room not found"})
			}
		})

		// Endpoint para enviar mensagem para uma sala específica (via HTTP)
		apiGroup.POST("/rooms/:room/message", func(c *hikari.Context) {
			roomName := c.Param("room")

			var msg ChatMessage
			if err := c.Bind(&msg); err != nil {
				c.JSON(http.StatusBadRequest, hikari.H{"error": "Invalid message format"})
				return
			}

			msg.Type = "message"
			msg.Room = roomName
			msg.Timestamp = time.Now()

			if hub, exists := app.GetWebSocketHub(roomName); exists {
				data, _ := json.Marshal(msg)
				hub.Broadcast(data)
				c.JSON(http.StatusOK, hikari.H{"status": "sent"})
			} else {
				c.JSON(http.StatusNotFound, hikari.H{"error": "Room not found"})
			}
		})
	}

	app.ListenAndServe()
}

// Handler para chat geral
func generalChatHandler(c *hikari.WSContext) {
	handleChatMessage(c, "general")
}

// Handler para chat de tecnologia
func techChatHandler(c *hikari.WSContext) {
	handleChatMessage(c, "tech")
}

// Handler para chat aleatório
func randomChatHandler(c *hikari.WSContext) {
	handleChatMessage(c, "random")
}

// Handler para chat VIP
func vipChatHandler(c *hikari.WSContext) {
	// Verificar se o usuário está autenticado
	if !c.GetBool("authenticated") {
		c.JSON(hikari.H{
			"type":  "error",
			"error": "Authentication required for VIP room",
		})
		return
	}

	username := c.GetString("username")
	c.Logger.Info("VIP user connected",
		zap.String("username", username),
		zap.String("conn_id", c.GetConnectionID()),
	)

	handleChatMessage(c, "vip")
}

// Handler genérico para mensagens de chat
func handleChatMessage(c *hikari.WSContext, roomName string) {
	if c.IsTextMessage() {
		var msg ChatMessage
		if err := c.Bind(&msg); err != nil {
			c.Logger.Error("Failed to parse message", zap.Error(err))
			c.JSON(hikari.H{
				"type":  "error",
				"error": "Invalid message format",
			})
			return
		}

		// Definir campos da mensagem
		msg.Room = roomName
		msg.Timestamp = time.Now()
		msg.ConnID = c.GetConnectionID()

		// Log da mensagem
		c.Logger.Info("Message received",
			zap.String("type", msg.Type),
			zap.String("username", msg.Username),
			zap.String("room", msg.Room),
			zap.String("conn_id", c.GetConnectionID()),
		)

		// Processar diferentes tipos de mensagem
		switch msg.Type {
		case "join":
			handleUserJoin(c, msg)
		case "leave":
			handleUserLeave(c, msg)
		case "message":
			handleMessage(c, msg)
		case "typing":
			handleTyping(c, msg)
		default:
			c.JSON(hikari.H{
				"type":  "error",
				"error": "Unknown message type",
			})
		}
	} else if c.IsBinaryMessage() {
		// Para mensagens binárias (ex: compartilhamento de arquivos)
		message := c.GetMessage() // Para texto, ou podemos usar um método para dados binários
		c.Logger.Info("Binary message received",
			zap.String("conn_id", c.GetConnectionID()),
			zap.Int("size", len(message)),
		)

		// Echo da mensagem binária para todos
		c.Broadcast([]byte(message))
	}
}

// Processar entrada de usuário
func handleUserJoin(c *hikari.WSContext, msg ChatMessage) {
	joinMsg := UserJoinLeave{
		Type:      "user_joined",
		Username:  msg.Username,
		Room:      msg.Room,
		Message:   msg.Username + " entrou na sala",
		Timestamp: time.Now(),
	}

	// Broadcast para todos na sala
	if err := c.BroadcastJSON(joinMsg); err != nil {
		c.Logger.Error("Failed to broadcast join message", zap.Error(err))
	}

	// Confirmar entrada para o usuário
	c.JSON(hikari.H{
		"type":    "joined",
		"room":    msg.Room,
		"message": "Bem-vindo à sala " + msg.Room,
	})
}

// Processar saída de usuário
func handleUserLeave(c *hikari.WSContext, msg ChatMessage) {
	leaveMsg := UserJoinLeave{
		Type:      "user_left",
		Username:  msg.Username,
		Room:      msg.Room,
		Message:   msg.Username + " saiu da sala",
		Timestamp: time.Now(),
	}

	// Broadcast para todos na sala
	if err := c.BroadcastJSON(leaveMsg); err != nil {
		c.Logger.Error("Failed to broadcast leave message", zap.Error(err))
	}
}

// Processar mensagem normal
func handleMessage(c *hikari.WSContext, msg ChatMessage) {
	// Validar mensagem
	if msg.Username == "" || msg.Message == "" {
		c.JSON(hikari.H{
			"type":  "error",
			"error": "Username and message are required",
		})
		return
	}

	// Limitar tamanho da mensagem
	if len(msg.Message) > 1000 {
		c.JSON(hikari.H{
			"type":  "error",
			"error": "Message too long (max 1000 characters)",
		})
		return
	}

	// Broadcast da mensagem para todos na sala
	if err := c.BroadcastJSON(msg); err != nil {
		c.Logger.Error("Failed to broadcast message", zap.Error(err))
		c.JSON(hikari.H{
			"type":  "error",
			"error": "Failed to send message",
		})
	}
}

// Processar indicador de digitação
func handleTyping(c *hikari.WSContext, msg ChatMessage) {
	typingMsg := hikari.H{
		"type":     "typing",
		"username": msg.Username,
		"room":     msg.Room,
		"conn_id":  c.GetConnectionID(),
	}

	// Broadcast para todos exceto o remetente
	if err := c.BroadcastJSON(typingMsg); err != nil {
		c.Logger.Error("Failed to broadcast typing indicator", zap.Error(err))
	}
}

// Middleware de autenticação simples para sala VIP
func authMiddleware(next hikari.HandlerFunc) hikari.HandlerFunc {
	return func(c *hikari.Context) {
		token := c.Query("token")

		if token == "" {
			c.JSON(http.StatusUnauthorized, hikari.H{"error": "Token required for VIP room"})
			return
		}

		// Validação simples do token (em produção, use JWT ou similar)
		if token != "vip123" {
			c.JSON(http.StatusUnauthorized, hikari.H{"error": "Invalid VIP token"})
			return
		}

		// Extrair username do token (simulado)
		username := c.Query("username")
		if username == "" {
			username = "VIP_User"
		}

		// Salvar informações no contexto
		c.Set("authenticated", true)
		c.Set("username", username)
		c.Set("user_level", "vip")

		c.Logger.Info("VIP user authenticated",
			zap.String("username", username),
		)

		next(c)
	}
}
