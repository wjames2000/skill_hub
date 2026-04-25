package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hpds/skill-hub/internal/model"
	"github.com/hpds/skill-hub/internal/repository"
	"github.com/hpds/skill-hub/pkg/errno"
	"github.com/hpds/skill-hub/pkg/logger"
	"github.com/hpds/skill-hub/pkg/response"
)

type AdminHandler struct {
	skillRepo    *repository.SkillRepo
	categoryRepo *repository.CategoryRepo
	userRepo     *repository.UserRepo
	favoriteRepo *repository.FavoriteRepo
	reviewRepo   *repository.ReviewRepo
}

func NewAdminHandler(
	skillRepo *repository.SkillRepo,
	categoryRepo *repository.CategoryRepo,
	userRepo *repository.UserRepo,
	favoriteRepo *repository.FavoriteRepo,
	reviewRepo *repository.ReviewRepo,
) *AdminHandler {
	return &AdminHandler{
		skillRepo:    skillRepo,
		categoryRepo: categoryRepo,
		userRepo:     userRepo,
		favoriteRepo: favoriteRepo,
		reviewRepo:   reviewRepo,
	}
}

func (h *AdminHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/admin/skills", h.ListSkills)
	rg.PUT("/admin/skills/:id/status", h.UpdateSkillStatus)
	rg.DELETE("/admin/skills/:id", h.DeleteSkill)
	rg.GET("/admin/reviews", h.ListAllReviews)
	rg.DELETE("/admin/reviews/:id", h.DeleteReview)
	rg.GET("/admin/users", h.ListUsers)
	rg.PUT("/admin/users/:id/role", h.UpdateUserRole)
	rg.GET("/admin/categories", h.ListCategories)
	rg.POST("/admin/categories", h.CreateCategory)
	rg.PUT("/admin/categories/:id", h.UpdateCategory)
	rg.DELETE("/admin/categories/:id", h.DeleteCategory)
	rg.GET("/admin/stats/dashboard", h.Dashboard)
	rg.GET("/admin/pending-review", h.ListPendingReviews)
	rg.PUT("/admin/skills/:id/approve", h.ApproveSkill)
	rg.PUT("/admin/skills/:id/reject", h.RejectSkill)
	rg.GET("/admin/logs", h.GetSystemLogs)
}

func (h *AdminHandler) ListSkills(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")
	statusStr := c.DefaultQuery("status", "-1")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	status, _ := strconv.Atoi(statusStr)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	sess := h.skillRepo.ListQuery()
	if status >= 0 {
		sess = sess.Where("status = ?", status)
	}
	sess = sess.Desc("created_at")

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

func (h *AdminHandler) UpdateSkillStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, errno.ParamError)
		return
	}

	var req struct {
		Status int `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errno.ParamError)
		return
	}

	if req.Status < 0 || req.Status > 3 {
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

	if err := h.skillRepo.SetStatus(id, req.Status); err != nil {
		response.Error(c, errno.DBError)
		return
	}

	response.Success(c, gin.H{"message": "status updated"})
}

func (h *AdminHandler) DeleteSkill(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, errno.ParamError)
		return
	}

	if err := h.skillRepo.Delete(id); err != nil {
		response.Error(c, errno.DBError)
		return
	}

	response.Success(c, gin.H{"message": "skill deleted"})
}

func (h *AdminHandler) ListAllReviews(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")
	skillIDStr := c.Query("skill_id")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var reviews []*model.Review
	var total int64
	var err error

	if skillIDStr != "" {
		skillID, _ := strconv.ParseInt(skillIDStr, 10, 64)
		reviews, total, err = h.reviewRepo.ListBySkill(skillID, page, pageSize)
	} else {
		reviews, total, err = h.reviewRepo.ListAll(page, pageSize)
	}
	if err != nil {
		response.Error(c, errno.DBError)
		return
	}

	response.Success(c, gin.H{
		"reviews": reviews,
		"total":   total,
		"page":    page,
		"size":    pageSize,
	})
}

func (h *AdminHandler) DeleteReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, errno.ParamError)
		return
	}

	review, err := h.reviewRepo.GetByID(id)
	if err != nil {
		response.Error(c, errno.DBError)
		return
	}
	if review == nil {
		response.Error(c, errno.NotFound)
		return
	}

	if err := h.reviewRepo.Delete(id); err != nil {
		response.Error(c, errno.DBError)
		return
	}

	avg, _ := h.reviewRepo.GetAvgScoreBySkill(review.SkillID)
	_ = h.skillRepo.UpdateScore(review.SkillID, avg)

	response.Success(c, gin.H{"message": "review deleted"})
}

func (h *AdminHandler) ListUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	users, total, err := h.userRepo.List(page, pageSize)
	if err != nil {
		response.Error(c, errno.DBError)
		return
	}

	response.Success(c, gin.H{
		"users": users,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

func (h *AdminHandler) UpdateUserRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, errno.ParamError)
		return
	}

	var req struct {
		Role string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errno.ParamError)
		return
	}

	if req.Role != "user" && req.Role != "admin" {
		response.Error(c, errno.ParamError)
		return
	}

	user, err := h.userRepo.GetByID(id)
	if err != nil {
		response.Error(c, errno.DBError)
		return
	}
	if user == nil {
		response.Error(c, errno.UserNotFound)
		return
	}

	user.Role = req.Role
	if err := h.userRepo.Update(user); err != nil {
		response.Error(c, errno.DBError)
		return
	}

	response.Success(c, gin.H{"message": "role updated"})
}

func (h *AdminHandler) ListCategories(c *gin.Context) {
	cats, err := h.categoryRepo.List()
	if err != nil {
		response.Error(c, errno.DBError)
		return
	}

	response.Success(c, gin.H{"categories": cats})
}

func (h *AdminHandler) CreateCategory(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Slug        string `json:"slug" binding:"required"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
		ParentID    int64  `json:"parent_id"`
		SortOrder   int    `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errno.ParamError)
		return
	}

	cat := &model.SkillCategory{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Icon:        req.Icon,
		ParentID:    req.ParentID,
		SortOrder:   req.SortOrder,
	}

	if err := h.categoryRepo.Create(cat); err != nil {
		response.Error(c, errno.DBError)
		return
	}

	response.Success(c, cat)
}

func (h *AdminHandler) UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, errno.ParamError)
		return
	}

	cat, err := h.categoryRepo.GetByID(id)
	if err != nil {
		response.Error(c, errno.DBError)
		return
	}
	if cat == nil {
		response.Error(c, errno.CategoryNotFound)
		return
	}

	var req struct {
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
		ParentID    int64  `json:"parent_id"`
		SortOrder   int    `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errno.ParamError)
		return
	}

	if req.Name != "" {
		cat.Name = req.Name
	}
	if req.Slug != "" {
		cat.Slug = req.Slug
	}
	if req.Description != "" {
		cat.Description = req.Description
	}
	if req.Icon != "" {
		cat.Icon = req.Icon
	}
	if req.SortOrder != 0 {
		cat.SortOrder = req.SortOrder
	}
	cat.ParentID = req.ParentID

	if err := h.categoryRepo.Update(cat); err != nil {
		response.Error(c, errno.DBError)
		return
	}

	response.Success(c, cat)
}

func (h *AdminHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, errno.ParamError)
		return
	}

	if err := h.categoryRepo.Delete(id); err != nil {
		response.Error(c, errno.DBError)
		return
	}

	response.Success(c, gin.H{"message": "category deleted"})
}

func (h *AdminHandler) Dashboard(c *gin.Context) {
	stats, err := h.skillRepo.GetStats()
	if err != nil {
		response.Error(c, errno.DBError)
		return
	}

	pendingCount, _ := h.skillRepo.CountByStatus(model.SkillStatusPending)
	activeCount, _ := h.skillRepo.CountByStatus(model.SkillStatusActive)
	disabledCount, _ := h.skillRepo.CountByStatus(model.SkillStatusDisabled)
	deprecatedCount, _ := h.skillRepo.CountByStatus(model.SkillStatusDeprecated)

	response.Success(c, gin.H{
		"total_skills":      stats.TotalSkills,
		"active_skills":     activeCount,
		"pending_skills":    pendingCount,
		"disabled_skills":   disabledCount,
		"deprecated_skills": deprecatedCount,
		"total_stars":       stats.TotalStars,
		"total_installs":    stats.TotalInstalls,
	})
}

// ListPendingReviews lists skills with status = 0 (pending).
func (h *AdminHandler) ListPendingReviews(c *gin.Context) {
	var skills []*model.Skill
	if err := h.skillRepo.ListQuery().
		Where("status = ?", model.SkillStatusPending).
		Desc("created_at").
		Find(&skills); err != nil {
		response.Error(c, errno.DBError)
		return
	}
	response.Success(c, skills)
}

// ApproveSkill sets a skill's status to 1 (active).
func (h *AdminHandler) ApproveSkill(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
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

	if err := h.skillRepo.SetStatus(id, model.SkillStatusActive); err != nil {
		response.Error(c, errno.DBError)
		return
	}
	response.Success(c, gin.H{"message": "skill approved"})
}

// RejectSkill sets a skill's status to 2 (disabled), with an optional reason.
func (h *AdminHandler) RejectSkill(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
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

	var req struct {
		Reason string `json:"reason"`
	}
	_ = c.ShouldBindJSON(&req) // reason is optional

	if err := h.skillRepo.SetStatus(id, model.SkillStatusDisabled); err != nil {
		response.Error(c, errno.DBError)
		return
	}
	response.Success(c, gin.H{"message": "skill rejected", "reason": req.Reason})
}

// GetSystemLogs returns recent system logs from the in-memory ring buffer.
// Accepts a "lines" query parameter (default 50, max 1000) to control the count.
func (h *AdminHandler) GetSystemLogs(c *gin.Context) {
	linesStr := c.DefaultQuery("lines", "50")
	lines, err := strconv.Atoi(linesStr)
	if err != nil || lines <= 0 {
		lines = 50
	}
	if lines > 1000 {
		lines = 1000
	}

	entries := logger.GetRecentLogs(lines)
	// Ensure we return an empty JSON array (not null) when no logs exist.
	if entries == nil {
		entries = []logger.LogEntry{}
	}
	response.Success(c, entries)
}
