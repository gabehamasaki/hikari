package hikari

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

// ResponseWriter wrapper to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     200, // Default status code
		written:        false,
	}
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.written {
		rw.statusCode = code
		rw.written = true
		rw.ResponseWriter.WriteHeader(code)
	}
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	if !rw.written {
		rw.WriteHeader(200) // Default to 200 if WriteHeader wasn't called
	}
	return rw.ResponseWriter.Write(data)
}

func (rw *responseWriter) StatusCode() int {
	return rw.statusCode
}

type Context struct {
	Writer  *responseWriter
	Request *http.Request
	Params  map[string]string
	Logger  *zap.Logger
}

func (c *Context) JSON(status int, v any) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(status)
	_ = json.NewEncoder(c.Writer).Encode(v)
}

func (c *Context) String(status int, format string, values ...any) {
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.WriteHeader(status)
	_, _ = c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) Param(key string) string {
	return c.Params[key]
}

func (c *Context) Wildcard() string {
	return c.Params["*"]
}

func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

func (c *Context) FormValue(key string) string {
	return c.Request.FormValue(key)
}

func (c *Context) Bind(v any) error {
	return json.NewDecoder(c.Request.Body).Decode(v)
}

func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) GetHeader(key string) string {
	return c.Writer.Header().Get(key)
}

func (c *Context) Method() string {
	return c.Request.Method
}

func (c *Context) Path() string {
	return c.Request.URL.Path
}

func (c *Context) Status(status int) {
	c.Writer.WriteHeader(status)
}

func (c *Context) GetStatus() int {
	return c.Writer.StatusCode()
}

func (c *Context) File(filePath string) {
	http.ServeFile(c.Writer, c.Request, filePath)
}
