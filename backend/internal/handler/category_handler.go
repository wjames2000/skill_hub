package handler

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hpds/skill-hub/internal/model"
	"github.com/hpds/skill-hub/internal/repository"
	"github.com/hpds/skill-hub/pkg/errno"
	"github.com/hpds/skill-hub/pkg/response"
)

type CategoryHandler struct {
	categoryRepo *repository.CategoryRepo
}

func NewCategoryHandler(categoryRepo *repository.CategoryRepo) *CategoryHandler {
	return &CategoryHandler{categoryRepo: categoryRepo}
}

func (h *CategoryHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/categories", h.ListCategories)
	rg.GET("/categories/:id", h.GetCategory)
	rg.POST("/categories", h.CreateCategory)
	rg.PUT("/categories/:id", h.UpdateCategory)
	rg.DELETE("/categories/:id", h.DeleteCategory)
}

func (h *CategoryHandler) ListCategories(c *gin.Context) {
	cats, err := h.categoryRepo.List()
	if err != nil {
		response.Error(c, errno.DBError)
		return
	}

	response.Success(c, gin.H{"categories": cats})
}

func (h *CategoryHandler) GetCategory(c *gin.Context) {
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

	response.Success(c, cat)
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		ZhName      string `json:"zh_name"`
		EnName      string `json:"en_name"`
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

	req.Slug = strings.ToLower(req.Slug)

	existing, _ := h.categoryRepo.GetBySlug(req.Slug)
	if existing != nil {
		response.Error(c, errno.CategoryExists)
		return
	}

	cat := &model.SkillCategory{
		Name:        req.Name,
		ZhName:      req.ZhName,
		EnName:      req.EnName,
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

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
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
		ZhName      string `json:"zh_name"`
		EnName      string `json:"en_name"`
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
	if req.ZhName != "" {
		cat.ZhName = req.ZhName
	}
	if req.EnName != "" {
		cat.EnName = req.EnName
	}
	if req.Slug != "" {
		req.Slug = strings.ToLower(req.Slug)
		existing, _ := h.categoryRepo.GetBySlug(req.Slug)
		if existing != nil && existing.ID != id {
			response.Error(c, errno.CategoryExists)
			return
		}
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

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
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

	if err := h.categoryRepo.Delete(id); err != nil {
		response.Error(c, errno.DBError)
		return
	}

	response.Success(c, gin.H{"message": "category deleted"})
}
