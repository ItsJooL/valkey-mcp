// Package rename_key implements the rename_key tool.
package rename_key

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

// Tool implements the rename_key functionality.
type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

// Input represents the input for rename_key tool.
type Input struct {
	Key    string `json:"key" jsonschema:"required,description=Current key name"`
	NewKey string `json:"new_key" jsonschema:"required,description=New key name"`
}

// Output represents the output of rename_key tool.
type Output struct {
	OldKey  string `json:"old_key"`
	NewKey  string `json:"new_key"`
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// NewTool creates a new rename_key tool.
func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"rename_key",
			"Rename a key to a new name",
			Input{},
		),
		client: client,
	}
}

func (t *Tool) Execute(ctx context.Context, input json.RawMessage) (interface{}, error) {
	var params Input
	if err := t.ParseInput(input, &params); err != nil {
		return nil, err
	}

	if params.Key == "" {
		return nil, fmt.Errorf("key cannot be empty")
	}
	if params.NewKey == "" {
		return nil, fmt.Errorf("new_key cannot be empty")
	}
	if params.Key == params.NewKey {
		return nil, fmt.Errorf("key and new_key must be different")
	}

	success, err := t.client.RenameKey(ctx, params.Key, params.NewKey)
	if err != nil {
		return nil, fmt.Errorf("failed to rename key from %q to %q: %w", params.Key, params.NewKey, err)
	}

	message := "Key renamed successfully"
	if !success {
		message = "Key does not exist or rename failed"
	}

	return Output{
		OldKey:  params.Key,
		NewKey:  params.NewKey,
		Success: success,
		Message: message,
	}, nil
}

// Init registers the tool with the registry.
func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
