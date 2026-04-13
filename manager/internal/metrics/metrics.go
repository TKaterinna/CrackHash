package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "manager_requests_total",
			Help: "Total number of HTTP requests received by manager, grouped by handler and status class",
		},
		[]string{"handler", "status_class"},
	)

	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "manager_request_duration_seconds",
			Help:    "Duration of HTTP request processing in seconds",
			Buckets: []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"handler"},
	)

	HealthStatus = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "manager_status",
			Help: "Health status of the manager (1=UP, 0=DOWN)",
		},
	)
)
