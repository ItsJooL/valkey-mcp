package xlen_stream

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
	Key string `json:"key" jsonschema:"required,description=Stream key"`
}

type Output struct {
	Count int64 `json:"count"`
}

func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool("xlen_stream", "Get number of entries in stream", Input{}),
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

	count, err := t.client.GetStreamLength(ctx, params.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to get stream length: %w", err)
	}

	return Output{Count: count}, nil
}

func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
