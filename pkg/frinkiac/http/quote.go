package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// APISearchResult represents a single result from the Frinkiac API search endpoint
type APISearchResult struct {
	ID        int    `json:"Id"`
	Episode   string `json:"Episode"` // e.g., S16E01
	Timestamp int    `json:"Timestamp"`
}

// validateEpisodeFormat checks if the episode format is valid
// Episode format is expected to be like "S16E01"
func validateEpisodeFormat(episode string) error {
	if len(episode) < 6 {
		return fmt.Errorf("invalid episode format: %s", episode)
	}
	return nil
}

// GetSeasonAndEpisode extracts season and episode numbers from an APISearchResult
// Episode format is expected to be like "S16E01"
func GetSeasonAndEpisode(result APISearchResult) (season int, episode int, err error) {
	if err := validateEpisodeFormat(result.Episode); err != nil {
		return 0, 0, err
	}

	// Extract season number from "S16" part
	seasonStr := result.Episode[1:3] // Skip 'S', get "16"
	season, err = strconv.Atoi(seasonStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid season number in episode %s: %w", result.Episode, err)
	}

	// Extract episode number from "E01" part
	episodeStr := result.Episode[4:6] // Skip "S16E", get "01"
	episode, err = strconv.Atoi(episodeStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid episode number in episode %s: %w", result.Episode, err)
	}

	return season, episode, nil
}

// GetImagePath constructs the image path for an APISearchResult
func GetImagePath(result APISearchResult) (string, error) {
	// Validate the episode format first
	if err := validateEpisodeFormat(result.Episode); err != nil {
		return "", err
	}

	return fmt.Sprintf("/img/%s/%d/medium.jpg", result.Episode, result.Timestamp), nil
}

// hasClass checks if a node has a specific class
func hasClass(n *html.Node, class string) bool {
	for _, attr := range n.Attr {
		if attr.Key == "class" {
			classes := strings.Fields(attr.Val)
			for _, c := range classes {
				if c == class {
					return true
				}
			}
			break
		}
	}
	return false
}

// GetQuote searches for a quote on Frinkiac and returns the results
func GetQuote(ctx context.Context, client *http.Client, config Config, quote string) ([]APISearchResult, error) {
	// Set up query parameters
	queryParams := url.Values{}
	queryParams.Set("q", quote)

	// Set up log context
	logContext := map[string]interface{}{
		"quote": quote,
	}

	// Make the request
	resp, err := doRequest(ctx, client, config, RequestOptions{
		Method:      http.MethodGet,
		Path:        "/api/search",
		QueryParams: queryParams,
		LogContext:  logContext,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse JSON response
	var apiResults []APISearchResult
	if err := json.NewDecoder(resp.Body).Decode(&apiResults); err != nil {
		return nil, fmt.Errorf("error decoding JSON response: %w", err)
	}

	return apiResults, nil
}
