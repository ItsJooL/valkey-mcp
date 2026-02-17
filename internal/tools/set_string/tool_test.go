package set_string

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetStringTool_Execute_Success(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{
		"key":   "test_key",
		"value": "test_value",
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	require.NotNil(t, result)

	output, ok := result.(Output)
	require.True(t, ok)
	assert.True(t, output.Success)
	assert.Equal(t, "test_key", output.Key)

	// Verify the value was actually set
	val, exists, _ := mockClient.GetString(ctx, "test_key")
	assert.True(t, exists)
	assert.Equal(t, "test_value", string(val))
}

func TestSetStringTool_Execute_WithTTL(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	ttl := int64(3600)
	input := map[string]interface{}{
		"key":         "test_key",
		"value":       "test_value",
		"ttl_seconds": ttl,
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	require.NotNil(t, result)

	output, ok := result.(Output)
	require.True(t, ok)
	assert.True(t, output.Success)

	// Verify TTL was set
	actualTTL, _ := mockClient.GetTTL(ctx, "test_key")
	assert.Equal(t, ttl, actualTTL)
}

func TestSetStringTool_Execute_NX_Success(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{
		"key":   "new_key",
		"value": "test_value",
		"nx":    true,
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	output, ok := result.(Output)
	require.True(t, ok)
	assert.True(t, output.Success)
}

func TestSetStringTool_Execute_NX_Fail(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	// Set initial value
	mockClient.SetString(ctx, "existing_key", "old_value", nil, false, false)

	// Try to set with NX
	input := map[string]interface{}{
		"key":   "existing_key",
		"value": "new_value",
		"nx":    true,
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	output, ok := result.(Output)
	require.True(t, ok)
	assert.False(t, output.Success)

	// Verify old value unchanged
	val, _, _ := mockClient.GetString(ctx, "existing_key")
	assert.Equal(t, "old_value", string(val))
}

func TestSetStringTool_Execute_EmptyKey(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{
		"key":   "",
		"value": "test_value",
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "key cannot be empty")
}

func TestSetStringTool_Metadata(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)

	assert.Equal(t, "set_string", tool.Name())
	assert.Contains(t, tool.Description(), "string")
	assert.NotNil(t, tool.InputSchema())
}
