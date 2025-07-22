package alpaca

import (
	"net/http"
	"runtime/metrics"
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
	req err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		c.Metrics.RecordError(endpoint, "request_creation_error")
		return err
	}

	req.Header.Set("APCA-API-KEY-ID", c.APIKey)
	req.Header.Set("APCA-API-SECRET-KEY", c.APISecret)
}
