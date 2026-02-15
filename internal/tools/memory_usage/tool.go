package memory_usage

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
	Key string `json:"key" jsonschema:"required,description=The key to get memory usage for"`
}

type Output struct {
	Bytes int64  `json:"bytes"`
	Key   string `json:"key"`
}

func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool("memory_usage", "Get memory usage of a key in bytes", Input{}),
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

	bytes, err := t.client.MemoryUsage(ctx, params.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to get memory usage: %w", err)
	}

	return Output{
		Bytes: bytes,
		Key:   params.Key,
	}, nil
}

func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
