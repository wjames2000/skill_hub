package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hpds/skill-hub/pkg/response"
	"github.com/stretchr/testify/assert"
)

func newTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func TestAuthMiddlewareNoSecret(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AuthMiddleware(""))

	r.GET("/protected", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "protected"})
	})

	t.Run("missing auth header returns error", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)
		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 10003, resp.Code) // Unauthorized
	})
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AuthMiddleware("test-secret"))

	r.GET("/protected", func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		assert.True(t, exists)
		role, exists := c.Get("role")
		assert.True(t, exists)
		c.JSON(200, gin.H{"user_id": userID, "role": role})
	})

	t.Run("valid token passes middleware", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer test-valid-token")
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("missing authorization header", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)
		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 10003, resp.Code)
	})

	t.Run("invalid token format", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "InvalidFormat")
		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 10003, resp.Code)
	})
}

func TestAdminMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("admin role passes", func(t *testing.T) {
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("role", "admin")
			c.Next()
		})
		r.Use(AdminMiddleware())

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/admin", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("user role is forbidden", func(t *testing.T) {
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("role", "user")
			c.Next()
		})
		r.Use(AdminMiddleware())

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/admin", nil)
		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 10004, resp.Code) // Forbidden
	})

	t.Run("missing role is forbidden", func(t *testing.T) {
		r := gin.New()
		r.Use(AdminMiddleware())

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/admin", nil)
		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 10004, resp.Code)
	})
}

func TestCORSMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(CORSMiddleware())

	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	t.Run("sets CORS headers", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "http://localhost:5173")
		r.ServeHTTP(w, req)

		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Authorization")
	})

	t.Run("handles OPTIONS preflight", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("OPTIONS", "/test", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, 204, w.Code)
	})
}

func TestRateLimiter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("within rate limit passes", func(t *testing.T) {
		r := gin.New()
		r.Use(RateLimiter(2, 1)) // 2 requests per second

		for i := 0; i < 2; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			r.ServeHTTP(w, req)
			assert.Equal(t, 200, w.Code)
		}
	})
}

func TestRecoveryMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RecoveryMiddleware())

	r.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	t.Run("recovers from panic", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/panic", nil)
		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 10001, resp.Code) // InternalError
	})
}

func TestLoggerMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(LoggerMiddleware())

	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	t.Run("logger middleware does not block", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
	})
}
