package client

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ExampleClient_GetQuote() {
	// Create a new client
	client := New()

	// Search for a quote
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	results, err := client.GetQuote(ctx, "garbage water")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Print the results
	for i, result := range results {
		fmt.Printf("Result %d:\n", i+1)
		fmt.Printf("  Season: %s\n", result.Season)
		fmt.Printf("  Episode: %s\n", result.Episode)
		fmt.Printf("  ID: %s\n", result.ID)
		fmt.Printf("  ImagePath: %s\n", result.ImagePath)
	}
}

func ExampleClient_GetScreenCap() {
	// Create a new client
	client := New()

	// Get a screen cap
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := client.GetScreenCap(ctx, "S09", "E22", "202334")
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

func TestClient(t *testing.T) {
	// This is just a placeholder test to ensure the package compiles
	// Real tests would make HTTP requests to the Frinkiac website or use mocks
	client := New()
	require.NotNil(t, client, "Failed to create client")
}

// TestParseAPIResponse tests parsing the API response from a file
func TestParseAPIResponse(t *testing.T) {
	// Read the test data file
	data, err := os.ReadFile("testdata/milhouse_api_response.json")
	require.NoError(t, err, "Failed to read test data file")

	// Parse the JSON data
	var apiResults []APISearchResult
	err = json.Unmarshal(data, &apiResults)
	require.NoError(t, err, "Failed to parse JSON data")

	// Verify that we have results
	require.NotEmpty(t, apiResults, "No results found in test data")

	// Verify that the first result is from season 16 episode 1
	firstResult := apiResults[0]
	assert.Equal(t, "S16E01", firstResult.Episode, "First result should be from S16E01")

	// Verify that several subsequent results are from season 10 episode 19
	s10e19Count := 0
	for i := 1; i < 10 && i < len(apiResults); i++ {
		if apiResults[i].Episode == "S10E19" {
			s10e19Count++
		}
	}
	assert.GreaterOrEqual(t, s10e19Count, 3, "Expected several subsequent results to be from S10E19")

	// Convert API results to QuoteResult objects (similar to what GetQuote does)
	results := make([]QuoteResult, 0, len(apiResults))
	for _, apiResult := range apiResults {
		// Extract season and episode from the format S16E01
		if len(apiResult.Episode) < 6 {
			t.Logf("Skipping result with invalid episode format: %s", apiResult.Episode)
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

	// Verify that the conversion worked correctly
	require.NotEmpty(t, results, "No results after conversion")

	// Verify that the first result is from season 16 episode 1
	assert.Equal(t, "S16", results[0].Season, "First result season should be S16")
	assert.Equal(t, "E01", results[0].Episode, "First result episode should be E01")

	// Verify that several subsequent results are from season 10 episode 19
	s10e19Count = 0
	for i := 1; i < 10 && i < len(results); i++ {
		if results[i].Season == "S10" && results[i].Episode == "E19" {
			s10e19Count++
		}
	}
	assert.GreaterOrEqual(t, s10e19Count, 3, "Expected several subsequent results to be from S10E19")
}

// TestParseAPICaptionResponse tests parsing the API caption response
func TestParseAPICaptionResponse(t *testing.T) {
	// Create a sample API caption response
	captionJSON := `{
		"Episode": {
			"Id": 671,
			"Key": "S16E01",
			"Season": 16,
			"EpisodeNumber": 1,
			"Title": "Treehouse of Horror XV",
			"Director": "David Silverman",
			"Writer": "Bill Odenkirk",
			"OriginalAirDate": "7-Nov-04",
			"WikiLink": "https://en.wikipedia.org/wiki/Treehouse_of_Horror_XV"
		},
		"Frame": {
			"Id": 2109531,
			"Episode": "S16E01",
			"Timestamp": 408242
		},
		"Subtitles": [
			{
				"Id": 172888,
				"RepresentativeTimestamp": 406573,
				"Episode": "S16E01",
				"StartTimestamp": 405767,
				"EndTimestamp": 407266,
				"Content": "( chuckling)",
				"Language": "en"
			},
			{
				"Id": 172889,
				"RepresentativeTimestamp": 408450,
				"Episode": "S16E01",
				"StartTimestamp": 407266,
				"EndTimestamp": 409400,
				"Content": "Everything's coming up Homer.",
				"Language": "en"
			},
			{
				"Id": 172890,
				"RepresentativeTimestamp": 410327,
				"Episode": "S16E01",
				"StartTimestamp": 409400,
				"EndTimestamp": 411767,
				"Content": "Yeah, well, the joke's on you, smart guy.",
				"Language": "en"
			}
		],
		"Nearby": [
			{
				"Id": 2109527,
				"Episode": "S16E01",
				"Timestamp": 407199
			},
			{
				"Id": 2109528,
				"Episode": "S16E01",
				"Timestamp": 407616
			},
			{
				"Id": 2109529,
				"Episode": "S16E01",
				"Timestamp": 408033
			},
			{
				"Id": 2109531,
				"Episode": "S16E01",
				"Timestamp": 408242
			},
			{
				"Id": 2109534,
				"Episode": "S16E01",
				"Timestamp": 408450
			},
			{
				"Id": 2109533,
				"Episode": "S16E01",
				"Timestamp": 408659
			},
			{
				"Id": 2109530,
				"Episode": "S16E01",
				"Timestamp": 408867
			}
		]
	}`

	// Parse the JSON data
	var apiCaption APICaption
	err := json.Unmarshal([]byte(captionJSON), &apiCaption)
	require.NoError(t, err, "Failed to parse JSON data")

	// Verify the episode details
	assert.Equal(t, 671, apiCaption.Episode.Id, "Episode ID should match")
	assert.Equal(t, "S16E01", apiCaption.Episode.Key, "Episode key should match")
	assert.Equal(t, 16, apiCaption.Episode.Season, "Season should match")
	assert.Equal(t, 1, apiCaption.Episode.EpisodeNumber, "Episode number should match")
	assert.Equal(t, "Treehouse of Horror XV", apiCaption.Episode.Title, "Title should match")

	// Verify the frame details
	assert.Equal(t, 2109531, apiCaption.Frame.Id, "Frame ID should match")
	assert.Equal(t, "S16E01", apiCaption.Frame.Episode, "Frame episode should match")
	assert.Equal(t, 408242, apiCaption.Frame.Timestamp, "Frame timestamp should match")

	// Verify the subtitles
	require.Len(t, apiCaption.Subtitles, 3, "Should have 3 subtitles")
	assert.Equal(t, "( chuckling)", apiCaption.Subtitles[0].Content, "First subtitle content should match")
	assert.Equal(t, "Everything's coming up Homer.", apiCaption.Subtitles[1].Content, "Second subtitle content should match")
	assert.Equal(t, "Yeah, well, the joke's on you, smart guy.", apiCaption.Subtitles[2].Content, "Third subtitle content should match")

	// Verify the nearby frames
	require.Len(t, apiCaption.Nearby, 7, "Should have 7 nearby frames")
	assert.Equal(t, 408242, apiCaption.Nearby[3].Timestamp, "Fourth nearby frame timestamp should match")

	// Test extracting caption text from subtitles
	var captionBuilder strings.Builder
	for i, subtitle := range apiCaption.Subtitles {
		if i > 0 {
			captionBuilder.WriteString(" ")
		}
		captionBuilder.WriteString(subtitle.Content)
	}
	caption := captionBuilder.String()
	assert.Equal(t, "( chuckling) Everything's coming up Homer. Yeah, well, the joke's on you, smart guy.", caption, "Combined caption should match")

	// Test constructing the image path
	imagePath := fmt.Sprintf("/img/%s/%d/medium.jpg", apiCaption.Frame.Episode, apiCaption.Frame.Timestamp)
	assert.Equal(t, "/img/S16E01/408242/medium.jpg", imagePath, "Image path should match")
}
