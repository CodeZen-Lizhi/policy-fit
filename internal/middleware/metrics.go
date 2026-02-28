package middleware

import (
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	totalRequests       int64
	totalFailures       int64
	totalProcessingTime int64
)

// Metrics 记录基础指标
func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		atomic.AddInt64(&totalRequests, 1)

		c.Next()

		if c.Writer.Status() >= 500 {
			atomic.AddInt64(&totalFailures, 1)
		}
		atomic.AddInt64(&totalProcessingTime, time.Since(start).Milliseconds())
	}
}

// MetricsSnapshotHandler 输出基础指标
func MetricsSnapshotHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		requests := atomic.LoadInt64(&totalRequests)
		failures := atomic.LoadInt64(&totalFailures)
		totalLatency := atomic.LoadInt64(&totalProcessingTime)
		avgLatency := int64(0)
		if requests > 0 {
			avgLatency = totalLatency / requests
		}

		failureRate := float64(0)
		if requests > 0 {
			failureRate = float64(failures) / float64(requests)
		}

		c.JSON(200, gin.H{
			"requests_total":         requests,
			"failures_total":         failures,
			"failure_rate":           failureRate,
			"avg_processing_ms":      avgLatency,
			"alert_failure_rate_gt":  0.15,
			"alert_triggered":        failureRate > 0.15,
			"reporting_time_unix_ms": time.Now().UnixMilli(),
		})
	}
}
