package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/rs/zerolog/log"
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

// RequestOptions contains options for making a request
type RequestOptions struct {
	Method      string
	Path        string
	QueryParams url.Values
	LogContext  map[string]interface{}
}

// doRequest makes an HTTP request and returns the response body
func (c *Client) doRequest(ctx context.Context, opts RequestOptions) (*http.Response, error) {
	// Construct URL
	u, err := url.Parse(fmt.Sprintf("%s%s", c.baseURL, opts.Path))
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %w", err)
	}

	// Add query parameters if provided
	if opts.QueryParams != nil {
		u.RawQuery = opts.QueryParams.Encode()
	}

	// Create request
	requestURL := u.String()

	// Create logger with context
	logEvent := log.Debug().Str("url", requestURL).Str("method", opts.Method)
	for k, v := range opts.LogContext {
		switch val := v.(type) {
		case string:
			logEvent = logEvent.Str(k, val)
		case int:
			logEvent = logEvent.Int(k, val)
		case bool:
			logEvent = logEvent.Bool(k, val)
		default:
			logEvent = logEvent.Interface(k, v)
		}
	}
	logEvent.Msg("sending request to frinkiac")

	// Create request
	req, err := http.NewRequestWithContext(ctx, opts.Method, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		// Read the response body for error details
		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()

		if readErr != nil {
			return nil, fmt.Errorf("unexpected status code: %d (failed to read response body: %v)",
				resp.StatusCode, readErr)
		}

		return nil, fmt.Errorf("unexpected status code: %d, body: %s",
			resp.StatusCode, string(body))
	}

	return resp, nil
}
