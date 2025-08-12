package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics interface {
	ObserveLatency(path, method, status string, duration float64)
	IncError(path, method, status string)
}

var _ Metrics = &CartMetrics{}

type CartMetrics struct {
	ResponseLatency *prometheus.HistogramVec
	ErrorsTotal     *prometheus.CounterVec
}

func RegisterMetrics() (*CartMetrics, error) {
	responseLatency := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_time_seconds",
			Help:    "HTTP response latencies in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method", "status"},
	)

	if err := prometheus.Register(responseLatency); err != nil {
		return nil, err
	}

	errorCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_failed_requests_total",
			Help: "Total failed HTTP requests",
		},
		[]string{"path", "method", "status"},
	)

	if err := prometheus.Register(errorCounter); err != nil {
		return nil, err
	}

	return &CartMetrics{
		ResponseLatency: responseLatency,
		ErrorsTotal:     errorCounter,
	}, nil
}

func (m *CartMetrics) ObserveLatency(path, method, status string, duration float64) {
	m.ResponseLatency.With(prometheus.Labels{
		"path":   path,
		"method": method,
		"status": status,
	}).Observe(duration)
}

func (m *CartMetrics) IncError(path, method, status string) {
	m.ErrorsTotal.With(prometheus.Labels{
		"path":   path,
		"method": method,
		"status": status,
	}).Inc()
}
