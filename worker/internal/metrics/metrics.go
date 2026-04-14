package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	TasksTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "worker_tasks_total",
			Help: "Total number of tasks processed by worker, grouped by result status",
		},
		[]string{"status"},
	)

	WorkerStatus = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "worker_status",
			Help: "Current worker status: 1 = busy (processing task), 0 = idle",
		},
	)

	TaskDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "worker_task_duration_seconds",
			Help:    "Duration of single task processing in seconds",
			Buckets: []float64{0.1, 0.5, 1, 2.5, 5, 10, 30, 60, 120},
		},
	)

	ActiveTasks = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "worker_active_tasks",
			Help: "Current number of concurrently processing tasks",
		},
	)
)
