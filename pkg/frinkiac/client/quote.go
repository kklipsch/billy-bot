package client

import (
	"context"
	"fmt"
	"io"
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

// GetQuote searches for a quote on Frinkiac and returns the results
func (c *Client) GetQuote(ctx context.Context, quote string) ([]QuoteResult, error) {
	// Construct URL with query parameter
	u, err := url.Parse(fmt.Sprintf("%s/", c.baseURL))
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %w", err)
	}

	q := u.Query()
	q.Set("q", quote)
	u.RawQuery = q.Encode()

	// Create request
	requestURL := u.String()
	log.Debug().Str("url", requestURL).Str("quote", quote).Msg("sending quote request to frinkiac")

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
	results, err := parseQuoteResults(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	log.Debug().Int("result_count", len(results)).Str("quote", quote).Msg("parsed quote results from frinkiac")
	return results, nil
}

// parseQuoteResults parses the HTML response from a quote search
func parseQuoteResults(body io.Reader) ([]QuoteResult, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %w", err)
	}

	var results []QuoteResult

	// Find all result divs
	var findResults func(*html.Node)
	findResults = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" && hasClass(n, "result") {
			// Found a result div, now look for the image tag
			var findImage func(*html.Node) *QuoteResult
			findImage = func(n *html.Node) *QuoteResult {
				if n.Type == html.ElementNode && n.Data == "img" {
					// Found an image tag, extract the src attribute
					var src string
					for _, attr := range n.Attr {
						if attr.Key == "src" {
							src = attr.Val
							break
						}
					}

					if src != "" {
						// Parse the image path to extract season, episode, and ID
						result := parseImagePath(src)
						if result != nil {
							return result
						}
					}
				}

				// Recursively search for image tags
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if result := findImage(c); result != nil {
						return result
					}
				}

				return nil
			}

			if result := findImage(n); result != nil {
				results = append(results, *result)
			}
		}

		// Recursively search for result divs
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findResults(c)
		}
	}

	findResults(doc)
	return results, nil
}

// parseImagePath parses an image path to extract season, episode, and ID
// Example path: /img/S09E22/202334/medium.jpg
func parseImagePath(path string) *QuoteResult {
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		return nil
	}

	// Extract season and episode from the format S09E22
	seasonEpisode := parts[2]
	if len(seasonEpisode) < 6 {
		return nil
	}

	season := seasonEpisode[:3]  // S09
	episode := seasonEpisode[3:] // E22
	id := parts[3]               // 202334

	return &QuoteResult{
		ImagePath: path,
		Season:    season,
		Episode:   episode,
		ID:        id,
	}
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
