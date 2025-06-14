package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ExampleGetQuote() {
	// Create a new client and config
	client := NewHTTPClient()
	config := DefaultConfig()

	// Search for a quote
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	results, err := GetQuote(ctx, client, config, "garbage water")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Print the results
	for i, result := range results {
		season, episode, err := GetSeasonAndEpisode(result.EpisodID)
		if err != nil {
			fmt.Printf("Error parsing season/episode for result %d: %v\n", i+1, err)
			continue
		}

		imagePath, err := GetImagePath(result.EpisodID, result.Timestamp)
		if err != nil {
			fmt.Printf("Error getting image path for result %d: %v\n", i+1, err)
			continue
		}

		fmt.Printf("Result %d:\n", i+1)
		fmt.Printf("  Season: %d\n", season)
		fmt.Printf("  Episode: %d\n", episode)
		fmt.Printf("  ID: %s\n", result.Timestamp)
		fmt.Printf("  ImagePath: %s\n", imagePath)
	}
}

func ExampleGetScreenCap() {
	// Create a new client and config
	client := NewHTTPClient()
	config := DefaultConfig()

	// Get a screen cap
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := GetScreenCap(ctx, client, config, 9, 22, Timestamp("202334"))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Print the result
	fmt.Printf("Season: %s\n", result.Season)
	fmt.Printf("Episode: %s\n", result.Episode)
	fmt.Printf("ID: %s\n", result.ID)
	fmt.Printf("ImagePath: %s\n", result.ImagePath)
	fmt.Printf("Caption: %s\n", result.Caption)
}

func TestHTTPClient(t *testing.T) {
	// This is just a placeholder test to ensure the package compiles
	// Real tests would make HTTP requests to the Frinkiac website or use mocks
	client := NewHTTPClient()
	require.NotNil(t, client, "Failed to create client")

	config := DefaultConfig()
	require.NotEmpty(t, config.BaseURL, "Config should have a base URL")
}

// TestDoRequest tests the doRequest function with a mock HTTP client
func TestDoRequest(t *testing.T) {
	// Create a mock HTTP client that returns a predefined response
	mockClient := &http.Client{
		Transport: &mockTransport{
			response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"test": "data"}`)),
			},
		},
	}

	// Create a config
	config := DefaultConfig()

	// Make a request
	ctx := context.Background()
	resp, err := doRequest(ctx, mockClient, config, RequestOptions{
		Method: http.MethodGet,
		Path:   "/test",
		LogContext: map[string]interface{}{
			"test": "value",
		},
	})

	// Verify the response
	require.NoError(t, err, "doRequest should not return an error")
	require.NotNil(t, resp, "doRequest should return a response")
	require.Equal(t, http.StatusOK, resp.StatusCode, "Response status code should be 200")

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "Failed to read response body")
	resp.Body.Close()

	// Verify the response body
	assert.Equal(t, `{"test": "data"}`, string(body), "Response body should match expected value")
}

// TestDoRequestError tests the doRequest function with an error response
func TestDoRequestError(t *testing.T) {
	// Create a mock HTTP client that returns an error response
	mockClient := &http.Client{
		Transport: &mockTransport{
			response: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(strings.NewReader(`{"error": "bad request"}`)),
			},
		},
	}

	// Create a config
	config := DefaultConfig()

	// Make a request
	ctx := context.Background()
	_, err := doRequest(ctx, mockClient, config, RequestOptions{
		Method: http.MethodGet,
		Path:   "/test",
	})

	// Verify the error
	require.Error(t, err, "doRequest should return an error for non-200 status codes")
	assert.Contains(t, err.Error(), "400", "Error should contain the status code")
	assert.Contains(t, err.Error(), "bad request", "Error should contain the response body")
}

// mockTransport is a mock http.RoundTripper that returns a predefined response
type mockTransport struct {
	response *http.Response
	err      error
}

func (m *mockTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return m.response, m.err
}

// TestParseAPIResponse tests parsing the API response from a file
func TestParseAPIResponse(t *testing.T) {
	// Read the test data file
	data, err := os.ReadFile("testdata/milhouse_api_response.json")
	require.NoError(t, err, "Failed to read test data file")

	// Parse the JSON data
	var apiResults []SearchResult
	err = json.Unmarshal(data, &apiResults)
	require.NoError(t, err, "Failed to parse JSON data")

	// Verify that we have results
	require.NotEmpty(t, apiResults, "No results found in test data")

	// Verify that the first result is from season 16 episode 1
	firstResult := apiResults[0]
	assert.Equal(t, EpisodeID("S16E01"), firstResult.EpisodID, "First result should be from S16E01")

	// Test GetSeasonAndEpisode function with the first result
	season, episode, err := GetSeasonAndEpisode(firstResult.EpisodID)
	require.NoError(t, err, "Failed to parse season and episode from first result")
	assert.Equal(t, 16, season, "First result season should be 16")
	assert.Equal(t, 1, episode, "First result episode should be 1")

	// Test GetImagePath function with the first result
	imagePath, err := GetImagePath(firstResult.EpisodID, firstResult.Timestamp)
	require.NoError(t, err, "Failed to get image path from first result")
	expectedPath := fmt.Sprintf("/img/S16E01/%s/medium.jpg", firstResult.Timestamp)
	assert.Equal(t, expectedPath, imagePath, "Image path should match expected format")

	// Verify that several subsequent results are from season 10 episode 19
	s10e19Count := 0
	for i := 1; i < 10 && i < len(apiResults); i++ {
		if apiResults[i].EpisodID == "S10E19" {
			s10e19Count++

			// Test utility functions on S10E19 results
			season, episode, err := GetSeasonAndEpisode(apiResults[i].EpisodID)
			require.NoError(t, err, "Failed to parse season and episode from S10E19 result")
			assert.Equal(t, 10, season, "S10E19 result season should be 10")
			assert.Equal(t, 19, episode, "S10E19 result episode should be 19")

			imagePath, err := GetImagePath(apiResults[i].EpisodID, apiResults[i].Timestamp)
			require.NoError(t, err, "Failed to get image path from S10E19 result")
			expectedPath := fmt.Sprintf("/img/S10E19/%s/medium.jpg", apiResults[i].Timestamp)
			assert.Equal(t, expectedPath, imagePath, "S10E19 image path should match expected format")
		}
	}
	assert.GreaterOrEqual(t, s10e19Count, 3, "Expected several subsequent results to be from S10E19")
}

// TestGetSeasonAndEpisode tests the GetSeasonAndEpisode function with various inputs
func TestGetSeasonAndEpisode(t *testing.T) {
	tests := []struct {
		name        string
		episodeID   EpisodeID
		expectError bool
		season      int
		episode     int
	}{
		{
			name:        "Valid S16E01",
			episodeID:   EpisodeID("S16E01"),
			expectError: false,
			season:      16,
			episode:     1,
		},
		{
			name:        "Valid S10E19",
			episodeID:   EpisodeID("S10E19"),
			expectError: false,
			season:      10,
			episode:     19,
		},
		{
			name:        "Invalid format - too short",
			episodeID:   EpisodeID("S16E"),
			expectError: true,
		},
		{
			name:        "Invalid format - no S",
			episodeID:   EpisodeID("16E01"),
			expectError: true,
		},
		{
			name:        "Invalid format - no E",
			episodeID:   EpisodeID("S1601"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			season, episode, err := GetSeasonAndEpisode(tt.episodeID)

			if tt.expectError {
				assert.Error(t, err, "Expected error for invalid input")
			} else {
				assert.NoError(t, err, "Expected no error for valid input")
				assert.Equal(t, tt.season, season, "Season should match expected value")
				assert.Equal(t, tt.episode, episode, "Episode should match expected value")
			}
		})
	}
}

// TestGetImagePath tests the GetImagePath function
func TestGetImagePath(t *testing.T) {
	tests := []struct {
		name        string
		episodeID   EpisodeID
		timestamp   Timestamp
		expectError bool
		expected    string
	}{
		{
			name:        "Valid result",
			episodeID:   "S16E01",
			timestamp:   "123456",
			expectError: false,
			expected:    "/img/S16E01/123456/medium.jpg",
		},
		{
			name:        "Invalid episode format",
			episodeID:   "S16E",
			timestamp:   "123456",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imagePath, err := GetImagePath(tt.episodeID, tt.timestamp)

			if tt.expectError {
				assert.Error(t, err, "Expected error for invalid input")
			} else {
				assert.NoError(t, err, "Expected no error for valid input")
				assert.Equal(t, tt.expected, imagePath, "Image path should match expected value")
			}
		})
	}
}
