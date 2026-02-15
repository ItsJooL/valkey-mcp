package dbsize

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
)

func TestDbsizeTool_Execute_Success(t *testing.T) {
	mockClient := client.NewMockClient()

	tool := NewTool(mockClient)
	input := `{}`

	result, err := tool.Execute(context.Background(), json.RawMessage(input))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	output, ok := result.(Output)
	if !ok {
		t.Fatalf("Expected Output type, got %T", result)
	}

	if output.Size < 0 {
		t.Errorf("Expected non-negative size, got %d", output.Size)
	}
}

func TestDbsizeTool_Execute_EmptyInput(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	input := `{}`

	result, err := tool.Execute(context.Background(), json.RawMessage(input))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	output, ok := result.(Output)
	if !ok {
		t.Fatalf("Expected Output type, got %T", result)
	}

	if output.Size < 0 {
		t.Errorf("Expected non-negative size, got %d", output.Size)
	}
}

func TestDbsizeTool_Execute_InvalidJSON(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	input := `{invalid`

	_, err := tool.Execute(context.Background(), json.RawMessage(input))
	if err == nil {
		t.Fatal("Expected error for invalid JSON, got nil")
	}
}

func TestDbsizeTool_Metadata(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)

	if tool.Name() != "dbsize" {
		t.Errorf("Expected name 'dbsize', got %q", tool.Name())
	}

	if tool.Description() == "" {
		t.Error("Expected non-empty description")
	}

	schema := tool.InputSchema()
	if schema == nil {
		t.Error("Expected non-nil input schema")
	}
}
