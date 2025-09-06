# Chat App - Exemplo WebSocket com Hikari

Este é um exemplo completo de aplicação de chat em tempo real usando WebSocket com o framework Hikari Go.

## ✨ Características Principais

- **💬 Chat Multi-Salas**: Quatro salas diferentes (General, Tecnologia, Aleatório, VIP)
- **🔐 Autenticação**: Sistema de tokens para sala VIP
- **📱 Interface Responsiva**: Design moderno que funciona em desktop e mobile
- **🔔 Notificações**: Sistema de notificações toast em tempo real
- **⌨️ Indicador de Digitação**: Mostra quando alguém está digitando
- **📊 Estatísticas**: Contadores de mensagens e tempo de conexão
- **🛡️ Middleware**: Demonstração de middleware de autenticação
- **🌐 API REST**: Endpoints HTTP para integração

## 🚀 Executando o Exemplo

### Pré-requisitos
- Go 1.19+
- Navegador moderno com suporte a WebSocket

### Passos

1. **Navegar para o diretório**:
   ```bash
   cd examples/chat-app
   ```

2. **Instalar dependências**:
   ```bash
   go mod tidy
   ```

3. **Executar o servidor**:
   ```bash
   go run main.go
   ```

4. **Abrir no navegador**:
   ```
   http://localhost:8080
   ```

## 🏛️ Arquitetura

### WebSocket Hubs
Cada sala de chat é um hub independente que gerencia suas próprias conexões:

```go
app.WebSocket("/ws/general", "general", generalChatHandler)
app.WebSocket("/ws/tech", "tech", techChatHandler)
app.WebSocket("/ws/random", "random", randomChatHandler)
app.WebSocket("/ws/vip", "vip", vipChatHandler, authMiddleware)
```

### Handler Pattern
```go
func generalChatHandler(c *hikari.WSContext) {
    if c.IsTextMessage() {
        var msg ChatMessage
        if err := c.Bind(&msg); err == nil {
            // Processar mensagem
            c.BroadcastJSON(msg)
        }
    }
}
```

## 📡 Protocolo WebSocket

### Tipos de Mensagem

#### 1. Entrada na Sala
```json
{
    "type": "join",
    "username": "João",
    "message": "entrou na sala"
}
```

#### 2. Mensagem de Chat
```json
{
    "type": "message",
    "username": "João",
    "message": "Olá pessoal!",
    "room": "general",
    "timestamp": "2025-09-06T15:30:00Z"
}
```

#### 3. Indicador de Digitação
```json
{
    "type": "typing",
    "username": "João",
    "room": "general"
}
```

#### 4. Saída da Sala
```json
{
    "type": "leave",
    "username": "João",
    "message": "saiu da sala"
}
```

### Respostas do Servidor

#### Confirmação de Entrada
```json
{
    "type": "joined",
    "room": "general",
    "message": "Bem-vindo à sala General"
}
```

#### Notificação de Usuário
```json
{
    "type": "user_joined",
    "username": "João",
    "room": "general",
    "message": "João entrou na sala",
    "timestamp": "2025-09-06T15:30:00Z"
}
```

#### Erro
```json
{
    "type": "error",
    "error": "Username and message are required"
}
```

## 🌐 API REST

### Endpoints Disponíveis

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `GET` | `/api/v1/rooms` | Lista salas disponíveis |
| `GET` | `/api/v1/rooms/{room}/stats` | Estatísticas da sala |
| `POST` | `/api/v1/rooms/{room}/message` | Enviar mensagem via HTTP |

### Exemplos de Uso

#### Listar Salas
```bash
curl http://localhost:8080/api/v1/rooms
```

**Resposta:**
```json
{
    "rooms": ["general", "tech", "random", "vip"]
}
```

#### Estatísticas da Sala
```bash
curl http://localhost:8080/api/v1/rooms/general/stats
```

**Resposta:**
```json
{
    "room": "general",
    "user_count": 3
}
```

#### Enviar Mensagem
```bash
curl -X POST http://localhost:8080/api/v1/rooms/general/message \
  -H "Content-Type: application/json" \
  -d '{
    "username": "Sistema",
    "message": "Mensagem via API"
  }'
```

## 🔐 Sistema VIP

A sala VIP demonstra o uso de middleware de autenticação:

### Acesso
- **Token**: `vip123`
- **URL**: `ws://localhost:8080/ws/vip?token=vip123&username=SeuNome`

### Middleware de Autenticação
```go
func authMiddleware(next hikari.HandlerFunc) hikari.HandlerFunc {
    return func(c *hikari.Context) {
        token := c.Query("token")
        if token != "vip123" {
            c.JSON(401, hikari.H{"error": "Invalid VIP token"})
            return
        }
        c.Set("authenticated", true)
        c.Set("username", c.Query("username"))
        next(c)
    }
}
```

## 🎨 Interface do Usuário

### Funcionalidades da Interface

- **📝 Painel de Login**: Seleção de nome de usuário e sala
- **💬 Interface de Chat**: Mensagens em tempo real, input com contador de caracteres
- **📊 Painel de Estatísticas**: Métricas em tempo real
- **🏠 Lista de Salas**: Visualização de salas e contagem de usuários
- **🔔 Notificações**: Sistema de toast para feedback
- **📱 Design Responsivo**: Otimizado para desktop e mobile

### Tecnologias Frontend

- **HTML5**: Estrutura semântica
- **CSS3**: Design moderno com gradientes e animações
- **JavaScript ES6+**: Cliente WebSocket nativo
- **WebSocket API**: Comunicação em tempo real

## 🛠️ Configuração WebSocket

```go
wsConfig := &hikari.WebSocketConfig{
    ReadBufferSize:    1024,           // Buffer de leitura
    WriteBufferSize:   1024,           // Buffer de escrita
    HandshakeTimeout:  10 * time.Second, // Timeout do handshake
    CheckOrigin:       func(r *http.Request) bool { return true }, // Validação CORS
    EnableCompression: true,           // Compressão habilitada
    PingInterval:      30 * time.Second, // Intervalo de ping
    PongTimeout:       60 * time.Second, // Timeout de pong
}
```

## 📊 Monitoramento e Logs

O servidor produz logs estruturados usando Zap:

```
INFO  Request started  method=GET path=/ws/general remote_addr=127.0.0.1:54321
INFO  WebSocket connection established  conn_id=conn_1725630000123456789_12345 hub=general
INFO  Message received  type=join username=João room=general conn_id=conn_1725630000123456789_12345
INFO  VIP user authenticated  username=Maria
```

### Métricas Disponíveis

- Número de conexões ativas por sala
- Mensagens enviadas/recebidas por usuário
- Tempo de conexão
- Erros de autenticação

## 🧪 Testando a Aplicação

### 1. Interface Web
Acesse `http://localhost:8080` e use a interface gráfica.

### 2. Ferramentas de Linha de Comando

#### wscat (recomendado)
```bash
# Instalar wscat
npm install -g wscat

# Conectar à sala geral
wscat -c ws://localhost:8080/ws/general

# Conectar à sala VIP
wscat -c "ws://localhost:8080/ws/vip?token=vip123&username=TestUser"
```

#### curl (API REST)
```bash
# Testar endpoints HTTP
curl http://localhost:8080/api/v1/rooms
curl http://localhost:8080/api/v1/rooms/general/stats
```

### 3. JavaScript Console
```javascript
const ws = new WebSocket('ws://localhost:8080/ws/general');

ws.onopen = () => {
    console.log('Conectado!');
    ws.send(JSON.stringify({
        type: 'join',
        username: 'TestUser',
        message: 'entrou na sala'
    }));
};

ws.onmessage = (event) => {
    console.log('Mensagem:', JSON.parse(event.data));
};
```

## 🎯 Funcionalidades Demonstradas

### Framework Hikari

1. **WebSocket Support**: Integração nativa com WebSocket
2. **Multiple Hubs**: Gestão de múltiplos hubs independentes
3. **Context Extension**: `WSContext` com métodos específicos para WebSocket
4. **Middleware System**: Middleware funcionando com WebSocket
5. **Graceful Shutdown**: Fechamento adequado de conexões
6. **Error Handling**: Tratamento robusto de erros
7. **Structured Logging**: Logs estruturados com Zap

### WebSocket Features

1. **Broadcasting**: Envio para todas as conexões do hub
2. **Direct Messaging**: Envio para conexões específicas
3. **Message Types**: Suporte para diferentes tipos de mensagem
4. **Connection Management**: Gestão automática do ciclo de vida
5. **Ping/Pong**: Keepalive automático
6. **Binary Support**: Suporte para mensagens binárias

## 🚀 Próximas Funcionalidades

- [ ] **Mensagens Privadas**: Chat direto entre usuários
- [ ] **Histórico**: Persistência de mensagens em banco de dados
- [ ] **Upload de Arquivos**: Compartilhamento de imagens/documentos
- [ ] **Reações**: Sistema de emojis e reações às mensagens
- [ ] **Moderação**: Sistema de moderação de conteúdo
- [ ] **Notificações Push**: Notificações desktop/mobile
- [ ] **Temas**: Suporte para tema escuro/claro
- [ ] **Internacionalização**: Suporte a múltiplos idiomas

## 🤝 Contribuindo

Este exemplo serve como referência para implementar WebSocket no Hikari. Sinta-se à vontade para:

- Reportar bugs
- Sugerir melhorias
- Implementar novas funcionalidades
- Melhorar a documentação

---

**Desenvolvido com ❤️ usando Hikari Go Framework**
