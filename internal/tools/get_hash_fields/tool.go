// Package get_hash_fields implements the get_hash_fields tool.
package get_hash_fields

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
	Key    string   `json:"key" jsonschema:"required"`
	Fields []string `json:"fields" jsonschema:"required,minItems=1"`
}

type Output struct {
	Key    string            `json:"key"`
	Fields map[string]string `json:"fields"`
	Count  int               `json:"count"`
}

func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool("get_hash_fields", "Get values for specific fields in a hash", Input{}),
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
	fields, err := t.client.GetMapFields(ctx, params.Key, params.Fields)
	if err != nil {
		return nil, fmt.Errorf("failed to get hash fields: %w", err)
	}
	return Output{Key: params.Key, Fields: fields, Count: len(fields)}, nil
}

func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
