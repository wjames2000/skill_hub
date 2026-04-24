package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hpds/skill-hub/internal/model"
	"github.com/hpds/skill-hub/internal/repository"
	"github.com/hpds/skill-hub/pkg/response"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func setupAuthTest() (*gin.Engine, *AuthHandler) {
	gin.SetMode(gin.TestMode)
	engine := newTestEngine()
	_ = engine.Sync2(&model.User{}, &model.APIKey{})
	userRepo := repository.NewUserRepo(engine)
	apiKeyRepo := repository.NewAPIKeyRepo(engine)

	hash, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
	_ = userRepo.Create(&model.User{
		Username:     "existinguser",
		Email:        "test@example.com",
		PasswordHash: string(hash),
		Role:         "user",
		Status:       1,
	})
	_ = userRepo.Create(&model.User{
		Username:     "adminuser",
		Email:        "admin@example.com",
		PasswordHash: string(hash),
		Role:         "admin",
		Status:       1,
	})

	h := NewAuthHandler(userRepo, apiKeyRepo, "test-secret", 72)
	r := gin.New()
	rg := r.Group("/api/v1")
	h.RegisterRoutes(rg)
	return r, h
}

func TestAuthHandler_Register(t *testing.T) {
	r, _ := setupAuthTest()

	t.Run("register success", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{"username":"newuser","email":"new@example.com","password":"newpass123"}`
		req, _ := http.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 0, resp.Code)
		data := resp.Data.(map[string]interface{})
		user := data["user"].(map[string]interface{})
		assert.Equal(t, "newuser", user["username"])
		assert.NotEmpty(t, data["token"])
	})

	t.Run("register duplicate username", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{"username":"existinguser","email":"other@example.com","password":"testpass"}`
		req, _ := http.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 20102, resp.Code) // UserExists
	})

	t.Run("register missing fields", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{"username":"test"}`
		req, _ := http.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 10002, resp.Code) // ParamError
	})
}

func TestAuthHandler_Login(t *testing.T) {
	r, _ := setupAuthTest()

	t.Run("login success with username", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{"username":"existinguser","password":"testpass"}`
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 0, resp.Code)
		data := resp.Data.(map[string]interface{})
		assert.NotEmpty(t, data["token"])
	})

	t.Run("login success with email", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{"username":"test@example.com","password":"testpass"}`
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 0, resp.Code)
		data := resp.Data.(map[string]interface{})
		assert.NotEmpty(t, data["token"])
	})

	t.Run("login wrong password", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{"username":"existinguser","password":"wrongpass"}`
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 20103, resp.Code) // InvalidPassword
	})

	t.Run("login user not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{"username":"nobody","password":"testpass"}`
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 20101, resp.Code) // UserNotFound
	})

	t.Run("login missing fields", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{}`
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 10002, resp.Code) // ParamError
	})
}

func TestAuthHandler_GitHubOAuth(t *testing.T) {
	r, _ := setupAuthTest()

	t.Run("github oauth returns redirect url", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/auth/github?client_id=test-client", nil)
		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 0, resp.Code)
		data := resp.Data.(map[string]interface{})
		redirect := data["redirect_url"].(string)
		assert.Contains(t, redirect, "https://github.com/login/oauth/authorize")
		assert.Contains(t, redirect, "test-client")
	})
}

func TestAuthHandler_GitHubCallback(t *testing.T) {
	r, _ := setupAuthTest()

	t.Run("github callback with new user", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{"code":"github_code_123"}`
		req, _ := http.NewRequest("POST", "/api/v1/auth/github/callback", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 0, resp.Code)
		data := resp.Data.(map[string]interface{})
		assert.NotEmpty(t, data["token"])
		user := data["user"].(map[string]interface{})
		assert.Contains(t, user["username"], "gh_")
	})

	t.Run("github callback without code", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{}`
		req, _ := http.NewRequest("POST", "/api/v1/auth/github/callback", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 10002, resp.Code) // ParamError
	})
}
