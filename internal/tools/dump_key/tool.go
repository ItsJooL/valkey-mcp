package dump_key

import (
	"context"
	"encoding/base64"
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
	Key string `json:"key" jsonschema:"required,description=Key to serialize"`
}

type Output struct {
	Serialized string `json:"serialized"`
	Size       int    `json:"size"`
}

func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool("dump_key", "Serialize value of key (returns base64-encoded serialization)", Input{}),
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

	serialized, err := t.client.DumpKey(ctx, params.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to dump key: %w", err)
	}

	// Encode binary data as base64 for JSON compatibility
	encoded := base64.StdEncoding.EncodeToString(serialized)

	return Output{
		Serialized: encoded,
		Size:       len(serialized),
	}, nil
}

func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
