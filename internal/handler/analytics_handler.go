package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhenglizhi/policy-fit/internal/service"
	"github.com/zhenglizhi/policy-fit/pkg/response"
)

// AnalyticsQueryService 埋点服务抽象
type AnalyticsQueryService interface {
	Track(
		ctx context.Context,
		eventName string,
		userID *int64,
		taskID *int64,
		properties map[string]interface{},
	) error
	FunnelByPeriod(ctx context.Context, period string) (map[string]int64, error)
	Dashboard(ctx context.Context, period string) (map[string]interface{}, error)
}

// AnalyticsHandler 埋点处理器
type AnalyticsHandler struct {
	service AnalyticsQueryService
	admin   *adminAuthorizer
}

// NewAnalyticsHandler 创建埋点处理器
func NewAnalyticsHandler(service AnalyticsQueryService, adminUserIDs []int64) *AnalyticsHandler {
	return &AnalyticsHandler{
		service: service,
		admin:   newAdminAuthorizer(adminUserIDs),
	}
}

// RegisterAnalyticsRoutes 注册埋点路由
func RegisterAnalyticsRoutes(r *gin.RouterGroup, h *AnalyticsHandler) {
	analytics := r.Group("/analytics")
	{
		analytics.POST("/events", h.TrackEvent)
		analytics.GET("/funnel", h.GetFunnel)
		analytics.GET("/overview", h.GetOverview)
	}
}

type trackEventRequest struct {
	EventName  string                 `json:"event_name" binding:"required,max=64"`
	TaskID     *int64                 `json:"task_id"`
	Properties map[string]interface{} `json:"properties"`
}

// TrackEvent 记录埋点事件
func (h *AnalyticsHandler) TrackEvent(c *gin.Context) {
	var req trackEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithStatus(c, http.StatusBadRequest, "PFIT-1001", "请检查输入参数")
		return
	}
	if req.TaskID != nil && *req.TaskID <= 0 {
		response.ErrorWithStatus(c, http.StatusBadRequest, "PFIT-1001", "请检查输入参数")
		return
	}

	userID := currentUserID(c)
	if userID <= 0 {
		response.ErrorWithStatus(c, http.StatusUnauthorized, "PFIT-1002", "请重新登录")
		return
	}
	uid := userID
	if err := h.service.Track(c.Request.Context(), req.EventName, &uid, req.TaskID, req.Properties); err != nil {
		if errors.Is(err, service.ErrInvalidArgument) {
			response.ErrorWithStatus(c, http.StatusBadRequest, "PFIT-1001", "请检查输入参数")
			return
		}
		response.ErrorWithStatus(c, http.StatusInternalServerError, "PFIT-1005", "系统繁忙，请稍后重试")
		return
	}
	response.Success(c, gin.H{"tracked": true})
}

// GetFunnel 查询漏斗
func (h *AnalyticsHandler) GetFunnel(c *gin.Context) {
	if !h.ensureAdmin(c) {
		return
	}
	period := c.DefaultQuery("period", "week")
	data, err := h.service.FunnelByPeriod(c.Request.Context(), period)
	if err != nil {
		if errors.Is(err, service.ErrInvalidArgument) {
			response.ErrorWithStatus(c, http.StatusBadRequest, "PFIT-1001", "请检查输入参数")
			return
		}
		response.ErrorWithStatus(c, http.StatusInternalServerError, "PFIT-1005", "系统繁忙，请稍后重试")
		return
	}
	response.Success(c, gin.H{
		"period": period,
		"funnel": data,
	})
}

// GetOverview 查询核心指标看板
func (h *AnalyticsHandler) GetOverview(c *gin.Context) {
	if !h.ensureAdmin(c) {
		return
	}
	period := c.DefaultQuery("period", "week")
	data, err := h.service.Dashboard(c.Request.Context(), period)
	if err != nil {
		if errors.Is(err, service.ErrInvalidArgument) {
			response.ErrorWithStatus(c, http.StatusBadRequest, "PFIT-1001", "请检查输入参数")
			return
		}
		response.ErrorWithStatus(c, http.StatusInternalServerError, "PFIT-1005", "系统繁忙，请稍后重试")
		return
	}
	response.Success(c, data)
}

func (h *AnalyticsHandler) ensureAdmin(c *gin.Context) bool {
	if h.admin.isAdmin(currentUserID(c)) {
		return true
	}
	response.ErrorWithStatus(c, http.StatusForbidden, "PFIT-1003", "无权限访问")
	return false
}
