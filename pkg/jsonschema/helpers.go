package jsonschema

// NewObjectSchema creates a new Schema for an object type
func NewObjectSchema(properties map[string]*Schema, required []string) *Schema {
	return &Schema{
		Type:       "object",
		Properties: properties,
		Required:   required,
	}
}

// NewArraySchema creates a new Schema for an array type
func NewArraySchema(items *Schema) *Schema {
	return &Schema{
		Type:  "array",
		Items: items,
	}
}

// NewStringSchema creates a new Schema for a string type
func NewStringSchema() *Schema {
	return &Schema{
		Type: "string",
	}
}

// NewNumberSchema creates a new Schema for a number type
func NewNumberSchema() *Schema {
	return &Schema{
		Type: "number",
	}
}

// NewIntegerSchema creates a new Schema for an integer type
func NewIntegerSchema() *Schema {
	return &Schema{
		Type: "integer",
	}
}

// NewBooleanSchema creates a new Schema for a boolean type
func NewBooleanSchema() *Schema {
	return &Schema{
		Type: "boolean",
	}
}
