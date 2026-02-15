package config_set

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
)

func TestConfigSetTool_Execute_Success(t *testing.T) {
	mockClient := client.NewMockClient()

	tool := NewTool(mockClient)
	input := `{"parameter": "maxmemory", "value": "1gb"}`

	result, err := tool.Execute(context.Background(), json.RawMessage(input))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	output, ok := result.(Output)
	if !ok {
		t.Fatalf("Expected Output type, got %T", result)
	}

	if !output.Success {
		t.Error("Expected success to be true")
	}

	if output.Message == "" {
		t.Error("Expected non-empty message")
	}
}

func TestConfigSetTool_Execute_EmptyParameter(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	input := `{"parameter": "", "value": "1gb"}`

	_, err := tool.Execute(context.Background(), json.RawMessage(input))
	if err == nil {
		t.Fatal("Expected error for empty parameter, got nil")
	}

	expectedError := "parameter cannot be empty"
	if err.Error() != expectedError {
		t.Errorf("Expected error %q, got %q", expectedError, err.Error())
	}
}

func TestConfigSetTool_Execute_EmptyValue(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	input := `{"parameter": "maxmemory", "value": ""}`

	_, err := tool.Execute(context.Background(), json.RawMessage(input))
	if err == nil {
		t.Fatal("Expected error for empty value, got nil")
	}

	expectedError := "value cannot be empty"
	if err.Error() != expectedError {
		t.Errorf("Expected error %q, got %q", expectedError, err.Error())
	}
}

func TestConfigSetTool_Execute_InvalidJSON(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	input := `{"parameter": "maxmemory"`

	_, err := tool.Execute(context.Background(), json.RawMessage(input))
	if err == nil {
		t.Fatal("Expected error for invalid JSON, got nil")
	}
}

func TestConfigSetTool_Metadata(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)

	if tool.Name() != "config_set" {
		t.Errorf("Expected name 'config_set', got %q", tool.Name())
	}

	if tool.Description() == "" {
		t.Error("Expected non-empty description")
	}

	schema := tool.InputSchema()
	if schema == nil {
		t.Error("Expected non-nil input schema")
	}
}
