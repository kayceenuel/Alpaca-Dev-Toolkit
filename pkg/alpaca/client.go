package alpaca

import (
	"alpaca-dev-toolkit/pkg/metrics"
	"fmt"
	"io"
	"log"
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

	remainingStr := resp.Header.Get("X-RateLimit-Remaining")

	if remainingStr != "" { // Check if header exists
		remaining, err := strconv.Atoi(remainingStr)
		if err != nil {
			c.Metrics.RecordError(endpoint, "rate_limit_parse_error")
			remaining = 0
		}

		if remaining < 20 {
			c.Metrics.RateLimitRemaining.WithLabelValues(endpoint).Inc()
			log.Printf(" Low rate limit on %s: %d calls remaining", endpoint, remaining)
		}
	}

	limitStr := resp.Header.Get("X-RateLimit-Limit")
	if limitStr != "" { // check if header exists
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
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

	fmt.Printf("Response time: %.0fms, Status: %d, Endpoint: %s\n",
		duration*1000, resp.StatusCode, endpoint)

	if resp.StatusCode == 200 {
		fmt.Printf("Success\n")
	} else {
		fmt.Printf(" Error body: %s\n", string(body))

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

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("\n--- Monitoring Cycle Started ---")
			for _, endpoint := range endpoints {
				c.MakeRequest(endpoint)
				time.Sleep(1 * time.Second)
			}
			fmt.Println("Monitoring Cycle Complete")
		}
	}
}
