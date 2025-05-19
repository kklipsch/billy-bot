package openrouter

import (
	"context"
	"net/http"

	"github.com/kklipsch/billy-bot/pkg/config"
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
	if result.Err != nil {
		return result.Err
	}

	return nil
}
