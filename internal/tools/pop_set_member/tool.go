// Package pop_set_member implements the pop_set_member tool.
package pop_set_member

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
	Count int64  `json:"count,omitempty" jsonschema:"description=Number of members to pop (default: 1)"`
}

type Output struct {
	Key     string `json:"key"`
	Members []any  `json:"members"`
	Count   int    `json:"count"`
}

func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool("pop_set_member", "Remove and return random members from a set", Input{}),
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
	if params.Count == 0 {
		params.Count = 1
	}
	raw, err := t.client.PopSet(ctx, params.Key, params.Count)
	if err != nil {
		return nil, fmt.Errorf("failed to pop from set: %w", err)
	}
	return Output{Key: params.Key, Members: base.SafeSlice(raw), Count: len(raw)}, nil
}

func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
