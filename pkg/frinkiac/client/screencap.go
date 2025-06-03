package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/html"
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
	// We need to convert the ID to a timestamp for the API
	// First, try to directly use the ID as a timestamp
	timestamp := id

	// If that doesn't work, try to search for the frame
	episodeKey := fmt.Sprintf("%s%s", season, episode)

	// Try to get the caption using the direct approach first
	u, err := url.Parse(fmt.Sprintf("%s/api/caption", c.baseURL))
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %w", err)
	}

	q := u.Query()
	q.Set("e", episodeKey)
	q.Set("t", timestamp)
	u.RawQuery = q.Encode()

	// Create request
	requestURL := u.String()
	log.Info().Str("url", requestURL).Str("season", season).Str("episode", episode).Str("timestamp", timestamp).Msg("sending screen cap request to frinkiac API")

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

	// If we get a 404, try to search for the frame using the episode
	if resp.StatusCode == http.StatusNotFound {
		log.Info().Str("season", season).Str("episode", episode).Str("id", id).Msg("direct caption lookup failed, trying search")

		// Try to get the frame by searching
		searchURL, err := url.Parse(fmt.Sprintf("%s/api/search", c.baseURL))
		if err != nil {
			return nil, fmt.Errorf("error parsing search URL: %w", err)
		}

		searchQ := searchURL.Query()
		searchQ.Set("q", episodeKey)
		searchURL.RawQuery = searchQ.Encode()

		searchReq, err := http.NewRequestWithContext(ctx, http.MethodGet, searchURL.String(), nil)
		if err != nil {
			return nil, fmt.Errorf("error creating search request: %w", err)
		}

		searchResp, err := c.httpClient.Do(searchReq)
		if err != nil {
			return nil, fmt.Errorf("error sending search request: %w", err)
		}
		defer searchResp.Body.Close()

		if searchResp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code from search API: %d", searchResp.StatusCode)
		}

		// Parse search results
		var searchResults []APISearchResult
		if err := json.NewDecoder(searchResp.Body).Decode(&searchResults); err != nil {
			return nil, fmt.Errorf("error decoding search response: %w", err)
		}

		// Find the frame with matching ID
		var foundTimestamp int
		idInt, err := strconv.Atoi(id)
		if err == nil {
			for _, result := range searchResults {
				if result.ID == idInt {
					foundTimestamp = result.Timestamp
					break
				}
			}
		}

		if foundTimestamp == 0 {
			// If we couldn't find the exact ID, use the first result
			if len(searchResults) > 0 {
				foundTimestamp = searchResults[0].Timestamp
			} else {
				return nil, fmt.Errorf("no frames found for episode %s", episodeKey)
			}
		}

		// Try again with the found timestamp
		u, err = url.Parse(fmt.Sprintf("%s/api/caption", c.baseURL))
		if err != nil {
			return nil, fmt.Errorf("error parsing URL: %w", err)
		}

		q = u.Query()
		q.Set("e", episodeKey)
		q.Set("t", fmt.Sprintf("%d", foundTimestamp))
		u.RawQuery = q.Encode()

		requestURL = u.String()
		log.Info().Str("url", requestURL).Str("season", season).Str("episode", episode).Int("timestamp", foundTimestamp).Msg("retrying screen cap request with found timestamp")

		req, err = http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %w", err)
		}

		resp, err = c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error sending request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}
	} else if resp.StatusCode != http.StatusOK {
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

// getScreenCapFromHTML gets a screen cap from Frinkiac using the HTML endpoint
func (c *Client) getScreenCapFromHTML(ctx context.Context, season, episode, id string) (*ScreenCapResult, error) {
	// Construct URL for the HTML endpoint
	requestURL := fmt.Sprintf("%s/caption/%s%s/%s", c.baseURL, season, episode, id)
	log.Info().Str("url", requestURL).Str("season", season).Str("episode", episode).Str("id", id).Msg("sending screen cap request to frinkiac HTML endpoint")

	// Create request
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
		log.Debug().Int("status_code", resp.StatusCode).Str("url", requestURL).Msg("unexpected status code from frinkiac HTML endpoint")
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse HTML to extract caption
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %w", err)
	}

	// Extract caption from HTML
	caption := extractCaptionFromHTML(doc)

	// Construct the image path directly
	imagePath := fmt.Sprintf("/img/%s%s/%s/medium.jpg", season, episode, id)

	result := &ScreenCapResult{
		ImagePath: imagePath,
		Caption:   caption,
		Season:    season,
		Episode:   episode,
		ID:        id,
	}

	log.Debug().Str("season", season).Str("episode", episode).Str("id", id).Str("caption", result.Caption).Msg("parsed screen cap result from frinkiac HTML endpoint")
	return result, nil
}

// extractCaptionFromHTML extracts the caption from the HTML document
func extractCaptionFromHTML(n *html.Node) string {
	// Look for the caption container which typically has a class like "caption-container" or similar
	var caption string
	var findCaption func(*html.Node)

	findCaption = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			// Check for caption container class
			for _, attr := range n.Attr {
				if attr.Key == "class" && strings.Contains(attr.Val, "caption") {
					// Found caption container, extract text
					var extractText func(*html.Node) string
					extractText = func(n *html.Node) string {
						if n.Type == html.TextNode {
							return n.Data
						}
						var result string
						for c := n.FirstChild; c != nil; c = c.NextSibling {
							result += extractText(c)
						}
						return result
					}

					caption = strings.TrimSpace(extractText(n))
					return
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findCaption(c)
		}
	}

	findCaption(n)
	return caption
}
