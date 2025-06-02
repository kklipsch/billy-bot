package openrouter

import (
	"context"
	"net/http"
)

// CompletionRequest represents a request to the OpenRouter completions API.
// Based on https://openrouter.ai/docs/api-reference/completion as of 2025-05-19.
type CompletionRequest struct {
	BaseRequest

	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// NewCompletionReq creates a new HTTP request for the OpenRouter completions API.
// It takes a context and a ChatCompletionRequest and returns an HTTP request ready to be sent.
func NewCompletionReq(ctx context.Context, request ChatCompletionRequest) (*http.Request, error) {
	return NewRequest(ctx, "POST", "completions", request)
}

// ChatCompletionRequest represents a request to the OpenRouter chat completions API.
// Based on https://openrouter.ai/docs/api-reference/chat-completion as of 2025-05-19.
type ChatCompletionRequest struct {
	BaseRequest
	ToolsEnabled          // the api documentation doesnt mention it but the tools does
	ResponseFormatEnabled // the api documentation doesnt mention it but the structured responses doc does

	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

// NewChatCompletionReq creates a new HTTP request for the OpenRouter chat completions API.
// It takes a context and a ChatCompletionRequest and returns an HTTP request ready to be sent.
func NewChatCompletionReq(ctx context.Context, request ChatCompletionRequest) (*http.Request, error) {
	return NewRequest(ctx, "POST", "chat/completions", request)
}

// BaseRequest contains common fields used in both CompletionRequest and ChatCompletionRequest.
// It includes parameters for controlling the model's behavior and output.
type BaseRequest struct {
	Models            []string           `json:"models,omitempty"`
	Provider          *ProviderRequest   `json:"provider,omitempty"`
	Reasoning         *ReasoningRequest  `json:"reasoning,omitempty"`
	Usage             *UsageRequest      `json:"usage,omitempty"`
	Transforms        []string           `json:"transforms,omitempty"`
	Stream            bool               `json:"stream,omitempty"`
	MaxTokens         *int               `json:"max_tokens,omitempty"`
	Temperature       *float64           `json:"temperature,omitempty"`
	Seed              *int64             `json:"seed,omitempty"`
	TopP              *float64           `json:"top_p,omitempty"`
	TopK              *int64             `json:"top_k,omitempty"`
	FrequencyPenalty  *float64           `json:"frequency_penalty,omitempty"`
	PresencePenalty   *float64           `json:"presence_penalty,omitempty"`
	RepititionPenalty *float64           `json:"repetition_penalty,omitempty"`
	LogitBias         map[string]float64 `json:"logit_bias,omitempty"`
	TopLogprobs       *int               `json:"top_logprobs,omitempty"`
	MinP              *float64           `json:"min_p,omitempty"`
	TopA              *float64           `json:"top_a,omitempty"`
}

// ChatMessage represents a message in a chat conversation with the AI model.
// It includes the role (e.g., "system", "user", "assistant") and the content of the message.
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`

	/*
		have seen fields refusal and reasoning but dont know the type
		as we've only received nulls
	*/
	FinishReason       string       `json:"finish_reason,omitempty"`
	NativeFinishReason string       `json:"native_finish_reason,omitempty"`
	Index              int          `json:"index,omitempty"`
	Message            *ChatMessage `json:"message,omitempty"`
	ToolCalls          []ToolCall   `json:"tool_calls,omitempty"`
}

// ToolCall represents a call to a tool by the AI model.
// It includes information about the tool being called and its parameters.
type ToolCall struct {
	Index    int       `json:"index,omitempty"`
	ID       string    `json:"id,omitempty"`
	Type     string    `json:"type,omitempty"`
	Function *Function `json:"function,omitempty"`
}

// ProviderRequest contains configuration options for selecting and routing between different AI providers.
// It allows for specifying provider preferences, fallback behavior, and pricing constraints.
type ProviderRequest struct {
	Sort string `json:"sort,omitempty"`

	// the api documentation doesn't mention the below 2025-05-21 but the Provider Routing page does
	Order             []string           `json:"order"`
	AllowFallback     bool               `json:"allow_fallback,omitempty"`
	RequireParameters bool               `json:"require_parameters,omitempty"`
	DataCollection    DataCollectionEnum `json:"data_collection,omitempty"`
	Only              []string           `json:"only"`
	Ignore            []string           `json:"ignore"`
	Quantizations     []string           `json:"quantizations"`
	MaxPrice          *MaxPrice          `json:"max_price,omitempty"`
}

// DataCollectionEnum represents the data collection policy for the request.
// It controls whether the provider is allowed to collect and use the data from the request.
type DataCollectionEnum string

const (
	// AllowDataCollection permits the provider to collect and use data from the request.
	AllowDataCollection DataCollectionEnum = "allow"
	// DenyDataCollection prevents the provider from collecting and using data from the request.
	DenyDataCollection DataCollectionEnum = "deny"
)

// MaxPrice defines price limits for different aspects of the API request.
// It allows setting maximum prices for the overall request, completion tokens, etc.
type MaxPrice struct {
	Price      *float64 `json:"price,omitempty"`
	Completion *float64 `json:"completion,omitempty"`
	Request    *float64 `json:"request,omitempty"`
	Image      *float64 `json:"image,omitempty"`
}

// ReasoningRequest configures the reasoning capabilities of the AI model.
// It allows controlling the effort level, token allocation, and whether to include reasoning in the response.
type ReasoningRequest struct {
	Effort    EffortEnum `json:"effort,omitempty"`
	MaxTokens *int       `json:"max_tokens,omitempty"`
	Exclude   bool       `json:"exclude,omitempty"`
}

// EffortEnum represents the level of reasoning effort the model should apply.
// It can be set to low, medium, or high depending on the complexity of the task.
type EffortEnum string

const (
	// EffortEnumLow indicates minimal reasoning effort, suitable for simple tasks.
	EffortEnumLow EffortEnum = "low"
	// EffortEnumMedium indicates moderate reasoning effort, suitable for average complexity tasks.
	EffortEnumMedium EffortEnum = "medium"
	// EffortEnumHigh indicates maximum reasoning effort, suitable for complex tasks.
	EffortEnumHigh EffortEnum = "high"
)

// UsageRequest configures whether to include token usage information in the response.
// When enabled, the response will include counts of prompt, completion, and total tokens.
type UsageRequest struct {
	Include bool `json:"include,omitempty"`
}

// CompletionResponse represents the response from the OpenRouter completions API.
// It contains the generated text and related metadata.
type CompletionResponse struct {
	ID      string            `json:"id,omitempty"`
	Choices []ChoicesResponse `json:"choices,omitempty"`
}

// ChoicesResponse represents a single choice in the completion response.
// This has been modified to include what we are receiving from the API, not just what is documented.
type ChoicesResponse struct {
	Text               string       `json:"text,omitempty"`
	Index              *int         `json:"index,omitempty"`
	FinishReason       string       `json:"finish_reason,omitempty"`
	NativeFinishReason string       `json:"native_finish_reason,omitempty"`
	Message            *ChatMessage `json:"message,omitempty"`
}

// ChatCompletionResponse represents the response from the OpenRouter chat completions API.
// It contains the generated messages and related metadata such as provider, model, and token usage.
type ChatCompletionResponse struct {
	ID       string                `json:"id,omitempty"`
	Provider string                `json:"provider,omitempty"`
	Model    string                `json:"model,omitempty"`
	Object   string                `json:"object,omitempty"`
	Created  *int                  `json:"created,omitempty"`
	Choices  []ChatChoicesResponse `json:"choices,omitempty"`
	Usage    *UsageResponse        `json:"usage,omitempty"`
}

// ChatChoicesResponse represents a single choice in the chat completion response.
// It contains the message generated by the AI model.
type ChatChoicesResponse struct {
	Message *ChatMessage `json:"message,omitempty"`
}

// UsageResponse contains information about token usage in the API request and response.
// It includes counts for prompt tokens, completion tokens, and the total number of tokens used.
type UsageResponse struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
