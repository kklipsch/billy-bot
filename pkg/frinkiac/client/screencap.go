package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
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

// GetScreenCap gets a screen cap from Frinkiac
func (c *Client) GetScreenCap(ctx context.Context, season, episode, id string) (*ScreenCapResult, error) {
	// Construct URL
	requestURL := fmt.Sprintf("%s/caption/%s%s/%s", c.baseURL, season, episode, id)
	log.Debug().Str("url", requestURL).Str("season", season).Str("episode", episode).Str("id", id).Msg("sending screen cap request to frinkiac")

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
		log.Debug().Int("status_code", resp.StatusCode).Str("url", requestURL).Msg("unexpected status code from frinkiac")
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse HTML response
	result, err := parseScreenCapResult(resp.Body, season, episode, id)
	if err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	log.Debug().Str("season", season).Str("episode", episode).Str("id", id).Str("caption", result.Caption).Msg("parsed screen cap result from frinkiac")
	return result, nil
}

// parseScreenCapResult parses the HTML response from a screen cap request
func parseScreenCapResult(body io.Reader, season, episode, id string) (*ScreenCapResult, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %w", err)
	}

	result := &ScreenCapResult{
		Season:  season,
		Episode: episode,
		ID:      id,
	}

	// Find the image path
	var findImage func(*html.Node) bool
	findImage = func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "img" {
			// Check if this is the main image
			for _, attr := range n.Attr {
				if attr.Key == "src" && strings.Contains(attr.Val, fmt.Sprintf("%s%s/%s", season, episode, id)) {
					result.ImagePath = attr.Val
					return true
				}
			}
		}

		// Recursively search for the image
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if findImage(c) {
				return true
			}
		}

		return false
	}

	// Find the caption
	var findCaption func(*html.Node) bool
	findCaption = func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "div" && hasClass(n, "caption") {
			var captionText strings.Builder
			var extractText func(*html.Node)
			extractText = func(n *html.Node) {
				if n.Type == html.TextNode {
					captionText.WriteString(n.Data)
				}
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					extractText(c)
				}
			}
			extractText(n)
			result.Caption = strings.TrimSpace(captionText.String())
			return true
		}

		// Recursively search for the caption
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if findCaption(c) {
				return true
			}
		}

		return false
	}

	findImage(doc)
	findCaption(doc)

	// If we didn't find an image path, construct one based on the season, episode, and ID
	if result.ImagePath == "" {
		result.ImagePath = fmt.Sprintf("/img/%s%s/%s/medium.jpg", season, episode, id)
	}

	return result, nil
}
