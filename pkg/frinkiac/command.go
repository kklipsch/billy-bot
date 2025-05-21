package frinkiac

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kklipsch/billy-bot/pkg/config"
	openrouter "github.com/kklipsch/billy-bot/pkg/openrouter"
)

var (
	FrinkiacTool = openrouter.Tool{
		Type: "function",
		Function: openrouter.Function{
			Name:        "frinkiac",
			Description: "Search frinkiac for a scene",
			Parameters: openrouter.Parameters{
				Type: "object",
				Properties: map[string]openrouter.Property{
					"quote": {
						Type:        "string",
						Description: "The quote to search for.",
					},
					"confidence": {
						Type:        "number",
						Description: "The confidence score of the quote.",
					},
					"character": {
						Type:        "string",
						Description: "The character who said the quote.",
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

	FrinkiacPrompt = openrouter.ChatMessage{
		Role: "system",
		Content: `You are a helpful assistant with encyclopedic knowledge of The Simpsons. 
		You have access to a website called frinkiac that can find scenes from the Simpsons based on the text used in closed captioning of the Simpsons.
		Your goal is to categorize a set of text and think of any Simpsons quotes that are relevant to the text that should be findable in frinkiac.
		Your output should be a list JSON objects with a confidence score from 0 to 1.0 and a quote that is a good search term for the frinkiac tool.
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
		ToolsEnabled: openrouter.ToolsEnabled{
			Tools:      []openrouter.Tool{FrinkiacTool},
			ToolChoice: "auto",
		},
	}

	req, err := openrouter.NewChatCompletionReq(ctx, request)
	result := openrouter.OpenRouterCall[openrouter.ChatCompletionResponse](ctx, apiKey, req, err, http.StatusOK)
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
