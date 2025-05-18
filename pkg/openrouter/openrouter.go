package openrouter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
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
	}

	requestJSON, err := json.Marshal(requestData)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(requestJSON))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("HTTP-Referer", "https://github.com/kklipsch/billy-bot")
	req.Header.Set("X-Title", "Billy Bot")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OpenRouter API returned status %d: %s", resp.StatusCode, string(body))
	}

	var openRouterResp OpenRouterResponse
	if err := json.Unmarshal(body, &openRouterResp); err != nil {
		return "", fmt.Errorf("error parsing response: %w", err)
	}

	if len(openRouterResp.Choices) == 0 {
		return "", fmt.Errorf("OpenRouter API returned no choices")
	}

	return openRouterResp.Choices[0].Message.Content, nil
}
