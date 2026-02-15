package xadd_stream

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
	Key    string            `json:"key" jsonschema:"required,description=Stream key"`
	ID     string            `json:"id" jsonschema:"description=Stream entry ID (* for auto-generate, defaults to *)"`
	Fields map[string]string `json:"fields" jsonschema:"required,description=Field-value pairs"`
}

type Output struct {
	ID string `json:"id"`
}

func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool("xadd_stream", "Add entry to stream with specified fields", Input{}),
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

	if len(params.Fields) == 0 {
		return nil, fmt.Errorf("fields cannot be empty")
	}

	// Default to auto-generate ID if not provided
	id := params.ID
	if id == "" {
		id = "*"
	}

	resultID, err := t.client.AddStream(ctx, params.Key, id, params.Fields)
	if err != nil {
		return nil, fmt.Errorf("failed to add stream entry: %w", err)
	}

	return Output{ID: resultID}, nil
}

func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
