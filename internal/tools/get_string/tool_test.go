package get_string

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetStringTool_Execute_Success(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	// Setup: Set a value in the mock
	mockClient.SetString(ctx, "test_key", "test_value", nil, false, false)

	// Execute
	input := map[string]interface{}{
		"key": "test_key",
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)

	output, ok := result.(Output)
	require.True(t, ok, "result should be Output type")
	assert.Equal(t, "test_key", output.Key)
	assert.Equal(t, "test_value", output.Value)
	assert.True(t, output.Exists)
}

func TestGetStringTool_Execute_KeyNotFound(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	// Execute with non-existent key
	input := map[string]interface{}{
		"key": "nonexistent_key",
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)

	output, ok := result.(Output)
	require.True(t, ok)
	assert.Equal(t, "nonexistent_key", output.Key)
	assert.Equal(t, "", output.Value)
	assert.False(t, output.Exists)
}

func TestGetStringTool_Execute_EmptyKey(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	// Execute with empty key
	input := map[string]interface{}{
		"key": "",
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "key cannot be empty")
}

func TestGetStringTool_Execute_InvalidInput(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	// Execute with invalid JSON
	invalidJSON := json.RawMessage(`{"invalid": 123}`)
	result, err := tool.Execute(ctx, invalidJSON)

	// Should still work but fail validation
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestGetStringTool_Metadata(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)

	assert.Equal(t, "get_string", tool.Name())
	assert.Contains(t, tool.Description(), "string")
	assert.NotNil(t, tool.InputSchema())
}
