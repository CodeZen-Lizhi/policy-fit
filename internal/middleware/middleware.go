package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhenglizhi/policy-fit/pkg/logger"
)

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		logger.Info("HTTP Request",
			"request_id", GetRequestID(c),
			"method", c.Request.Method,
			"path", path,
			"query", sanitizeQuery(query),
			"status", status,
			"latency", latency.String(),
			"ip", maskIP(c.ClientIP()),
		)
	}
}

// CORS 跨域中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func sanitizeQuery(query string) string {
	if query == "" {
		return ""
	}

	lowered := strings.ToLower(query)
	if strings.Contains(lowered, "token=") ||
		strings.Contains(lowered, "authorization=") ||
		strings.Contains(lowered, "password=") {
		return "[REDACTED]"
	}
	if len(query) > 120 {
		return query[:120] + "..."
	}
	return query
}

func maskIP(ip string) string {
	if strings.Count(ip, ".") == 3 {
		parts := strings.Split(ip, ".")
		return parts[0] + "." + parts[1] + ".*.*"
	}
	return ip
}
