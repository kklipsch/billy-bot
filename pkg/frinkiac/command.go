package frinkiac

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kklipsch/billy-bot/pkg/config"
	"github.com/kklipsch/billy-bot/pkg/frinkiac/client"
	"github.com/kklipsch/billy-bot/pkg/jsonschema"
	openrouter "github.com/kklipsch/billy-bot/pkg/openrouter"
)

var (
	FrinkiacResponseSchema = jsonschema.NewArraySchema(
		jsonschema.NewObjectSchema(
			map[string]*jsonschema.Schema{
				"quote":      jsonschema.NewStringSchema(),
				"confidence": jsonschema.NewNumberSchema(),
				"character":  jsonschema.NewStringSchema(),
				"season":     jsonschema.NewIntegerSchema(),
				"episode":    jsonschema.NewIntegerSchema(),
			},
			[]string{"quote", "confidence"},
		),
	)

	FrinkiacPrompt = openrouter.ChatMessage{
		Role: "system",
		Content: `You are a helpful assistant with encyclopedic knowledge of The Simpsons. 
		You have access to a website called frinkiac that can find scenes from the Simpsons based on the text used in closed captioning of the Simpsons.
		Your goal is to categorize a set of text and think of any Simpsons quotes that are relevant to the text that should be findable in frinkiac.
		Your output should be a list JSON quote objects with a confidence score from 0 to 1.0 and a quote that is a good search term for the frinkiac tool.
		If you can identify the season and episode number, include those as well.
		You should sort the list by confidence score in descending order.`,
	}
)

// Command represents the CLI command for OpenRouter
type Command struct {
	Prompt string `arg:"" help:"The prompt to send to the AI model."`
	Model  string `default:"openrouter/auto" help:"The model to use."`
	APIKey string `name:"api-key" short:"k" help:"OpenRouter API key. If not provided, OPENROUTER_API_KEY env var is used."`
}

// QuoteResponse represents a quote response from the AI model
type QuoteResponse struct {
	Quote      string  `json:"quote"`
	Confidence float64 `json:"confidence"`
	Character  string  `json:"character,omitempty"`
	Season     int     `json:"season,omitempty"`
	Episode    int     `json:"episode,omitempty"`
}

// Run executes the OpenRouter command
func (o *Command) Run(ctx context.Context) error {
	apiKey, err := config.GetFlagOrEnvVar(o.APIKey, "OPENROUTER_API_KEY")
	if err != nil {
		return err
	}

	request := openrouter.ChatCompletionRequest{
		Model: o.Model,
		Messages: []openrouter.ChatMessage{
			FrinkiacPrompt,
			{Role: "user", Content: o.Prompt},
		},
		ResponseFormatEnabled: openrouter.NewResponseFormatEnabled("quotes", FrinkiacResponseSchema),
		BaseRequest: openrouter.BaseRequest{
			Provider: &openrouter.ProviderRequest{
				RequireParameters: true,
			},
		},
	}

	req, err := openrouter.NewChatCompletionReq(ctx, request)
	result := openrouter.OpenRouterCall[openrouter.ChatCompletionResponse](ctx, apiKey, req, err, http.StatusOK)
	if result.Err != nil {
		return result.Err
	}

	// Parse the response to extract quotes
	var quotes []QuoteResponse
	if err := json.Unmarshal([]byte(result.Body), &quotes); err != nil {
		// Try to extract from the choices if direct unmarshaling fails
		if len(result.Result.Choices) > 0 {
			content := result.Result.Choices[0].Message.Content
			if content != "" {
				if err := json.Unmarshal([]byte(content), &quotes); err != nil {
					return fmt.Errorf("error parsing quotes from response content: %w", err)
				}
			}
		} else {
			return fmt.Errorf("error parsing quotes from response: %w", err)
		}
	}

	if len(quotes) == 0 {
		return fmt.Errorf("no quotes found in response")
	}

	// Create a Frinkiac client
	frinkiacClient := client.New()

	// Process each quote
	fmt.Println("Quotes found:")
	for i, quote := range quotes {
		fmt.Printf("%d. %s (confidence: %.2f)\n", i+1, quote.Quote, quote.Confidence)

		// If season and episode are provided, try to get the screen cap directly
		if quote.Season > 0 && quote.Episode > 0 {
			season := fmt.Sprintf("S%02d", quote.Season)
			episode := fmt.Sprintf("E%02d", quote.Episode)

			fmt.Printf("   Season %s, Episode %s provided by AI\n", season, episode)

			// We don't have the ID, so we need to search for the quote first
			results, err := frinkiacClient.GetQuote(ctx, quote.Quote)
			if err != nil {
				fmt.Printf("   Error searching for quote: %v\n", err)
				continue
			}

			if len(results) > 0 {
				// Use the first result
				result := results[0]
				fmt.Printf("   Found screen cap: Season %s, Episode %s, ID %s\n", result.Season, result.Episode, result.ID)

				// Get the screen cap
				screenCap, err := frinkiacClient.GetScreenCap(ctx, result.Season, result.Episode, result.ID)
				if err != nil {
					fmt.Printf("   Error getting screen cap: %v\n", err)
					continue
				}

				fmt.Printf("   Caption: %s\n", screenCap.Caption)
				fmt.Printf("   Image URL: %s%s\n", client.BaseURL, screenCap.ImagePath)
			} else {
				fmt.Println("   No screen caps found for this quote")
			}
		} else {
			// Search for the quote
			results, err := frinkiacClient.GetQuote(ctx, quote.Quote)
			if err != nil {
				fmt.Printf("   Error searching for quote: %v\n", err)
				continue
			}

			if len(results) > 0 {
				// Use the first result
				result := results[0]
				fmt.Printf("   Found screen cap: Season %s, Episode %s, ID %s\n", result.Season, result.Episode, result.ID)

				// Get the screen cap
				screenCap, err := frinkiacClient.GetScreenCap(ctx, result.Season, result.Episode, result.ID)
				if err != nil {
					fmt.Printf("   Error getting screen cap: %v\n", err)
					continue
				}

				fmt.Printf("   Caption: %s\n", screenCap.Caption)
				fmt.Printf("   Image URL: %s%s\n", client.BaseURL, screenCap.ImagePath)
			} else {
				fmt.Println("   No screen caps found for this quote")
			}
		}

		fmt.Println()
	}

	return nil
}
