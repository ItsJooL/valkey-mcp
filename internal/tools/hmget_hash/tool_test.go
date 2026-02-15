package hmget_hash

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHmgetHash_Execute_Success(t *testing.T) {
	mockClient := client.NewMockClient()
	mockClient.SetMap(context.Background(), "test_hash", map[string]string{"field1": "value1", "field2": "value2", "field3": "value3"})
	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{"key": "test_hash", "fields": []string{"field1", "field2"}}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	output, ok := result.(Output)
	require.True(t, ok)
	assert.Equal(t, "test_hash", output.Key)
	assert.Equal(t, "value1", output.Result["field1"])
	assert.Equal(t, "value2", output.Result["field2"])
}

func TestHmgetHash_Execute_PartialFields(t *testing.T) {
	mockClient := client.NewMockClient()
	mockClient.SetMap(context.Background(), "test_hash", map[string]string{"field1": "value1", "field2": "value2", "field3": "value3"})
	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{"key": "test_hash", "fields": []string{"field1", "nonexistent"}}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	output, ok := result.(Output)
	require.True(t, ok)
	assert.Equal(t, "test_hash", output.Key)
	assert.Equal(t, "value1", output.Result["field1"])
}

func TestHmgetHash_Execute_EmptyFields(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{"key": "test_hash", "fields": []string{}}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.Error(t, err)
	assert.Nil(t, result)
}
