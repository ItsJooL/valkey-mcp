package eval_script

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
)

func TestEvalScriptTool_Execute_Success(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	input := `{"script": "return 1", "keys": [], "args": []}`

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

func TestEvalScriptTool_Execute_EmptyScript(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	input := `{"script": ""}`

	_, err := tool.Execute(context.Background(), json.RawMessage(input))
	if err == nil {
		t.Fatal("Expected error for empty script, got nil")
	}
}

func TestEvalScriptTool_Metadata(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)

	if tool.Name() != "eval_script" {
		t.Errorf("Expected name 'eval_script', got %q", tool.Name())
	}
}
