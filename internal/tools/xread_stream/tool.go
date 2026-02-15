package xread_stream

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
	Key   string `json:"key" jsonschema:"required,description=Stream key"`
	ID    string `json:"id" jsonschema:"required,description=Start ID ($ for new entries, 0 for first)"`
	Count int64  `json:"count" jsonschema:"description=Maximum entries to return (0 for all)"`
}

type Output struct {
	Entries []map[string]string `json:"entries"`
	Count   int64               `json:"count"`
}

func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool("xread_stream", "Read entries from stream starting at ID", Input{}),
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

	if params.ID == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	entries, err := t.client.ReadStream(ctx, params.Key, params.ID, params.Count)
	if err != nil {
		return nil, fmt.Errorf("failed to read stream: %w", err)
	}

	return Output{
		Entries: entries,
		Count:   int64(len(entries)),
	}, nil
}

func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
