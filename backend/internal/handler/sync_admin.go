package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hpds/skill-hub/internal/service"
	"github.com/hpds/skill-hub/pkg/errno"
	"github.com/hpds/skill-hub/pkg/response"
)

type SyncAdminHandler struct {
	syncService *service.SyncService
}

func NewSyncAdminHandler(syncService *service.SyncService) *SyncAdminHandler {
	return &SyncAdminHandler{syncService: syncService}
}

func (h *SyncAdminHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/sync/status", h.SyncStatus)
	rg.POST("/sync/trigger", h.TriggerSync)
	rg.GET("/sync/tasks", h.SyncTaskList)
	rg.GET("/stats", h.AdminStats)
}

func (h *SyncAdminHandler) SyncStatus(c *gin.Context) {
	taskIDStr := c.Query("task_id")
	if taskIDStr != "" {
		taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
		if err != nil {
			response.Error(c, errno.ParamError)
			return
		}
		task, err := h.syncService.GetSyncStatus(c.Request.Context(), taskID)
		if err != nil {
			response.Error(c, errno.InternalError)
			return
		}
		if task == nil {
			response.Error(c, errno.NotFound)
			return
		}
		response.Success(c, task)
		return
	}

	running, err := h.syncService.GetRunningTask(c.Request.Context())
	if err != nil {
		response.Error(c, errno.InternalError)
		return
	}
	if running != nil {
		response.Success(c, gin.H{"running": running})
		return
	}

	latest, err := h.syncService.GetLatestSyncTask(c.Request.Context(), "full")
	if err != nil {
		response.Error(c, errno.InternalError)
		return
	}

	result := gin.H{
		"running": nil,
		"latest":  latest,
	}

	stats, err := h.syncService.GetSkillStats(c.Request.Context())
	if err == nil {
		result["stats"] = stats
	}

	response.Success(c, result)
}

func (h *SyncAdminHandler) TriggerSync(c *gin.Context) {
	var req struct {
		Type     string `json:"type" binding:"required"`
		Strategy string `json:"strategy"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errno.ParamError)
		return
	}

	if req.Type != "full" && req.Type != "incremental" {
		response.Error(c, errno.ParamError)
		return
	}

	var task interface{}
	var err error

	switch req.Type {
	case "full":
		if req.Strategy == "" {
			req.Strategy = "topic"
		}
		task, err = h.syncService.TriggerFullSync(c.Request.Context(), req.Strategy)
	case "incremental":
		task, err = h.syncService.TriggerIncrementalSync(c.Request.Context())
	}

	if err != nil {
		response.Error(c, errno.InternalError)
		return
	}

	response.Success(c, task)
}

func (h *SyncAdminHandler) SyncTaskList(c *gin.Context) {
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

	tasks, total, err := h.syncService.ListSyncTasks(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, errno.InternalError)
		return
	}

	response.Success(c, gin.H{
		"tasks": tasks,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

func (h *SyncAdminHandler) AdminStats(c *gin.Context) {
	stats, err := h.syncService.GetSkillStats(c.Request.Context())
	if err != nil {
		response.Error(c, errno.InternalError)
		return
	}

	running, _ := h.syncService.GetRunningTask(c.Request.Context())

	response.Success(c, gin.H{
		"stats":   stats,
		"running": running != nil,
	})
}

func parsePagination(c *gin.Context) (page, pageSize int) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	page, _ = strconv.Atoi(pageStr)
	pageSize, _ = strconv.Atoi(pageSizeStr)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return
}

func parseIDParam(c *gin.Context) (int64, error) {
	idStr := c.Param("id")
	return strconv.ParseInt(idStr, 10, 64)
}
