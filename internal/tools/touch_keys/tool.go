package touch_keys

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
	Keys []string `json:"keys" jsonschema:"required,description=List of keys to update access time for"`
}

type Output struct {
	Count   int64    `json:"count"`
	Keys    []string `json:"keys"`
	Updated int64    `json:"updated"`
}

func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool("touch_keys", "Update access time for multiple keys in Valkey", Input{}),
		client:   client,
	}
}

func (t *Tool) Execute(ctx context.Context, input json.RawMessage) (interface{}, error) {
	var params Input
	if err := t.ParseInput(input, &params); err != nil {
		return nil, err
	}

	if len(params.Keys) == 0 {
		return nil, fmt.Errorf("keys list cannot be empty")
	}

	count, err := t.client.TouchKeys(ctx, params.Keys)
	if err != nil {
		return nil, fmt.Errorf("failed to touch keys: %w", err)
	}

	return Output{
		Count:   int64(len(params.Keys)),
		Keys:    params.Keys,
		Updated: count,
	}, nil
}

func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
