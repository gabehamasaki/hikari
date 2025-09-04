# File Upload Example

Complete file upload system with multiple upload methods, file management and static serving using Hikari-Go.

**Language / Idioma:** [English](README.md) | [Português Brasil](README.pt-BR.md)

## Features

- 📤 Single file upload
- 📤📤 Multiple file upload (batch)
- 📋 File listing and metadata
- 📥 File download with proper headers
- 🗑️ File deletion
- 🌐 Static file serving
- ✅ File type validation
- 📏 File size limits (10MB per file)
- 🔒 Security checks for static serving
- 🏗️ Organized route groups structure
- 🩺 Health check and system info
- 📊 Upload statistics

## How to run

```bash
cd examples/file-upload
go run main.go
```

The server will start at `http://localhost:8082`

## API Structure

The API uses grouped routes for organized file operations:

```
/                    → API information
/api/v1/
├── /health          → Health check
├── /info            → System information & stats
├── /files/
│   ├── GET /        → List all files
│   ├── GET /:id     → Get file information
│   └── DELETE /:id  → Delete file
├── /upload/
│   ├── POST /       → Upload single file
│   └── POST /multiple → Upload multiple files
└── GET /download/:id → Download file
/static/*            → Serve uploaded files directly
```

## File Storage

- **Upload Directory:** `./uploads`
- **Naming:** Files are renamed with unique IDs: `file_<timestamp><extension>`
- **Metadata:** Stored in memory (ID, original name, size, content type, upload time)
- **Max Size:** 10MB per file, 50MB total for batch uploads

## Endpoints

### GET /
Returns API information and available endpoints.

### GET /api/v1/health
Health check endpoint that verifies upload directory status.

**Example Response:**
```json
{
  "status": "healthy",
  "upload_directory": "./uploads",
  "directory_exists": true,
  "files_count": 5,
  "timestamp": "2025-01-01T10:00:00Z"
}
```

### GET /api/v1/info
Returns system information and upload statistics.

**Example Response:**
```json
{
  "service": "file-upload",
  "version": "1.0.0",
  "total_files": 10,
  "total_size": 52428800,
  "upload_dir": "./uploads"
}
```

### File Operations

#### POST /api/v1/upload
Uploads a single file.

**Content-Type:** `multipart/form-data`

**Form Fields:**
- `file`: The file to upload

**Example:**
```bash
curl -X POST http://localhost:8082/api/v1/upload \
  -F "file=@document.pdf"
```

**Response:**
```json
{
  "message": "File uploaded successfully",
  "file": {
    "id": "file_1234567890123456789",
    "name": "document.pdf",
    "size": 1048576,
    "content_type": "application/pdf",
    "uploaded_at": "2025-01-01T10:00:00Z",
    "path": "file_1234567890123456789.pdf"
  },
  "download_url": "/api/v1/download/file_1234567890123456789",
  "static_url": "/static/file_1234567890123456789.pdf"
}
```

#### POST /api/v1/upload/multiple
Uploads multiple files in a single request.

**Content-Type:** `multipart/form-data`

**Form Fields:**
- `files`: Multiple files to upload

**Example:**
```bash
curl -X POST http://localhost:8082/api/v1/upload/multiple \
  -F "files=@file1.txt" \
  -F "files=@file2.jpg" \
  -F "files=@file3.pdf"
```

**Response:**
```json
{
  "message": "Processed 3 files",
  "uploaded_files": [
    {
      "file": { /* file info */ },
      "download_url": "/api/v1/download/file_id_1",
      "static_url": "/static/file_name_1.txt"
    }
  ],
  "uploaded_count": 2,
  "total_count": 3,
  "errors": [
    {
      "file": "large_file.zip",
      "error": "File too large (max 10MB)"
    }
  ],
  "error_count": 1
}
```

#### GET /api/v1/files
Lists all uploaded files with metadata.

**Example Response:**
```json
{
  "files": [
    {
      "id": "file_1234567890123456789",
      "name": "document.pdf",
      "size": 1048576,
      "content_type": "application/pdf",
      "uploaded_at": "2025-01-01T10:00:00Z",
      "path": "file_1234567890123456789.pdf"
    }
  ],
  "count": 1
}
```

#### GET /api/v1/files/:id
Returns information about a specific file.

**Example:**
```bash
curl http://localhost:8082/api/v1/files/file_1234567890123456789
```

#### GET /api/v1/download/:id
Downloads a file with proper headers and content disposition.

**Example:**
```bash
curl -O http://localhost:8082/api/v1/download/file_1234567890123456789
```

#### DELETE /api/v1/files/:id
Deletes a file from both storage and metadata.

**Example:**
```bash
curl -X DELETE http://localhost:8082/api/v1/files/file_1234567890123456789
```

### Static File Serving

#### GET /static/*
Serves uploaded files directly via static serving.

**Security Features:**
- Path traversal protection
- Upload directory confinement
- File existence validation

**Example:**
```bash
curl http://localhost:8082/static/file_1234567890123456789.pdf
```

## Code Structure

### Route Groups

```go
// API v1 group
v1Group := app.Group("/api/v1")
{
    // File management routes
    filesGroup := v1Group.Group("/files")
    {
        filesGroup.GET("/", listFiles)
        filesGroup.GET("/:id", getFileInfo)
        filesGroup.DELETE("/:id", deleteFile)
    }

    // Upload routes
    uploadGroup := v1Group.Group("/upload")
    {
        uploadGroup.POST("/", uploadFile)
        uploadGroup.POST("/multiple", uploadMultipleFiles)
    }

    // Direct endpoints
    v1Group.GET("/download/:id", downloadFile)
    v1Group.GET("/health", healthCheck)
    v1Group.GET("/info", systemInfo)
}

// Static serving (outside API versioning)
app.GET("/static/*", serveStatic)
```

### File Type Validation

Supported file types:
- **Images:** JPEG, PNG, GIF, WebP
- **Documents:** PDF, Word, Excel
- **Text:** Plain text, CSV, JSON, XML
- **Archives:** ZIP
- **Generic:** Any text/* content type

### Security Features

- File size limits (10MB per file)
- Content type validation
- Path traversal protection
- Secure file naming
- Upload directory isolation

## Testing

Use the provided HTTP test file:
```
examples/file-upload/requests/test-sequence.http
```

The test sequence includes:
1. API information and health checks
2. Single file upload
3. Multiple file upload
4. File listing and information
5. File download (both methods)
6. File deletion
7. Error scenarios and edge cases

## Error Handling

The API provides detailed error responses:

- **400 Bad Request:** Invalid file, size limit exceeded, no file provided
- **404 Not Found:** File not found in metadata or disk
- **403 Forbidden:** Path traversal attempt
- **500 Internal Server Error:** File system errors
- **207 Multi-Status:** Partial success in batch uploads
