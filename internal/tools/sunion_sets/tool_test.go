package sunion_sets

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSunionSets_Execute_Success(t *testing.T) {
	mockClient := client.NewMockClient()
	mockClient.AddSet(context.Background(), "set1", []string{"a", "b", "c"})
	mockClient.AddSet(context.Background(), "set2", []string{"c", "d", "e"})
	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{"keys": []string{"set1", "set2"}}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	output, ok := result.(Output)
	require.True(t, ok)
	assert.Greater(t, output.Count, int64(0))
	assert.NotNil(t, output.Members)
	assert.Len(t, output.Members, int(output.Count))
}

func TestSunionSets_Execute_EmptyKeys(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{"keys": []string{}}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "keys cannot be empty")
}

func TestSunionSets_Execute_NonexistentSets(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{"keys": []string{"nonexistent1", "nonexistent2"}}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	output, ok := result.(Output)
	require.True(t, ok)
	assert.Equal(t, int64(0), output.Count)
	assert.Empty(t, output.Members)
}

func TestSunionSets_Execute_MultipleSetUnion(t *testing.T) {
	mockClient := client.NewMockClient()
	mockClient.AddSet(context.Background(), "set1", []string{"a", "b"})
	mockClient.AddSet(context.Background(), "set2", []string{"c", "d"})
	mockClient.AddSet(context.Background(), "set3", []string{"e", "f"})
	tool := NewTool(mockClient)
	ctx := context.Background()

	input := map[string]interface{}{"keys": []string{"set1", "set2", "set3"}}
	inputJSON, _ := json.Marshal(input)
	result, err := tool.Execute(ctx, inputJSON)

	require.NoError(t, err)
	output, ok := result.(Output)
	require.True(t, ok)
	assert.Equal(t, int64(6), output.Count)
	assert.Len(t, output.Members, 6)
}
