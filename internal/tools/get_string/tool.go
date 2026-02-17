// Package get_string implements the get_string tool.
package get_string

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

// Tool implements the get_string functionality.
type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

// Input represents the input for get_string tool.
type Input struct {
	Key string `json:"key" jsonschema:"required,description=Key to retrieve"`
}

// Output represents the output of get_string tool.
type Output struct {
	Key    string `json:"key"`
	Value  any    `json:"value"`
	Exists bool   `json:"exists"`
}

// NewTool creates a new get_string tool.
func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"get_string",
			"Get a string value from Valkey by key",
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

	raw, exists, err := t.client.GetString(ctx, params.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to get string for key %q: %w", params.Key, err)
	}

	return Output{
		Key:    params.Key,
		Value:  base.SafeValue(raw),
		Exists: exists,
	}, nil
}

// Init registers the tool with the registry.
func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
