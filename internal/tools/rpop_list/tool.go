// Package rpop_list implements the rpop_list tool.
package rpop_list

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

// Tool implements the rpop_list functionality.
type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

// Input represents the input for rpop_list tool.
type Input struct {
	Key   string `json:"key" jsonschema:"required,description=List key"`
	Count int64  `json:"count,omitempty" jsonschema:"description=Number of elements to pop (default: 1)"`
}

// Output represents the output of rpop_list tool.
type Output struct {
	Key      string   `json:"key"`
	Elements []string `json:"elements"`
	Count    int      `json:"count"`
}

// NewTool creates a new rpop_list tool.
func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"rpop_list",
			"Remove and return elements from the right (tail) of a list",
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

	if params.Count == 0 {
		params.Count = 1
	}

	elements, err := t.client.PopList(ctx, params.Key, params.Count, true)
	if err != nil {
		return nil, fmt.Errorf("failed to pop from list %q: %w", params.Key, err)
	}

	return Output{
		Key:      params.Key,
		Elements: elements,
		Count:    len(elements),
	}, nil
}

// Init registers the tool with the registry.
func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
