package keys_by_pattern

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeysByPattern_Execute_Success(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	// Setup test keys
	mockClient.SetString(ctx, "test:key1", "value1", nil, false, false)
	mockClient.SetString(ctx, "test:key2", "value2", nil, false, false)

	input := map[string]interface{}{"pattern": "*"}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	output, ok := result.(Output)
	require.True(t, ok)
	assert.NotNil(t, output.Keys)
	assert.GreaterOrEqual(t, output.Count, 0)
}

func TestKeysByPattern_Execute_NoKeysFound(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{"pattern": "nonexistent:*"}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	output, ok := result.(Output)
	require.True(t, ok)
	assert.Equal(t, 0, output.Count)
	assert.Equal(t, 0, len(output.Keys))
}

func TestKeysByPattern_Execute_EmptyPattern(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{"pattern": ""}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "pattern cannot be empty")
}
