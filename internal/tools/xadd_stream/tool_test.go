package xadd_stream

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestXAddStream_Execute_Success(t *testing.T) {
	mockClient := &client.MockValkeyClient{}
	mockClient.AddStreamFunc = func(ctx context.Context, key, id string, fields map[string]string) (string, error) {
		return "1609459200000-0", nil
	}

	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{
		"key":    "mystream",
		"id":     "*",
		"fields": map[string]string{"name": "alice", "age": "30"},
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	output, ok := result.(Output)
	require.True(t, ok)
	assert.Equal(t, "1609459200000-0", output.ID)
}

func TestXAddStream_Execute_WithExplicitID(t *testing.T) {
	mockClient := &client.MockValkeyClient{}
	mockClient.AddStreamFunc = func(ctx context.Context, key, id string, fields map[string]string) (string, error) {
		return id, nil
	}

	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{
		"key":    "mystream",
		"id":     "1609459200000-5",
		"fields": map[string]string{"message": "hello"},
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	output, ok := result.(Output)
	require.True(t, ok)
	assert.Equal(t, "1609459200000-5", output.ID)
}

func TestXAddStream_Execute_EmptyKey(t *testing.T) {
	mockClient := &client.MockValkeyClient{}
	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{
		"key":    "",
		"id":     "*",
		"fields": map[string]string{"field": "value"},
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "key cannot be empty")
}

func TestXAddStream_Execute_AutoGenerateID(t *testing.T) {
	mockClient := &client.MockValkeyClient{}
	mockClient.AddStreamFunc = func(ctx context.Context, key, id string, fields map[string]string) (string, error) {
		assert.Equal(t, "*", id)
		return "1609459200000-0", nil
	}
	
	tool := NewTool(mockClient)
	ctx := context.Background()

	// Empty ID should auto-generate with "*"
	input := map[string]interface{}{
		"key":    "mystream",
		"fields": map[string]string{"field": "value"},
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	output, ok := result.(Output)
	require.True(t, ok)
	assert.Equal(t, "1609459200000-0", output.ID)
}

func TestXAddStream_Execute_EmptyFields(t *testing.T) {
	mockClient := &client.MockValkeyClient{}
	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{
		"key":    "mystream",
		"id":     "*",
		"fields": map[string]string{},
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "fields cannot be empty")
}

func TestXAddStream_Execute_MultipleFields(t *testing.T) {
	mockClient := &client.MockValkeyClient{}
	mockClient.AddStreamFunc = func(ctx context.Context, key, id string, fields map[string]string) (string, error) {
		assert.Equal(t, 3, len(fields))
		return "1609459200000-1", nil
	}

	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{
		"key": "mystream",
		"id":  "*",
		"fields": map[string]string{
			"user":   "bob",
			"action": "login",
			"ip":     "192.168.1.1",
		},
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	output, ok := result.(Output)
	require.True(t, ok)
	assert.Equal(t, "1609459200000-1", output.ID)
}
