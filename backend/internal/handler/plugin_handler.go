package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hpds/skill-hub/internal/model"
	"github.com/hpds/skill-hub/internal/repository"
	"github.com/hpds/skill-hub/pkg/errno"
	mls "github.com/hpds/skill-hub/pkg/meilisearch"
	"github.com/hpds/skill-hub/pkg/response"
)

type PluginHandler struct {
	skillRepo    *repository.SkillRepo
	categoryRepo *repository.CategoryRepo
	meiliCli     *mls.Client
	favoriteRepo *repository.FavoriteRepo
	reviewRepo   *repository.ReviewRepo
}

func NewPluginHandler(
	skillRepo *repository.SkillRepo,
	categoryRepo *repository.CategoryRepo,
	meiliCli *mls.Client,
	favoriteRepo *repository.FavoriteRepo,
	reviewRepo *repository.ReviewRepo,
) *PluginHandler {
	return &PluginHandler{
		skillRepo:    skillRepo,
		categoryRepo: categoryRepo,
		meiliCli:     meiliCli,
		favoriteRepo: favoriteRepo,
		reviewRepo:   reviewRepo,
	}
}

func (h *PluginHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/plugin/hot", h.HotSkills)
	rg.GET("/plugin/download", h.DownloadSkill)
	rg.GET("/plugin/recommend", h.RecommendSkills)
	rg.POST("/plugin/sync/status", h.SyncStatus)
	rg.GET("/plugin/categories", h.ListCategories)
	rg.GET("/plugin/skills", h.ListSkills)
}

func (h *PluginHandler) HotSkills(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	limit := 20
	if n, err := strconv.Atoi(limitStr); err == nil && n > 0 && n <= 100 {
		limit = n
	}

	sess := h.skillRepo.ListQuery()
	sess = sess.Where("status = ?", model.SkillStatusActive).Desc("installs").Desc("stars")

	var skills []*model.Skill
	if err := sess.Limit(limit, 0).Find(&skills); err != nil {
		response.Error(c, errno.DBError)
		return
	}

	response.Success(c, gin.H{"skills": skills})
}

func (h *PluginHandler) DownloadSkill(c *gin.Context) {
	idStr := c.Query("id")
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

	if err := h.skillRepo.IncrementInstalls(id); err != nil {
		response.Error(c, errno.DBError)
		return
	}

	response.Success(c, gin.H{
		"skill":        skill,
		"download_url": skill.Repository,
	})
}

func (h *PluginHandler) RecommendSkills(c *gin.Context) {
	category := c.Query("category")
	limitStr := c.DefaultQuery("limit", "10")
	limit := 10
	if n, err := strconv.Atoi(limitStr); err == nil && n > 0 && n <= 50 {
		limit = n
	}

	sess := h.skillRepo.ListQuery()
	sess = sess.Where("status = ?", model.SkillStatusActive)

	if category != "" {
		sess = sess.Where("category = ?", category)
	}

	sess = sess.Desc("score").Desc("stars")

	var skills []*model.Skill
	if err := sess.Limit(limit, 0).Find(&skills); err != nil {
		response.Error(c, errno.DBError)
		return
	}

	response.Success(c, gin.H{"skills": skills})
}

func (h *PluginHandler) SyncStatus(c *gin.Context) {
	var req struct {
		InstalledIDs []int64 `json:"installed_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errno.ParamError)
		return
	}

	type SkillSyncStatus struct {
		ID        int64  `json:"id"`
		Name      string `json:"name"`
		Version   string `json:"version"`
		UpdatedAt string `json:"updated_at"`
		IsRemoved bool   `json:"is_removed"`
	}

	results := make([]SkillSyncStatus, 0, len(req.InstalledIDs))
	for _, id := range req.InstalledIDs {
		skill, err := h.skillRepo.GetByID(id)
		if err != nil || skill == nil || skill.Status == model.SkillStatusDeprecated {
			results = append(results, SkillSyncStatus{
				ID:        id,
				IsRemoved: true,
			})
			continue
		}
		results = append(results, SkillSyncStatus{
			ID:        skill.ID,
			Name:      skill.Name,
			Version:   skill.Version,
			UpdatedAt: skill.UpdatedAt.Format("2006-01-02T15:04:05Z"),
			IsRemoved: false,
		})
	}

	response.Success(c, gin.H{"skills": results})
}

func (h *PluginHandler) ListCategories(c *gin.Context) {
	cats, err := h.categoryRepo.List()
	if err != nil {
		response.Error(c, errno.DBError)
		return
	}

	response.Success(c, gin.H{"categories": cats})
}

func (h *PluginHandler) ListSkills(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "50")
	category := c.Query("category")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 200 {
		pageSize = 50
	}

	sess := h.skillRepo.ListQuery()
	sess = sess.Where("status = ?", model.SkillStatusActive)
	if category != "" {
		sess = sess.Where("category = ?", category)
	}
	sess = sess.Desc("stars")

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
