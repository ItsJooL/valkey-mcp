package hvals_hash

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHvalsHash_Execute_Success(t *testing.T) {
	mockClient := client.NewMockClient()
	mockClient.SetMap(context.Background(), "test_hash", map[string]string{"field1": "value1", "field2": "value2", "field3": "value3"})
	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{"key": "test_hash"}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	output, ok := result.(Output)
	require.True(t, ok)
	assert.Equal(t, "test_hash", output.Key)
	assert.Len(t, output.Result, 3)
	assert.Contains(t, output.Result, "value1")
	assert.Contains(t, output.Result, "value2")
	assert.Contains(t, output.Result, "value3")
}

func TestHvalsHash_Execute_KeyNotFound(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{"key": "nonexistent_hash"}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	output, ok := result.(Output)
	require.True(t, ok)
	assert.Equal(t, "nonexistent_hash", output.Key)
	assert.Empty(t, output.Result)
}

func TestHvalsHash_Execute_EmptyKey(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{"key": ""}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.Error(t, err)
	assert.Nil(t, result)
}
