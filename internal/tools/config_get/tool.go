package config_get

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

type Input struct {
	Parameter string `json:"parameter" jsonschema:"required,description=Configuration parameter name to retrieve"`
}

type Output struct {
	Parameters map[string]string `json:"parameters" jsonschema:"description=Configuration parameter values"`
}

type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"config_get",
			"Get Redis/Valkey server configuration parameters",
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

	if params.Parameter == "" {
		return nil, fmt.Errorf("parameter cannot be empty")
	}

	config, err := t.client.ConfigGet(ctx, params.Parameter)
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	return Output{
		Parameters: config,
	}, nil
}

func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
