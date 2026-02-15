package keys_by_pattern

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
	Pattern string `json:"pattern" jsonschema:"required,description=Key pattern to match (supports wildcards like * and ?)"`
}

type Output struct {
	Keys  []string `json:"keys"`
	Count int      `json:"count"`
}

func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool("keys_by_pattern", "Get keys matching a pattern in Valkey", Input{}),
		client:   client,
	}
}

func (t *Tool) Execute(ctx context.Context, input json.RawMessage) (interface{}, error) {
	var params Input
	if err := t.ParseInput(input, &params); err != nil {
		return nil, err
	}

	if params.Pattern == "" {
		return nil, fmt.Errorf("pattern cannot be empty")
	}

	keys, err := t.client.KeysByPattern(ctx, params.Pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to get keys by pattern: %w", err)
	}

	return Output{
		Keys:  keys,
		Count: len(keys),
	}, nil
}

func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
