package hlen_hash

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
	Key string `json:"key" jsonschema:"required,description=Hash key"`
}

type Output struct {
	Key    string `json:"key"`
	Result int64  `json:"result"`
}

func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool("hlen_hash", "Get the number of fields in a hash", Input{}),
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

	result, err := t.client.GetMapLength(ctx, params.Key)
	if err != nil {
		return nil, fmt.Errorf("operation failed: %w", err)
	}

	return Output{Key: params.Key, Result: result}, nil
}

func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
