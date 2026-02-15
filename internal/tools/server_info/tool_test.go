package server_info

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerInfoTool_Execute_Success(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	require.NotNil(t, result)

	output, ok := result.(Output)
	require.True(t, ok, "result should be Output type")
	assert.NotNil(t, output.Info)
	assert.NotEmpty(t, output.Info)
}

func TestServerInfoTool_Execute_EmptyInput(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	result, err := tool.Execute(ctx, json.RawMessage(`{}`))

	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestServerInfoTool_Execute_WithError(t *testing.T) {
	mockClient := client.NewMockClient()
	mockClient.GetServerInfoError = assert.AnError
	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestServerInfoTool_Metadata(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)

	assert.Equal(t, "server_info", tool.Name())
	assert.Contains(t, tool.Description(), "server")
	// InputSchema can be nil for tools with no required input
}
