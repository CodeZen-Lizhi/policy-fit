package handler

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhenglizhi/policy-fit/internal/domain"
	"github.com/zhenglizhi/policy-fit/internal/service"
	"github.com/zhenglizhi/policy-fit/pkg/response"
)

// RuleConfigService 规则配置服务抽象
type RuleConfigService interface {
	Publish(
		ctx context.Context,
		content map[string]interface{},
		changelog string,
		actorID int64,
	) (*domain.RuleConfigVersion, error)
	Rollback(ctx context.Context, version string, actorID int64) error
	SetGrayMode(ctx context.Context, version string, enabled bool, actorID int64) error
	GetActive(ctx context.Context) (*domain.RuleConfigVersion, error)
	ListVersions(ctx context.Context, limit int) ([]domain.RuleConfigVersion, error)
}

// RuleConfigHandler 规则配置管理处理器
type RuleConfigHandler struct {
	service RuleConfigService
	admin   *adminAuthorizer
}

// NewRuleConfigHandler 创建规则配置管理处理器
func NewRuleConfigHandler(service RuleConfigService, adminUserIDs []int64) *RuleConfigHandler {
	return &RuleConfigHandler{
		service: service,
		admin:   newAdminAuthorizer(adminUserIDs),
	}
}

// RegisterRuleConfigRoutes 注册规则管理路由
func RegisterRuleConfigRoutes(r *gin.RouterGroup, h *RuleConfigHandler) {
	adminRules := r.Group("/admin/rules")
	{
		adminRules.GET("/active", h.GetActive)
		adminRules.GET("/versions", h.ListVersions)
		adminRules.POST("/publish", h.Publish)
		adminRules.POST("/rollback", h.Rollback)
		adminRules.POST("/gray", h.SetGrayMode)
	}
}

type publishRuleRequest struct {
	Content   map[string]interface{} `json:"content"`
	Changelog string                 `json:"changelog" binding:"required,max=2000"`
}

type rollbackRuleRequest struct {
	Version string `json:"version" binding:"required,max=64"`
}

type grayRuleRequest struct {
	Version string `json:"version" binding:"required,max=64"`
	Enabled bool   `json:"enabled"`
}

// GetActive 获取当前生效规则
func (h *RuleConfigHandler) GetActive(c *gin.Context) {
	if !h.ensureAdmin(c) {
		return
	}
	cfg, err := h.service.GetActive(c.Request.Context())
	if err != nil {
		handleRuleConfigError(c, err)
		return
	}
	response.Success(c, cfg)
}

// ListVersions 查询规则版本
func (h *RuleConfigHandler) ListVersions(c *gin.Context) {
	if !h.ensureAdmin(c) {
		return
	}
	limit := 20
	if raw := c.Query("limit"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil || v <= 0 || v > 100 {
			response.ErrorWithStatus(c, http.StatusBadRequest, "PFIT-1001", "请检查输入参数")
			return
		}
		limit = v
	}
	items, err := h.service.ListVersions(c.Request.Context(), limit)
	if err != nil {
		handleRuleConfigError(c, err)
		return
	}
	response.Success(c, gin.H{"items": items})
}

// Publish 发布规则版本
func (h *RuleConfigHandler) Publish(c *gin.Context) {
	if !h.ensureAdmin(c) {
		return
	}
	var req publishRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil || len(req.Content) == 0 {
		response.ErrorWithStatus(c, http.StatusBadRequest, "PFIT-1001", "请检查输入参数")
		return
	}
	cfg, err := h.service.Publish(c.Request.Context(), req.Content, req.Changelog, currentUserID(c))
	if err != nil {
		handleRuleConfigError(c, err)
		return
	}
	response.Success(c, cfg)
}

// Rollback 回滚规则版本
func (h *RuleConfigHandler) Rollback(c *gin.Context) {
	if !h.ensureAdmin(c) {
		return
	}
	var req rollbackRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithStatus(c, http.StatusBadRequest, "PFIT-1001", "请检查输入参数")
		return
	}
	if err := h.service.Rollback(c.Request.Context(), req.Version, currentUserID(c)); err != nil {
		handleRuleConfigError(c, err)
		return
	}
	response.Success(c, gin.H{
		"version": req.Version,
		"active":  true,
	})
}

// SetGrayMode 设置灰度开关
func (h *RuleConfigHandler) SetGrayMode(c *gin.Context) {
	if !h.ensureAdmin(c) {
		return
	}
	var req grayRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithStatus(c, http.StatusBadRequest, "PFIT-1001", "请检查输入参数")
		return
	}
	if err := h.service.SetGrayMode(c.Request.Context(), req.Version, req.Enabled, currentUserID(c)); err != nil {
		handleRuleConfigError(c, err)
		return
	}
	response.Success(c, gin.H{
		"version": req.Version,
		"is_gray": req.Enabled,
	})
}

func (h *RuleConfigHandler) ensureAdmin(c *gin.Context) bool {
	if h.admin.isAdmin(currentUserID(c)) {
		return true
	}
	response.ErrorWithStatus(c, http.StatusForbidden, "PFIT-1003", "无权限访问")
	return false
}

func handleRuleConfigError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidArgument):
		response.ErrorWithStatus(c, http.StatusBadRequest, "PFIT-1001", "请检查输入参数")
	case errors.Is(err, sql.ErrNoRows):
		response.ErrorWithStatus(c, http.StatusNotFound, "PFIT-5004", "规则版本不存在")
	default:
		response.ErrorWithStatus(c, http.StatusInternalServerError, "PFIT-1005", "系统繁忙，请稍后重试")
	}
}
