package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// Timestamp represents a Frinkiac timestamp as a string
type Timestamp string

// UnmarshalJSON implements json.Unmarshaler for Timestamp
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	// Remove quotes if present (string value)
	if len(data) > 0 && data[0] == '"' && data[len(data)-1] == '"' {
		*t = Timestamp(data[1 : len(data)-1])
		return nil
	}
	// Convert number to string
	*t = Timestamp(string(data))
	return nil
}

// EpisodeID represents a Simpsons episode identifier in the format S##E##
type EpisodeID string

// ErrInvalidEpisodeFormat is returned when an episode format is invalid
var ErrInvalidEpisodeFormat = errors.New("invalid episode format")

// SearchResult represents a single result from the Frinkiac API search endpoint
type SearchResult struct {
	ID        int       `json:"Id"`
	EpisodID  EpisodeID `json:"Episode"` // e.g., S16E01
	Timestamp Timestamp `json:"Timestamp"`
}

// validateEpisodeFormat checks if the episode format is valid
// Episode format is expected to be like "S16E01"
func validateEpisodeFormat(episode string) error {
	if len(episode) < 6 {
		return fmt.Errorf("%w: %s", ErrInvalidEpisodeFormat, episode)
	}
	return nil
}

// GetSeasonAndEpisode extracts season and episode numbers from an EpisodeID
// Episode format is expected to be like "S16E01"
func GetSeasonAndEpisode(episodeID EpisodeID) (season int, episode int, err error) {
	episodeStr := string(episodeID)
	if err := validateEpisodeFormat(episodeStr); err != nil {
		return 0, 0, err
	}

	// Extract season number from "S16" part
	seasonStr := episodeStr[1:3] // Skip 'S', get "16"
	season, err = strconv.Atoi(seasonStr)
	if err != nil {
		return 0, 0, fmt.Errorf("%w: invalid season number in episode %s: %w", ErrInvalidEpisodeFormat, episodeStr, err)
	}

	// Extract episode number from "E01" part
	episodeNumStr := episodeStr[4:6] // Skip "S16E", get "01"
	episode, err = strconv.Atoi(episodeNumStr)
	if err != nil {
		return 0, 0, fmt.Errorf("%w: invalid episode number in episode %s: %w", ErrInvalidEpisodeFormat, episodeStr, err)
	}

	return season, episode, nil
}

// GetImagePath constructs the image path for a SearchResult
func GetImagePath(result SearchResult) (string, error) {
	// Validate the episode format first
	episodeStr := string(result.EpisodID)
	if err := validateEpisodeFormat(episodeStr); err != nil {
		return "", err
	}

	return fmt.Sprintf("/img/%s/%s/medium.jpg", result.EpisodID, result.Timestamp), nil
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
func GetQuote(ctx context.Context, client *http.Client, config Config, quote string) ([]SearchResult, error) {
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
	var apiResults []SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&apiResults); err != nil {
		return nil, fmt.Errorf("error decoding JSON response: %w", err)
	}

	return apiResults, nil
}
