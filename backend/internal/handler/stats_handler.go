package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hpds/skill-hub/internal/model"
	"github.com/hpds/skill-hub/internal/repository"
	"github.com/hpds/skill-hub/pkg/errno"
	"github.com/hpds/skill-hub/pkg/response"
)

type StatsHandler struct {
	skillRepo    *repository.SkillRepo
	categoryRepo *repository.CategoryRepo
	favoriteRepo *repository.FavoriteRepo
}

func NewStatsHandler(skillRepo *repository.SkillRepo, categoryRepo *repository.CategoryRepo, favoriteRepo *repository.FavoriteRepo) *StatsHandler {
	return &StatsHandler{skillRepo: skillRepo, categoryRepo: categoryRepo, favoriteRepo: favoriteRepo}
}

func (h *StatsHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/stats", h.GetStats)
	rg.GET("/stats/top-skills", h.TopSkills)
}

func (h *StatsHandler) GetStats(c *gin.Context) {
	stats, err := h.skillRepo.GetStats()
	if err != nil {
		response.Error(c, errno.DBError)
		return
	}

	categories, err := h.categoryRepo.List()
	if err != nil {
		response.Error(c, errno.DBError)
		return
	}

	categoryStats := make([]gin.H, 0)
	for _, cat := range categories {
		count, _ := h.categoryRepo.GetSkillCountByCategory(cat.Name)
		categoryStats = append(categoryStats, gin.H{
			"id":    cat.ID,
			"name":  cat.Name,
			"slug":  cat.Slug,
			"icon":  cat.Icon,
			"count": count,
		})
	}

	response.Success(c, gin.H{
		"total_skills":   stats.TotalSkills,
		"active_skills":  stats.ActiveSkills,
		"total_stars":    stats.TotalStars,
		"total_installs": stats.TotalInstalls,
		"categories":     categoryStats,
	})
}

func (h *StatsHandler) TopSkills(c *gin.Context) {
	sort := c.DefaultQuery("sort", "stars")
	limitStr := c.DefaultQuery("limit", "10")
	limit := 10
	if n, err := strconv.Atoi(limitStr); err == nil && n > 0 && n <= 50 {
		limit = n
	}

	sess := h.skillRepo.ListQuery().Where("status = ?", model.SkillStatusActive)

	switch sort {
	case "installs":
		sess = sess.Desc("installs")
	case "score":
		sess = sess.Desc("score")
	default:
		sess = sess.Desc("stars")
	}

	var skills []*model.Skill
	if err := sess.Limit(limit, 0).Find(&skills); err != nil {
		response.Error(c, errno.DBError)
		return
	}

	response.Success(c, gin.H{
		"sort":   sort,
		"skills": skills,
	})
}
