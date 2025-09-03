# Todo App Example

Uma API REST simples para gerenciamento de tarefas usando Hikari-Go.

## Características

- CRUD completo para tarefas
- Filtragem por status (completed/pending)
- Middleware de CORS personalizado
- Validação de dados
- Estrutura JSON organizada

## Como executar

```bash
cd examples/todo-app
go run main.go
```

O servidor será iniciado em `http://localhost:8080`

## Endpoints

### GET /
Retorna informações sobre a API e lista de endpoints disponíveis.

### GET /todos
Lista todas as tarefas.

**Query Parameters:**
- `status`: `completed` ou `pending` para filtrar tarefas

**Exemplo:**
```bash
curl http://localhost:8080/todos
curl http://localhost:8080/todos?status=completed
```

### GET /todos/:id
Retorna uma tarefa específica por ID.

**Exemplo:**
```bash
curl http://localhost:8080/todos/1
```

### POST /todos
Cria uma nova tarefa.

**Body:**
```json
{
  "title": "Nova tarefa",
  "content": "Descrição da tarefa"
}
```

**Exemplo:**
```bash
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"Estudar Go","content":"Aprender sobre goroutines"}'
```

### PUT /todos/:id
Atualiza uma tarefa existente.

**Body:**
```json
{
  "title": "Título atualizado",
  "content": "Conteúdo atualizado",
  "completed": true
}
```

### DELETE /todos/:id
Remove uma tarefa.

**Exemplo:**
```bash
curl -X DELETE http://localhost:8080/todos/1
```

### PATCH /todos/:id/toggle
Alterna o status de conclusão de uma tarefa.

**Exemplo:**
```bash
curl -X PATCH http://localhost:8080/todos/1/toggle
```

## Funcionalidades Demonstradas

- **Routing**: Diferentes métodos HTTP e parâmetros de rota
- **JSON Binding**: Deserialização automática de JSON
- **Query Parameters**: Filtragem usando query strings
- **Custom Middleware**: Middleware CORS personalizado
- **Error Handling**: Validação e tratamento de erros
- **Response Formatting**: Respostas JSON estruturadas
