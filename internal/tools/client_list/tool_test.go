package client_list

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientListTool_Execute_Success(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	require.NotNil(t, result)

	output, ok := result.(Output)
	require.True(t, ok)
	assert.GreaterOrEqual(t, output.ClientCount, 0)
	assert.NotEmpty(t, output.RawInfo)
}

func TestClientListTool_Execute_EmptyInput(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	result, err := tool.Execute(ctx, json.RawMessage(`{}`))

	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestClientListTool_Metadata(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)

	assert.Equal(t, "client_list", tool.Name())
	assert.Contains(t, tool.Description(), "client")
	assert.NotNil(t, tool.InputSchema())
}
