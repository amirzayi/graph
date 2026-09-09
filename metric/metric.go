package metric

import "github.com/prometheus/client_golang/prometheus"

var (
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	TotalTasks = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "total_tasks",
			Help: "Number of total tasks registered on the app",
		},
	)
)

func init() {
	prometheus.MustRegister(RequestsTotal, RequestDuration, TotalTasks)
}
