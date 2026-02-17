package sinter_sets

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
	Keys []string `json:"keys" jsonschema:"required,description=Set keys to intersect"`
}

type Output struct {
	Members []any `json:"members"`
	Count   int64 `json:"count"`
}

func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool("sinter_sets", "Get the intersection of multiple sets", Input{}),
		client:   client,
	}
}

func (t *Tool) Execute(ctx context.Context, input json.RawMessage) (interface{}, error) {
	var params Input
	if err := t.ParseInput(input, &params); err != nil {
		return nil, err
	}

	if len(params.Keys) == 0 {
		return nil, fmt.Errorf("keys cannot be empty")
	}

	raw, err := t.client.SetIntersection(ctx, params.Keys)
	if err != nil {
		return nil, fmt.Errorf("set intersection operation failed: %w", err)
	}

	return Output{
		Members: base.SafeSlice(raw),
		Count:   int64(len(raw)),
	}, nil
}

func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
