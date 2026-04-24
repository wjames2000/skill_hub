package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/hpds/skill-hub/pkg/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var baseURL string

func TestMain(m *testing.M) {
	baseURL = os.Getenv("SKILL_HUB_TEST_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	baseURL = strings.TrimRight(baseURL, "/")
	os.Exit(m.Run())
}

type apiClient struct {
	baseURL string
	token   string
}

func newClient() *apiClient {
	return &apiClient{baseURL: baseURL}
}

func (c *apiClient) setToken(tok string) {
	c.token = tok
}

func (c *apiClient) do(method, path string, body interface{}) (*http.Response, error) {
	var buf bytes.Buffer
	if body != nil {
		b, _ := json.Marshal(body)
		buf = *bytes.NewBuffer(b)
	}
	req, err := http.NewRequest(method, c.baseURL+path, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	return http.DefaultClient.Do(req)
}

func (c *apiClient) get(path string) (*http.Response, error) {
	return c.do("GET", path, nil)
}

func (c *apiClient) post(path string, body interface{}) (*http.Response, error) {
	return c.do("POST", path, body)
}

func parseResponse(resp *http.Response) (*response.Response, error) {
	var r response.Response
	err := json.NewDecoder(resp.Body).Decode(&r)
	return &r, err
}

func TestIntegration_SkillsAPI(t *testing.T) {
	client := newClient()

	t.Run("GET /api/v1/skills returns skills list", func(t *testing.T) {
		resp, err := client.get("/api/v1/skills")
		require.NoError(t, err)
		defer resp.Body.Close()

		r, err := parseResponse(resp)
		require.NoError(t, err)
		assert.Equal(t, 0, r.Code, "expected success, got: %v", r.Message)
	})

	t.Run("GET /api/v1/skills/:id returns skill detail", func(t *testing.T) {
		resp, err := client.get("/api/v1/skills/1")
		require.NoError(t, err)
		defer resp.Body.Close()

		r, err := parseResponse(resp)
		require.NoError(t, err)
		assert.Equal(t, 0, r.Code)
	})

	t.Run("GET /api/v1/skills/999 returns not found", func(t *testing.T) {
		resp, err := client.get("/api/v1/skills/999")
		require.NoError(t, err)
		defer resp.Body.Close()

		r, err := parseResponse(resp)
		require.NoError(t, err)
		assert.Equal(t, 30101, r.Code) // SkillNotFound
	})

	t.Run("POST /api/v1/skills/search returns results", func(t *testing.T) {
		resp, err := client.post("/api/v1/skills/search", map[string]string{
			"query": "excel",
		})
		require.NoError(t, err)
		defer resp.Body.Close()

		r, err := parseResponse(resp)
		require.NoError(t, err)
		assert.Equal(t, 0, r.Code)
	})

	t.Run("POST /api/v1/skills/search with empty query returns error", func(t *testing.T) {
		resp, err := client.post("/api/v1/skills/search", map[string]string{})
		require.NoError(t, err)
		defer resp.Body.Close()

		r, err := parseResponse(resp)
		require.NoError(t, err)
		assert.Equal(t, 10002, r.Code)
	})
}

func TestIntegration_CategoryAPI(t *testing.T) {
	client := newClient()

	t.Run("GET /api/v1/categories returns categories", func(t *testing.T) {
		resp, err := client.get("/api/v1/categories")
		require.NoError(t, err)
		defer resp.Body.Close()

		r, err := parseResponse(resp)
		require.NoError(t, err)
		assert.Equal(t, 0, r.Code)
	})
}

func TestIntegration_AuthAPI(t *testing.T) {
	client := newClient()

	t.Run("POST /api/v1/auth/register creates new user", func(t *testing.T) {
		username := fmt.Sprintf("inttest_%d", 12345)
		resp, err := client.post("/api/v1/auth/register", map[string]string{
			"username": username,
			"email":    username + "@test.com",
			"password": "testpass123",
		})
		require.NoError(t, err)
		defer resp.Body.Close()

		r, err := parseResponse(resp)
		require.NoError(t, err)
		assert.Equal(t, 0, r.Code, "register failed: %v", r.Message)
	})

	t.Run("POST /api/v1/auth/login with valid credentials", func(t *testing.T) {
		resp, err := client.post("/api/v1/auth/login", map[string]string{
			"username": "admin",
			"password": "admin123",
		})
		require.NoError(t, err)
		defer resp.Body.Close()

		r, err := parseResponse(resp)
		require.NoError(t, err)
		assert.Equal(t, 0, r.Code)
	})
}

func TestIntegration_RouterAPI(t *testing.T) {
	client := newClient()

	t.Run("POST /api/v1/router/match returns matched skills", func(t *testing.T) {
		resp, err := client.post("/api/v1/router/match", map[string]interface{}{
			"query": "analyze excel data",
			"top_k": 5,
		})
		require.NoError(t, err)
		defer resp.Body.Close()

		r, err := parseResponse(resp)
		require.NoError(t, err)
		assert.Equal(t, 0, r.Code, "match failed: %v", r.Message)
	})

	t.Run("POST /api/v1/router/match with empty query returns error", func(t *testing.T) {
		resp, err := client.post("/api/v1/router/match", map[string]interface{}{
			"query": "",
		})
		require.NoError(t, err)
		defer resp.Body.Close()

		r, err := parseResponse(resp)
		require.NoError(t, err)
		assert.Equal(t, 10002, r.Code)
	})

	t.Run("POST /api/v1/router/execute with skill", func(t *testing.T) {
		resp, err := client.post("/api/v1/router/execute", map[string]interface{}{
			"query":    "analyze data",
			"skill_id": 1,
		})
		require.NoError(t, err)
		defer resp.Body.Close()

		r, err := parseResponse(resp)
		require.NoError(t, err)
		assert.Equal(t, 0, r.Code, "execute failed: %v", r.Message)
	})
}

func TestIntegration_HealthCheck(t *testing.T) {
	client := newClient()

	t.Run("GET /api/v1/health returns ok", func(t *testing.T) {
		resp, err := client.get("/api/v1/health")
		require.NoError(t, err)
		defer resp.Body.Close()

		r, err := parseResponse(resp)
		require.NoError(t, err)
		assert.Equal(t, 0, r.Code)
	})
}
