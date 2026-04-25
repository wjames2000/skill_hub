package handler

import (
	"strconv"

	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/hpds/skill-hub/internal/model"
	"github.com/hpds/skill-hub/internal/repository"
	"github.com/hpds/skill-hub/pkg/errno"
	mls "github.com/hpds/skill-hub/pkg/meilisearch"
	"github.com/hpds/skill-hub/pkg/response"
)

type SkillHandler struct {
	skillRepo    *repository.SkillRepo
	categoryRepo *repository.CategoryRepo
	meiliCli     *mls.Client
}

func NewSkillHandler(skillRepo *repository.SkillRepo, categoryRepo *repository.CategoryRepo, meiliCli *mls.Client) *SkillHandler {
	return &SkillHandler{skillRepo: skillRepo, categoryRepo: categoryRepo, meiliCli: meiliCli}
}

func (h *SkillHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/skills", h.ListSkills)
	rg.GET("/skills/trending", h.ListTrendingSkills)
	rg.GET("/skills/latest", h.ListLatestSkills)
	rg.GET("/skills/:id", h.GetSkill)
	rg.POST("/skills/search", h.SearchSkills)
}

func (h *SkillHandler) ListSkills(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")
	category := c.Query("category")
	sort := c.DefaultQuery("sort", "stars")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	sess := h.skillRepo.ListQuery()

	if category != "" {
		sess = sess.Where("category = ?", category)
	}

	sess = sess.Where("status = ?", model.SkillStatusActive)

	switch sort {
	case "installs":
		sess = sess.Desc("installs")
	case "created_at":
		sess = sess.Desc("created_at")
	case "score":
		sess = sess.Desc("score")
	case "name":
		sess = sess.Asc("name")
	default:
		sess = sess.Desc("stars")
	}

	var skills []*model.Skill
	total, err := sess.Limit(pageSize, (page-1)*pageSize).FindAndCount(&skills)
	if err != nil {
		response.Error(c, errno.DBError)
		return
	}

	response.Success(c, gin.H{
		"skills": skills,
		"total":  total,
		"page":   page,
		"size":   pageSize,
	})
}

func (h *SkillHandler) ListTrendingSkills(c *gin.Context) {
	c.Request.URL.RawQuery = "sort=installs&page_size=6"
	h.ListSkills(c)
}

func (h *SkillHandler) ListLatestSkills(c *gin.Context) {
	c.Request.URL.RawQuery = "sort=created_at&page_size=10"
	h.ListSkills(c)
}

func (h *SkillHandler) GetSkill(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, errno.ParamError)
		return
	}

	skill, err := h.skillRepo.GetByID(id)
	if err != nil {
		response.Error(c, errno.DBError)
		return
	}
	if skill == nil {
		response.Error(c, errno.NotFound)
		return
	}

	response.Success(c, skill)
}

func (h *SkillHandler) SearchSkills(c *gin.Context) {
	var req struct {
		Query    string `json:"query" binding:"required"`
		Page     int    `json:"page"`
		PageSize int    `json:"page_size"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errno.ParamError)
		return
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	if h.meiliCli == nil {
		skills, err := h.skillRepo.SearchByName(req.Query, req.PageSize)
		if err != nil {
			response.Error(c, errno.DBError)
			return
		}
		response.Success(c, gin.H{
			"skills": skills,
			"total":  len(skills),
			"page":   req.Page,
			"size":   req.PageSize,
		})
		return
	}

	resp, err := h.meiliCli.Search("skills", req.Query, int64(req.PageSize))
	if err != nil {
		response.Error(c, errno.InternalError)
		return
	}

	var skills []*model.Skill
	for _, hit := range resp.Hits {
		rawID, ok := hit["id"]
		if !ok {
			continue
		}
		var idFloat float64
		if data, err := rawID.MarshalJSON(); err == nil {
			_ = json.Unmarshal(data, &idFloat)
		}
		if idFloat <= 0 {
			continue
		}
		skill, err := h.skillRepo.GetByID(int64(idFloat))
		if err != nil || skill == nil {
			continue
		}
		skills = append(skills, skill)
	}

	response.Success(c, gin.H{
		"skills": skills,
		"total":  resp.EstimatedTotalHits,
		"page":   req.Page,
		"size":   req.PageSize,
	})
}
