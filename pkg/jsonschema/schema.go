package jsonschema

// Schema represents a JSON Schema structure with type safety
type Schema struct {
	Name                 string             `json:"name,omitempty"`
	Type                 string             `json:"type,omitempty"`
	Strict               bool               `json:"strict,omitempty"`
	Properties           map[string]*Schema `json:"properties,omitempty"`
	Required             []string           `json:"required,omitempty"`
	Items                *Schema            `json:"items,omitempty"`
	AdditionalProperties interface{}        `json:"additionalProperties,omitempty"`
	Enum                 []interface{}      `json:"enum,omitempty"`
	AllOf                []*Schema          `json:"allOf,omitempty"`
	AnyOf                []*Schema          `json:"anyOf,omitempty"`
	OneOf                []*Schema          `json:"oneOf,omitempty"`
	Not                  *Schema            `json:"not,omitempty"`
	Definitions          map[string]*Schema `json:"definitions,omitempty"`
	Title                string             `json:"title,omitempty"`
	Description          string             `json:"description,omitempty"`
	Default              interface{}        `json:"default,omitempty"`
	Format               string             `json:"format,omitempty"`
	Ref                  string             `json:"$ref,omitempty"`

	// Validation keywords for numeric types
	MultipleOf       *float64 `json:"multipleOf,omitempty"`
	Maximum          *float64 `json:"maximum,omitempty"`
	ExclusiveMaximum *bool    `json:"exclusiveMaximum,omitempty"`
	Minimum          *float64 `json:"minimum,omitempty"`
	ExclusiveMinimum *bool    `json:"exclusiveMinimum,omitempty"`

	// Validation keywords for strings
	MaxLength *int   `json:"maxLength,omitempty"`
	MinLength *int   `json:"minLength,omitempty"`
	Pattern   string `json:"pattern,omitempty"`

	// Validation keywords for arrays
	MaxItems    *int  `json:"maxItems,omitempty"`
	MinItems    *int  `json:"minItems,omitempty"`
	UniqueItems *bool `json:"uniqueItems,omitempty"`

	// Validation keywords for objects
	MaxProperties     *int                   `json:"maxProperties,omitempty"`
	MinProperties     *int                   `json:"minProperties,omitempty"`
	PatternProperties map[string]*Schema     `json:"patternProperties,omitempty"`
	Dependencies      map[string]interface{} `json:"dependencies,omitempty"`
}
