package handler

import (
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhenglizhi/policy-fit/internal/domain"
	"github.com/zhenglizhi/policy-fit/internal/middleware"
	"github.com/zhenglizhi/policy-fit/internal/service"
	"github.com/zhenglizhi/policy-fit/pkg/response"
)

// TaskService 任务服务抽象
type TaskService interface {
	CreateTask(ctx context.Context, userID int64, requestID string) (*domain.AnalysisTask, error)
	GetTask(ctx context.Context, taskID int64, actorID int64) (*domain.AnalysisTask, error)
	RunTask(ctx context.Context, taskID int64, actorID int64) error
	DeleteTask(ctx context.Context, taskID int64, actorID int64) error
}

// DocumentService 文档服务抽象
type DocumentService interface {
	CreateDocument(
		ctx context.Context,
		taskID int64,
		actorID int64,
		docType domain.DocumentType,
		fileName string,
		storageKey string,
	) (*domain.Document, error)
}

// FindingService 风险服务抽象
type FindingService interface {
	ListByTask(ctx context.Context, taskID int64) ([]domain.RiskFinding, error)
	Summarize(findings []domain.RiskFinding) map[string]int
}

// StorageService 存储服务抽象
type StorageService interface {
	SaveUploadedFile(
		ctx context.Context,
		taskID int64,
		docType domain.DocumentType,
		fileHeader *multipart.FileHeader,
	) (string, error)
}

// ExportService 导出服务抽象
type ExportService interface {
	ExportMarkdown(ctx context.Context, taskID int64, actorID int64, lang string) ([]byte, error)
	ExportPDF(ctx context.Context, taskID int64, actorID int64, lang string) ([]byte, error)
}

// ComparisonService 对比服务抽象
type ComparisonService interface {
	CompareTasks(ctx context.Context, fromTaskID int64, toTaskID int64, actorID int64) (*service.ComparisonResult, error)
}

// AnalyticsService 埋点服务抽象
type AnalyticsService interface {
	Track(
		ctx context.Context,
		eventName string,
		userID *int64,
		taskID *int64,
		properties map[string]interface{},
	) error
}

// TaskHandler 任务处理器
type TaskHandler struct {
	taskSvc     TaskService
	documentSvc DocumentService
	findingSvc  FindingService
	storageSvc  StorageService
	exportSvc   ExportService
	compareSvc  ComparisonService
	analytics   AnalyticsService
}

// NewTaskHandler 创建任务处理器
func NewTaskHandler(
	taskSvc TaskService,
	documentSvc DocumentService,
	findingSvc FindingService,
	storageSvc StorageService,
	exportSvc ExportService,
	compareSvc ComparisonService,
	analytics AnalyticsService,
) *TaskHandler {
	return &TaskHandler{
		taskSvc:     taskSvc,
		documentSvc: documentSvc,
		findingSvc:  findingSvc,
		storageSvc:  storageSvc,
		exportSvc:   exportSvc,
		compareSvc:  compareSvc,
		analytics:   analytics,
	}
}

// RegisterTaskRoutes 注册任务路由
func RegisterTaskRoutes(r *gin.RouterGroup, h *TaskHandler) {
	tasks := r.Group("/tasks")
	{
		tasks.POST("", h.CreateTask)
		tasks.GET("/:id", h.GetTask)
		tasks.POST("/:id/documents", h.UploadDocument)
		tasks.POST("/:id/run", h.RunTask)
		tasks.GET("/:id/findings", h.GetFindings)
		tasks.GET("/:id/export", h.ExportReport)
		tasks.GET("/compare", h.CompareReports)
		tasks.DELETE("/:id", h.DeleteTask)
	}
}

type createTaskRequest struct {
	RequestID string `json:"request_id" binding:"omitempty,max=64"`
}

// CreateTask 创建任务
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req createTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithStatus(c, http.StatusBadRequest, "PFIT-1001", "请检查输入参数")
		return
	}

	task, err := h.taskSvc.CreateTask(c.Request.Context(), currentUserID(c), req.RequestID)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, gin.H{
		"task_id": task.ID,
		"status":  task.Status,
	})
	h.trackEvent(c, "task_created", &task.ID, map[string]interface{}{"status": task.Status})
}

// GetTask 获取任务详情
func (h *TaskHandler) GetTask(c *gin.Context) {
	taskID, ok := parseTaskID(c)
	if !ok {
		return
	}

	task, err := h.taskSvc.GetTask(c.Request.Context(), taskID, currentUserID(c))
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, task)
}

// UploadDocument 上传文档
func (h *TaskHandler) UploadDocument(c *gin.Context) {
	taskID, ok := parseTaskID(c)
	if !ok {
		return
	}

	docType := domain.DocumentType(c.PostForm("docType"))
	switch docType {
	case domain.DocTypeReport, domain.DocTypePolicy, domain.DocTypeDisclosure:
	default:
		response.ErrorWithStatus(c, http.StatusBadRequest, "PFIT-1001", "请检查输入参数")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		response.ErrorWithStatus(c, http.StatusBadRequest, "PFIT-1001", "请检查输入参数")
		return
	}

	storageKey, err := h.storageSvc.SaveUploadedFile(c.Request.Context(), taskID, docType, file)
	if err != nil {
		handleError(c, err)
		return
	}

	doc, err := h.documentSvc.CreateDocument(c.Request.Context(), taskID, currentUserID(c), docType, file.Filename, storageKey)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, gin.H{
		"document_id": doc.ID,
		"storage_key": doc.StorageKey,
	})
	h.trackEvent(c, "document_uploaded", &taskID, map[string]interface{}{"doc_type": docType})
}

// RunTask 运行任务
func (h *TaskHandler) RunTask(c *gin.Context) {
	taskID, ok := parseTaskID(c)
	if !ok {
		return
	}

	if err := h.taskSvc.RunTask(c.Request.Context(), taskID, currentUserID(c)); err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, gin.H{
		"status": "pending",
	})
	h.trackEvent(c, "task_run", &taskID, map[string]interface{}{})
}

// GetFindings 获取风险发现
func (h *TaskHandler) GetFindings(c *gin.Context) {
	taskID, ok := parseTaskID(c)
	if !ok {
		return
	}

	findings, err := h.findingSvc.ListByTask(c.Request.Context(), taskID)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, gin.H{
		"summary":           h.findingSvc.Summarize(findings),
		"findings":          findings,
		"retrieval_sources": buildRetrievalSourceMap(findings),
	})
	h.trackEvent(c, "report_viewed", &taskID, map[string]interface{}{})
}

// ExportReport 导出报告
func (h *TaskHandler) ExportReport(c *gin.Context) {
	taskID, ok := parseTaskID(c)
	if !ok {
		return
	}

	format := c.DefaultQuery("format", "md")
	lang := c.DefaultQuery("lang", "zh-CN")
	var (
		content     []byte
		contentType string
		fileName    string
		err         error
	)
	switch format {
	case "md":
		content, err = h.exportSvc.ExportMarkdown(c.Request.Context(), taskID, currentUserID(c), lang)
		contentType = "text/markdown; charset=utf-8"
		fileName = "report.md"
	case "pdf":
		content, err = h.exportSvc.ExportPDF(c.Request.Context(), taskID, currentUserID(c), lang)
		contentType = "application/pdf"
		fileName = "report.pdf"
	default:
		response.ErrorWithStatus(c, http.StatusBadRequest, "PFIT-1001", "请检查输入参数")
		return
	}
	if err != nil {
		handleError(c, err)
		return
	}

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Data(http.StatusOK, contentType, content)
	h.trackEvent(c, "report_exported", &taskID, map[string]interface{}{"format": format})
}

// CompareReports 对比历史报告
func (h *TaskHandler) CompareReports(c *gin.Context) {
	fromTaskID, err := strconv.ParseInt(c.Query("from"), 10, 64)
	if err != nil || fromTaskID <= 0 {
		response.ErrorWithStatus(c, http.StatusBadRequest, "PFIT-1001", "请检查输入参数")
		return
	}
	toTaskID, err := strconv.ParseInt(c.Query("to"), 10, 64)
	if err != nil || toTaskID <= 0 {
		response.ErrorWithStatus(c, http.StatusBadRequest, "PFIT-1001", "请检查输入参数")
		return
	}

	result, err := h.compareSvc.CompareTasks(c.Request.Context(), fromTaskID, toTaskID, currentUserID(c))
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, result)
}

// DeleteTask 删除任务
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	taskID, ok := parseTaskID(c)
	if !ok {
		return
	}

	if err := h.taskSvc.DeleteTask(c.Request.Context(), taskID, currentUserID(c)); err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
	h.trackEvent(c, "task_deleted", &taskID, map[string]interface{}{})
}

func parseTaskID(c *gin.Context) (int64, bool) {
	raw := c.Param("id")
	taskID, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || taskID <= 0 {
		response.ErrorWithStatus(c, http.StatusBadRequest, "PFIT-1001", "请检查输入参数")
		return 0, false
	}
	return taskID, true
}

func currentUserID(c *gin.Context) int64 {
	if userID := middleware.GetUserID(c); userID > 0 {
		return userID
	}
	raw := c.GetHeader("X-User-ID")
	if raw == "" {
		return 1
	}
	userID, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || userID <= 0 {
		return 1
	}
	return userID
}

func handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidArgument):
		response.ErrorWithStatus(c, http.StatusBadRequest, "PFIT-1001", "请检查输入参数")
	case errors.Is(err, service.ErrTaskNotFound):
		response.ErrorWithStatus(c, http.StatusNotFound, "PFIT-5001", "任务不存在或已被删除")
	case errors.Is(err, service.ErrForbidden):
		response.ErrorWithStatus(c, http.StatusForbidden, "PFIT-1003", "无权限访问")
	case errors.Is(err, service.ErrTaskStatusConflict):
		response.ErrorWithStatus(c, http.StatusConflict, "PFIT-5002", "当前任务状态不支持此操作")
	case errors.Is(err, service.ErrRequiredDocumentsMissing):
		response.ErrorWithStatus(c, http.StatusBadRequest, "PFIT-5003", "请上传体检报告与保险条款后再开始分析")
	case errors.Is(err, service.ErrUnsupportedFileType):
		response.ErrorWithStatus(c, http.StatusBadRequest, "PFIT-2001", "仅支持 PDF 格式文件")
	case errors.Is(err, service.ErrFileTooLarge):
		response.ErrorWithStatus(c, http.StatusBadRequest, "PFIT-2002", "文件大小超过 30MB，请压缩后重传")
	default:
		response.ErrorWithStatus(c, http.StatusInternalServerError, "PFIT-1005", "系统繁忙，请稍后重试")
	}
}

func (h *TaskHandler) trackEvent(c *gin.Context, eventName string, taskID *int64, properties map[string]interface{}) {
	if h.analytics == nil {
		return
	}
	userID := currentUserID(c)
	if userID <= 0 {
		return
	}
	uid := userID
	_ = h.analytics.Track(c.Request.Context(), eventName, &uid, taskID, properties)
}

func buildRetrievalSourceMap(findings []domain.RiskFinding) map[string][]domain.Evidence {
	result := make(map[string][]domain.Evidence, len(findings))
	for _, finding := range findings {
		if len(finding.PolicyEvidence) == 0 {
			continue
		}
		limit := len(finding.PolicyEvidence)
		if limit > 3 {
			limit = 3
		}
		result[finding.Topic] = append([]domain.Evidence{}, finding.PolicyEvidence[:limit]...)
	}
	return result
}
