package openrouter

import (
	"context"
	"net/http"
)

// per https://openrouter.ai/docs/api-reference/completion 2025-05-19
type CompletionRequest struct {
	BaseRequest

	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

func NewCompletionReq(ctx context.Context, request ChatCompletionRequest) (*http.Request, error) {
	return NewRequest(ctx, "POST", "completions", request)
}

// per https://openrouter.ai/docs/api-reference/chat-completion 2025-05-19
type ChatCompletionRequest struct {
	BaseRequest
	ToolsEnabled          // the api documentation doesnt mention it but the tools does
	ResponseFormatEnabled // the api documentation doesnt mention it but the structured responses doc does

	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

func NewChatCompletionReq(ctx context.Context, request ChatCompletionRequest) (*http.Request, error) {
	return NewRequest(ctx, "POST", "chat/completions", request)
}

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

type ToolCall struct {
	Index    int       `json:"index,omitempty"`
	ID       string    `json:"id,omitempty"`
	Type     string    `json:"type,omitempty"`
	Function *Function `json:"function,omitempty"`
}

type ProviderRequest struct {
	Sort string `json:"sort,omitempty"`
}

type ReasoningRequest struct {
	Effort    EffortEnum `json:"effort,omitempty"`
	MaxTokens *int       `json:"max_tokens,omitempty"`
	Exclude   bool       `json:"exclude,omitempty"`
}

type EffortEnum string

const (
	EffortEnumLow    EffortEnum = "low"
	EffortEnumMedium EffortEnum = "medium"
	EffortEnumHigh   EffortEnum = "high"
)

type UsageRequest struct {
	Include bool `json:"include,omitempty"`
}

type CompletionResponse struct {
	ID      string            `json:"id,omitempty"`
	Choices []ChoicesResponse `json:"choices,omitempty"`
}

// modified to include what we are receiving from the API not just what is documented
type ChoicesResponse struct {
	Text               string       `json:"text,omitempty"`
	Index              *int         `json:"index,omitempty"`
	FinishReason       string       `json:"finish_reason,omitempty"`
	NativeFinishReason string       `json:"native_finish_reason,omitempty"`
	Message            *ChatMessage `json:"message,omitempty"`
}

type ChatCompletionResponse struct {
	ID       string                `json:"id,omitempty"`
	Provider string                `json:"provider,omitempty"`
	Model    string                `json:"model,omitempty"`
	Object   string                `json:"object,omitempty"`
	Created  *int                  `json:"created,omitempty"`
	Choices  []ChatChoicesResponse `json:"choices,omitempty"`
	Usage    *UsageResponse        `json:"usage,omitempty"`
}

type ChatChoicesResponse struct {
	Message *ChatMessage `json:"message,omitempty"`
}

type UsageResponse struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
