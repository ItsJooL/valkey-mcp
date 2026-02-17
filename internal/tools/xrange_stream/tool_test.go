package xrange_stream

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestXRangeStream_Execute_Success(t *testing.T) {
	mockClient := &client.MockValkeyClient{}
	mockClient.GetStreamRangeFunc = func(ctx context.Context, key, start, end string, count int64) ([]client.StreamEntry, error) {
		return []client.StreamEntry{
			{ID: "1-0", FieldValues: map[string][]byte{"name": []byte("alice"), "age": []byte("30")}},
			{ID: "2-0", FieldValues: map[string][]byte{"name": []byte("bob"), "age": []byte("25")}},
		}, nil
	}

	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{
		"key":   "mystream",
		"start": "-",
		"end":   "+",
		"count": int64(10),
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	output, ok := result.(Output)
	require.True(t, ok)
	assert.Equal(t, int64(2), output.Count)
	assert.Len(t, output.Entries, 2)
	assert.Equal(t, "alice", output.Entries[0]["name"])
}

func TestXRangeStream_Execute_WithIDRange(t *testing.T) {
	mockClient := &client.MockValkeyClient{}
	mockClient.GetStreamRangeFunc = func(ctx context.Context, key, start, end string, count int64) ([]client.StreamEntry, error) {
		assert.Equal(t, "1609459200000-0", start)
		assert.Equal(t, "1609459200000-100", end)
		return []client.StreamEntry{
			{ID: "1609459200000-0", FieldValues: map[string][]byte{"value": []byte("test")}},
		}, nil
	}

	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{
		"key":   "mystream",
		"start": "1609459200000-0",
		"end":   "1609459200000-100",
		"count": int64(5),
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	output, ok := result.(Output)
	require.True(t, ok)
	assert.Equal(t, int64(1), output.Count)
}

func TestXRangeStream_Execute_EmptyKey(t *testing.T) {
	mockClient := &client.MockValkeyClient{}
	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{
		"key":   "",
		"start": "-",
		"end":   "+",
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "key cannot be empty")
}

func TestXRangeStream_Execute_EmptyStart(t *testing.T) {
	mockClient := &client.MockValkeyClient{}
	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{
		"key":   "mystream",
		"start": "",
		"end":   "+",
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "start cannot be empty")
}

func TestXRangeStream_Execute_EmptyEnd(t *testing.T) {
	mockClient := &client.MockValkeyClient{}
	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{
		"key":   "mystream",
		"start": "-",
		"end":   "",
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "end cannot be empty")
}

func TestXRangeStream_Execute_NoEntries(t *testing.T) {
	mockClient := &client.MockValkeyClient{}
	mockClient.GetStreamRangeFunc = func(ctx context.Context, key, start, end string, count int64) ([]client.StreamEntry, error) {
		return []client.StreamEntry{}, nil
	}

	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{
		"key":   "emptystream",
		"start": "-",
		"end":   "+",
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	output, ok := result.(Output)
	require.True(t, ok)
	assert.Equal(t, int64(0), output.Count)
	assert.Len(t, output.Entries, 0)
}
