package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gabehamasaki/hikari-go/pkg/hikari"
)

type FileInfo struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	UploadedAt  time.Time `json:"uploaded_at"`
	Path        string    `json:"path"`
}

// In-memory storage for file metadata
var files = make(map[string]*FileInfo)
var uploadDir = "./uploads"

func main() {
	app := hikari.New(":8082")

	// Create uploads directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create uploads directory: %v", err))
	}

	// Middleware for CORS
	app.Use(func(next hikari.HandlerFunc) hikari.HandlerFunc {
		return func(c *hikari.Context) {
			c.SetHeader("Access-Control-Allow-Origin", "*")
			c.SetHeader("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
			c.SetHeader("Access-Control-Allow-Headers", "Content-Type")

			if c.Method() == "OPTIONS" {
				c.Status(http.StatusOK)
				return
			}

			next(c)
		}
	})

	// Routes
	app.GET("/", homePage)
	app.GET("/files", listFiles)
	app.GET("/files/:id", getFileInfo)
	app.GET("/download/:id", downloadFile)
	app.POST("/upload", uploadFile)
	app.POST("/upload/multiple", uploadMultipleFiles)
	app.DELETE("/files/:id", deleteFile)
	app.GET("/health", healthCheck)

	// Static file serving for uploads directory
	app.GET("/static/*", serveStatic)

	fmt.Println("📁 File Upload Server running on http://localhost:8082")
	fmt.Println("📂 Upload directory: " + uploadDir)
	app.ListenAndServe()
}

func homePage(c *hikari.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "File Upload API",
		"version": "1.0.0",
		"features": []string{
			"Single file upload",
			"Multiple file upload",
			"File download",
			"File listing",
			"File deletion",
			"Static file serving",
		},
		"endpoints": map[string]string{
			"POST /upload":          "Upload single file",
			"POST /upload/multiple": "Upload multiple files",
			"GET /files":            "List all files",
			"GET /files/:id":        "Get file information",
			"GET /download/:id":     "Download file",
			"DELETE /files/:id":     "Delete file",
			"GET /static/*":         "Serve static files",
		},
		"max_file_size": "10MB per file",
	})
}

func uploadFile(c *hikari.Context) {
	// Parse multipart form with max memory of 10MB
	err := c.Request.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Unable to parse form",
		})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "No file uploaded or invalid file",
		})
		return
	}
	defer file.Close()

	// Validate file size (max 10MB)
	if header.Size > 10<<20 {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "File too large (max 10MB)",
		})
		return
	}

	// Validate file type (basic validation)
	contentType := header.Header.Get("Content-Type")
	if !isAllowedFileType(contentType) {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "File type not allowed",
		})
		return
	}

	// Generate unique filename
	fileID := generateFileID()
	fileExt := filepath.Ext(header.Filename)
	fileName := fileID + fileExt
	filePath := filepath.Join(uploadDir, fileName)

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Unable to create file",
		})
		return
	}
	defer dst.Close()

	// Copy uploaded file to destination
	size, err := io.Copy(dst, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Unable to save file",
		})
		return
	}

	// Store file metadata
	fileInfo := &FileInfo{
		ID:          fileID,
		Name:        header.Filename,
		Size:        size,
		ContentType: contentType,
		UploadedAt:  time.Now(),
		Path:        fileName,
	}
	files[fileID] = fileInfo

	c.JSON(http.StatusCreated, map[string]interface{}{
		"message":      "File uploaded successfully",
		"file":         fileInfo,
		"download_url": fmt.Sprintf("/download/%s", fileID),
		"static_url":   fmt.Sprintf("/static/%s", fileName),
	})
}

func uploadMultipleFiles(c *hikari.Context) {
	// Parse multipart form with max memory of 50MB
	err := c.Request.ParseMultipartForm(50 << 20) // 50MB total
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Unable to parse form",
		})
		return
	}

	form := c.Request.MultipartForm
	uploadedFiles := form.File["files"]

	if len(uploadedFiles) == 0 {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "No files uploaded",
		})
		return
	}

	var results []map[string]interface{}
	var errors []map[string]string

	for i, header := range uploadedFiles {
		file, err := header.Open()
		if err != nil {
			errors = append(errors, map[string]string{
				"file":  header.Filename,
				"error": "Unable to open file",
			})
			continue
		}

		// Validate file size
		if header.Size > 10<<20 {
			file.Close()
			errors = append(errors, map[string]string{
				"file":  header.Filename,
				"error": "File too large (max 10MB)",
			})
			continue
		}

		// Validate file type
		contentType := header.Header.Get("Content-Type")
		if !isAllowedFileType(contentType) {
			file.Close()
			errors = append(errors, map[string]string{
				"file":  header.Filename,
				"error": "File type not allowed",
			})
			continue
		}

		// Generate unique filename
		fileID := generateFileID() + "_" + strconv.Itoa(i)
		fileExt := filepath.Ext(header.Filename)
		fileName := fileID + fileExt
		filePath := filepath.Join(uploadDir, fileName)

		// Create destination file
		dst, err := os.Create(filePath)
		if err != nil {
			file.Close()
			errors = append(errors, map[string]string{
				"file":  header.Filename,
				"error": "Unable to create file",
			})
			continue
		}

		// Copy uploaded file to destination
		size, err := io.Copy(dst, file)
		dst.Close()
		file.Close()

		if err != nil {
			os.Remove(filePath) // Clean up on error
			errors = append(errors, map[string]string{
				"file":  header.Filename,
				"error": "Unable to save file",
			})
			continue
		}

		// Store file metadata
		fileInfo := &FileInfo{
			ID:          fileID,
			Name:        header.Filename,
			Size:        size,
			ContentType: contentType,
			UploadedAt:  time.Now(),
			Path:        fileName,
		}
		files[fileID] = fileInfo

		results = append(results, map[string]interface{}{
			"file":         fileInfo,
			"download_url": fmt.Sprintf("/download/%s", fileID),
			"static_url":   fmt.Sprintf("/static/%s", fileName),
		})
	}

	response := map[string]interface{}{
		"message":        fmt.Sprintf("Processed %d files", len(uploadedFiles)),
		"uploaded_files": results,
		"uploaded_count": len(results),
		"total_count":    len(uploadedFiles),
	}

	if len(errors) > 0 {
		response["errors"] = errors
		response["error_count"] = len(errors)
	}

	statusCode := http.StatusCreated
	if len(errors) == len(uploadedFiles) {
		statusCode = http.StatusBadRequest
	} else if len(errors) > 0 {
		statusCode = http.StatusPartialContent
	}

	c.JSON(statusCode, response)
}

func listFiles(c *hikari.Context) {
	var fileList []*FileInfo
	for _, fileInfo := range files {
		fileList = append(fileList, fileInfo)
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"files": fileList,
		"count": len(fileList),
	})
}

func getFileInfo(c *hikari.Context) {
	fileID := c.Param("id")
	fileInfo, exists := files[fileID]
	if !exists {
		c.JSON(http.StatusNotFound, map[string]string{
			"error": "File not found",
		})
		return
	}

	// Check if file still exists on disk
	filePath := filepath.Join(uploadDir, fileInfo.Path)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		delete(files, fileID) // Remove from memory if file doesn't exist
		c.JSON(http.StatusNotFound, map[string]string{
			"error": "File not found on disk",
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"file":         fileInfo,
		"download_url": fmt.Sprintf("/download/%s", fileID),
		"static_url":   fmt.Sprintf("/static/%s", fileInfo.Path),
	})
}

func downloadFile(c *hikari.Context) {
	fileID := c.Param("id")
	fileInfo, exists := files[fileID]
	if !exists {
		c.JSON(http.StatusNotFound, map[string]string{
			"error": "File not found",
		})
		return
	}

	filePath := filepath.Join(uploadDir, fileInfo.Path)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		delete(files, fileID)
		c.JSON(http.StatusNotFound, map[string]string{
			"error": "File not found on disk",
		})
		return
	}

	// Set headers for file download
	c.SetHeader("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileInfo.Name))
	c.SetHeader("Content-Type", fileInfo.ContentType)
	c.SetHeader("Content-Length", strconv.FormatInt(fileInfo.Size, 10))

	// Serve file
	http.ServeFile(c.Writer, c.Request, filePath)
}

func deleteFile(c *hikari.Context) {
	fileID := c.Param("id")
	fileInfo, exists := files[fileID]
	if !exists {
		c.JSON(http.StatusNotFound, map[string]string{
			"error": "File not found",
		})
		return
	}

	filePath := filepath.Join(uploadDir, fileInfo.Path)

	// Delete file from disk
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Unable to delete file",
		})
		return
	}

	// Remove from memory
	delete(files, fileID)

	c.JSON(http.StatusOK, map[string]string{
		"message": "File deleted successfully",
	})
}

func serveStatic(c *hikari.Context) {
	// Get the file path from the wildcard parameter
	filePath := c.Wildcard()

	if filePath == "" {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "No file specified",
		})
		return
	}

	fullPath := filepath.Join(uploadDir, filePath)

	// Security check: ensure the path is within the upload directory
	absUploadDir, _ := filepath.Abs(uploadDir)
	absFullPath, _ := filepath.Abs(fullPath)
	if !strings.HasPrefix(absFullPath, absUploadDir) {
		c.JSON(http.StatusForbidden, map[string]string{
			"error": "Access denied",
		})
		return
	}

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, map[string]string{
			"error": "File not found",
		})
		return
	}

	// Serve the file
	c.File(fullPath)
}

func healthCheck(c *hikari.Context) {
	// Check if upload directory exists and is writable
	uploadInfo, err := os.Stat(uploadDir)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, map[string]interface{}{
			"status":           "unhealthy",
			"upload_directory": "not accessible",
			"error":            err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status":           "healthy",
		"upload_directory": uploadDir,
		"directory_exists": uploadInfo.IsDir(),
		"files_count":      len(files),
		"timestamp":        time.Now(),
	})
}

// Helper functions
func isAllowedFileType(contentType string) bool {
	allowedTypes := []string{
		"image/jpeg",
		"image/png",
		"image/gif",
		"image/webp",
		"text/plain",
		"text/csv",
		"application/pdf",
		"application/json",
		"application/xml",
		"application/zip",
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/vnd.ms-excel",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	}

	for _, allowedType := range allowedTypes {
		if contentType == allowedType {
			return true
		}
	}

	// Allow text files
	if strings.HasPrefix(contentType, "text/") {
		return true
	}

	return false
}

func generateFileID() string {
	return fmt.Sprintf("file_%d", time.Now().UnixNano())
}
