package exists_key

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
	Keys []string `json:"keys" jsonschema:"required,minItems=1,description=Array of keys to check for existence"`
}

type Output struct {
	Count int `json:"count" jsonschema:"description=Number of keys that exist"`
}

func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool("exists_key", "Check if a key exists in Valkey", Input{}),
		client:   client,
	}
}

func (t *Tool) Execute(ctx context.Context, input json.RawMessage) (interface{}, error) {
	var params Input
	if err := t.ParseInput(input, &params); err != nil {
		return nil, err
	}

	if len(params.Keys) == 0 {
		return nil, fmt.Errorf("keys array cannot be empty")
	}

	existsMap, err := t.client.ExistsKeys(ctx, params.Keys)
	if err != nil {
		return nil, fmt.Errorf("failed to check key existence: %w", err)
	}

	count := 0
	for _, exists := range existsMap {
		if exists {
			count++
		}
	}

	return Output{Count: count}, nil
}

func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
