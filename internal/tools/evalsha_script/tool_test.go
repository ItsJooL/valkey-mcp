package evalsha_script

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
)

func TestEvalshaScriptTool_Execute_Success(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	input := `{"sha": "abc123", "keys": [], "args": []}`

	result, err := tool.Execute(context.Background(), json.RawMessage(input))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	output, ok := result.(Output)
	if !ok {
		t.Fatalf("Expected Output type, got %T", result)
	}

	if output.Result == nil {
		t.Error("Expected non-nil result")
	}
}

func TestEvalshaScriptTool_Execute_EmptySHA(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	input := `{"sha": ""}`

	_, err := tool.Execute(context.Background(), json.RawMessage(input))
	if err == nil {
		t.Fatal("Expected error for empty SHA, got nil")
	}
}

func TestEvalshaScriptTool_Metadata(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)

	if tool.Name() != "evalsha_script" {
		t.Errorf("Expected name 'evalsha_script', got %q", tool.Name())
	}
}
