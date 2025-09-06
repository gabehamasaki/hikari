.PHONY: help todo-app user-management file-upload run-all clean

# Default target
help:
	@echo "Makefile para executar exemplos do Hikari-Go"
	@echo ""
	@echo "Exemplos disponíveis:"
	@echo "  todo-app         - Executa o exemplo de aplicação TODO"
	@echo "  user-management  - Executa o exemplo de gerenciamento de usuários"
	@echo "  file-upload      - Executa o exemplo de upload de arquivos"
	@echo "  chat-app         - Executa o exemplo de chat WebSocket"
	@echo "  run-all          - Executa todos os exemplos concorrentemente"
	@echo ""
	@echo "Outros comandos:"
	@echo "  clean           - Remove arquivos temporários"
	@echo "  help            - Mostra esta ajuda"

# Exemplos da pasta examples/
todo-app:
	@echo "🚀 Executando exemplo: TODO App"
	@cd ./examples/todo-app && go run .

user-management:
	@echo "🚀 Executando exemplo: User Management"
	@cd ./examples/user-management && go run .

file-upload:
	@echo "🚀 Executando exemplo: File Upload"
	@cd ./examples/file-upload && go run .

chat-app:
	@echo "🚀 Executando exemplo: Chat WebSocket"
	@cd ./examples/chat-app && go run .

# Executa todos os exemplos concorrentemente
run-all:
	@echo "🚀 Executando todos os exemplos concorrentemente..."
	@npx concurrently --kill-others-on-fail --prefix-colors "cyan,magenta,yellow,green" --names "TODO,USER,FILE,CHAT" \
		"cd ./examples/todo-app && go run ." \
		"cd ./examples/user-management && go run ." \
		"cd ./examples/file-upload && go run ." \
		"cd ./examples/chat-app && go run ."

# Comando para limpar arquivos temporários
clean:
	@echo "🧹 Limpando arquivos temporários..."
	@find . -name "*.log" -delete
	@find . -name "*.tmp" -delete
	@echo "✅ Limpeza concluída"

# Comando para executar todos os exemplos (apenas builds para verificar se compilam)
check-all:
	@echo "🔍 Verificando se todos os exemplos compilam..."
	@cd ./examples/todo-app && go build -o /tmp/todo-app .
	@cd ./examples/user-management && go build -o /tmp/user-management .
	@cd ./examples/file-upload && go build -o /tmp/file-upload .
	@cd ./examples/chat-app && go build -o /tmp/chat-app .
	@echo "✅ Todos os exemplos compilam corretamente"
	@rm -f /tmp/todo-app /tmp/user-management /tmp/file-upload /tmp/chat-app
