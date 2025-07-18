package main

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
	go func() {
		for {
			// Record the duration in the histogram 
			duration := time.Since(start).Seconds()
			requestDurationHistogram.Observe(duration)

			errorCounter.Inc() // only if an error happened

		}
	}()
}

// Create a counter metrics
var (
	errorCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	}), 

	requestDurationHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "http_request_duration_seconds", 
		Help: "Histogram of response time for HTTP requests", 
		Buckets: prometheus.DefBuckets, // Default buckets: [0.005, 0.01, 0.025, ..., 10.0]
	}),
)

func main() {
	//register the metrics 
	prometheus.MustRegister(errorCounter)
	prometheus.MustRegister(requestDurationHistogram)

	recordMetrics()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
