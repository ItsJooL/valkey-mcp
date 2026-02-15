package xrange_stream

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
	Start string `json:"start" jsonschema:"required,description=Start ID (- for first entry)"`
	End   string `json:"end" jsonschema:"required,description=End ID (+ for last entry)"`
	Count int64  `json:"count" jsonschema:"description=Maximum entries to return (0 for all)"`
}

type Output struct {
	Entries []map[string]string `json:"entries"`
	Count   int64               `json:"count"`
}

func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool("xrange_stream", "Get stream entries in ID range", Input{}),
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

	if params.Start == "" {
		return nil, fmt.Errorf("start cannot be empty")
	}

	if params.End == "" {
		return nil, fmt.Errorf("end cannot be empty")
	}

	entries, err := t.client.GetStreamRange(ctx, params.Key, params.Start, params.End, params.Count)
	if err != nil {
		return nil, fmt.Errorf("failed to get stream range: %w", err)
	}

	return Output{
		Entries: entries,
		Count:   int64(len(entries)),
	}, nil
}

func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
