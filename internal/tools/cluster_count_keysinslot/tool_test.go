package cluster_count_keysinslot

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
)

func TestClusterCountKeysinslotTool_Execute_Success(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	input := `{"slot": 100}`

	result, err := tool.Execute(context.Background(), json.RawMessage(input))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	output, ok := result.(Output)
	if !ok {
		t.Fatalf("Expected Output type, got %T", result)
	}

	if output.Count < 0 {
		t.Errorf("Expected non-negative count, got %d", output.Count)
	}
}

func TestClusterCountKeysinslotTool_Execute_InvalidSlot(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	input := `{"slot": 20000}`

	_, err := tool.Execute(context.Background(), json.RawMessage(input))
	if err == nil {
		t.Fatal("Expected error for invalid slot, got nil")
	}
}

func TestClusterCountKeysinslotTool_Metadata(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)

	if tool.Name() != "cluster_count_keysinslot" {
		t.Errorf("Expected name 'cluster_count_keysinslot', got %q", tool.Name())
	}
}
