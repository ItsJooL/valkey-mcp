package object_encoding

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
	Key string `json:"key" jsonschema:"required,description=The key to get encoding type for"`
}

type Output struct {
	Encoding string `json:"encoding"`
	Key      string `json:"key"`
}

func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool("object_encoding", "Get the encoding type of a key's value in Valkey", Input{}),
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

	encoding, err := t.client.ObjectEncoding(ctx, params.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to get object encoding: %w", err)
	}

	return Output{
		Encoding: encoding,
		Key:      params.Key,
	}, nil
}

func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
