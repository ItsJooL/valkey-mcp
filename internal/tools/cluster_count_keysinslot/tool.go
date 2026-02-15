package cluster_count_keysinslot

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

type Input struct {
	Slot int64 `json:"slot" jsonschema:"required,description=Hash slot number (0-16383)"`
}

type Output struct {
	Count int64 `json:"count" jsonschema:"description=Number of keys in the slot"`
}

type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"cluster_count_keysinslot",
			"Count the number of keys in a specific hash slot",
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

	if params.Slot < 0 || params.Slot > 16383 {
		return nil, fmt.Errorf("slot must be between 0 and 16383")
	}

	count, err := t.client.CountKeysInSlot(ctx, params.Slot)
	if err != nil {
		return nil, fmt.Errorf("failed to count keys in slot: %w", err)
	}

	return Output{
		Count: count,
	}, nil
}

func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
