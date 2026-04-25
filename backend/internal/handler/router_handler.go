package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hpds/skill-hub/internal/service"
	"github.com/hpds/skill-hub/pkg/errno"
	"github.com/hpds/skill-hub/pkg/response"
)

type RouterHandler struct {
	svc *service.RouterService
}

func NewRouterHandler(svc *service.RouterService) *RouterHandler {
	return &RouterHandler{svc: svc}
}

func (h *RouterHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/router/match", h.Match)
	rg.POST("/router/execute", h.Execute)
	rg.POST("/router/feedback", h.Feedback)
	rg.GET("/router/logs", h.ListLogs)
}

func (h *RouterHandler) Match(c *gin.Context) {
	var req service.MatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errno.ParamError)
		return
	}
	if req.Query == "" {
		response.Error(c, errno.ParamError)
		return
	}

	userID, _ := c.Get("user_id")
	if id, ok := userID.(int64); ok {
		req.UserID = id
	}

	result, err := h.svc.Match(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, errno.InternalError)
		return
	}

	response.Success(c, result)
}

func (h *RouterHandler) Execute(c *gin.Context) {
	var req service.ExecuteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errno.ParamError)
		return
	}
	if req.Query == "" || req.SkillID <= 0 {
		response.Error(c, errno.ParamError)
		return
	}

	userID, _ := c.Get("user_id")
	if id, ok := userID.(int64); ok {
		req.UserID = id
	}

	result, err := h.svc.Execute(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, errno.InternalError)
		return
	}

	response.Success(c, result)
}

func (h *RouterHandler) Feedback(c *gin.Context) {
	var req service.FeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errno.ParamError)
		return
	}
	if req.LogID <= 0 || req.Score < 1 || req.Score > 5 {
		response.Error(c, errno.ParamError)
		return
	}

	if err := h.svc.SubmitFeedback(c.Request.Context(), &req); err != nil {
		response.Error(c, errno.InternalError)
		return
	}

	response.Success(c, gin.H{"message": "feedback submitted"})
}

func (h *RouterHandler) ListLogs(c *gin.Context) {
	logs, err := h.svc.ListLogs(c.Request.Context())
	if err != nil {
		response.Error(c, errno.InternalError)
		return
	}
	response.Success(c, gin.H{"logs": logs})
}
