package get_string_range

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTool_Execute_Success(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	// Setup: create initial string if needed
	mockClient.SetString(ctx, "test_key", "initial_value", nil, false, false)

	input := map[string]interface{}{
		"key": "test_key",
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestTool_Execute_EmptyKey(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{
		"key": "",
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestTool_Metadata(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)

	assert.Equal(t, "get_string_range", tool.Name())
	assert.NotEmpty(t, tool.Description())
}
