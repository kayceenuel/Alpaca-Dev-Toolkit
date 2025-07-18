package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
	go func() {
		for {
			requestCounter := prometheus.NewCounter(
				prometheus.CounterOpts{
					Name: "http_requests_total",
					Help: "Total number of HTTP requests",
				},
			)

			// Create a Histogram metric
			requestDurationHistogram := prometheus.NewHistogram(
				prometheus.HistogramOpts{
					Name: "http_requests_duration_seconds",
					Help: "Histrogram of response time for HTTP requests",
				},
			)
		}
	}()
}

func main() {
	recordMetrics()

	http.Handle("/main", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
