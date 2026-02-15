// Package get_string_range implements the get_string_range tool.
package get_string_range

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

type Input struct {
	Key   string `json:"key" jsonschema:"required"`
	Start int64  `json:"start" jsonschema:"required"`
	End   int64  `json:"end" jsonschema:"required"`
}

type Output struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool("get_string_range", "Get a substring of a string by start and end index", Input{}),
		client:   client,
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
	value, err := t.client.GetRange(ctx, params.Key, params.Start, params.End)
	if err != nil {
		return nil, fmt.Errorf("failed to get string range: %w", err)
	}
	return Output{Key: params.Key, Value: value}, nil
}

func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
