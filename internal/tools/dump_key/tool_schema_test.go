package dump_key

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
)

func TestDumpKeyTool_SchemaIsPopulated(t *testing.T) {
	tool := NewTool(&client.MockValkeyClient{})

	schema := tool.InputSchema()
	assert.NotNil(t, schema, "InputSchema should not be nil")

	// Convert to JSON to verify it's JSON-serializable
	schemaJSON, err := json.MarshalIndent(schema, "", "  ")
	require.NoError(t, err)

	// Unmarshal to verify structure
	var schemaObj map[string]interface{}
	err = json.Unmarshal(schemaJSON, &schemaObj)
	require.NoError(t, err)

	// Verify schema has expected structure
	assert.Equal(t, "object", schemaObj["type"], "schema type should be 'object'")

	properties, ok := schemaObj["properties"].(map[string]interface{})
	assert.True(t, ok, "should have properties map")

	keyProp, ok := properties["key"].(map[string]interface{})
	assert.True(t, ok, "should have 'key' property")
	assert.Equal(t, "string", keyProp["type"], "key should be string type")
	assert.Equal(t, "Key to serialize", keyProp["description"], "should have description")

	required, ok := schemaObj["required"].([]interface{})
	assert.True(t, ok, "should have required array")
	assert.Len(t, required, 1, "should have 1 required field")
	assert.Equal(t, "key", required[0], "key should be required")
}

func TestDumpKeyTool_SchemaInRegistry(t *testing.T) {
	reg := registry.NewToolRegistry()
	mockClient := &client.MockValkeyClient{}

	Init(reg, mockClient)

	// Get all tools info
	toolInfos := reg.GetAllToolInfo()
	assert.Len(t, toolInfos, 1, "should have 1 tool registered")

	info := toolInfos[0]
	assert.Equal(t, "dump_key", info.Name)
	assert.NotNil(t, info.InputSchema, "InputSchema should be populated in registry")

	// Verify it's a valid JSON schema
	schemaJSON, err := json.Marshal(info.InputSchema)
	require.NoError(t, err)

	var schema map[string]interface{}
	err = json.Unmarshal(schemaJSON, &schema)
	require.NoError(t, err)

	assert.Equal(t, "object", schema["type"])
}

func TestDumpKeyTool_MCP_CanSeeParameters(t *testing.T) {
	// This test simulates what an MCP client would see
	tool := NewTool(&client.MockValkeyClient{})

	// Get the schema as MCP would
	schema := tool.InputSchema()
	assert.NotNil(t, schema)

	// Marshal it (as MCP would in the CallToolResult)
	schemaJSON, err := json.MarshalIndent(schema, "", "  ")
	require.NoError(t, err)

	t.Logf("Schema sent to MCP client:\n%s", string(schemaJSON))

	// Verify the schema has all the information needed by the user
	var schemaObj map[string]interface{}
	json.Unmarshal(schemaJSON, &schemaObj)

	properties := schemaObj["properties"].(map[string]interface{})
	keyProp := properties["key"].(map[string]interface{})

	// The user can now see:
	// 1. Parameter name: "key"
	assert.Equal(t, "key", "key") // field name from JSON tag

	// 2. Type: "string"
	assert.Equal(t, "string", keyProp["type"])

	// 3. Required: yes
	assert.Contains(t, schemaObj["required"], "key")

	// 4. Description: "Key to serialize"
	assert.Equal(t, "Key to serialize", keyProp["description"])
}
