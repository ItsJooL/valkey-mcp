// Package append_string implements the append_string tool.
package append_string

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

// Tool implements the append_string functionality.
type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

// Input represents the input for append_string tool.
type Input struct {
	Key   string `json:"key" jsonschema:"required,description=Key to append to"`
	Value string `json:"value" jsonschema:"required,description=Value to append"`
}

// Output represents the output of append_string tool.
type Output struct {
	Key           string `json:"key"`
	NewLength     int64  `json:"new_length" jsonschema:"description=Length of the string after append"`
	AppendedValue string `json:"appended_value"`
}

// NewTool creates a new append_string tool.
func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"append_string",
			"Append a value to a string",
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
	if params.Value == "" {
		return nil, fmt.Errorf("value cannot be empty")
	}

	newLength, err := t.client.AppendString(ctx, params.Key, params.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to append to key %q: %w", params.Key, err)
	}

	return Output{
		Key:           params.Key,
		NewLength:     newLength,
		AppendedValue: params.Value,
	}, nil
}

// Init registers the tool with the registry.
func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
