package client

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"
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
	if client == nil {
		t.Fatal("Failed to create client")
	}
}

// TestParseAPIResponse tests parsing the API response from a file
func TestParseAPIResponse(t *testing.T) {
	// Read the test data file
	data, err := os.ReadFile("testdata/milhouse_api_response.json")
	if err != nil {
		t.Fatalf("Failed to read test data file: %v", err)
	}

	// Parse the JSON data
	var apiResults []APISearchResult
	if err := json.Unmarshal(data, &apiResults); err != nil {
		t.Fatalf("Failed to parse JSON data: %v", err)
	}

	// Verify that we have results
	if len(apiResults) == 0 {
		t.Fatal("No results found in test data")
	}

	// Verify that the first result is from season 16 episode 1
	firstResult := apiResults[0]
	if firstResult.Episode != "S16E01" {
		t.Errorf("Expected first result to be from S16E01, got %s", firstResult.Episode)
	}

	// Verify that several subsequent results are from season 10 episode 19
	s10e19Count := 0
	for i := 1; i < 10 && i < len(apiResults); i++ {
		if apiResults[i].Episode == "S10E19" {
			s10e19Count++
		}
	}
	if s10e19Count < 3 {
		t.Errorf("Expected several subsequent results to be from S10E19, got %d", s10e19Count)
	}

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
		id := fmt.Sprintf("%d", apiResult.Id)

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
	if len(results) == 0 {
		t.Fatal("No results after conversion")
	}

	// Verify that the first result is from season 16 episode 1
	if results[0].Season != "S16" || results[0].Episode != "E01" {
		t.Errorf("Expected first result to be from S16E01, got %s%s", results[0].Season, results[0].Episode)
	}

	// Verify that several subsequent results are from season 10 episode 19
	s10e19Count = 0
	for i := 1; i < 10 && i < len(results); i++ {
		if results[i].Season == "S10" && results[i].Episode == "E19" {
			s10e19Count++
		}
	}
	if s10e19Count < 3 {
		t.Errorf("Expected several subsequent results to be from S10E19, got %d", s10e19Count)
	}
}
