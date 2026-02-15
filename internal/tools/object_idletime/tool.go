package object_idletime

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

type Input struct {
	Key string `json:"key" jsonschema:"required,description=Key to check idle time for"`
}

type Output struct {
	IdleTime int64 `json:"idle_time" jsonschema:"description=Idle time in seconds since last access"`
}

type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"object_idletime",
			"Get the idle time (time since last access) of a key in seconds",
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

	if params.Key == "" {
		return nil, fmt.Errorf("key cannot be empty")
	}

	idleTime, err := t.client.ObjectIdletime(ctx, params.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to get object idle time: %w", err)
	}

	return Output{
		IdleTime: idleTime,
	}, nil
}

func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
