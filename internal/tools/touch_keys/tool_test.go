package touch_keys

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTouchKeys_Execute_Success(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	// Setup test keys
	mockClient.SetString(ctx, "key1", "value1", nil, false, false)
	mockClient.SetString(ctx, "key2", "value2", nil, false, false)
	mockClient.SetString(ctx, "key3", "value3", nil, false, false)

	input := map[string]interface{}{"keys": []string{"key1", "key2", "key3"}}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	output, ok := result.(Output)
	require.True(t, ok)
	assert.Equal(t, int64(3), output.Count)
	assert.GreaterOrEqual(t, output.Updated, int64(0))
	assert.Equal(t, 3, len(output.Keys))
}

func TestTouchKeys_Execute_PartialKeyNotFound(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	// Setup only one key
	mockClient.SetString(ctx, "existing_key", "value", nil, false, false)

	input := map[string]interface{}{"keys": []string{"existing_key", "nonexistent_key"}}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	output, ok := result.(Output)
	require.True(t, ok)
	assert.Equal(t, int64(2), output.Count)
	assert.GreaterOrEqual(t, output.Updated, int64(0))
}

func TestTouchKeys_Execute_EmptyKeysList(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{"keys": []string{}}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "keys list cannot be empty")
}
