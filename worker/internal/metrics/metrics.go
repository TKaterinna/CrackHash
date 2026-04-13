package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	TasksTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "worker_tasks_total",
			Help: "Total number of tasks processed by worker",
		},
		[]string{"status"},
	)

	TasksInProgress = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "worker_tasks_in_progress",
			Help: "Number of tasks currently being processed (0 or 1)",
		},
	)

	TaskDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "worker_task_duration_seconds",
			Help:    "Duration of task processing in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)
)
