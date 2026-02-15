package restore_key

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRestoreKey_Execute_Success(t *testing.T) {
	mockClient := &client.MockValkeyClient{}
	mockClient.RestoreKeyFunc = func(ctx context.Context, key string, ttl int64, serialized []byte) (bool, error) {
		assert.Equal(t, "mykey", key)
		assert.Equal(t, int64(1000), ttl)
		return true, nil
	}

	tool := NewTool(mockClient)
	ctx := context.Background()

	testData := []byte("hello")
	encoded := base64.StdEncoding.EncodeToString(testData)

	input := map[string]interface{}{
		"key":        "mykey",
		"ttl":        int64(1000),
		"serialized": encoded,
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	output, ok := result.(Output)
	require.True(t, ok)
	assert.True(t, output.Success)
	assert.Equal(t, "Key restored successfully", output.Message)
}

func TestRestoreKey_Execute_NoTTL(t *testing.T) {
	mockClient := &client.MockValkeyClient{}
	mockClient.RestoreKeyFunc = func(ctx context.Context, key string, ttl int64, serialized []byte) (bool, error) {
		assert.Equal(t, int64(0), ttl)
		return true, nil
	}

	tool := NewTool(mockClient)
	ctx := context.Background()

	testData := []byte{0x00, 0x05, 0x68, 0x65, 0x6c, 0x6c, 0x6f}
	encoded := base64.StdEncoding.EncodeToString(testData)

	input := map[string]interface{}{
		"key":        "restorekey",
		"ttl":        int64(0),
		"serialized": encoded,
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	output, ok := result.(Output)
	require.True(t, ok)
	assert.True(t, output.Success)
}

func TestRestoreKey_Execute_EmptyKey(t *testing.T) {
	mockClient := &client.MockValkeyClient{}
	tool := NewTool(mockClient)
	ctx := context.Background()

	encoded := base64.StdEncoding.EncodeToString([]byte("data"))
	input := map[string]interface{}{
		"key":        "",
		"ttl":        int64(0),
		"serialized": encoded,
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "key cannot be empty")
}

func TestRestoreKey_Execute_EmptySerializedData(t *testing.T) {
	mockClient := &client.MockValkeyClient{}
	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{
		"key":        "mykey",
		"ttl":        int64(0),
		"serialized": "",
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "either 'serialized' or 'serialized_value' must be provided")
}

func TestRestoreKey_Execute_InvalidBase64(t *testing.T) {
	mockClient := &client.MockValkeyClient{}
	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{
		"key":        "mykey",
		"ttl":        int64(0),
		"serialized": "not-valid-base64!!!",
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to decode base64")
}

func TestRestoreKey_Execute_LargeBinaryData(t *testing.T) {
	mockClient := &client.MockValkeyClient{}
	mockClient.RestoreKeyFunc = func(ctx context.Context, key string, ttl int64, serialized []byte) (bool, error) {
		assert.Equal(t, 1024, len(serialized))
		return true, nil
	}

	tool := NewTool(mockClient)
	ctx := context.Background()

	// Create 1KB of binary data
	largeData := make([]byte, 1024)
	for i := 0; i < len(largeData); i++ {
		largeData[i] = byte(i % 256)
	}
	encoded := base64.StdEncoding.EncodeToString(largeData)

	input := map[string]interface{}{
		"key":        "largekey",
		"ttl":        int64(5000),
		"serialized": encoded,
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	output, ok := result.(Output)
	require.True(t, ok)
	assert.True(t, output.Success)
}
