// Package get_list_length implements the get_list_length tool.
package get_list_length

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

// Tool implements the get_list_length functionality.
type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

// Input represents the input for get_list_length tool.
type Input struct {
	Key string `json:"key" jsonschema:"required,description=List key"`
}

// Output represents the output of get_list_length tool.
type Output struct {
	Key    string `json:"key"`
	Length int64  `json:"length"`
	Exists bool   `json:"exists"`
}

// NewTool creates a new get_list_length tool.
func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"get_list_length",
			"Get the number of elements in a list",
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

	length, err := t.client.GetListLength(ctx, params.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to get list length for key %q: %w", params.Key, err)
	}

	return Output{
		Key:    params.Key,
		Length: length,
		Exists: length > 0,
	}, nil
}

// Init registers the tool with the registry.
func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
