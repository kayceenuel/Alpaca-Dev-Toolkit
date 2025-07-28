package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct {
	RequestDuration    *prometheus.HistogramVec
	RequestsTotal      *prometheus.CounterVec
	ErrorsTotal        *prometheus.CounterVec
	RateLimitRemaining *prometheus.GaugeVec
	RateLimitWarnings  *prometheus.CounterVec
}

func NewCollector() *Collector {
	collector := &Collector{
		RequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "alpaca_api_request_duration_seconds",
				Help:    "Histogram of response time for Alpaca API requests",
				Buckets: []float64{0.1, 0.25, 0.5, 1.0, 2.5, 5.0, 10.0},
			},
			[]string{"endpoint", "status_code"},
		),

		RequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "alpaca_api_requests_total",
				Help: "Total number of Alpaca API requests",
			},
			[]string{"endpoint", "status_code"},
		),

		ErrorsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "alpaca_api_errors_total",
				Help: "Total number of Alpaca API errors",
			},
			[]string{"endpoint", "error_type"},
		),

		RateLimitRemaining: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "alpaca_rate_limit_remaining",
				Help: "How many Alpaca Calls you have left this window",
			},
			[]string{"endpoint"},
		),

		RateLimitWarnings: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "alpaca_rate_limit_warning_total",
				Help: "Times we warned that rate limit was low",
			},
			[]string{"endpoint"},
		),
	}

	// Register metrics
	prometheus.MustRegister(collector.RequestDuration)
	prometheus.MustRegister(collector.RequestsTotal)
	prometheus.MustRegister(collector.ErrorsTotal)
	prometheus.MustRegister(collector.RateLimitRemaining)
	prometheus.MustRegister(collector.RateLimitWarnings)

	return collector
}

func (c *Collector) RecordRequest(endpoint string, statusCode string, duration float64) {
	c.RequestDuration.WithLabelValues(endpoint, statusCode).Observe(duration)
	c.RequestsTotal.WithLabelValues(endpoint, statusCode).Inc()
}

func (c *Collector) RecordError(endpoint string, errorType string) {
	c.ErrorsTotal.WithLabelValues(endpoint, errorType).Inc()
}
