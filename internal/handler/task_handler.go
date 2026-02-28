package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhenglizhi/policy-fit/pkg/response"
)

// TaskHandler 任务处理器
type TaskHandler struct {
	// TODO: 注入 service
}

// NewTaskHandler 创建任务处理器
func NewTaskHandler() *TaskHandler {
	return &TaskHandler{}
}

// RegisterTaskRoutes 注册任务路由
func RegisterTaskRoutes(r *gin.RouterGroup) {
	h := NewTaskHandler()

	tasks := r.Group("/tasks")
	{
		tasks.POST("", h.CreateTask)
		tasks.GET("/:id", h.GetTask)
		tasks.POST("/:id/documents", h.UploadDocument)
		tasks.POST("/:id/run", h.RunTask)
		tasks.GET("/:id/findings", h.GetFindings)
		tasks.DELETE("/:id", h.DeleteTask)
	}
}

// CreateTask 创建任务
func (h *TaskHandler) CreateTask(c *gin.Context) {
	// TODO: 实现
	response.Success(c, gin.H{"task_id": 1})
}

// GetTask 获取任务详情
func (h *TaskHandler) GetTask(c *gin.Context) {
	// TODO: 实现
	response.Success(c, gin.H{"status": "pending"})
}

// UploadDocument 上传文档
func (h *TaskHandler) UploadDocument(c *gin.Context) {
	// TODO: 实现
	response.Success(c, gin.H{"document_id": 1})
}

// RunTask 运行任务
func (h *TaskHandler) RunTask(c *gin.Context) {
	// TODO: 实现
	response.Success(c, gin.H{"status": "running"})
}

// GetFindings 获取风险发现
func (h *TaskHandler) GetFindings(c *gin.Context) {
	// TODO: 实现
	response.Success(c, gin.H{"findings": []interface{}{}})
}

// DeleteTask 删除任务
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	// TODO: 实现
	c.Status(http.StatusNoContent)
}
