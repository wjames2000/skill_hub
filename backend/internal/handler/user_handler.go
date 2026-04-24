package handler

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hpds/skill-hub/internal/model"
	"github.com/hpds/skill-hub/internal/repository"
	"github.com/hpds/skill-hub/pkg/errno"
	"github.com/hpds/skill-hub/pkg/logger"
	"github.com/hpds/skill-hub/pkg/response"
)

type UserHandler struct {
	userRepo     *repository.UserRepo
	favoriteRepo *repository.FavoriteRepo
	reviewRepo   *repository.ReviewRepo
	apiKeyRepo   *repository.APIKeyRepo
	skillRepo    *repository.SkillRepo
}

func NewUserHandler(
	userRepo *repository.UserRepo,
	favoriteRepo *repository.FavoriteRepo,
	reviewRepo *repository.ReviewRepo,
	apiKeyRepo *repository.APIKeyRepo,
	skillRepo *repository.SkillRepo,
) *UserHandler {
	return &UserHandler{
		userRepo:     userRepo,
		favoriteRepo: favoriteRepo,
		reviewRepo:   reviewRepo,
		apiKeyRepo:   apiKeyRepo,
		skillRepo:    skillRepo,
	}
}

func (h *UserHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/user/profile", h.GetProfile)
	rg.PUT("/user/profile", h.UpdateProfile)

	rg.GET("/user/favorites", h.ListFavorites)
	rg.POST("/user/favorites", h.AddFavorite)
	rg.DELETE("/user/favorites/:skill_id", h.RemoveFavorite)
	rg.GET("/user/favorites/check/:skill_id", h.CheckFavorite)

	rg.GET("/user/reviews", h.ListMyReviews)
	rg.POST("/user/reviews", h.AddReview)
	rg.PUT("/user/reviews/:id", h.UpdateReview)
	rg.DELETE("/user/reviews/:id", h.DeleteReview)

	rg.GET("/user/api-keys", h.ListAPIKeys)
	rg.POST("/user/api-keys", h.CreateAPIKey)
	rg.DELETE("/user/api-keys/:id", h.RevokeAPIKey)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, ok := userID.(int64)
	if !ok || id <= 0 {
		response.Error(c, errno.Unauthorized)
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

	favCount, _ := h.favoriteRepo.CountByUser(id)

	response.Success(c, gin.H{
		"id":             user.ID,
		"username":       user.Username,
		"email":          user.Email,
		"avatar_url":     user.AvatarURL,
		"bio":            user.Bio,
		"role":           user.Role,
		"github_id":      user.GitHubID,
		"favorite_count": favCount,
		"created_at":     user.CreatedAt,
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, ok := userID.(int64)
	if !ok || id <= 0 {
		response.Error(c, errno.Unauthorized)
		return
	}

	var req struct {
		AvatarURL string `json:"avatar_url"`
		Bio       string `json:"bio"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errno.ParamError)
		return
	}

	user, err := h.userRepo.GetByID(id)
	if err != nil || user == nil {
		response.Error(c, errno.UserNotFound)
		return
	}

	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	user.Bio = req.Bio

	if err := h.userRepo.Update(user); err != nil {
		response.Error(c, errno.DBError)
		return
	}

	response.Success(c, gin.H{"message": "profile updated"})
}

func (h *UserHandler) AddFavorite(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, ok := userID.(int64)
	if !ok || uid <= 0 {
		response.Error(c, errno.Unauthorized)
		return
	}

	var req struct {
		SkillID int64 `json:"skill_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errno.ParamError)
		return
	}

	exists, err := h.favoriteRepo.Exists(uid, req.SkillID)
	if err != nil {
		response.Error(c, errno.DBError)
		return
	}
	if exists {
		response.Error(c, errno.AlreadyFavorited)
		return
	}

	fav := &model.Favorite{
		UserID:  uid,
		SkillID: req.SkillID,
	}
	if err := h.favoriteRepo.Create(fav); err != nil {
		response.Error(c, errno.DBError)
		return
	}

	response.Success(c, gin.H{"message": "favorited"})
}

func (h *UserHandler) RemoveFavorite(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, ok := userID.(int64)
	if !ok || uid <= 0 {
		response.Error(c, errno.Unauthorized)
		return
	}

	skillIDStr := c.Param("skill_id")
	skillID, err := strconv.ParseInt(skillIDStr, 10, 64)
	if err != nil {
		response.Error(c, errno.ParamError)
		return
	}

	if err := h.favoriteRepo.Delete(uid, skillID); err != nil {
		response.Error(c, errno.DBError)
		return
	}

	response.Success(c, gin.H{"message": "unfavorited"})
}

func (h *UserHandler) ListFavorites(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, ok := userID.(int64)
	if !ok || uid <= 0 {
		response.Error(c, errno.Unauthorized)
		return
	}

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

	favs, total, err := h.favoriteRepo.ListByUser(uid, page, pageSize)
	if err != nil {
		response.Error(c, errno.DBError)
		return
	}

	var skills []*model.Skill
	for _, fav := range favs {
		skill, err := h.skillRepo.GetByID(fav.SkillID)
		if err != nil || skill == nil {
			continue
		}
		skills = append(skills, skill)
	}

	response.Success(c, gin.H{
		"skills": skills,
		"total":  total,
		"page":   page,
		"size":   pageSize,
	})
}

func (h *UserHandler) CheckFavorite(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, ok := userID.(int64)
	if !ok || uid <= 0 {
		response.Success(c, gin.H{"favorited": false})
		return
	}

	skillIDStr := c.Param("skill_id")
	skillID, err := strconv.ParseInt(skillIDStr, 10, 64)
	if err != nil {
		response.Error(c, errno.ParamError)
		return
	}

	exists, err := h.favoriteRepo.Exists(uid, skillID)
	if err != nil {
		response.Error(c, errno.DBError)
		return
	}

	response.Success(c, gin.H{"favorited": exists})
}

func (h *UserHandler) AddReview(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, ok := userID.(int64)
	if !ok || uid <= 0 {
		response.Error(c, errno.Unauthorized)
		return
	}

	var req struct {
		SkillID int64  `json:"skill_id" binding:"required"`
		Score   int    `json:"score" binding:"required"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errno.ParamError)
		return
	}
	if req.Score < 1 || req.Score > 5 {
		response.Error(c, errno.ParamError)
		return
	}

	skill, err := h.skillRepo.GetByID(req.SkillID)
	if err != nil {
		response.Error(c, errno.DBError)
		return
	}
	if skill == nil {
		response.Error(c, errno.NotFound)
		return
	}

	review := &model.Review{
		UserID:  uid,
		SkillID: req.SkillID,
		Score:   req.Score,
		Content: req.Content,
	}
	if err := h.reviewRepo.Create(review); err != nil {
		response.Error(c, errno.DBError)
		return
	}

	avg, _ := h.reviewRepo.GetAvgScoreBySkill(req.SkillID)
	count, _ := h.reviewRepo.CountBySkill(req.SkillID)
	_ = h.skillRepo.UpdateScore(req.SkillID, avg)

	response.Success(c, gin.H{
		"review":       review,
		"avg_score":    avg,
		"review_count": count,
	})
}

func (h *UserHandler) UpdateReview(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, ok := userID.(int64)
	if !ok || uid <= 0 {
		response.Error(c, errno.Unauthorized)
		return
	}

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
	if review == nil || review.UserID != uid {
		response.Error(c, errno.Forbidden)
		return
	}

	var req struct {
		Score   int    `json:"score"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errno.ParamError)
		return
	}
	if req.Score < 1 || req.Score > 5 {
		response.Error(c, errno.ParamError)
		return
	}

	review.Score = req.Score
	review.Content = req.Content
	if err := h.reviewRepo.Update(review); err != nil {
		response.Error(c, errno.DBError)
		return
	}

	avg, _ := h.reviewRepo.GetAvgScoreBySkill(review.SkillID)
	_ = h.skillRepo.UpdateScore(review.SkillID, avg)

	response.Success(c, gin.H{"review": review})
}

func (h *UserHandler) DeleteReview(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, ok := userID.(int64)
	if !ok || uid <= 0 {
		response.Error(c, errno.Unauthorized)
		return
	}

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
	if review == nil || review.UserID != uid {
		response.Error(c, errno.Forbidden)
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

func (h *UserHandler) ListMyReviews(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, ok := userID.(int64)
	if !ok || uid <= 0 {
		response.Error(c, errno.Unauthorized)
		return
	}

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

	reviews, total, err := h.reviewRepo.ListByUser(uid, page, pageSize)
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

func (h *UserHandler) ListAPIKeys(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, ok := userID.(int64)
	if !ok || uid <= 0 {
		response.Error(c, errno.Unauthorized)
		return
	}

	keys, err := h.apiKeyRepo.ListByUser(uid)
	if err != nil {
		response.Error(c, errno.DBError)
		return
	}

	response.Success(c, gin.H{"api_keys": keys})
}

func (h *UserHandler) CreateAPIKey(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, ok := userID.(int64)
	if !ok || uid <= 0 {
		response.Error(c, errno.Unauthorized)
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errno.ParamError)
		return
	}

	keyBytes := make([]byte, 32)
	if _, err := rand.Read(keyBytes); err != nil {
		logger.Error("generate api key", logger.ErrorField(err))
		response.Error(c, errno.InternalError)
		return
	}
	keyStr := "sk_" + hex.EncodeToString(keyBytes)

	expiresAt := time.Now().Add(365 * 24 * time.Hour)

	apiKey := &model.APIKey{
		UserID:    uid,
		Name:      req.Name,
		Key:       keyStr,
		ExpiresAt: &expiresAt,
	}

	if err := h.apiKeyRepo.Create(apiKey); err != nil {
		response.Error(c, errno.DBError)
		return
	}

	response.Success(c, gin.H{
		"api_key": apiKey,
		"raw_key": keyStr,
	})
}

func (h *UserHandler) RevokeAPIKey(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, ok := userID.(int64)
	if !ok || uid <= 0 {
		response.Error(c, errno.Unauthorized)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, errno.ParamError)
		return
	}

	key, err := h.apiKeyRepo.GetByID(id)
	if err != nil {
		response.Error(c, errno.DBError)
		return
	}
	if key == nil || key.UserID != uid {
		response.Error(c, errno.APIKeyNotFound)
		return
	}

	if err := h.apiKeyRepo.Revoke(id); err != nil {
		response.Error(c, errno.DBError)
		return
	}

	response.Success(c, gin.H{"message": "api key revoked"})
}
