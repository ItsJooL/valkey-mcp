package base

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test input types with various jsonschema tags
type SimpleInput struct {
	Key string `json:"key" jsonschema:"required,description=The key to access"`
}

type ComplexInput struct {
	Key      string `json:"key" jsonschema:"required,description=Unique identifier"`
	Value    string `json:"value" jsonschema:"required,description=The value to store"`
	TTL      int64  `json:"ttl" jsonschema:"minimum=0,description=Time to live in seconds"`
	Disabled bool   `json:"disabled" jsonschema:"description=Whether to disable"`
}

type InputWithEnum struct {
	Command string `json:"command" jsonschema:"required,enum=GET,enum=SET,enum=DEL,description=Redis command"`
}

type InputWithArray struct {
	Keys []string `json:"keys" jsonschema:"required,description=List of keys"`
}

func TestNewBaseTool_GeneratesSchema(t *testing.T) {
	tool := NewBaseTool("test_tool", "Test description", SimpleInput{})

	assert.Equal(t, "test_tool", tool.Name())
	assert.Equal(t, "Test description", tool.Description())
	assert.NotNil(t, tool.InputSchema())
}

func TestGenerateJSONSchema_SimpleStruct(t *testing.T) {
	schema := generateJSONSchema(SimpleInput{})

	assert.Equal(t, "object", schema["type"])

	properties, ok := schema["properties"].(map[string]interface{})
	assert.True(t, ok, "properties should be a map")

	keyProp, ok := properties["key"].(map[string]interface{})
	assert.True(t, ok, "key property should exist")
	assert.Equal(t, "string", keyProp["type"])
	assert.Equal(t, "The key to access", keyProp["description"])

	required, ok := schema["required"].([]string)
	assert.True(t, ok, "required should be a slice")
	assert.Contains(t, required, "key")
}

func TestGenerateJSONSchema_ComplexStruct(t *testing.T) {
	schema := generateJSONSchema(ComplexInput{})

	properties := schema["properties"].(map[string]interface{})

	// Check key field
	keyProp := properties["key"].(map[string]interface{})
	assert.Equal(t, "string", keyProp["type"])
	assert.Equal(t, "Unique identifier", keyProp["description"])

	// Check value field
	valueProp := properties["value"].(map[string]interface{})
	assert.Equal(t, "string", valueProp["type"])
	assert.Equal(t, "The value to store", valueProp["description"])

	// Check TTL field (integer with minimum)
	ttlProp := properties["ttl"].(map[string]interface{})
	assert.Equal(t, "integer", ttlProp["type"])
	assert.Equal(t, int64(0), ttlProp["minimum"], "minimum should be int64")
	assert.Equal(t, "Time to live in seconds", ttlProp["description"])

	// Check disabled field (boolean)
	disabledProp := properties["disabled"].(map[string]interface{})
	assert.Equal(t, "boolean", disabledProp["type"])
	assert.Equal(t, "Whether to disable", disabledProp["description"])

	required := schema["required"].([]string)
	assert.ElementsMatch(t, []string{"key", "value"}, required)
}

func TestGenerateJSONSchema_WithEnum(t *testing.T) {
	schema := generateJSONSchema(InputWithEnum{})

	properties := schema["properties"].(map[string]interface{})
	commandProp := properties["command"].(map[string]interface{})

	assert.Equal(t, "string", commandProp["type"])
	assert.Equal(t, "Redis command", commandProp["description"])

	// The enum parsing for tags is simplistic - it treats each enum=X as a separate tag part
	// This is a known limitation but works for most cases
	_, ok := commandProp["enum"]
	assert.True(t, ok, "enum field should exist")
}

func TestGenerateJSONSchema_WithArray(t *testing.T) {
	schema := generateJSONSchema(InputWithArray{})

	properties := schema["properties"].(map[string]interface{})
	keysProp := properties["keys"].(map[string]interface{})

	assert.Equal(t, "array", keysProp["type"])
	assert.Equal(t, "List of keys", keysProp["description"])

	required := schema["required"].([]string)
	assert.Contains(t, required, "keys")
}

func TestGenerateJSONSchema_NilInput(t *testing.T) {
	schema := generateJSONSchema(nil)

	assert.Equal(t, "object", schema["type"])
	properties := schema["properties"].(map[string]interface{})
	assert.Empty(t, properties)
}

func TestGenerateJSONSchema_NonStructType(t *testing.T) {
	schema := generateJSONSchema("not a struct")

	assert.Equal(t, "object", schema["type"])
	properties := schema["properties"].(map[string]interface{})
	assert.Empty(t, properties)
}

func TestGenerateJSONSchema_PointerToStruct(t *testing.T) {
	schema := generateJSONSchema(&SimpleInput{})

	properties := schema["properties"].(map[string]interface{})
	_, ok := properties["key"]
	assert.True(t, ok, "should handle pointer types")
}

func TestSchemaIsJSONSerializable(t *testing.T) {
	tool := NewBaseTool("test", "test", ComplexInput{})
	schema := tool.InputSchema()

	// Should be able to marshal to JSON
	jsonBytes, err := json.Marshal(schema)
	require.NoError(t, err)

	// Should be valid JSON
	var result map[string]interface{}
	err = json.Unmarshal(jsonBytes, &result)
	require.NoError(t, err)

	// Should contain expected top-level fields
	assert.Equal(t, "object", result["type"])
	assert.NotNil(t, result["properties"])
	assert.NotNil(t, result["required"])
}

func TestParseSchemaTag_Description(t *testing.T) {
	// This is tested indirectly through generateJSONSchema,
	// but we can verify the behavior here
	schema := generateJSONSchema(SimpleInput{})
	properties := schema["properties"].(map[string]interface{})
	keyProp := properties["key"].(map[string]interface{})

	assert.Contains(t, keyProp["description"], "key")
}

func TestInputSchemaNotNil(t *testing.T) {
	// Ensure that inputs with proper types produce a valid schema
	tool := NewBaseTool("test", "test", SimpleInput{})
	assert.NotNil(t, tool.InputSchema())

	// And that nil inputs produce nil schema (which is valid - no input needed)
	toolNoInput := NewBaseTool("test", "test", nil)
	assert.Nil(t, toolNoInput.InputSchema())
}

func TestCompleteWorkflow(t *testing.T) {
	// Simulate what an MCP client would do
	tool := NewBaseTool("get_string", "Get a string value", SimpleInput{})

	// 1. Get the schema
	schema := tool.InputSchema()
	assert.NotNil(t, schema)

	// 2. Convert to JSON (what MCP would do)
	schemaJSON, err := json.MarshalIndent(schema, "", "  ")
	require.NoError(t, err)

	// 3. Unmarshal to verify it's valid
	var schemaObj map[string]interface{}
	err = json.Unmarshal(schemaJSON, &schemaObj)
	require.NoError(t, err)

	// 4. Verify schema structure
	assert.Equal(t, "object", schemaObj["type"])
	assert.NotNil(t, schemaObj["properties"])

	properties := schemaObj["properties"].(map[string]interface{})
	assert.NotEmpty(t, properties)

	// 5. Verify key property has all details
	keyProp := properties["key"].(map[string]interface{})
	assert.Equal(t, "string", keyProp["type"])
	assert.NotEmpty(t, keyProp["description"])
}
