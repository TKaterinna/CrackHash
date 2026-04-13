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

		endpoint := c.FullPath()
		if endpoint == "" {
			endpoint = c.Request.URL.Path
		}

		status := c.Writer.Status()

		metrics.RequestsTotal.WithLabelValues(endpoint, string(rune(status))).Inc()

		metrics.RequestDuration.WithLabelValues(endpoint).Observe(elapsed)
	}
}
