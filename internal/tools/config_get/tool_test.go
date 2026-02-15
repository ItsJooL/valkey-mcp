package config_get

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
)

func TestConfigGetTool_Execute_Success(t *testing.T) {
	mockClient := client.NewMockClient()

	tool := NewTool(mockClient)
	input := `{"parameter": "maxmemory"}`

	result, err := tool.Execute(context.Background(), json.RawMessage(input))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	output, ok := result.(Output)
	if !ok {
		t.Fatalf("Expected Output type, got %T", result)
	}

	if output.Parameters == nil {
		t.Error("Expected non-nil parameters map")
	}
}

func TestConfigGetTool_Execute_EmptyParameter(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	input := `{"parameter": ""}`

	_, err := tool.Execute(context.Background(), json.RawMessage(input))
	if err == nil {
		t.Fatal("Expected error for empty parameter, got nil")
	}

	expectedError := "parameter cannot be empty"
	if err.Error() != expectedError {
		t.Errorf("Expected error %q, got %q", expectedError, err.Error())
	}
}

func TestConfigGetTool_Execute_InvalidJSON(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	input := `{"parameter": }`

	_, err := tool.Execute(context.Background(), json.RawMessage(input))
	if err == nil {
		t.Fatal("Expected error for invalid JSON, got nil")
	}
}

func TestConfigGetTool_Metadata(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)

	if tool.Name() != "config_get" {
		t.Errorf("Expected name 'config_get', got %q", tool.Name())
	}

	if tool.Description() == "" {
		t.Error("Expected non-empty description")
	}

	schema := tool.InputSchema()
	if schema == nil {
		t.Error("Expected non-nil input schema")
	}
}
