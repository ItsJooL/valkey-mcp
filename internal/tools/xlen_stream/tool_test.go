package xlen_stream

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestXLenStream_Execute_Success(t *testing.T) {
	mockClient := &client.MockValkeyClient{}
	mockClient.GetStreamLengthFunc = func(ctx context.Context, key string) (int64, error) {
		return int64(42), nil
	}

	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{"key": "mystream"}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	output, ok := result.(Output)
	require.True(t, ok)
	assert.Equal(t, int64(42), output.Count)
}

func TestXLenStream_Execute_EmptyStream(t *testing.T) {
	mockClient := &client.MockValkeyClient{}
	mockClient.GetStreamLengthFunc = func(ctx context.Context, key string) (int64, error) {
		return int64(0), nil
	}

	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{"key": "emptystream"}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	output, ok := result.(Output)
	require.True(t, ok)
	assert.Equal(t, int64(0), output.Count)
}

func TestXLenStream_Execute_NonexistentStream(t *testing.T) {
	mockClient := &client.MockValkeyClient{}
	mockClient.GetStreamLengthFunc = func(ctx context.Context, key string) (int64, error) {
		return int64(0), nil
	}

	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{"key": "nonexistent"}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	output, ok := result.(Output)
	require.True(t, ok)
	assert.Equal(t, int64(0), output.Count)
}

func TestXLenStream_Execute_EmptyKey(t *testing.T) {
	mockClient := &client.MockValkeyClient{}
	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{"key": ""}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "key cannot be empty")
}

func TestXLenStream_Execute_LargeStream(t *testing.T) {
	mockClient := &client.MockValkeyClient{}
	mockClient.GetStreamLengthFunc = func(ctx context.Context, key string) (int64, error) {
		return int64(1000000), nil
	}

	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{"key": "largestream"}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	output, ok := result.(Output)
	require.True(t, ok)
	assert.Equal(t, int64(1000000), output.Count)
}
