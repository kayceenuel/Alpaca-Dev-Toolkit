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

}
