package openrouter

// the api documentation at https://openrouter.ai/docs/api-reference/completion per 2025-05-19
// does not include tool support but https://openrouter.ai/docs/features/tool-calling
// implies the following
type ToolsEnabled struct {
	Tools      []Tool `json:"tools,omitempty"`
	ToolChoice string `json:"tool_choice,omitempty"`
}

type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

type Function struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Parameters  Parameters `json:"parameters"`
}

type Parameters struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required,omitempty"`
}

type Property struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
}
