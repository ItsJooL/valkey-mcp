package sdiff_sets

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSdiffSets_Execute_Success(t *testing.T) {
	mockClient := client.NewMockClient()
	mockClient.AddSet(context.Background(), "set1", []string{"a", "b", "c", "d"})
	mockClient.AddSet(context.Background(), "set2", []string{"b", "c"})
	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{
		"keys": []string{"set1", "set2"},
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	output, ok := result.(Output)
	require.True(t, ok)
	assert.Greater(t, output.Count, int64(0))
	assert.NotNil(t, output.Members)
	assert.Len(t, output.Members, int(output.Count))
}

func TestSdiffSets_Execute_InsufficientKeys(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{
		"keys": []string{"set1"},
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "at least 2 keys required")
}

func TestSdiffSets_Execute_EmptyKeys(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{
		"keys": []string{},
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "at least 2 keys required")
}

func TestSdiffSets_Execute_MultipleOtherKeys(t *testing.T) {
	mockClient := client.NewMockClient()
	mockClient.AddSet(context.Background(), "set1", []string{"a", "b", "c", "d", "e"})
	mockClient.AddSet(context.Background(), "set2", []string{"b", "c"})
	mockClient.AddSet(context.Background(), "set3", []string{"d"})
	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{
		"keys": []string{"set1", "set2", "set3"},
	}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	output, ok := result.(Output)
	require.True(t, ok)
	assert.Greater(t, output.Count, int64(0))
	assert.NotNil(t, output.Members)
}
