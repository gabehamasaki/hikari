package hikari

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

type WSContext struct {
	*Context
	connection  *WebSocketConnection
	messageType int
	data        []byte
}

// Send envia uma mensagem através desta conexão WebSocket
func (wsc *WSContext) Send(data []byte) {
	wsc.connection.Send(data)
}

// SendJSON envia uma mensagem JSON através desta conexão
func (wsc *WSContext) JSON(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	wsc.Send(data)
	return nil
}

// SendString envia uma string através desta conexão
func (wsc *WSContext) String(message string) {
	wsc.Send([]byte(message))
}

// Broadcast envia mensagem para todas as conexões do hub
func (wsc *WSContext) Broadcast(data []byte) {
	wsc.connection.hub.Broadcast(data)
}

// BroadcastJSON envia mensagem JSON para todas as conexões do hub
func (wsc *WSContext) BroadcastJSON(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	wsc.Broadcast(data)
	return nil
}

// BroadcastString envia string para todas as conexões do hub
func (wsc *WSContext) BroadcastString(message string) {
	wsc.Broadcast([]byte(message))
}

// SendToConnection envia mensagem para uma conexão específica do hub
func (wsc *WSContext) SendToConnection(connID string, data []byte) bool {
	return wsc.connection.hub.SendToConnection(connID, data)
}

// GetConnectionID retorna o ID desta conexão
func (wsc *WSContext) GetConnectionID() string {
	return wsc.connection.id
}

// GetHubName retorna o nome do hub desta conexão
func (wsc *WSContext) GetHubName() string {
	return wsc.connection.hub.name
}

// IsTextMessage verifica se a mensagem é do tipo texto
func (wsc *WSContext) IsTextMessage() bool {
	return wsc.messageType == websocket.TextMessage
}

// IsBinaryMessage verifica se a mensagem é do tipo binário
func (wsc *WSContext) IsBinaryMessage() bool {
	return wsc.messageType == websocket.BinaryMessage
}

// GetMessage retorna a mensagem como string (se for texto)
func (wsc *WSContext) GetMessage() string {
	if wsc.IsTextMessage() {
		return string(wsc.data)
	}
	return ""
}

// BindMessage faz bind da mensagem JSON para uma estrutura
func (wsc *WSContext) Bind(v interface{}) error {
	if !wsc.IsTextMessage() {
		return fmt.Errorf("message is not text type")
	}
	return json.Unmarshal(wsc.data, v)
}
