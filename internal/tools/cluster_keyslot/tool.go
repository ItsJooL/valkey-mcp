package cluster_keyslot

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

type Input struct {
	Key string `json:"key" jsonschema:"required,description=Key to get the hash slot for"`
}

type Output struct {
	Slot int64 `json:"slot" jsonschema:"description=Hash slot number for the key"`
}

type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"cluster_keyslot",
			"Get the hash slot for a key in a Redis/Valkey cluster",
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

	slot, err := t.client.GetKeySlot(ctx, params.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to get key slot: %w", err)
	}

	return Output{
		Slot: slot,
	}, nil
}

func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
