package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhenglizhi/policy-fit/internal/i18n"
)

// Response 统一响应结构
type Response struct {
	Code      string      `json:"code,omitempty"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
	Timestamp string      `json:"timestamp,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	requestID := c.Writer.Header().Get("X-Request-ID")
	if requestID == "" {
		requestID = c.GetHeader("X-Request-ID")
	}
	c.JSON(http.StatusOK, gin.H{
		"data":       data,
		"request_id": requestID,
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
	})
}

// Error 错误响应
func Error(c *gin.Context, code string, message string) {
	requestID := c.Writer.Header().Get("X-Request-ID")
	if requestID == "" {
		requestID = c.GetHeader("X-Request-ID")
	}
	c.JSON(http.StatusBadRequest, Response{
		Code:      code,
		Message:   localizeMessage(c, code, message),
		RequestID: requestID,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// ErrorWithStatus 带状态码的错误响应
func ErrorWithStatus(c *gin.Context, status int, code string, message string) {
	requestID := c.Writer.Header().Get("X-Request-ID")
	if requestID == "" {
		requestID = c.GetHeader("X-Request-ID")
	}
	c.JSON(status, Response{
		Code:      code,
		Message:   localizeMessage(c, code, message),
		RequestID: requestID,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

func localizeMessage(c *gin.Context, code, fallback string) string {
	return i18n.TranslateErrorMessage(code, c.GetHeader("Accept-Language"), fallback)
}
