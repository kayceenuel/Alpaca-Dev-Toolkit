package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	http.Handle("/main", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
