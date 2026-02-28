package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhenglizhi/policy-fit/internal/domain"
	"github.com/zhenglizhi/policy-fit/internal/service"
)

func TestRuleConfigHandlerAdminAPIs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := &fakeRuleConfigService{
		active: &domain.RuleConfigVersion{
			ID:        1,
			Version:   "v1",
			Changelog: "init",
			Content:   map[string]interface{}{"topics": []string{"hypertension"}},
			IsActive:  true,
			CreatedAt: time.Now(),
		},
	}

	h := NewRuleConfigHandler(svc, []int64{99})
	router := gin.New()
	v1 := router.Group("/api/v1")
	RegisterRuleConfigRoutes(v1, h)

	// Non-admin forbidden
	forbiddenRec := httptest.NewRecorder()
	forbiddenReq := httptest.NewRequest(http.MethodGet, "/api/v1/admin/rules/active", nil)
	forbiddenReq.Header.Set("X-User-ID", "1")
	router.ServeHTTP(forbiddenRec, forbiddenReq)
	if forbiddenRec.Code != http.StatusForbidden {
		t.Fatalf("forbidden status = %d", forbiddenRec.Code)
	}

	// Publish
	publishRec := httptest.NewRecorder()
	publishBody := bytes.NewBufferString(`{"changelog":"add bp rule","content":{"rules":[{"topic":"hypertension"}]}}`)
	publishReq := httptest.NewRequest(http.MethodPost, "/api/v1/admin/rules/publish", publishBody)
	publishReq.Header.Set("Content-Type", "application/json")
	publishReq.Header.Set("X-User-ID", "99")
	router.ServeHTTP(publishRec, publishReq)
	if publishRec.Code != http.StatusOK {
		t.Fatalf("publish status = %d, body = %s", publishRec.Code, publishRec.Body.String())
	}
	if svc.publishCount != 1 {
		t.Fatalf("publish count = %d", svc.publishCount)
	}

	// List versions
	listRec := httptest.NewRecorder()
	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/admin/rules/versions?limit=10", nil)
	listReq.Header.Set("X-User-ID", "99")
	router.ServeHTTP(listRec, listReq)
	if listRec.Code != http.StatusOK {
		t.Fatalf("list versions status = %d, body = %s", listRec.Code, listRec.Body.String())
	}

	// Gray mode
	grayRec := httptest.NewRecorder()
	grayReq := httptest.NewRequest(http.MethodPost, "/api/v1/admin/rules/gray", bytes.NewBufferString(`{"version":"v2","enabled":true}`))
	grayReq.Header.Set("Content-Type", "application/json")
	grayReq.Header.Set("X-User-ID", "99")
	router.ServeHTTP(grayRec, grayReq)
	if grayRec.Code != http.StatusOK {
		t.Fatalf("gray status = %d, body = %s", grayRec.Code, grayRec.Body.String())
	}
	if svc.grayVersion != "v2" || !svc.grayEnabled {
		t.Fatalf("unexpected gray args: version=%s enabled=%v", svc.grayVersion, svc.grayEnabled)
	}

	// Rollback
	rollbackRec := httptest.NewRecorder()
	rollbackReq := httptest.NewRequest(http.MethodPost, "/api/v1/admin/rules/rollback", bytes.NewBufferString(`{"version":"v1"}`))
	rollbackReq.Header.Set("Content-Type", "application/json")
	rollbackReq.Header.Set("X-User-ID", "99")
	router.ServeHTTP(rollbackRec, rollbackReq)
	if rollbackRec.Code != http.StatusOK {
		t.Fatalf("rollback status = %d, body = %s", rollbackRec.Code, rollbackRec.Body.String())
	}
	if svc.rollbackVersion != "v1" {
		t.Fatalf("rollback version = %s", svc.rollbackVersion)
	}

	// Active
	activeRec := httptest.NewRecorder()
	activeReq := httptest.NewRequest(http.MethodGet, "/api/v1/admin/rules/active", nil)
	activeReq.Header.Set("X-User-ID", "99")
	router.ServeHTTP(activeRec, activeReq)
	if activeRec.Code != http.StatusOK {
		t.Fatalf("active status = %d, body = %s", activeRec.Code, activeRec.Body.String())
	}

	var payload struct {
		Data domain.RuleConfigVersion `json:"data"`
	}
	if err := json.Unmarshal(activeRec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal active response: %v", err)
	}
	if payload.Data.Version == "" {
		t.Fatalf("active version should not be empty")
	}
}

type fakeRuleConfigService struct {
	active          *domain.RuleConfigVersion
	versions        []domain.RuleConfigVersion
	publishCount    int
	rollbackVersion string
	grayVersion     string
	grayEnabled     bool
}

func (f *fakeRuleConfigService) Publish(
	_ context.Context,
	content map[string]interface{},
	changelog string,
	_ int64,
) (*domain.RuleConfigVersion, error) {
	if len(content) == 0 || changelog == "" {
		return nil, service.ErrInvalidArgument
	}
	f.publishCount++
	cfg := &domain.RuleConfigVersion{
		ID:        int64(f.publishCount + 1),
		Version:   "v2",
		Changelog: changelog,
		Content:   content,
		IsActive:  true,
		CreatedAt: time.Now(),
	}
	f.active = cfg
	f.versions = append([]domain.RuleConfigVersion{*cfg}, f.versions...)
	return cfg, nil
}

func (f *fakeRuleConfigService) Rollback(_ context.Context, version string, _ int64) error {
	if version == "" {
		return service.ErrInvalidArgument
	}
	f.rollbackVersion = version
	return nil
}

func (f *fakeRuleConfigService) SetGrayMode(_ context.Context, version string, enabled bool, _ int64) error {
	if version == "" {
		return service.ErrInvalidArgument
	}
	f.grayVersion = version
	f.grayEnabled = enabled
	return nil
}

func (f *fakeRuleConfigService) GetActive(_ context.Context) (*domain.RuleConfigVersion, error) {
	if f.active == nil {
		return nil, service.ErrInvalidArgument
	}
	return f.active, nil
}

func (f *fakeRuleConfigService) ListVersions(_ context.Context, limit int) ([]domain.RuleConfigVersion, error) {
	if limit <= 0 {
		return nil, service.ErrInvalidArgument
	}
	if len(f.versions) == 0 && f.active != nil {
		f.versions = []domain.RuleConfigVersion{*f.active}
	}
	return f.versions, nil
}
