// Package lrange_list implements the lrange_list tool.
package lrange_list

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

// Tool implements the lrange_list functionality.
type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

// Input represents the input for lrange_list tool.
type Input struct {
	Key   string `json:"key" jsonschema:"required,description=List key"`
	Start int64  `json:"start" jsonschema:"required,description=Start index (0-based, negative for from-end)"`
	Stop  int64  `json:"stop" jsonschema:"required,description=Stop index (inclusive, negative for from-end)"`
}

// Output represents the output of lrange_list tool.
type Output struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
	Count  int      `json:"count"`
}

// NewTool creates a new lrange_list tool.
func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"lrange_list",
			"Get a range of elements from a list",
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

	values, err := t.client.GetListRange(ctx, params.Key, params.Start, params.Stop)
	if err != nil {
		return nil, fmt.Errorf("failed to get list range for key %q: %w", params.Key, err)
	}

	return Output{
		Key:    params.Key,
		Values: values,
		Count:  len(values),
	}, nil
}

// Init registers the tool with the registry.
func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
