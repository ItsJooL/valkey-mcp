package expire_key

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
)

func TestExpireKeyTool_SchemaWithNumericConstraints(t *testing.T) {
	tool := NewTool(&client.MockValkeyClient{})

	schema := tool.InputSchema()
	assert.NotNil(t, schema)

	// Convert to JSON and back to verify MCP compatibility
	schemaJSON, err := json.MarshalIndent(schema, "", "  ")
	require.NoError(t, err, "schema should be JSON serializable")

	t.Logf("Schema sent to MCP client:\n%s", string(schemaJSON))

	// Unmarshal to verify structure
	var schemaObj map[string]interface{}
	err = json.Unmarshal(schemaJSON, &schemaObj)
	require.NoError(t, err, "schema should be valid JSON")

	// Verify schema structure
	assert.Equal(t, "object", schemaObj["type"])

	properties := schemaObj["properties"].(map[string]interface{})

	// Verify seconds field has numeric minimum
	secondsProp := properties["seconds"].(map[string]interface{})
	assert.Equal(t, "integer", secondsProp["type"])
	assert.Equal(t, "Seconds until expiration", secondsProp["description"])

	// The key assertion: minimum should be numeric (int64), not string
	minimum := secondsProp["minimum"]
	require.NotNil(t, minimum, "minimum constraint should exist")

	// Check it's numeric (will be float64 after JSON unmarshal, but original is int64)
	switch v := minimum.(type) {
	case float64:
		assert.Equal(t, float64(1), v, "minimum should be numeric value 1")
	case int64:
		assert.Equal(t, int64(1), v, "minimum should be numeric value 1")
	default:
		t.Fatalf("minimum should be numeric, got %T: %v", v, v)
	}
}
