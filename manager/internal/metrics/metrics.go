package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "manager_requests_total",
			Help: "Total number of HTTP requests received by manager",
		},
		[]string{"endpoint", "status"},
	)

	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "manager_request_duration_seconds",
			Help:    "Duration of HTTP request processing in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint"},
	)

	HealthStatus = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "manager_health_status",
			Help: "Health status of the manager (1=UP, 0=DOWN)",
		},
	)
)
