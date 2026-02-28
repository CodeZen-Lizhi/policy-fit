package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhenglizhi/policy-fit/internal/domain"
	"github.com/zhenglizhi/policy-fit/internal/service"
)

func TestTaskHandlerCoreAPIs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	taskSvc := &fakeTaskService{
		tasks: map[int64]*domain.AnalysisTask{},
	}
	docSvc := &fakeDocumentService{
		documents: map[int64][]domain.Document{},
	}
	findingSvc := &fakeFindingService{}
	storageSvc := &fakeStorageService{}
	exportSvc := &fakeExportService{}
	compareSvc := &fakeCompareService{}
	analyticsSvc := &fakeAnalyticsService{}

	h := NewTaskHandler(taskSvc, docSvc, findingSvc, storageSvc, exportSvc, compareSvc, analyticsSvc)
	router := gin.New()
	v1 := router.Group("/api/v1")
	RegisterTaskRoutes(v1, h)

	// 1. CreateTask
	createRec := httptest.NewRecorder()
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", strings.NewReader(`{"request_id":"req-1"}`))
	createReq.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(createRec, createReq)
	if createRec.Code != http.StatusOK {
		t.Fatalf("CreateTask status = %d, body = %s", createRec.Code, createRec.Body.String())
	}

	taskID := parseTaskIDFromResponse(t, createRec.Body.Bytes())

	// 2. GetTask
	getRec := httptest.NewRecorder()
	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/"+strconv.FormatInt(taskID, 10), nil)
	router.ServeHTTP(getRec, getReq)
	if getRec.Code != http.StatusOK {
		t.Fatalf("GetTask status = %d, body = %s", getRec.Code, getRec.Body.String())
	}

	// 3. UploadDocument (report)
	uploadReportRec := httptest.NewRecorder()
	uploadReportReq := newUploadRequest(
		t,
		"/api/v1/tasks/"+strconv.FormatInt(taskID, 10)+"/documents",
		"docType",
		"report",
		"file",
		"report.pdf",
		[]byte("%PDF-1.4 report"),
	)
	router.ServeHTTP(uploadReportRec, uploadReportReq)
	if uploadReportRec.Code != http.StatusOK {
		t.Fatalf("UploadDocument(report) status = %d, body = %s", uploadReportRec.Code, uploadReportRec.Body.String())
	}

	// 4. UploadDocument (policy)
	uploadPolicyRec := httptest.NewRecorder()
	uploadPolicyReq := newUploadRequest(
		t,
		"/api/v1/tasks/"+strconv.FormatInt(taskID, 10)+"/documents",
		"docType",
		"policy",
		"file",
		"policy.pdf",
		[]byte("%PDF-1.4 policy"),
	)
	router.ServeHTTP(uploadPolicyRec, uploadPolicyReq)
	if uploadPolicyRec.Code != http.StatusOK {
		t.Fatalf("UploadDocument(policy) status = %d, body = %s", uploadPolicyRec.Code, uploadPolicyRec.Body.String())
	}

	// 5. RunTask
	runRec := httptest.NewRecorder()
	runReq := httptest.NewRequest(http.MethodPost, "/api/v1/tasks/"+strconv.FormatInt(taskID, 10)+"/run", nil)
	router.ServeHTTP(runRec, runReq)
	if runRec.Code != http.StatusOK {
		t.Fatalf("RunTask status = %d, body = %s", runRec.Code, runRec.Body.String())
	}

	// 6. GetFindings
	findingsRec := httptest.NewRecorder()
	findingsReq := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/"+strconv.FormatInt(taskID, 10)+"/findings", nil)
	router.ServeHTTP(findingsRec, findingsReq)
	if findingsRec.Code != http.StatusOK {
		t.Fatalf("GetFindings status = %d, body = %s", findingsRec.Code, findingsRec.Body.String())
	}

	// 7. Export
	exportRec := httptest.NewRecorder()
	exportReq := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/"+strconv.FormatInt(taskID, 10)+"/export?format=md", nil)
	router.ServeHTTP(exportRec, exportReq)
	if exportRec.Code != http.StatusOK {
		t.Fatalf("Export status = %d, body = %s", exportRec.Code, exportRec.Body.String())
	}

	// 8. DeleteTask
	deleteRec := httptest.NewRecorder()
	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/v1/tasks/"+strconv.FormatInt(taskID, 10), nil)
	router.ServeHTTP(deleteRec, deleteReq)
	if deleteRec.Code != http.StatusNoContent {
		t.Fatalf("DeleteTask status = %d, body = %s", deleteRec.Code, deleteRec.Body.String())
	}

	// 9. Compare
	compareRec := httptest.NewRecorder()
	compareReq := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/compare?from=1&to=2", nil)
	router.ServeHTTP(compareRec, compareReq)
	if compareRec.Code != http.StatusOK {
		t.Fatalf("Compare status = %d, body = %s", compareRec.Code, compareRec.Body.String())
	}
}

func parseTaskIDFromResponse(t *testing.T, raw []byte) int64 {
	t.Helper()

	var payload struct {
		Data struct {
			TaskID int64 `json:"task_id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	return payload.Data.TaskID
}

func newUploadRequest(
	t *testing.T,
	url string,
	docTypeKey string,
	docTypeValue string,
	fileKey string,
	fileName string,
	content []byte,
) *http.Request {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	if err := writer.WriteField(docTypeKey, docTypeValue); err != nil {
		t.Fatalf("write field: %v", err)
	}
	part, err := writer.CreateFormFile(fileKey, fileName)
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := part.Write(content); err != nil {
		t.Fatalf("write form file content: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, url, &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

type fakeTaskService struct {
	tasks  map[int64]*domain.AnalysisTask
	nextID int64
}

func (f *fakeTaskService) CreateTask(_ context.Context, userID int64, requestID string) (*domain.AnalysisTask, error) {
	f.nextID++
	task := &domain.AnalysisTask{
		ID:        f.nextID,
		UserID:    userID,
		RequestID: requestID,
		Status:    domain.TaskStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	f.tasks[task.ID] = task
	return task, nil
}

func (f *fakeTaskService) GetTask(_ context.Context, taskID int64, _ int64) (*domain.AnalysisTask, error) {
	task, ok := f.tasks[taskID]
	if !ok {
		return nil, service.ErrTaskNotFound
	}
	return task, nil
}

func (f *fakeTaskService) RunTask(_ context.Context, taskID int64, _ int64) error {
	if _, ok := f.tasks[taskID]; !ok {
		return service.ErrTaskNotFound
	}
	return nil
}

func (f *fakeTaskService) DeleteTask(_ context.Context, taskID int64, _ int64) error {
	if _, ok := f.tasks[taskID]; !ok {
		return service.ErrTaskNotFound
	}
	delete(f.tasks, taskID)
	return nil
}

type fakeDocumentService struct {
	nextID    int64
	documents map[int64][]domain.Document
}

func (f *fakeDocumentService) CreateDocument(
	_ context.Context,
	taskID int64,
	_ int64,
	docType domain.DocumentType,
	fileName string,
	storageKey string,
) (*domain.Document, error) {
	f.nextID++
	doc := domain.Document{
		ID:         f.nextID,
		TaskID:     taskID,
		DocType:    docType,
		FileName:   fileName,
		StorageKey: storageKey,
	}
	f.documents[taskID] = append(f.documents[taskID], doc)
	return &doc, nil
}

type fakeFindingService struct{}

func (f *fakeFindingService) ListByTask(_ context.Context, _ int64) ([]domain.RiskFinding, error) {
	return []domain.RiskFinding{}, nil
}

func (f *fakeFindingService) Summarize(_ []domain.RiskFinding) map[string]int {
	return map[string]int{"red": 0, "yellow": 0, "green": 0}
}

type fakeStorageService struct{}

func (f *fakeStorageService) SaveUploadedFile(
	_ context.Context,
	taskID int64,
	docType domain.DocumentType,
	fileHeader *multipart.FileHeader,
) (string, error) {
	return "task/" + strconv.FormatInt(taskID, 10) + "/" + string(docType) + "/" + fileHeader.Filename, nil
}

type fakeExportService struct{}

func (f *fakeExportService) ExportMarkdown(context.Context, int64, int64, string) ([]byte, error) {
	return []byte("# report"), nil
}

func (f *fakeExportService) ExportPDF(context.Context, int64, int64, string) ([]byte, error) {
	return []byte("%PDF"), nil
}

type fakeCompareService struct{}

func (f *fakeCompareService) CompareTasks(context.Context, int64, int64, int64) (*service.ComparisonResult, error) {
	return &service.ComparisonResult{
		FromTaskID: 1,
		ToTaskID:   2,
		Summary:    "diff",
	}, nil
}

type fakeAnalyticsService struct{}

func (f *fakeAnalyticsService) Track(
	context.Context,
	string,
	*int64,
	*int64,
	map[string]interface{},
) error {
	return nil
}
