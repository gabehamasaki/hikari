# Chat App - Exemplo WebSocket

Este é um exemplo completo de aplicação de chat em tempo real usando WebSocket com o framework Hikari.

## ✨ Funcionalidades

- **Chat Multi-Salas**: Suporte para múltiplas salas de chat (General, Tech, Random, VIP)
- **Autenticação VIP**: Sala VIP com sistema de tokens
- **Interface Responsiva**: Interface moderna e responsiva
- **Notificações em Tempo Real**: Sistema de notificações toast
- **Indicador de Digitação**: Mostra quando alguém está digitando
- **Estatísticas**: Contadores de mensagens e tempo de conexão
- **Middleware Support**: Demonstração de middleware de autenticação

## 🚀 Como Executar

1. **Instalar Dependências**:
   ```bash
   cd examples/chat-app
   go mod tidy
   ```

2. **Executar o Servidor**:
   ```bash
   go run main.go
   ```

3. **Abrir no Navegador**:
   ```
   http://localhost:8080
   ```

## 🏗️ Estrutura do Projeto

```
chat-app/
├── main.go              # Servidor principal com WebSocket
├── static/
│   ├── index.html       # Interface do chat
│   ├── style.css        # Estilos CSS
│   └── app.js           # JavaScript client-side
├── requests/
│   └── test-requests.http # Exemplos de requisições HTTP
└── README.md            # Esta documentação
```

## 📡 WebSocket Endpoints

### Chat Rooms

| Endpoint | Hub | Descrição |
|----------|-----|-----------|
| `ws://localhost:8080/ws/general` | `general` | Chat geral |
| `ws://localhost:8080/ws/tech` | `tech` | Chat de tecnologia |
| `ws://localhost:8080/ws/random` | `random` | Chat aleatório |
| `ws://localhost:8080/ws/vip` | `vip` | Chat VIP (requer token) |

### Mensagens WebSocket

#### Entrada na Sala
```json
{
    "type": "join",
    "username": "João",
    "message": "entrou na sala"
}
```

#### Mensagem de Chat
```json
{
    "type": "message",
    "username": "João",
    "message": "Olá pessoal!"
}
```

#### Indicador de Digitação
```json
{
    "type": "typing",
    "username": "João"
}
```

#### Saída da Sala
```json
{
    "type": "leave",
    "username": "João",
    "message": "saiu da sala"
}
```

## 🌐 API Endpoints

### GET `/api/v1/rooms`
Lista todas as salas disponíveis.

**Resposta:**
```json
{
    "rooms": ["general", "tech", "random", "vip"]
}
```

### GET `/api/v1/rooms/{room}/stats`
Estatísticas de uma sala específica.

**Resposta:**
```json
{
    "room": "general",
    "user_count": 5,
    "messages_sent": 0
}
```

### POST `/api/v1/rooms/{room}/message`
Envia mensagem para uma sala via HTTP.

**Body:**
```json
{
    "username": "Sistema",
    "message": "Mensagem via API"
}
```

## 🔐 Autenticação VIP

Para acessar a sala VIP, use:
- **Token**: `vip123`
- **URL**: `ws://localhost:8080/ws/vip?token=vip123&username=SeuNome`

O middleware de autenticação valida o token e permite acesso à sala exclusiva.

## 💡 Exemplos de Uso

### Cliente JavaScript Básico

```javascript
const ws = new WebSocket('ws://localhost:8080/ws/general');

ws.onopen = () => {
    ws.send(JSON.stringify({
        type: 'join',
        username: 'TestUser',
        message: 'entrou na sala'
    }));
};

ws.onmessage = (event) => {
    const message = JSON.parse(event.data);
    console.log('Mensagem recebida:', message);
};

// Enviar mensagem
ws.send(JSON.stringify({
    type: 'message',
    username: 'TestUser',
    message: 'Olá chat!'
}));
```

### Exemplo com cURL (HTTP API)

```bash
# Obter estatísticas de uma sala
curl http://localhost:8080/api/v1/rooms/general/stats

# Enviar mensagem via HTTP
curl -X POST http://localhost:8080/api/v1/rooms/general/message \
  -H "Content-Type: application/json" \
  -d '{
    "username": "Sistema",
    "message": "Mensagem via API HTTP"
  }'
```

## 🎨 Interface

A aplicação possui uma interface moderna com:

- **Painel de Login**: Seleção de nome e sala
- **Chat Interface**: Mensagens, input, indicadores
- **Sidebar**: Estatísticas e lista de salas
- **Notificações**: Sistema de toast para feedback
- **Design Responsivo**: Funciona em desktop e mobile

## 🔧 Configurações WebSocket

```go
wsConfig := &hikari.WebSocketConfig{
    ReadBufferSize:    1024,
    WriteBufferSize:   1024,
    HandshakeTimeout:  10 * time.Second,
    CheckOrigin:       func(r *http.Request) bool { return true },
    EnableCompression: true,
    PingInterval:      30 * time.Second,
    PongTimeout:       60 * time.Second,
}
```

## 📝 Logs

O servidor produz logs estruturados para:

- Conexões WebSocket estabelecidas/fechadas
- Mensagens recebidas/enviadas
- Erros de autenticação
- Estatísticas das salas
- Requests HTTP

## 🌟 Funcionalidades Avançadas

- **Multiple Hubs**: Cada sala opera em um hub independente
- **Broadcast**: Mensagens enviadas para todos os usuários da sala
- **Direct Messages**: Suporte para mensagens diretas (futuro)
- **Graceful Shutdown**: Fechamento adequado das conexões
- **Error Handling**: Tratamento robusto de erros
- **Middleware**: Sistema de middleware para autenticação

## 🚀 Próximas Funcionalidades

- [ ] Mensagens diretas entre usuários
- [ ] Histórico de mensagens
- [ ] Upload de arquivos
- [ ] Emojis e reações
- [ ] Notificações desktop
- [ ] Temas escuro/claro
