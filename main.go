package main

import (
	"alpaca-dev-toolkit/pkg/alpaca"
	"alpaca-dev-toolkit/pkg/metrics"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	//create a logger with JSON handler
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("Starting Alpaca Performance Monitor")

	// load env file
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file: %v", err)
		os.Exit(1)
	}

	apiKey := os.Getenv("APCA_API_KEY_ID")
	apiSecret := os.Getenv("APCA_API_SECRET_KEY")

	if apiKey == "" || apiSecret == "" {
		slog.Error("API credentials not set in environment variables")
		os.Exit(1)
	}

	slog.Info("Environment loaded successfully")

	// Initialize metrics
	metricsCollector := metrics.NewCollector()
	slog.Info("Metrics collector initalized")

	// initialize Alpaca client
	alpacaClient := alpaca.NewClient(apiKey, apiSecret, metricsCollector)
	slog.Info("Alpaca client initialized")

	// start monitoring
	go alpacaClient.StartMonitoring(30 * time.Second)
	slog.Info("Background monitoring started", "interval", "30s")

	// serve Prometheus metrics
	http.Handle("/metrics", promhttp.Handler())

	//health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"alpaca-monitor"}`))
	})

	// Startup complete
	slog.Info("server starting",
		"port", 2112,
		"metrics_endpoint", "http://locahost:2112/metrics",
		"health_endpoint", "http://localhost:2112/health")

	// Start server
	if err := http.ListenAndServe(":2112", nil); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
