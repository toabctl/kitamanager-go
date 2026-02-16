package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/eenemeene/kitamanager-go/internal/ctxkeys"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestStructuredLogger(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		path       string
		query      string
		setUserID  bool
		setError   bool
	}{
		{
			name:       "successful request",
			statusCode: http.StatusOK,
			path:       "/api/test",
		},
		{
			name:       "client error",
			statusCode: http.StatusBadRequest,
			path:       "/api/test",
		},
		{
			name:       "server error",
			statusCode: http.StatusInternalServerError,
			path:       "/api/test",
		},
		{
			name:       "request with query",
			statusCode: http.StatusOK,
			path:       "/api/test",
			query:      "foo=bar",
		},
		{
			name:       "request with user ID",
			statusCode: http.StatusOK,
			path:       "/api/test",
			setUserID:  true,
		},
		{
			name:       "request with error",
			statusCode: http.StatusBadRequest,
			path:       "/api/test",
			setError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(StructuredLogger())
			r.GET("/api/test", func(c *gin.Context) {
				if tt.setUserID {
					c.Set(ctxkeys.UserID, uint(1))
				}
				if tt.setError {
					_ = c.Error(gin.Error{
						Err:  http.ErrBodyNotAllowed,
						Type: gin.ErrorTypePublic,
					})
				}
				c.Status(tt.statusCode)
			})

			path := tt.path
			if tt.query != "" {
				path += "?" + tt.query
			}

			req, _ := http.NewRequest("GET", path, nil)
			req.Header.Set("User-Agent", "test-agent")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.statusCode {
				t.Errorf("expected status %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestStructuredLogger_NotFound(t *testing.T) {
	r := gin.New()
	r.Use(StructuredLogger())
	r.GET("/api/exists", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/api/notfound", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}
