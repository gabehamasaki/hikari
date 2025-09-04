package hikari

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Context struct {
	context.Context
	Writer  *responseWriter
	Request *http.Request
	Params  map[string]string
	Logger  *zap.Logger

	storage      map[string]interface{}
	mutexStorage sync.RWMutex
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

func (c *Context) Set(key string, value interface{}) {
	c.mutexStorage.Lock()
	defer c.mutexStorage.Unlock()

	c.storage[key] = value
}

func (c *Context) Get(key string) (interface{}, bool) {
	c.mutexStorage.RLock()
	defer c.mutexStorage.RUnlock()

	value, exists := c.storage[key]
	return value, exists
}

func (c *Context) MustGet(key string) interface{} {
	c.mutexStorage.RLock()
	defer c.mutexStorage.RUnlock()

	value, exists := c.storage[key]
	if !exists {
		c.Logger.Error("Key not found in context storage", zap.String("key", key))
		return nil
	}
	return value
}

func (c *Context) GetString(key string) string {
	if value, exists := c.Get(key); exists {
		if s, ok := value.(string); ok {
			return s
		}
	}
	return ""
}

func (c *Context) GetInt(key string) int {
	if value, exists := c.Get(key); exists {
		if i, ok := value.(int); ok {
			return i
		}
	}
	return 0
}

func (c *Context) GetBool(key string) bool {
	if value, exists := c.Get(key); exists {
		if b, ok := value.(bool); ok {
			return b
		}
	}
	return false
}

func (c *Context) Keys() []string {
	c.mutexStorage.RLock()
	defer c.mutexStorage.RUnlock()

	if c.storage == nil {
		return []string{}
	}

	keys := make([]string, 0, len(c.storage))
	for key := range c.storage {
		keys = append(keys, key)
	}
	return keys
}

func (c *Context) WithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.Context, timeout)
}

func (c *Context) WithCancel() (context.Context, context.CancelFunc) {
	return context.WithCancel(c.Context)
}

func (c *Context) WithValue(key, value interface{}) context.Context {
	return context.WithValue(c.Context, key, value)
}

func (c *Context) Value(key interface{}) interface{} {
	return c.Context.Value(key)
}

func (c *Context) Done() <-chan struct{} {
	return c.Context.Done()
}

func (c *Context) Err() error {
	return c.Context.Err()
}
