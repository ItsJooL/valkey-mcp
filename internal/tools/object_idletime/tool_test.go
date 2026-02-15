package object_idletime

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
)

func TestObjectIdletimeTool_Execute_Success(t *testing.T) {
	mockClient := client.NewMockClient()

	tool := NewTool(mockClient)
	input := `{"key": "test-key"}`

	result, err := tool.Execute(context.Background(), json.RawMessage(input))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	output, ok := result.(Output)
	if !ok {
		t.Fatalf("Expected Output type, got %T", result)
	}

	// MockClient returns 0 for idle time by default
	if output.IdleTime < 0 {
		t.Errorf("Expected non-negative idle time, got %d", output.IdleTime)
	}
}

func TestObjectIdletimeTool_Execute_EmptyKey(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	input := `{"key": ""}`

	_, err := tool.Execute(context.Background(), json.RawMessage(input))
	if err == nil {
		t.Fatal("Expected error for empty key, got nil")
	}

	expectedError := "key cannot be empty"
	if err.Error() != expectedError {
		t.Errorf("Expected error %q, got %q", expectedError, err.Error())
	}
}

func TestObjectIdletimeTool_Execute_InvalidJSON(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	input := `{"key": }`

	_, err := tool.Execute(context.Background(), json.RawMessage(input))
	if err == nil {
		t.Fatal("Expected error for invalid JSON, got nil")
	}
}

func TestObjectIdletimeTool_Metadata(t *testing.T) {
	mockClient := &client.MockClient{}
	tool := NewTool(mockClient)

	if tool.Name() != "object_idletime" {
		t.Errorf("Expected name 'object_idletime', got %q", tool.Name())
	}

	if tool.Description() == "" {
		t.Error("Expected non-empty description")
	}

	schema := tool.InputSchema()
	if schema == nil {
		t.Error("Expected non-nil input schema")
	}
}
