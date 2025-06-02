package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/html"
)

// QuoteResult represents a single result from a quote search
type QuoteResult struct {
	ImagePath string // e.g., /img/S09E22/202334/medium.jpg
	Season    string // e.g., S09
	Episode   string // e.g., E22
	ID        string // e.g., 202334
}

// APISearchResult represents a single result from the Frinkiac API search endpoint
type APISearchResult struct {
	ID        int    `json:"Id"`
	Episode   string `json:"Episode"` // e.g., S16E01
	Timestamp int    `json:"Timestamp"`
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
func (c *Client) GetQuote(ctx context.Context, quote string) ([]QuoteResult, error) {
	// Construct URL with query parameter for the API endpoint
	u, err := url.Parse(fmt.Sprintf("%s/api/search", c.baseURL))
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %w", err)
	}

	q := u.Query()
	q.Set("q", quote)
	u.RawQuery = q.Encode()

	// Create request
	requestURL := u.String()
	log.Debug().Str("url", requestURL).Str("quote", quote).Msg("sending quote request to frinkiac API")

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
	var apiResults []APISearchResult
	if err := json.NewDecoder(resp.Body).Decode(&apiResults); err != nil {
		return nil, fmt.Errorf("error decoding JSON response: %w", err)
	}

	// Convert API results to QuoteResult objects
	results := make([]QuoteResult, 0, len(apiResults))
	for _, apiResult := range apiResults {
		// Extract season and episode from the format S16E01
		if len(apiResult.Episode) < 6 {
			log.Debug().Str("episode", apiResult.Episode).Msg("invalid episode format")
			continue
		}

		season := apiResult.Episode[:3]  // S16
		episode := apiResult.Episode[3:] // E01
		id := fmt.Sprintf("%d", apiResult.ID)

		// Construct the image path
		imagePath := fmt.Sprintf("/img/%s/%s/medium.jpg", apiResult.Episode, id)

		results = append(results, QuoteResult{
			ImagePath: imagePath,
			Season:    season,
			Episode:   episode,
			ID:        id,
		})
	}

	log.Debug().Int("result_count", len(results)).Str("quote", quote).Msg("parsed quote results from frinkiac API")
	return results, nil
}
