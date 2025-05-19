package openrouter

import (
	"context"
	"fmt"
	"net/http"
	"os"
)

// Command represents the CLI command for OpenRouter
type Command struct {
	Prompt string `arg:"" help:"The prompt to send to the AI model."`
	Model  string `default:"openrouter/auto" help:"The model to use."`
	APIKey string `name:"api-key" short:"k" help:"OpenRouter API key. If not provided, OPENROUTER_API_KEY env var is used."`
}

// Run executes the OpenRouter command
func (o *Command) Run(ctx context.Context) error {
	// Get API key from flag or environment variable
	apiKey := o.APIKey
	if apiKey == "" {
		apiKey = os.Getenv("OPENROUTER_API_KEY")
		if apiKey == "" {
			return fmt.Errorf("API key not provided and OPENROUTER_API_KEY environment variable not set")
		}
		fmt.Println("Using API key from OPENROUTER_API_KEY environment variable")
	} else {
		fmt.Println("Using API key from command line flag")
	}

	fmt.Printf("Sending prompt to %s: %s\n", o.Model, o.Prompt)

	request := ChatCompletionRequest{
		Model:    o.Model,
		Messages: []ChatMessage{{Role: "user", Content: o.Prompt}},
		ToolsEnabled: ToolsEnabled{
			Tools: []Tool{
				{
					Type: "function",
					Function: Function{
						Name:        "get_code",
						Description: "Get code from a file",
						Parameters: Parameters{
							Type: "object",
							Properties: map[string]Property{
								"file_path": {
									Type:        "string",
									Description: "Path to the file",
								},
								"line_number": {
									Type:        "integer",
									Description: "Line number to retrieve",
								},
							},
							Required: []string{"file_path", "line_number"},
						},
					},
				},
			},
			ToolChoice: "auto",
		},
	}

	req, err := NewChatCompletionReq(ctx, request)
	result := OpenRouterCall[ChatCompletionResponse](ctx, apiKey, req, err, http.StatusOK)

	fmt.Println("Response from OpenRouter:")
	fmt.Println(result.Body)
	if result.Err != nil {
		return fmt.Errorf("error: %w", result.Err)
	}

	return nil
}
