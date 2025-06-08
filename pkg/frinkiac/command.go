package frinkiac

import (
	"context"
	"fmt"

	"github.com/kklipsch/billy-bot/pkg/config"
	"github.com/kklipsch/billy-bot/pkg/frinkiac/ai"
	"github.com/kklipsch/billy-bot/pkg/frinkiac/http"
)

// Command represents the CLI command for OpenRouter
type Command struct {
	Prompt string `arg:"" help:"The prompt to send to the AI model."`
	Model  string `default:"openrouter/auto" help:"The model to use."`
	APIKey string `name:"api-key" short:"k" help:"OpenRouter API key. If not provided, OPENROUTER_API_KEY env var is used."`
}

// Run executes the OpenRouter command
func (o *Command) Run(ctx context.Context) error {
	apiKey, err := config.GetFlagOrEnvVar(o.APIKey, "OPENROUTER_API_KEY")
	if err != nil {
		return err
	}

	quotes, err := ai.GetCandidateQuotes(ctx, o.Prompt, apiKey)
	if err != nil {
		return err
	}

	// Create a Frinkiac HTTP client and config
	client := http.NewHTTPClient()
	config := http.DefaultConfig()

	// Process each quote
	fmt.Println("Quotes found:")
	for i, quote := range quotes {
		fmt.Printf("%d. %s (confidence: %.2f) [S%02d E%02d]\n", i+1, quote.Quote, quote.Confidence, quote.Season, quote.Episode)

		// Search for the quote
		results, err := http.GetQuote(ctx, client, config, quote.Quote)
		if err != nil {
			fmt.Printf("   Error searching for quote: %v\n", err)
			continue
		}

		if len(results) > 0 {
			// Use the first result
			result := results[0]
			fmt.Printf("   Found screen cap: Season %s, Episode %s, ID %s\n", result.Season, result.Episode, result.ID)

			// Get the screen cap
			screenCap, err := http.GetScreenCap(ctx, client, config, result.Season, result.Episode, result.ID)
			if err != nil {
				fmt.Printf("   Error getting screen cap: %v\n", err)
				continue
			}

			fmt.Printf("   Caption: %s\n", screenCap.Caption)
			fmt.Printf("   Image URL: %s%s\n", http.BaseURL, screenCap.ImagePath)
		} else {
			fmt.Println("   No screen caps found for this quote")
		}

		fmt.Println()
	}

	return nil
}
