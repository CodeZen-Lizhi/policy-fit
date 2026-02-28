package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/zhenglizhi/policy-fit/internal/service"
)

func TestAnalyticsHandlerAPIs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := &fakeAnalyticsQueryService{}
	h := NewAnalyticsHandler(svc, []int64{99})
	router := gin.New()
	v1 := router.Group("/api/v1")
	RegisterAnalyticsRoutes(v1, h)

	// Track event
	trackRec := httptest.NewRecorder()
	trackReq := httptest.NewRequest(http.MethodPost, "/api/v1/analytics/events", bytes.NewBufferString(`{"event_name":"report_viewed","task_id":11,"properties":{"source":"web"}}`))
	trackReq.Header.Set("Content-Type", "application/json")
	trackReq.Header.Set("X-User-ID", "2")
	router.ServeHTTP(trackRec, trackReq)
	if trackRec.Code != http.StatusOK {
		t.Fatalf("track event status = %d, body = %s", trackRec.Code, trackRec.Body.String())
	}
	if svc.trackCount != 1 {
		t.Fatalf("track count = %d", svc.trackCount)
	}

	// Non-admin should be forbidden
	forbiddenRec := httptest.NewRecorder()
	forbiddenReq := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/funnel?period=week", nil)
	forbiddenReq.Header.Set("X-User-ID", "2")
	router.ServeHTTP(forbiddenRec, forbiddenReq)
	if forbiddenRec.Code != http.StatusForbidden {
		t.Fatalf("funnel forbidden status = %d", forbiddenRec.Code)
	}

	// Admin overview
	overviewRec := httptest.NewRecorder()
	overviewReq := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/overview?period=month", nil)
	overviewReq.Header.Set("X-User-ID", "99")
	router.ServeHTTP(overviewRec, overviewReq)
	if overviewRec.Code != http.StatusOK {
		t.Fatalf("overview status = %d, body = %s", overviewRec.Code, overviewRec.Body.String())
	}

	var payload struct {
		Data map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(overviewRec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal overview: %v", err)
	}
	if payload.Data["period"] != "month" {
		t.Fatalf("unexpected period: %v", payload.Data["period"])
	}
}

type fakeAnalyticsQueryService struct {
	trackCount int
}

func (f *fakeAnalyticsQueryService) Track(
	_ context.Context,
	eventName string,
	userID *int64,
	taskID *int64,
	properties map[string]interface{},
) error {
	if eventName == "" || userID == nil || *userID <= 0 {
		return service.ErrInvalidArgument
	}
	if taskID != nil && *taskID <= 0 {
		return service.ErrInvalidArgument
	}
	if properties == nil {
		properties = map[string]interface{}{}
	}
	f.trackCount++
	return nil
}

func (f *fakeAnalyticsQueryService) FunnelByPeriod(_ context.Context, period string) (map[string]int64, error) {
	if period != "all" && period != "week" && period != "month" {
		return nil, service.ErrInvalidArgument
	}
	return map[string]int64{
		"task_created":    10,
		"task_completed":  8,
		"report_viewed":   6,
		"report_exported": 2,
	}, nil
}

func (f *fakeAnalyticsQueryService) Dashboard(_ context.Context, period string) (map[string]interface{}, error) {
	if period != "all" && period != "week" && period != "month" {
		return nil, service.ErrInvalidArgument
	}
	return map[string]interface{}{
		"period":          period,
		"task_created":    int64(10),
		"task_completed":  int64(8),
		"completion_rate": 0.8,
	}, nil
}
