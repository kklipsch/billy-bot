package client

import (
	"context"
	"fmt"
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
