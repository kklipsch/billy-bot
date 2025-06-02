package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/rs/zerolog/log"
)

// ScreenCapResult represents the result of a screen cap request
type ScreenCapResult struct {
	ImagePath string
	Caption   string
	Season    string
	Episode   string
	ID        string
}

// APICaption represents the response from the Frinkiac API caption endpoint
type APICaption struct {
	Episode struct {
		ID              int    `json:"Id"`
		Key             string `json:"Key"`
		Season          int    `json:"Season"`
		EpisodeNumber   int    `json:"EpisodeNumber"`
		Title           string `json:"Title"`
		Director        string `json:"Director"`
		Writer          string `json:"Writer"`
		OriginalAirDate string `json:"OriginalAirDate"`
		WikiLink        string `json:"WikiLink"`
	} `json:"Episode"`
	Frame struct {
		ID        int    `json:"Id"`
		Episode   string `json:"Episode"`
		Timestamp int    `json:"Timestamp"`
	} `json:"Frame"`
	Subtitles []struct {
		ID                      int    `json:"Id"`
		RepresentativeTimestamp int    `json:"RepresentativeTimestamp"`
		Episode                 string `json:"Episode"`
		StartTimestamp          int    `json:"StartTimestamp"`
		EndTimestamp            int    `json:"EndTimestamp"`
		Content                 string `json:"Content"`
		Language                string `json:"Language"`
	} `json:"Subtitles"`
	Nearby []struct {
		ID        int    `json:"Id"`
		Episode   string `json:"Episode"`
		Timestamp int    `json:"Timestamp"`
	} `json:"Nearby"`
}

// GetScreenCap gets a screen cap from Frinkiac
func (c *Client) GetScreenCap(ctx context.Context, season, episode, id string) (*ScreenCapResult, error) {
	// First try the JSON API endpoint
	result, err := c.getScreenCapFromAPI(ctx, season, episode, id)
	if err != nil {
		// If the API endpoint fails, fall back to the HTML endpoint
		log.Info().Str("season", season).Str("episode", episode).Str("id", id).Msg("API endpoint failed, falling back to HTML endpoint")
		return c.getScreenCapFromHTML(ctx, season, episode, id)
	}
	return result, nil
}

// getScreenCapFromAPI gets a screen cap from Frinkiac using the JSON API
func (c *Client) getScreenCapFromAPI(ctx context.Context, season, episode, id string) (*ScreenCapResult, error) {
	// Combine season and episode for the 'e' parameter (e.g., S16E01)
	episodeKey := fmt.Sprintf("%s%s", season, episode)

	// Set up query parameters
	queryParams := url.Values{}
	queryParams.Set("e", episodeKey)
	queryParams.Set("t", id)

	// Set up log context
	logContext := map[string]interface{}{
		"season":  season,
		"episode": episode,
		"id":      id,
	}

	// Make the request
	resp, err := c.doRequest(ctx, RequestOptions{
		Method:      http.MethodGet,
		Path:        "/api/caption",
		QueryParams: queryParams,
		LogContext:  logContext,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse JSON response
	var apiCaption APICaption
	if err := json.NewDecoder(resp.Body).Decode(&apiCaption); err != nil {
		return nil, fmt.Errorf("error decoding JSON response: %w", err)
	}

	// Extract caption text from subtitles
	var captionBuilder strings.Builder
	for _, subtitle := range apiCaption.Subtitles {
		if captionBuilder.Len() > 0 {
			captionBuilder.WriteString(" ")
		}
		captionBuilder.WriteString(subtitle.Content)
	}
	caption := captionBuilder.String()

	// Construct the image path
	imagePath := fmt.Sprintf("/img/%s/%d/medium.jpg", apiCaption.Frame.Episode, apiCaption.Frame.Timestamp)

	result := &ScreenCapResult{
		ImagePath: imagePath,
		Caption:   caption,
		Season:    season,
		Episode:   episode,
		ID:        id,
	}

	log.Debug().Str("season", season).Str("episode", episode).Str("id", id).Str("caption", result.Caption).Msg("parsed screen cap result from frinkiac API")
	return result, nil
}

// getScreenCapFromHTML gets a screen cap from Frinkiac using the HTML endpoint
func (c *Client) getScreenCapFromHTML(ctx context.Context, season, episode, id string) (*ScreenCapResult, error) {
	// Set up log context
	logContext := map[string]interface{}{
		"season":  season,
		"episode": episode,
		"id":      id,
		"type":    "HTML",
	}

	// Make the request
	path := fmt.Sprintf("/caption/%s%s/%s", season, episode, id)
	resp, err := c.doRequest(ctx, RequestOptions{
		Method:     http.MethodGet,
		Path:       path,
		LogContext: logContext,
	})
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	// Construct the image path directly
	imagePath := fmt.Sprintf("/img/%s%s/%s/medium.jpg", season, episode, id)

	// Create a result with an empty caption
	result := &ScreenCapResult{
		ImagePath: imagePath,
		Caption:   "",
		Season:    season,
		Episode:   episode,
		ID:        id,
	}

	log.Debug().Str("season", season).Str("episode", episode).Str("id", id).Str("caption", result.Caption).Msg("created screen cap result from frinkiac HTML endpoint")
	return result, nil
}
