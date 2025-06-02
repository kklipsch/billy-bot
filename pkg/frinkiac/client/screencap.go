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
		Id              int    `json:"Id"`
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
		Id        int    `json:"Id"`
		Episode   string `json:"Episode"`
		Timestamp int    `json:"Timestamp"`
	} `json:"Frame"`
	Subtitles []struct {
		Id                      int    `json:"Id"`
		RepresentativeTimestamp int    `json:"RepresentativeTimestamp"`
		Episode                 string `json:"Episode"`
		StartTimestamp          int    `json:"StartTimestamp"`
		EndTimestamp            int    `json:"EndTimestamp"`
		Content                 string `json:"Content"`
		Language                string `json:"Language"`
	} `json:"Subtitles"`
	Nearby []struct {
		Id        int    `json:"Id"`
		Episode   string `json:"Episode"`
		Timestamp int    `json:"Timestamp"`
	} `json:"Nearby"`
}

// GetScreenCap gets a screen cap from Frinkiac using the JSON API
func (c *Client) GetScreenCap(ctx context.Context, season, episode, id string) (*ScreenCapResult, error) {
	// Construct URL with query parameters for the API endpoint
	u, err := url.Parse(fmt.Sprintf("%s/api/caption", c.baseURL))
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %w", err)
	}

	// Combine season and episode for the 'e' parameter (e.g., S16E01)
	episodeKey := fmt.Sprintf("%s%s", season, episode)

	q := u.Query()
	q.Set("e", episodeKey)
	q.Set("t", id)
	u.RawQuery = q.Encode()

	// Create request
	requestURL := u.String()
	log.Debug().Str("url", requestURL).Str("season", season).Str("episode", episode).Str("id", id).Msg("sending screen cap request to frinkiac API")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Debug().Int("status_code", resp.StatusCode).Str("url", requestURL).Msg("unexpected status code from frinkiac API")
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

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
