package handler

import (
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhenglizhi/policy-fit/internal/domain"
	"github.com/zhenglizhi/policy-fit/internal/service"
)

func BenchmarkTaskHandlerCreateTask(b *testing.B) {
	gin.SetMode(gin.TestMode)

	taskSvc := &benchTaskService{}
	docSvc := &benchDocumentService{}
	findingSvc := &benchFindingService{}
	storageSvc := &benchStorageService{}
	exportSvc := &benchExportService{}
	compareSvc := &benchCompareService{}
	analyticsSvc := &fakeAnalyticsService{}

	h := NewTaskHandler(taskSvc, docSvc, findingSvc, storageSvc, exportSvc, compareSvc, analyticsSvc)
	router := gin.New()
	v1 := router.Group("/api/v1")
	RegisterTaskRoutes(v1, h)

	body := `{"request_id":"bench"}`
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(rec, req)
			if rec.Code != http.StatusOK {
				b.Fatalf("unexpected status: %d", rec.Code)
			}
		}
	})
}

type benchTaskService struct {
	nextID atomic.Int64
}

func (s *benchTaskService) CreateTask(_ context.Context, userID int64, requestID string) (*domain.AnalysisTask, error) {
	id := s.nextID.Add(1)
	return &domain.AnalysisTask{
		ID:        id,
		UserID:    userID,
		RequestID: requestID,
		Status:    domain.TaskStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (s *benchTaskService) GetTask(_ context.Context, taskID int64, actorID int64) (*domain.AnalysisTask, error) {
	return &domain.AnalysisTask{ID: taskID, UserID: actorID, Status: domain.TaskStatusSuccess}, nil
}

func (s *benchTaskService) RunTask(context.Context, int64, int64) error {
	return nil
}

func (s *benchTaskService) DeleteTask(context.Context, int64, int64) error {
	return nil
}

type benchDocumentService struct{}

func (s *benchDocumentService) CreateDocument(context.Context, int64, int64, domain.DocumentType, string, string) (*domain.Document, error) {
	return &domain.Document{ID: 1}, nil
}

type benchFindingService struct{}

func (s *benchFindingService) ListByTask(context.Context, int64) ([]domain.RiskFinding, error) {
	return nil, nil
}

func (s *benchFindingService) Summarize([]domain.RiskFinding) map[string]int {
	return map[string]int{"red": 0, "yellow": 0, "green": 0}
}

type benchStorageService struct{}

func (s *benchStorageService) SaveUploadedFile(context.Context, int64, domain.DocumentType, *multipart.FileHeader) (string, error) {
	return "", nil
}

type benchExportService struct{}

func (s *benchExportService) ExportMarkdown(context.Context, int64, int64, string) ([]byte, error) {
	return []byte(""), nil
}

func (s *benchExportService) ExportPDF(context.Context, int64, int64, string) ([]byte, error) {
	return []byte(""), nil
}

type benchCompareService struct{}

func (s *benchCompareService) CompareTasks(context.Context, int64, int64, int64) (*service.ComparisonResult, error) {
	return &service.ComparisonResult{}, nil
}
