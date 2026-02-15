// Package incr_hash_field implements the incr_hash_field tool.
package incr_hash_field

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
	Key    string `json:"key" jsonschema:"required"`
	Field  string `json:"field" jsonschema:"required"`
	Amount int64  `json:"amount,omitempty" jsonschema:"description=Amount to increment (default: 1)"`
}

type Output struct {
	Key      string `json:"key"`
	Field    string `json:"field"`
	NewValue int64  `json:"new_value"`
}

func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool("incr_hash_field", "Increment a numeric field in a hash", Input{}),
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
	if params.Field == "" {
		return nil, fmt.Errorf("field cannot be empty")
	}
	if params.Amount == 0 {
		params.Amount = 1
	}
	value, err := t.client.IncrementMapField(ctx, params.Key, params.Field, params.Amount)
	if err != nil {
		return nil, fmt.Errorf("failed to increment hash field: %w", err)
	}
	return Output{Key: params.Key, Field: params.Field, NewValue: value}, nil
}

func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
