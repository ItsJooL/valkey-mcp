// Package persist_key implements the persist_key tool.
package persist_key

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

// Tool implements the persist_key functionality.
type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

// Input represents the input for persist_key tool.
type Input struct {
	Key string `json:"key" jsonschema:"required,description=Key to persist"`
}

// Output represents the output of persist_key tool.
type Output struct {
	Key     string `json:"key"`
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// NewTool creates a new persist_key tool.
func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"persist_key",
			"Remove the expiration timeout from a key (make it persistent)",
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

	success, err := t.client.PersistKey(ctx, params.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to persist key %q: %w", params.Key, err)
	}

	message := "Key made persistent"
	if !success {
		message = "Key does not exist or has no expiration"
	}

	return Output{
		Key:     params.Key,
		Success: success,
		Message: message,
	}, nil
}

// Init registers the tool with the registry.
func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
