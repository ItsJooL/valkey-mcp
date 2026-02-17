// Package get_list_index implements the get_list_index tool.
package get_list_index

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

// Tool implements the get_list_index functionality.
type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

// Input represents the input for get_list_index tool.
type Input struct {
	Key   string `json:"key" jsonschema:"required,description=List key"`
	Index int64  `json:"index" jsonschema:"required,description=Index (0-based, negative for from-end)"`
}

// Output represents the output of get_list_index tool.
type Output struct {
	Key    string `json:"key"`
	Index  int64  `json:"index"`
	Value  any    `json:"value"`
	Exists bool   `json:"exists"`
}

// NewTool creates a new get_list_index tool.
func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"get_list_index",
			"Get an element from a list by index",
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

	raw, exists, err := t.client.GetListIndex(ctx, params.Key, params.Index)
	if err != nil {
		return nil, fmt.Errorf("failed to get list element for key %q: %w", params.Key, err)
	}

	return Output{
		Key:    params.Key,
		Index:  params.Index,
		Value:  base.SafeValue(raw),
		Exists: exists,
	}, nil
}

// Init registers the tool with the registry.
func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
