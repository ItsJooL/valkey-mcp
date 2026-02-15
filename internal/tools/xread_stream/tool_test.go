package xread_stream

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestXReadStream_Execute_Success(t *testing.T) {
	mockClient := &client.MockValkeyClient{}
	mockClient.ReadStreamFunc = func(ctx context.Context, key, id string, count int64) ([]map[string]string, error) {
		return []map[string]string{
			{"message": "hello"},
			{"message": "world"},
		}, nil
	}

	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{
		"key":   "mystream",
		"id":    "0",
		"count": int64(10),
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	output, ok := result.(Output)
	require.True(t, ok)
	assert.Equal(t, int64(2), output.Count)
	assert.Len(t, output.Entries, 2)
	assert.Equal(t, "hello", output.Entries[0]["message"])
}

func TestXReadStream_Execute_WithSpecificID(t *testing.T) {
	mockClient := &client.MockValkeyClient{}
	mockClient.ReadStreamFunc = func(ctx context.Context, key, id string, count int64) ([]map[string]string, error) {
		assert.Equal(t, "1609459200000-0", id)
		return []map[string]string{
			{"data": "test"},
		}, nil
	}

	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{
		"key":   "mystream",
		"id":    "1609459200000-0",
		"count": int64(5),
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	output, ok := result.(Output)
	require.True(t, ok)
	assert.Equal(t, int64(1), output.Count)
}

func TestXReadStream_Execute_EmptyKey(t *testing.T) {
	mockClient := &client.MockValkeyClient{}
	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{
		"key":   "",
		"id":    "0",
		"count": int64(10),
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "key cannot be empty")
}

func TestXReadStream_Execute_EmptyID(t *testing.T) {
	mockClient := &client.MockValkeyClient{}
	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{
		"key":   "mystream",
		"id":    "",
		"count": int64(10),
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "id cannot be empty")
}

func TestXReadStream_Execute_NoEntries(t *testing.T) {
	mockClient := &client.MockValkeyClient{}
	mockClient.ReadStreamFunc = func(ctx context.Context, key, id string, count int64) ([]map[string]string, error) {
		return []map[string]string{}, nil
	}

	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{
		"key":   "mystream",
		"id":    "$",
		"count": int64(10),
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	output, ok := result.(Output)
	require.True(t, ok)
	assert.Equal(t, int64(0), output.Count)
	assert.Len(t, output.Entries, 0)
}
