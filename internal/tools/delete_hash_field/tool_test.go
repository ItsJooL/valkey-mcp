package delete_hash_field

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

	mockClient.SetMap(ctx, "test_hash", map[string]string{
		"field1": "value1",
		"field2": "value2",
	})

	input := map[string]interface{}{
		"key":    "test_hash",
		"fields": []string{"field1"},
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
		"key":    "",
		"fields": []string{"field1"},
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestTool_Metadata(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)

	assert.Equal(t, "delete_hash_field", tool.Name())
	assert.NotEmpty(t, tool.Description())
}
