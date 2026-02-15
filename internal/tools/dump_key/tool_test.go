package dump_key

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDumpKey_Execute_Success(t *testing.T) {
	mockClient := &client.MockValkeyClient{}
	testData := []byte{0x00, 0x05, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x09}
	mockClient.DumpKeyFunc = func(ctx context.Context, key string) ([]byte, error) {
		return testData, nil
	}

	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{"key": "mykey"}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	output, ok := result.(Output)
	require.True(t, ok)
	assert.Equal(t, len(testData), output.Size)
	// Verify base64 encoding is correct
	decoded, _ := base64.StdEncoding.DecodeString(output.Serialized)
	assert.Equal(t, testData, decoded)
}

func TestDumpKey_Execute_StringValue(t *testing.T) {
	mockClient := &client.MockValkeyClient{}
	mockClient.DumpKeyFunc = func(ctx context.Context, key string) ([]byte, error) {
		// Simulate dumped string "hello"
		return []byte("hello"), nil
	}

	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{"key": "stringkey"}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	output, ok := result.(Output)
	require.True(t, ok)
	assert.Equal(t, 5, output.Size)
	decoded, _ := base64.StdEncoding.DecodeString(output.Serialized)
	assert.Equal(t, []byte("hello"), decoded)
}

func TestDumpKey_Execute_EmptyKey(t *testing.T) {
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

func TestDumpKey_Execute_NonexistentKey(t *testing.T) {
	mockClient := &client.MockValkeyClient{}
	mockClient.DumpKeyFunc = func(ctx context.Context, key string) ([]byte, error) {
		// Simulate nil for non-existent key
		return nil, nil
	}

	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{"key": "nonexistent"}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	output, ok := result.(Output)
	require.True(t, ok)
	assert.Equal(t, 0, output.Size)
	assert.Equal(t, "", output.Serialized)
}

func TestDumpKey_Execute_LargeBinaryData(t *testing.T) {
	mockClient := &client.MockValkeyClient{}
	largeData := make([]byte, 1024)
	for i := 0; i < len(largeData); i++ {
		largeData[i] = byte(i % 256)
	}

	mockClient.DumpKeyFunc = func(ctx context.Context, key string) ([]byte, error) {
		return largeData, nil
	}

	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{"key": "binarykey"}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	output, ok := result.(Output)
	require.True(t, ok)
	assert.Equal(t, 1024, output.Size)
	decoded, _ := base64.StdEncoding.DecodeString(output.Serialized)
	assert.Equal(t, largeData, decoded)
}
