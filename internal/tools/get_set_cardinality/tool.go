// Package get_set_cardinality implements the get_set_cardinality tool.
package get_set_cardinality

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

// Tool implements the get_set_cardinality functionality.
type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

// Input represents the input for get_set_cardinality tool.
type Input struct {
	Key string `json:"key" jsonschema:"required,description=Set key"`
}

// Output represents the output of get_set_cardinality tool.
type Output struct {
	Key         string `json:"key"`
	Cardinality int64  `json:"cardinality" jsonschema:"description=Number of members in the set"`
	Exists      bool   `json:"exists"`
}

// NewTool creates a new get_set_cardinality tool.
func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"get_set_cardinality",
			"Get the number of members in a set (cardinality)",
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

	cardinality, err := t.client.GetSetSize(ctx, params.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to get set cardinality for key %q: %w", params.Key, err)
	}

	return Output{
		Key:         params.Key,
		Cardinality: cardinality,
		Exists:      cardinality > 0,
	}, nil
}

// Init registers the tool with the registry.
func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
