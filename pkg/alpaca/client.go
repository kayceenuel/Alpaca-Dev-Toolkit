package alpaca

import (
	"alpaca-dev-toolkit/pkg/metrics"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	APIKEY     string
	APISecret  string
	HTTPClient *http.Client
	Metrics    *metrics.Collector
}

func NewClient(apiKey, apiSecret string, metricsCollector *metrics.Collector) *Client {
	return &Client{
		APIKEY:     apiKey,
		APISecret:  apiSecret,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		Metrics:    metricsCollector,
	}
}

func (c *Client) MakeRequest(endpoint string) error {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		c.Metrics.RecordError(endpoint, "request_creation_error")
		return err
	}

	req.Header.Set("APCA-API-KEY-ID", c.APIKEY)
	req.Header.Set("APCA-API-SECRET-KEY", c.APISecret)

	start := time.Now()
	resp, err := c.HTTPClient.Do(req)
	duration := time.Since(start).Seconds()

	if err != nil {
		c.Metrics.RecordError(endpoint, "network_error")
		return err
	}

	defer resp.Body.Close()

	var remaining int

	remainingStr := resp.Header.Get("X-RateLimit-Remaining")
	if remainingStr != "" {
		var err error
		remaining, err = strconv.Atoi(remainingStr)
		if err != nil {
			c.Metrics.RecordError(endpoint, "rate_limit_parse_error")
			remaining = 0
		}

		c.Metrics.RateLimitRemaining.WithLabelValues(endpoint).Set(float64(remaining))

		if remaining < 20 {
			slog.Warn("Low rate limit detected",
				"endpoint", endpoint,
				"remaining", remaining)
		}
	}

	limitStr := resp.Header.Get("X-RateLimit-Limit")
	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err == nil {
			c.Metrics.RateLimitLimit.WithLabelValues(endpoint).Set(float64(limit))
		}
	}

	statusCode := fmt.Sprintf("%d", resp.StatusCode)
	c.Metrics.RecordRequest(endpoint, statusCode, duration)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.Metrics.RecordError(endpoint, "body_read_error")
		return err
	}

	if resp.StatusCode >= 400 {
		c.Metrics.RecordError(endpoint, "http_error")
	}

	// Fixed with slog calls
	slog.Info("API request completed",
		"endpoint", endpoint,
		"duration_ms", duration*1000,
		"status_code", resp.StatusCode,
		"rate_limit_remaining", remaining)

	if resp.StatusCode == 200 {
		slog.Debug("Request succeeded",
			"endpoint", endpoint,
			"status_code", resp.StatusCode)
	} else {
		slog.Error("Request failed",
			"endpoint", endpoint,
			"status_code", resp.StatusCode,
			"response_body", string(body))
	}

	return nil
}

func (c *Client) StartMonitoring(interval time.Duration) {
	endpoints := []string{
		"https://paper-api.alpaca.markets/v2/account",
		"https://paper-api.alpaca.markets/v2/positions",
		"https://paper-api.alpaca.markets/v2/orders",
		"https://paper-api.alpaca.markets/v2/assets",
	}

	slog.Info("Starting monitoring loop",
		"interval", interval,
		"endpoints", len(endpoints))

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			slog.Info("Monitoring cycle started")
			for _, endpoint := range endpoints {
				c.MakeRequest(endpoint)
				time.Sleep(1 * time.Second)
			}
			slog.Info("Monitoring cycle completed")
		}
	}
}
