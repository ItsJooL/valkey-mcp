// Package get_hash implements the get_hash tool.
package get_hash

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

// Tool implements the get_hash functionality.
type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

// Input represents the input for get_hash tool.
type Input struct {
	Key string `json:"key" jsonschema:"required,description=Hash key"`
}

// Output represents the output of get_hash tool.
type Output struct {
	Key        string         `json:"key"`
	Fields     map[string]any `json:"fields"`
	FieldCount int            `json:"field_count"`
	Exists     bool           `json:"exists"`
}

// NewTool creates a new get_hash tool.
func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"get_hash",
			"Get all fields and values of a hash",
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

	raw, err := t.client.GetMap(ctx, params.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to get hash for key %q: %w", params.Key, err)
	}

	return Output{
		Key:        params.Key,
		Fields:     base.SafeMap(raw),
		FieldCount: len(raw),
		Exists:     len(raw) > 0,
	}, nil
}

// Init registers the tool with the registry.
func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
