package openrouter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
<<<<<<< HEAD
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

	response, err := o.callOpenRouter(ctx, apiKey, o.Model, o.Prompt)
	if err != nil {
		return err
	}

	fmt.Println("Response from OpenRouter:")
	fmt.Println(response)
	return nil
}

// ChatMessage represents a message in a chat conversation
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenRouterRequest represents a request to the OpenRouter API
type OpenRouterRequest struct {
	Model      string        `json:"model"`
	Messages   []ChatMessage `json:"messages"`
	Tools      []Tool        `json:"tools,omitempty"`
	ToolChoice string        `json:"tool_choice,omitempty"`
}

type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

// Function represents a function that can be called by the model
type Function struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Parameters  Parameters `json:"parameters"`
}

// Parameters defines the schema for function parameters
type Parameters struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required,omitempty"`
}

// Property defines a property in a parameter schema
type Property struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
}

// ToolChoice specifies how the model should use tools
type ToolChoice struct {
	Type     string `json:"type,omitempty"`
	Function struct {
		Name string `json:"name,omitempty"`
	} `json:"function,omitempty"`
}

// OpenRouterResponse represents a response from the OpenRouter API
type OpenRouterResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// callOpenRouter sends a request to the OpenRouter API
func (o *Command) callOpenRouter(ctx context.Context, apiKey, model, prompt string) (string, error) {
	requestData := OpenRouterRequest{
		Model: model,
		Messages: []ChatMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
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
	}

	requestJSON, err := json.Marshal(requestData)
=======
	"slices"
	"strings"
)

func OpenRouterCall[T any](ctx context.Context, apiKey string, req *http.Request, err error, allowedStatus ...int) Response[T] {
>>>>>>> 9de49fa (move command into its own fiel)
	if err != nil {
		return Response[T]{Err: fmt.Errorf("error creating request: %w", err)}
	}

	AddDefaultHeaders(apiKey, req)
	resp, err := http.DefaultClient.Do(req)
	return FromResponse[T](ctx, resp, err, allowedStatus...)
}

func AddDefaultHeaders(APIKey string, req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+APIKey)
	req.Header.Set("HTTP-Referer", "https://github.com/kklipsch/billy-bot")
	req.Header.Set("X-Title", "Billy Bot")
}

type Response[T any] struct {
	Body   string
	Err    error
	Result T
}

func FromResponse[T any](ctx context.Context, resp *http.Response, err error, allowedStatus ...int) (oresp Response[T]) {
	oresp = Response[T]{}

	if err != nil {
		oresp.Err = fmt.Errorf("error sending request: %w", err)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		oresp.Err = fmt.Errorf("error reading response body: %w", err)
		return
	}
	defer resp.Body.Close()

	oresp.Body = strings.TrimSpace(string(body))

	if !slices.Contains(allowedStatus, resp.StatusCode) {
		oresp.Err = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		return
	}

	if err = json.Unmarshal(body, &oresp.Result); err != nil {
		oresp.Err = fmt.Errorf("error unmarshaling response: %w", err)
	}

	return oresp
}
