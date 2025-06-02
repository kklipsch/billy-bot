package client

import (
	"net/http"
	"time"
)

const (
	// BaseURL is the base URL for the Frinkiac website
	BaseURL = "https://frinkiac.com"
)

// Client represents a client for interacting with the Frinkiac website
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// New creates a new Frinkiac client
func New(opts ...Option) *Client {
	c := &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: BaseURL,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Option is a function that configures a Client
type Option func(*Client)

// WithHTTPClient sets the HTTP client to use
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithBaseURL sets the base URL to use
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}
