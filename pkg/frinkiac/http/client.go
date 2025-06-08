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

// Config holds configuration for making Frinkiac HTTP requests
type Config struct {
	BaseURL string
}

// DefaultConfig returns a default configuration for Frinkiac requests
func DefaultConfig() Config {
	return Config{
		BaseURL: BaseURL,
	}
}

// NewHTTPClient creates a new HTTP client with appropriate timeout for Frinkiac
func NewHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
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
func doRequest(ctx context.Context, client *http.Client, config Config, opts RequestOptions) (*http.Response, error) {
	// Construct URL
	u, err := url.Parse(fmt.Sprintf("%s%s", config.BaseURL, opts.Path))
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
	resp, err := client.Do(req)
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
