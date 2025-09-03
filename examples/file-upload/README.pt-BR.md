# File Upload Example

Sistema completo de upload e download de arquivos usando Hikari-Go.

**Language / Idioma:** [English](README.md) | [Português Brasil](README.pt-BR.md)

## Características

- Upload de arquivo único
- Upload de múltiplos arquivos
- Download de arquivos
- Listagem de arquivos
- Remoção de arquivos
- Servir arquivos estáticos
- Validação de tipo e tamanho de arquivo
- Health check do sistema

## Como executar

```bash
cd examples/file-upload
go run main.go
```

O servidor será iniciado em `http://localhost:8082`

**Nota:** O diretório `uploads` será criado automaticamente.

## Limites e Restrições

- **Tamanho máximo por arquivo:** 10MB
- **Tamanho máximo total (múltiplos arquivos):** 50MB
- **Tipos permitidos:**
  - Imagens: JPEG, PNG, GIF, WebP
  - Documentos: PDF, DOC, DOCX, XLS, XLSX
  - Texto: TXT, CSV, JSON, XML
  - Arquivos: ZIP
  - Todos os tipos text/*

## Endpoints

### GET /
Informações sobre a API e endpoints disponíveis.

### POST /upload
Upload de um único arquivo.

**Content-Type:** `multipart/form-data`
**Form Field:** `file`

**Exemplo com curl:**
```bash
curl -X POST http://localhost:8082/upload \
  -F "file=@example.jpg"
```

**Response:**
```json
{
  "message": "File uploaded successfully",
  "file": {
    "id": "file_1693834567123456789",
    "name": "example.jpg",
    "size": 245760,
    "content_type": "image/jpeg",
    "uploaded_at": "2023-09-04T15:30:00Z",
    "path": "file_1693834567123456789.jpg"
  },
  "download_url": "/download/file_1693834567123456789",
  "static_url": "/static/file_1693834567123456789.jpg"
}
```

### POST /upload/multiple
Upload de múltiplos arquivos.

**Content-Type:** `multipart/form-data`
**Form Field:** `files` (multiple)

**Exemplo com curl:**
```bash
curl -X POST http://localhost:8082/upload/multiple \
  -F "files=@file1.jpg" \
  -F "files=@file2.pdf" \
  -F "files=@file3.txt"
```

### GET /files
Lista todos os arquivos enviados.

**Exemplo:**
```bash
curl http://localhost:8082/files
```

### GET /files/:id
Obtém informações de um arquivo específico.

**Exemplo:**
```bash
curl http://localhost:8082/files/file_1693834567123456789
```

### GET /download/:id
Faz download de um arquivo.

**Exemplo:**
```bash
curl -O http://localhost:8082/download/file_1693834567123456789
```

### DELETE /files/:id
Remove um arquivo.

**Exemplo:**
```bash
curl -X DELETE http://localhost:8082/files/file_1693834567123456789
```

### GET /static/:filename
Serve arquivos estáticos diretamente.

**Exemplo:**
```bash
curl http://localhost:8082/static/file_1693834567123456789.jpg
```

### GET /health
Verifica o status do sistema.

**Exemplo:**
```bash
curl http://localhost:8082/health
```

## Exemplos de Uso Completo

### 1. Upload de uma Imagem
```bash
# Upload
curl -X POST http://localhost:8082/upload \
  -F "file=@photo.jpg"

# Resposta
{
  "message": "File uploaded successfully",
  "file": {
    "id": "file_1693834567123456789",
    "name": "photo.jpg",
    "size": 245760,
    "content_type": "image/jpeg",
    "uploaded_at": "2023-09-04T15:30:00Z",
    "path": "file_1693834567123456789.jpg"
  },
  "download_url": "/download/file_1693834567123456789",
  "static_url": "/static/file_1693834567123456789.jpg"
}
```

### 2. Upload de Múltiplos Arquivos
```bash
curl -X POST http://localhost:8082/upload/multiple \
  -F "files=@document.pdf" \
  -F "files=@image.png" \
  -F "files=@data.csv"
```

### 3. Listar Todos os Arquivos
```bash
curl http://localhost:8082/files
```

### 4. Download de Arquivo
```bash
# Por ID (com nome original)
curl -O -J http://localhost:8082/download/file_1693834567123456789

# Direto por path estático
curl -O http://localhost:8082/static/file_1693834567123456789.jpg
```

### 5. Verificar Informações de um Arquivo
```bash
curl http://localhost:8082/files/file_1693834567123456789
```

### 6. Remover um Arquivo
```bash
curl -X DELETE http://localhost:8082/files/file_1693834567123456789
```

## Exemplo HTML Simples

Crie um arquivo `upload.html` para testar na web:

```html
<!DOCTYPE html>
<html>
<head>
    <title>File Upload Test</title>
</head>
<body>
    <h1>Single File Upload</h1>
    <form action="http://localhost:8082/upload" method="post" enctype="multipart/form-data">
        <input type="file" name="file" required>
        <button type="submit">Upload</button>
    </form>

    <h1>Multiple File Upload</h1>
    <form action="http://localhost:8082/upload/multiple" method="post" enctype="multipart/form-data">
        <input type="file" name="files" multiple required>
        <button type="submit">Upload Multiple</button>
    </form>
</body>
</html>
```

## Funcionalidades Demonstradas

- **File Upload**: Upload de arquivos únicos e múltiplos
- **File Validation**: Validação de tipo e tamanho de arquivo
- **File Serving**: Servir arquivos estáticos e downloads
- **Error Handling**: Tratamento de erros de upload e validação
- **Metadata Storage**: Armazenamento de informações dos arquivos
- **Security**: Validação de paths para evitar directory traversal
- **Health Check**: Monitoramento do status do sistema
- **CORS Support**: Suporte para requisições cross-origin

## Estrutura de Arquivos

```
uploads/
├── file_1693834567123456789.jpg
├── file_1693834567123456790_0.pdf
├── file_1693834567123456791_1.txt
└── ...
```

## Melhorias para Produção

Para usar em produção, considere implementar:

- Persistência dos metadados em banco de dados
- Armazenamento em cloud (AWS S3, Google Cloud Storage)
- Autenticação e autorização
- Rate limiting para uploads
- Scan de vírus/malware
- Compressão de imagens
- Thumbnails automáticos
- Cleanup automático de arquivos antigos
- Logs detalhados de uploads
- Backup automático
