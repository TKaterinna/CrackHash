package middleware

import (
	"time"

	"github.com/TKaterinna/CrackHash/manager/internal/metrics"
	"github.com/gin-gonic/gin"
)

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		elapsed := time.Since(start).Seconds()
		statusCode := c.Writer.Status()

		handler := getHandlerByPath(c.Request.URL.Path)
		if handler == "" {
			return
		}

		statusClass := normalizeStatusCode(statusCode)

		metrics.RequestsTotal.WithLabelValues(handler, statusClass).Inc()
		metrics.RequestDuration.WithLabelValues(handler).Observe(elapsed)
	}
}

func getHandlerByPath(path string) string {
	switch path {
	case "/api/hash/crack":
		return "crack"
	case "/api/hash/status":
		return "status"
	default:
		return ""
	}
}

func normalizeStatusCode(code int) string {
	switch {
	case code >= 500:
		return "5xx"
	case code >= 400:
		return "4xx"
	case code >= 300:
		return "3xx"
	case code >= 200:
		return "2xx"
	default:
		return "1xx"
	}
}
