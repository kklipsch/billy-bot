package openrouter

// ToolsEnabled provides configuration for enabling AI model tool usage capabilities.
// While not documented in the main API reference as of 2025-05-19, this functionality
// is described in https://openrouter.ai/docs/features/tool-calling.
type ToolsEnabled struct {
	Tools      []Tool `json:"tools,omitempty"`
	ToolChoice string `json:"tool_choice,omitempty"`
}

// Tool defines a tool that can be called by the AI model.
// It specifies the type of tool and the function it provides.
type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

// Function describes a callable function that the AI model can use.
// It includes the function name, description, arguments, and parameter specifications.
type Function struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Arguments   string     `json:"arguments"`
	Parameters  Parameters `json:"parameters"`
}

// Parameters defines the structure of function parameters for a tool.
// It specifies the parameter type, properties, and which properties are required.
type Parameters struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required,omitempty"`
}

// Property defines a single property within a function parameter.
// It includes the property type, description, and optional enumeration of allowed values.
type Property struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
}
