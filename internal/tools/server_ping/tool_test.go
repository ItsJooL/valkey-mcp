package server_ping

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerPingTool_Execute_Success(t *testing.T) {
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
	assert.True(t, output.Alive)
	assert.GreaterOrEqual(t, output.LatencyMs, 0.0)
}

func TestServerPingTool_Execute_WithError(t *testing.T) {
	mockClient := client.NewMockClient()
	mockClient.PingError = assert.AnError
	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	// Note: This tool returns Output even on error, with Alive=false
	require.NoError(t, err)
	require.NotNil(t, result)

	output, ok := result.(Output)
	require.True(t, ok)
	assert.False(t, output.Alive)
	assert.NotEmpty(t, output.Message)
}

func TestServerPingTool_Execute_EmptyInput(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	result, err := tool.Execute(ctx, json.RawMessage(`{}`))

	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestServerPingTool_Metadata(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)

	assert.Equal(t, "server_ping", tool.Name())
	assert.NotEmpty(t, tool.Description())
	// InputSchema can be nil for tools with no required input
}
