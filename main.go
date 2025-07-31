package main

import (
	"alpaca-dev-toolkit/pkg/alpaca"
	"alpaca-dev-toolkit/pkg/metrics"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// load env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	API_Key := os.Getenv("APCA_API_KEY_ID")
	API_Secret := os.Getenv("APCA_API_SECRET_KEY")

	if API_Key == "" || API_Secret == "" {
		log.Fatal("API credentials not set in environment variables")
	}

	// Initialize metrics
	metricsCollector := metrics.NewCollector()

	// initialize Alpaca client
	alpacaClient := alpaca.NewClient(API_Key, API_Secret, metricsCollector)

	// start monitoring
	go alpacaClient.StartMonitoring(30 * time.Second)

	// serve Prometheus metrics
	http.Handle("/metrics", promhttp.Handler())

	//health check endpoint
	http.HandleFunc("health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	log.Println("Alpaca Performance Monitor started!")
	log.Println("Metrics available at: http://localhost:2112/metrics")
	log.Fatal(http.ListenAndServe(":2112", nil))
}
