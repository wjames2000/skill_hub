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
	"github.com/hpds/skill-hub/pkg/errno"
	mls "github.com/hpds/skill-hub/pkg/meilisearch"
	"github.com/hpds/skill-hub/pkg/response"
	"github.com/stretchr/testify/assert"
	"xorm.io/xorm"
)

func newTestEngine() *xorm.Engine {
	engine, err := xorm.NewEngine("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	engine.ShowSQL(false)
	_ = engine.Sync2(&model.Skill{}, &model.SkillCategory{})
	return engine
}

func seedSkills(engine *xorm.Engine) {
	skills := []*model.Skill{
		{Name: "excel-trend-analyzer", DisplayName: "Excel Trend Analyzer", Description: "Analyze Excel data trends", Repository: "https://github.com/test/excel", Stars: 100, Status: model.SkillStatusActive, Category: "data-analysis"},
		{Name: "ppt-generator", DisplayName: "PPT Generator", Description: "Generate PPT from data", Repository: "https://github.com/test/ppt", Stars: 200, Status: model.SkillStatusActive, Category: "presentation"},
		{Name: "archived-skill", DisplayName: "Archived", Description: "Old skill", Repository: "https://github.com/test/old", Stars: 50, Status: model.SkillStatusDeprecated, Category: "other"},
	}
	for _, s := range skills {
		_, _ = engine.Insert(s)
	}
}

func setupSkillTest() (*gin.Engine, *SkillHandler) {
	gin.SetMode(gin.TestMode)
	engine := newTestEngine()
	seedSkills(engine)
	skillRepo := repository.NewSkillRepo(engine)
	catRepo := repository.NewCategoryRepo(engine)
	h := NewSkillHandler(skillRepo, catRepo, nil)
	r := gin.New()
	rg := r.Group("/api/v1")
	h.RegisterRoutes(rg)
	return r, h
}

func TestSkillHandler_ListSkills(t *testing.T) {
	r, _ := setupSkillTest()

	t.Run("list all skills (active only)", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/skills", nil)
		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 0, resp.Code)
		data := resp.Data.(map[string]interface{})
		skills := data["skills"].([]interface{})
		total := data["total"].(float64)
		assert.Equal(t, 2, len(skills)) // only active ones
		assert.Equal(t, float64(2), total)
	})

	t.Run("list with category filter", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/skills?category=presentation", nil)
		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 0, resp.Code)
		data := resp.Data.(map[string]interface{})
		skills := data["skills"].([]interface{})
		assert.Equal(t, 1, len(skills))
	})

	t.Run("list with sort by name asc", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/skills?sort=name", nil)
		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 0, resp.Code)
		data := resp.Data.(map[string]interface{})
		skills := data["skills"].([]interface{})
		first := skills[0].(map[string]interface{})
		assert.Equal(t, "excel-trend-analyzer", first["name"])
	})

	t.Run("list with invalid page defaults to 1", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/skills?page=0&page_size=1", nil)
		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 0, resp.Code)
		data := resp.Data.(map[string]interface{})
		assert.Equal(t, float64(1), data["page"])
	})

	t.Run("list with page_size capped at 100", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/skills?page_size=200", nil)
		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 0, resp.Code)
		data := resp.Data.(map[string]interface{})
		assert.Equal(t, float64(20), data["size"])
	})
}

func TestSkillHandler_GetSkill(t *testing.T) {
	r, _ := setupSkillTest()

	t.Run("get existing skill", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/skills/1", nil)
		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 0, resp.Code)
		data := resp.Data.(map[string]interface{})
		assert.Equal(t, "excel-trend-analyzer", data["name"])
	})

	t.Run("get non-existing skill returns error", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/skills/999", nil)
		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, errno.NotFound.Code, resp.Code)
	})

	t.Run("get skill with invalid id", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/skills/abc", nil)
		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, errno.ParamError.Code, resp.Code)
	})
}

func TestSkillHandler_SearchSkills(t *testing.T) {
	r, _ := setupSkillTest()

	t.Run("search with empty query", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{}`
		req, _ := http.NewRequest("POST", "/api/v1/skills/search", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, errno.ParamError.Code, resp.Code)
	})

	t.Run("search with valid query (fallback to Like)", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{"query":"excel"}`
		req, _ := http.NewRequest("POST", "/api/v1/skills/search", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 0, resp.Code)
		data := resp.Data.(map[string]interface{})
		skills := data["skills"].([]interface{})
		assert.True(t, len(skills) >= 1)
	})

	t.Run("search with no results", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{"query":"zzzzzznotexist"}`
		req, _ := http.NewRequest("POST", "/api/v1/skills/search", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 0, resp.Code)
		data := resp.Data.(map[string]interface{})
		skills := data["skills"].([]interface{})
		assert.Equal(t, 0, len(skills))
	})
}

func TestNewSkillHandler(t *testing.T) {
	engine := newTestEngine()
	skillRepo := repository.NewSkillRepo(engine)
	catRepo := repository.NewCategoryRepo(engine)
	h := NewSkillHandler(skillRepo, catRepo, nil)
	assert.NotNil(t, h)
	assert.Equal(t, skillRepo, h.skillRepo)
	assert.Nil(t, h.meiliCli)
}

func TestSkillHandler_SearchWithMeili(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := newTestEngine()
	seedSkills(engine)
	skillRepo := repository.NewSkillRepo(engine)
	catRepo := repository.NewCategoryRepo(engine)
	meiliCli := &mls.Client{}
	h := NewSkillHandler(skillRepo, catRepo, meiliCli)
	r := gin.New()
	rg := r.Group("/api/v1")
	h.RegisterRoutes(rg)

	t.Run("search with meili returns internal error on meili failure", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{"query":"excel"}`
		req, _ := http.NewRequest("POST", "/api/v1/skills/search", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, errno.InternalError.Code, resp.Code)
	})
}
