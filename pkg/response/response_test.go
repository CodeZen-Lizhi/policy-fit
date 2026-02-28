package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestErrorWithStatusLocalized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.GET("/err", func(c *gin.Context) {
		ErrorWithStatus(c, http.StatusBadRequest, "PFIT-1001", "请检查输入参数")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/err", nil)
	req.Header.Set("Accept-Language", "en-US")
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status=%d", rec.Code)
	}

	var payload Response
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if payload.Message != "Invalid request parameters" {
		t.Fatalf("unexpected localized message: %s", payload.Message)
	}
}
