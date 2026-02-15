// Package get_hash_field implements the get_hash_field tool.
package get_hash_field

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

// Tool implements the get_hash_field functionality.
type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

// Input represents the input for get_hash_field tool.
type Input struct {
	Key   string `json:"key" jsonschema:"required,description=Hash key"`
	Field string `json:"field" jsonschema:"required,description=Field name"`
}

// Output represents the output of get_hash_field tool.
type Output struct {
	Key    string `json:"key"`
	Field  string `json:"field"`
	Value  string `json:"value"`
	Exists bool   `json:"exists"`
}

// NewTool creates a new get_hash_field tool.
func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"get_hash_field",
			"Get the value of a specific field in a hash",
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
	if params.Field == "" {
		return nil, fmt.Errorf("field cannot be empty")
	}

	value, exists, err := t.client.GetMapField(ctx, params.Key, params.Field)
	if err != nil {
		return nil, fmt.Errorf("failed to get hash field for key %q: %w", params.Key, err)
	}

	return Output{
		Key:    params.Key,
		Field:  params.Field,
		Value:  value,
		Exists: exists,
	}, nil
}

// Init registers the tool with the registry.
func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
