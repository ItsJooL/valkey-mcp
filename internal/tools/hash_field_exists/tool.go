// Package hash_field_exists implements the hash_field_exists tool.
package hash_field_exists

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
	Field string `json:"field" jsonschema:"required"`
}

type Output struct {
	Key    string `json:"key"`
	Field  string `json:"field"`
	Exists bool   `json:"exists"`
}

func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool("hash_field_exists", "Check if a field exists in a hash", Input{}),
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
	exists, err := t.client.MapFieldExists(ctx, params.Key, params.Field)
	if err != nil {
		return nil, fmt.Errorf("failed to check hash field: %w", err)
	}
	return Output{Key: params.Key, Field: params.Field, Exists: exists}, nil
}

func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
