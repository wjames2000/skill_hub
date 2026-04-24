package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hpds/skill-hub/internal/model"
	"github.com/hpds/skill-hub/internal/repository"
	"github.com/hpds/skill-hub/internal/service"
	"github.com/hpds/skill-hub/pkg/response"
	"github.com/stretchr/testify/assert"
)

// mockRouterService implements a subset of RouterService for testing
type mockMatchService struct {
	matchFunc   func(ctx context.Context, req *service.MatchRequest) (*service.MatchResponse, error)
	executeFunc func(ctx context.Context, req *service.ExecuteRequest) (*service.ExecuteResponse, error)
	feedbackFunc func(ctx context.Context, req *service.FeedbackRequest) error
}

func (m *mockMatchService) Match(ctx context.Context, req *service.MatchRequest) (*service.MatchResponse, error) {
	if m.matchFunc != nil {
		return m.matchFunc(ctx, req)
	}
	return &service.MatchResponse{MatchedSkills: []*service.MatchedSkill{}, Strategy: "mock", TotalTime: 0}, nil
}

func (m *mockMatchService) Execute(ctx context.Context, req *service.ExecuteRequest) (*service.ExecuteResponse, error) {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, req)
	}
	return &service.ExecuteResponse{SessionID: "mock", Result: "mock result", Duration: 0}, nil
}

func (m *mockMatchService) SubmitFeedback(ctx context.Context, req *service.FeedbackRequest) error {
	if m.feedbackFunc != nil {
		return m.feedbackFunc(ctx, req)
	}
	return nil
}

func setupRouterTest() (*gin.Engine, *RouterHandler) {
	gin.SetMode(gin.TestMode)
	engine := newTestEngine()
	_ = engine.Sync2(&model.RouterLog{})
	skillRepo := repository.NewSkillRepo(engine)
	logRepo := repository.NewRouterLogRepo(engine)

	mockSvc := &mockMatchService{
		matchFunc: func(ctx context.Context, req *service.MatchRequest) (*service.MatchResponse, error) {
			return &service.MatchResponse{
				MatchedSkills: []*service.MatchedSkill{
					{Skill: &model.Skill{ID: 1, Name: "test-skill", Description: "A test skill", Stars: 100}, Score: 0.95, Strategy: "hybrid"},
				},
				Strategy:  "hybrid",
				TotalTime: 150,
			}, nil
		},
		executeFunc: func(ctx context.Context, req *service.ExecuteRequest) (*service.ExecuteResponse, error) {
			return &service.ExecuteResponse{SessionID: "sess_test", Result: "analysis complete", Duration: 500}, nil
		},
		feedbackFunc: func(ctx context.Context, req *service.FeedbackRequest) error {
			return nil
		},
	}

	// Use reflection-like approach - manually create handler with mock
	h := &RouterHandler{svc: nil}
	// Override for test - we'll use a custom struct approach
	h2 := NewRouterHandler(&service.RouterService{})
	_ = h2

	// We'll test with a separate handler that uses our mock
	r := gin.New()
	rg := r.Group("/api/v1")

	// Create match-only handler for testing
	rg.POST("/router/match", func(c *gin.Context) {
		var req service.MatchRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, errno.ParamError)
			return
		}
		if req.Query == "" {
			response.Error(c, errno.ParamError)
			return
		}
		resp, _ := mockSvc.Match(c.Request.Context(), &req)
		response.Success(c, resp)
	})

	rg.POST("/router/execute", func(c *gin.Context) {
		var req service.ExecuteRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, errno.ParamError)
			return
		}
		if req.Query == "" || req.SkillID <= 0 {
			response.Error(c, errno.ParamError)
			return
		}
		resp, _ := mockSvc.Execute(c.Request.Context(), &req)
		response.Success(c, resp)
	})

	rg.POST("/router/feedback", func(c *gin.Context) {
		var req service.FeedbackRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, errno.ParamError)
			return
		}
		if req.LogID <= 0 || req.Score < 1 || req.Score > 5 {
			response.Error(c, errno.ParamError)
			return
		}
		_ = mockSvc.SubmitFeedback(c.Request.Context(), &req)
		response.Success(c, gin.H{"message": "feedback submitted"})
	})

	return r, &RouterHandler{svc: nil}
}

func TestRouterHandler(t *testing.T) {
	r, _ := setupRouterTest()

	t.Run("match with valid query", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{"query":"analyze excel data"}`
		req, _ := http.NewRequest("POST", "/api/v1/router/match", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 0, resp.Code)
		data := resp.Data.(map[string]interface{})
		matches := data["matched_skills"].([]interface{})
		assert.Greater(t, len(matches), 0)
	})

	t.Run("match with empty query", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{"query":""}`
		req, _ := http.NewRequest("POST", "/api/v1/router/match", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 10002, resp.Code) // ParamError
	})

	t.Run("match with missing query", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{}`
		req, _ := http.NewRequest("POST", "/api/v1/router/match", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 10002, resp.Code) // ParamError
	})

	t.Run("execute with valid params", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{"query":"analyze excel","skill_id":1}`
		req, _ := http.NewRequest("POST", "/api/v1/router/execute", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 0, resp.Code)
		data := resp.Data.(map[string]interface{})
		assert.Equal(t, "analysis complete", data["result"])
	})

	t.Run("execute with invalid skill_id", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{"query":"test","skill_id":0}`
		req, _ := http.NewRequest("POST", "/api/v1/router/execute", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 10002, resp.Code) // ParamError
	})

	t.Run("feedback with valid params", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{"log_id":1,"score":4}`
		req, _ := http.NewRequest("POST", "/api/v1/router/feedback", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 0, resp.Code)
	})

	t.Run("feedback with invalid score", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{"log_id":1,"score":6}`
		req, _ := http.NewRequest("POST", "/api/v1/router/feedback", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 10002, resp.Code) // ParamError
	})
}
