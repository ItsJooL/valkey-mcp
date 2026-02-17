package sdiff_sets

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
	Keys []string `json:"keys" jsonschema:"required,minItems=2,description=Array of set keys - first key is subtracted from, remaining keys are subtracted"`
}

type Output struct {
	Members []any `json:"members"`
	Count   int64 `json:"count"`
}

func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool("sdiff_sets", "Get the difference of sets (members in first set but not in others)", Input{}),
		client:   client,
	}
}

func (t *Tool) Execute(ctx context.Context, input json.RawMessage) (interface{}, error) {
	var params Input
	if err := t.ParseInput(input, &params); err != nil {
		return nil, err
	}

	if len(params.Keys) < 2 {
		return nil, fmt.Errorf("at least 2 keys required for set difference")
	}

	firstKey := params.Keys[0]
	otherKeys := params.Keys[1:]

	raw, err := t.client.SetDifference(ctx, firstKey, otherKeys)
	if err != nil {
		return nil, fmt.Errorf("set difference operation failed: %w", err)
	}

	return Output{
		Members: base.SafeSlice(raw),
		Count:   int64(len(raw)),
	}, nil
}

func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
