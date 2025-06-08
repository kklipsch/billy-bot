package frinkiac

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kklipsch/billy-bot/pkg/jsonschema"
	openrouter "github.com/kklipsch/billy-bot/pkg/openrouter"
)

var (
	// quotesResponseSchema defines the JSON schema for validating quote responses from the AI
	quotesResponseSchema = jsonschema.NewArraySchema(
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

	// quotesPrompt is the system prompt used to instruct the AI model about Simpsons quotes
	quotesPrompt = openrouter.ChatMessage{
		Role: "system",
		Content: `You are a helpful assistant with encyclopedic knowledge of The Simpsons. 
		You have access to a website called frinkiac that can find scenes from the Simpsons based on the text used in closed captioning of the Simpsons.
		Your goal is to categorize a set of text and think of any Simpsons quotes that are relevant to the text that should be findable in frinkiac.
		Your output should be a list JSON quote objects with a confidence score from 0 to 1.0 and a quote that is a good search term for the frinkiac tool.
		If you can identify the season and episode number, include those as well.
		You should sort the list by confidence score in descending order.`,
	}
)

// QuoteResponse represents a quote response from the AI model
type QuoteResponse struct {
	Quote      string  `json:"quote"`
	Confidence float64 `json:"confidence"`
	Character  string  `json:"character,omitempty"`
	Season     int     `json:"season,omitempty"`
	Episode    int     `json:"episode,omitempty"`
}

// GetCandidateQuotes fetches candidate Simpson quotes for a given prompt using OpenRouter AI
func GetCandidateQuotes(ctx context.Context, prompt, apiKey string) ([]QuoteResponse, error) {
	request := openrouter.ChatCompletionRequest{
		Model: "openrouter/auto",
		Messages: []openrouter.ChatMessage{
			quotesPrompt,
			{Role: "user", Content: prompt},
		},
		ResponseFormatEnabled: openrouter.NewResponseFormatEnabled("quotes", quotesResponseSchema),
		BaseRequest: openrouter.BaseRequest{
			Provider: &openrouter.ProviderRequest{
				RequireParameters: true,
			},
		},
	}

	req, err := openrouter.NewChatCompletionReq(ctx, request)
	result := openrouter.Call[openrouter.ChatCompletionResponse](ctx, apiKey, req, err, http.StatusOK)
	if result.Err != nil {
		return nil, result.Err
	}

	// Parse the response to extract quotes
	var quotes []QuoteResponse
	if err := json.Unmarshal([]byte(result.Body), &quotes); err != nil {
		// Try to extract from the choices if direct unmarshaling fails
		if len(result.Result.Choices) > 0 {
			content := result.Result.Choices[0].Message.Content
			if content != "" {
				if err := json.Unmarshal([]byte(content), &quotes); err != nil {
					return nil, fmt.Errorf("error parsing quotes from response content: %w", err)
				}
			}
		} else {
			return nil, fmt.Errorf("error parsing quotes from response: %w", err)
		}
	}

	if len(quotes) == 0 {
		return nil, fmt.Errorf("no quotes found in response")
	}

	return quotes, nil
}