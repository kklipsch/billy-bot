package openrouter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kklipsch/billy-bot/pkg/config"
)

var (
	FrinkiacTool = Tool{
		Type: "function",
		Function: Function{
			Name:        "frinkiac",
			Description: "Search frinkiac for a scene",
			Parameters: Parameters{
				Type: "object",
				Properties: map[string]Property{
					"quote": {
						Type:        "string",
						Description: "The quote to search for.",
					},
					"confidence": {
						Type:        "number",
						Description: "The confidence score of the quote.",
					},
					"season": {
						Type:        "integer",
						Description: "The season number.",
					},
					"episode": {
						Type:        "integer",
						Description: "The episode number.",
					},
				},
				Required: []string{"quote"},
			},
		},
	}

	FrinkiacPrompt = ChatMessage{
		Role: "system",
		Content: `You are a helpful assistant with encyclopedic knowledge of The Simpsons. 
		You have access to a tool called frinkiac that can find scenes from the Simppsons based on quotes.
		Your goal is to categorize a set of text and think of any Simpsons quotes that are relevant to the text.
		Your output should be a list JSON objects with a confidence score from 0 to 1.0 and a quote that is a good search term for the frinkiac tool.
		If you can identify the season and episode number, include those as well.
		You should sort the list by confidence score in descending order.`,
	}
)

// Frinkiac represents the CLI command for OpenRouter
type Frinkiac struct {
	Prompt string `arg:"" help:"The prompt to send to the AI model."`
	Model  string `default:"openrouter/auto" help:"The model to use."`
	APIKey string `name:"api-key" short:"k" help:"OpenRouter API key. If not provided, OPENROUTER_API_KEY env var is used."`
}

// Run executes the OpenRouter command
func (o *Frinkiac) Run(ctx context.Context) error {
	apiKey, err := config.GetFlagOrEnvVar(o.APIKey, "OPENROUTER_API_KEY")
	if err != nil {
		return err
	}

	request := ChatCompletionRequest{
		Model: o.Model,
		Messages: []ChatMessage{
			FrinkiacPrompt,
			{Role: "user", Content: o.Prompt},
		},
		ToolsEnabled: ToolsEnabled{
			Tools:      []Tool{FrinkiacTool},
			ToolChoice: "auto",
		},
	}

	req, err := NewChatCompletionReq(ctx, request)
	result := OpenRouterCall[ChatCompletionResponse](ctx, apiKey, req, err, http.StatusOK)
	if result.Err != nil {
		return result.Err
	}

	fmt.Println(result.Body)
	/*

		if len(result.Result.Choices) == 0 {
			return fmt.Errorf("no choices returned from OpenRouter")
		}

		for _, choice := range result.Result.Choices {
			fmt.Printf("Choice: %v\n", choice.Message)
			for _, toolCall := range choice.Message.ToolCalls {
				if toolCall.Type != "function" {
					continue
				}

				if toolCall.Function.Name != FrinkiacTool.Function.Name {
					continue
				}

				fmt.Printf("Tool call: %s\n", toolCall.Function.Name)
				fmt.Printf("Arguments: %s\n", toolCall.Function.Arguments)
			}
		}
	*/

	return nil
}
