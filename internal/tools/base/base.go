package base

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// BaseTool provides common functionality for all tools.
type BaseTool struct {
	name        string
	description string
	inputType   interface{}
	schema      interface{}
}

// NewBaseTool creates a new base tool.
func NewBaseTool(name, description string, inputType interface{}) BaseTool {
	bt := BaseTool{
		name:        name,
		description: description,
		inputType:   inputType,
	}
	// Generate JSON schema from input type if provided
	if inputType != nil {
		bt.schema = generateJSONSchema(inputType)
	}
	return bt
}

// Name returns the tool name.
func (b BaseTool) Name() string {
	return b.name
}

// Description returns the tool description.
func (b BaseTool) Description() string {
	return b.description
}

// InputSchema returns the input schema.
func (b BaseTool) InputSchema() interface{} {
	return b.schema
}

// generateJSONSchema creates a JSON schema from a struct type.
func generateJSONSchema(inputType interface{}) map[string]interface{} {
	schema := map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
	}

	t := reflect.TypeOf(inputType)
	if t == nil {
		return schema
	}

	// Handle pointer types
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Only process struct types
	if t.Kind() != reflect.Struct {
		return schema
	}

	properties := schema["properties"].(map[string]interface{})
	required := []string{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		// Extract JSON field name (before comma)
		fieldName := strings.Split(jsonTag, ",")[0]
		if fieldName == "" {
			fieldName = field.Name
		}

		// Parse jsonschema tag for metadata
		schemaTag := field.Tag.Get("jsonschema")
		fieldSchema := parseSchemaTag(field, schemaTag)

		properties[fieldName] = fieldSchema

		// Check if field is required
		if strings.Contains(schemaTag, "required") {
			required = append(required, fieldName)
		}
	}

	if len(required) > 0 {
		schema["required"] = required
	}

	return schema
}

// parseSchemaTag parses jsonschema struct tag and returns field schema.
func parseSchemaTag(field reflect.StructField, tag string) map[string]interface{} {
	fieldSchema := map[string]interface{}{}

	// Determine type
	fieldType := getJSONSchemaType(field.Type)
	if fieldType != "" {
		fieldSchema["type"] = fieldType
	}

	// Parse comma-separated tag values
	parts := strings.Split(tag, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)

		// Skip the 'required' keyword - it's handled separately
		if part == "required" {
			continue
		}

		// Parse description
		if strings.HasPrefix(part, "description=") {
			description := strings.TrimPrefix(part, "description=")
			fieldSchema["description"] = description
		}

		// Parse enum values
		if strings.HasPrefix(part, "enum=") {
			enumStr := strings.TrimPrefix(part, "enum=")
			// Simple enum parsing: "enum=value1,value2" or "enum=[value1,value2]"
			values := strings.Split(strings.Trim(enumStr, "[]"), ",")
			for i, v := range values {
				values[i] = strings.TrimSpace(v)
			}
			fieldSchema["enum"] = values
		}

		// Parse minimum (convert to number)
		if strings.HasPrefix(part, "minimum=") {
			minStr := strings.TrimPrefix(part, "minimum=")
			if num, err := parseNumber(minStr); err == nil {
				fieldSchema["minimum"] = num
			}
		}

		// Parse maximum (convert to number)
		if strings.HasPrefix(part, "maximum=") {
			maxStr := strings.TrimPrefix(part, "maximum=")
			if num, err := parseNumber(maxStr); err == nil {
				fieldSchema["maximum"] = num
			}
		}

		// Parse minLength (convert to number)
		if strings.HasPrefix(part, "minLength=") {
			lenStr := strings.TrimPrefix(part, "minLength=")
			if num, err := parseNumber(lenStr); err == nil {
				fieldSchema["minLength"] = num
			}
		}

		// Parse maxLength (convert to number)
		if strings.HasPrefix(part, "maxLength=") {
			lenStr := strings.TrimPrefix(part, "maxLength=")
			if num, err := parseNumber(lenStr); err == nil {
				fieldSchema["maxLength"] = num
			}
		}
	}

	return fieldSchema
}

// getJSONSchemaType maps Go types to JSON Schema types.
func getJSONSchemaType(t reflect.Type) string {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "integer"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "boolean"
	case reflect.Slice, reflect.Array:
		return "array"
	case reflect.Map:
		return "object"
	default:
		return ""
	}
}

// parseNumber converts a string to either int or float64 for JSON schema constraints.
func parseNumber(s string) (interface{}, error) {
	// Try parsing as integer first
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i, nil
	}
	// Fall back to float
	return strconv.ParseFloat(s, 64)
}

// ParseInput is a helper to parse JSON input into a typed struct.
func (b BaseTool) ParseInput(input json.RawMessage, target interface{}) error {
	if len(input) == 0 {
		return nil
	}
	if err := json.Unmarshal(input, target); err != nil {
		return fmt.Errorf("invalid input format: %w", err)
	}
	return nil
}

// Execute must be implemented by concrete tools.
func (b BaseTool) Execute(ctx context.Context, input json.RawMessage) (interface{}, error) {
	return nil, fmt.Errorf("execute not implemented for tool %s", b.name)
}
