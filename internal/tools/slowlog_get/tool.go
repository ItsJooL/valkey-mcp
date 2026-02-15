package slowlog_get

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

type Input struct {
	Count int64 `json:"count" jsonschema:"description=Number of slowlog entries to retrieve (0 for all)"`
}

type Output struct {
	Entries []map[string]interface{} `json:"entries" jsonschema:"description=Slowlog entries"`
	Count   int                      `json:"count" jsonschema:"description=Number of entries returned"`
}

type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"slowlog_get",
			"Get slow query log entries from Redis/Valkey server",
			Input{},
		),
		client: client,
	}
}

func (t *Tool) Execute(ctx context.Context, input json.RawMessage) (interface{}, error) {
	var params Input
	if err := t.ParseInput(input, &params); err != nil {
		return nil, err
	}

	entries, err := t.client.GetSlowlog(ctx, params.Count)
	if err != nil {
		return nil, fmt.Errorf("failed to get slowlog: %w", err)
	}

	return Output{
		Entries: entries,
		Count:   len(entries),
	}, nil
}

func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
