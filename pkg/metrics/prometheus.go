package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Collector struct {
	RequestDuration *prometheus.HistogramVec
	RequestsTotal   *prometheus.CounterVec
	ErrorsTotal     *prometheus.CounterVec
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
				Name:"alpaca_api_requests_total",
				Help: "Total number of Alpaca API requests",
			}, 
			[]string{"endpoint", "status_code"}, 
		),

	}
}

func (c *Collector) {

	NewCollector()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
