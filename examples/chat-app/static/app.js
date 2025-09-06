class HikariChat {
    constructor() {
        this.ws = null;
        this.username = '';
        this.currentRoom = '';
        this.isConnected = false;
        this.messagesSent = 0;
        this.messagesReceived = 0;
        this.connectionStartTime = null;
        this.typingTimer = null;
        this.isTyping = false;

        this.initializeElements();
        this.attachEventListeners();
        this.updateRoomStats();

        // Atualizar tempo de conexão a cada segundo
        setInterval(() => this.updateConnectionTime(), 1000);
    }

    initializeElements() {
        // Login elements
        this.usernameInput = document.getElementById('username');
        this.roomSelect = document.getElementById('roomSelect');
        this.vipTokenInput = document.getElementById('vipToken');
        this.tokenGroup = document.getElementById('tokenGroup');
        this.connectBtn = document.getElementById('connectBtn');
        this.loginPanel = document.getElementById('loginPanel');

        // Chat elements
        this.chatPanel = document.getElementById('chatPanel');
        this.messagesContainer = document.getElementById('messages');
        this.messageInput = document.getElementById('messageInput');
        this.sendBtn = document.getElementById('sendBtn');
        this.disconnectBtn = document.getElementById('disconnectBtn');
        this.clearBtn = document.getElementById('clearBtn');
        this.typingIndicator = document.getElementById('typingIndicator');
        this.charCount = document.getElementById('charCount');

        // Status elements
        this.currentRoomSpan = document.getElementById('currentRoom');
        this.userCountSpan = document.getElementById('userCount');
        this.connectionStatus = document.getElementById('connectionStatus');
        this.messagesSentSpan = document.getElementById('messagesSent');
        this.messagesReceivedSpan = document.getElementById('messagesReceived');
        this.connectionTimeSpan = document.getElementById('connectionTime');
        this.wsStatusSpan = document.getElementById('wsStatus');

        // Toast
        this.toast = document.getElementById('toast');
    }

    attachEventListeners() {
        // Login events
        this.roomSelect.addEventListener('change', () => {
            this.tokenGroup.style.display = this.roomSelect.value === 'vip' ? 'block' : 'none';
        });

        this.connectBtn.addEventListener('click', () => this.connect());
        this.disconnectBtn.addEventListener('click', () => this.disconnect());

        // Chat events
        this.sendBtn.addEventListener('click', () => this.sendMessage());
        this.clearBtn.addEventListener('click', () => this.clearMessages());

        this.messageInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault();
                this.sendMessage();
            }
        });

        this.messageInput.addEventListener('input', () => {
            this.updateCharCount();
            this.handleTyping();
        });

        // Enter key no username
        this.usernameInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                this.connect();
            }
        });

        // Room switching
        document.querySelectorAll('.room-item').forEach(item => {
            item.addEventListener('click', () => {
                const room = item.dataset.room;
                if (room !== this.currentRoom && this.isConnected) {
                    this.switchRoom(room);
                }
            });
        });
    }

    connect() {
        const username = this.usernameInput.value.trim();
        const room = this.roomSelect.value;

        if (!username) {
            this.showToast('Por favor, digite um nome de usuário', 'error');
            this.usernameInput.focus();
            return;
        }

        if (username.length < 2) {
            this.showToast('Nome de usuário deve ter pelo menos 2 caracteres', 'error');
            return;
        }

        this.username = username;
        this.currentRoom = room;

        // Construir URL do WebSocket
        let wsUrl = `ws://localhost:8080/ws/${room}`;

        if (room === 'vip') {
            const token = this.vipTokenInput.value.trim();
            if (!token) {
                this.showToast('Token VIP é obrigatório', 'error');
                this.vipTokenInput.focus();
                return;
            }
            wsUrl += `?token=${encodeURIComponent(token)}&username=${encodeURIComponent(username)}`;
        }

        this.connectBtn.disabled = true;
        this.connectBtn.textContent = 'Conectando...';

        try {
            this.ws = new WebSocket(wsUrl);
            this.setupWebSocketEvents();
        } catch (error) {
            this.showToast('Erro ao conectar: ' + error.message, 'error');
            this.connectBtn.disabled = false;
            this.connectBtn.textContent = 'Conectar';
        }
    }

    setupWebSocketEvents() {
        this.ws.onopen = () => {
            this.isConnected = true;
            this.connectionStartTime = new Date();

            this.loginPanel.style.display = 'none';
            this.chatPanel.style.display = 'flex';

            this.updateConnectionStatus(true);
            this.updateRoomDisplay();
            this.showToast(`Conectado à sala ${this.getRoomName(this.currentRoom)}`, 'success');

            // Enviar mensagem de entrada
            this.sendWebSocketMessage({
                type: 'join',
                username: this.username,
                message: 'entrou na sala'
            });

            this.messageInput.focus();
            this.updateStats();
        };

        this.ws.onmessage = (event) => {
            try {
                const message = JSON.parse(event.data);
                this.handleWebSocketMessage(message);
            } catch (error) {
                console.error('Erro ao processar mensagem:', error);
            }
        };

        this.ws.onerror = (error) => {
            console.error('WebSocket error:', error);
            this.showToast('Erro na conexão WebSocket', 'error');
        };

        this.ws.onclose = (event) => {
            this.isConnected = false;
            this.updateConnectionStatus(false);

            if (event.code === 1006) {
                this.showToast('Conexão perdida inesperadamente', 'error');
            } else if (event.code === 1000) {
                this.showToast('Desconectado com sucesso', 'info');
            } else {
                this.showToast(`Conexão fechada (código: ${event.code})`, 'warning');
            }

            // Reset UI se não foi desconexão intencional
            if (event.code !== 1000) {
                setTimeout(() => {
                    this.resetToLogin();
                }, 2000);
            }
        };
    }

    handleWebSocketMessage(message) {
        this.messagesReceived++;

        switch (message.type) {
            case 'message':
                this.displayChatMessage(message);
                break;
            case 'user_joined':
                this.displaySystemMessage(message.message, 'join');
                break;
            case 'user_left':
                this.displaySystemMessage(message.message, 'leave');
                break;
            case 'joined':
                this.displaySystemMessage(message.message, 'success');
                break;
            case 'typing':
                this.displayTypingIndicator(message);
                break;
            case 'error':
                this.displaySystemMessage(message.error, 'error');
                break;
            default:
                console.log('Mensagem desconhecida:', message);
        }

        this.updateStats();
    }

    sendMessage() {
        const message = this.messageInput.value.trim();

        if (!message) return;
        if (!this.isConnected) {
            this.showToast('Não conectado ao chat', 'error');
            return;
        }

        this.sendWebSocketMessage({
            type: 'message',
            username: this.username,
            message: message
        });

        this.messageInput.value = '';
        this.updateCharCount();
        this.messagesSent++;
        this.updateStats();
    }

    sendWebSocketMessage(message) {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify(message));
        }
    }

    displayChatMessage(message) {
        const messageEl = document.createElement('div');
        messageEl.className = `message ${message.username === this.username ? 'own' : 'other'}`;

        const time = new Date(message.timestamp).toLocaleTimeString('pt-BR', {
            hour: '2-digit',
            minute: '2-digit'
        });

        messageEl.innerHTML = `
            <div class="message-header">
                <strong>${this.escapeHtml(message.username)}</strong>
            </div>
            <div class="message-content">${this.escapeHtml(message.message)}</div>
            <div class="message-time">${time}</div>
        `;

        this.messagesContainer.appendChild(messageEl);
        this.scrollToBottom();
    }

    displaySystemMessage(message, type = 'system') {
        const messageEl = document.createElement('div');
        messageEl.className = `message ${type}`;
        messageEl.innerHTML = `<div class="message-content">${this.escapeHtml(message)}</div>`;

        this.messagesContainer.appendChild(messageEl);
        this.scrollToBottom();
    }

    displayTypingIndicator(message) {
        if (message.username === this.username) return;

        this.typingIndicator.textContent = `${message.username} está digitando...`;

        clearTimeout(this.typingTimer);
        this.typingTimer = setTimeout(() => {
            this.typingIndicator.textContent = '';
        }, 3000);
    }

    handleTyping() {
        if (!this.isConnected || this.isTyping) return;

        this.isTyping = true;
        this.sendWebSocketMessage({
            type: 'typing',
            username: this.username
        });

        setTimeout(() => {
            this.isTyping = false;
        }, 2000);
    }

    disconnect() {
        if (this.ws) {
            this.sendWebSocketMessage({
                type: 'leave',
                username: this.username,
                message: 'saiu da sala'
            });

            this.ws.close(1000);
        }

        this.resetToLogin();
    }

    resetToLogin() {
        this.isConnected = false;
        this.ws = null;
        this.connectionStartTime = null;

        this.chatPanel.style.display = 'none';
        this.loginPanel.style.display = 'block';

        this.connectBtn.disabled = false;
        this.connectBtn.textContent = 'Conectar';

        this.messagesContainer.innerHTML = '';
        this.messageInput.value = '';
        this.typingIndicator.textContent = '';

        this.updateConnectionStatus(false);
        this.updateStats();

        // Remover classe active das salas
        document.querySelectorAll('.room-item').forEach(item => {
            item.classList.remove('active');
        });
    }

    switchRoom(newRoom) {
        if (newRoom === this.currentRoom) return;

        // Implementação futura para trocar de sala sem desconectar
        this.showToast('Funcionalidade de troca de sala em desenvolvimento', 'info');
    }

    clearMessages() {
        this.messagesContainer.innerHTML = '';
        this.showToast('Mensagens limpas', 'info');
    }

    updateCharCount() {
        const length = this.messageInput.value.length;
        this.charCount.textContent = `${length}/1000`;

        if (length > 900) {
            this.charCount.className = 'char-count danger';
        } else if (length > 800) {
            this.charCount.className = 'char-count warning';
        } else {
            this.charCount.className = 'char-count';
        }
    }

    updateConnectionStatus(connected) {
        this.connectionStatus.textContent = connected ? '🟢 Conectado' : '🔴 Desconectado';
        this.connectionStatus.className = `status ${connected ? 'connected' : 'disconnected'}`;
        this.wsStatusSpan.textContent = connected ? 'Conectado' : 'Desconectado';
    }

    updateRoomDisplay() {
        const roomName = this.getRoomName(this.currentRoom);
        this.currentRoomSpan.textContent = `Sala: ${roomName}`;

        // Atualizar classe active nas salas
        document.querySelectorAll('.room-item').forEach(item => {
            item.classList.toggle('active', item.dataset.room === this.currentRoom);
        });
    }

    updateConnectionTime() {
        if (this.connectionStartTime && this.isConnected) {
            const diff = new Date() - this.connectionStartTime;
            const minutes = Math.floor(diff / 60000);
            const seconds = Math.floor((diff % 60000) / 1000);
            this.connectionTimeSpan.textContent = `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
        } else {
            this.connectionTimeSpan.textContent = '00:00';
        }
    }

    updateStats() {
        this.messagesSentSpan.textContent = this.messagesSent;
        this.messagesReceivedSpan.textContent = this.messagesReceived;
    }

    async updateRoomStats() {
        try {
            const rooms = ['general', 'tech', 'random', 'vip'];

            for (const room of rooms) {
                try {
                    const response = await fetch(`/api/v1/rooms/${room}/stats`);
                    if (response.ok) {
                        const stats = await response.json();
                        const roomEl = document.querySelector(`[data-room="${room}"] .room-users`);
                        if (roomEl) {
                            roomEl.textContent = `${stats.user_count} usuários`;
                        }
                    }
                } catch (error) {
                    console.error(`Erro ao obter stats da sala ${room}:`, error);
                }
            }
        } catch (error) {
            console.error('Erro ao atualizar estatísticas das salas:', error);
        }

        // Atualizar a cada 10 segundos
        setTimeout(() => this.updateRoomStats(), 10000);
    }

    getRoomName(room) {
        const roomNames = {
            'general': '💬 General',
            'tech': '💻 Tecnologia',
            'random': '🎲 Aleatório',
            'vip': '👑 VIP'
        };
        return roomNames[room] || room;
    }

    scrollToBottom() {
        this.messagesContainer.scrollTop = this.messagesContainer.scrollHeight;
    }

    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    showToast(message, type = 'info') {
        this.toast.textContent = message;
        this.toast.className = `toast ${type} show`;

        setTimeout(() => {
            this.toast.classList.remove('show');
        }, 4000);
    }
}

// Inicializar o chat quando a página carregar
document.addEventListener('DOMContentLoaded', () => {
    new HikariChat();
});
