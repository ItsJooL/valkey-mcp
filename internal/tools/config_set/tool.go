package config_set

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

type Input struct {
	Parameter string `json:"parameter" jsonschema:"required,description=Configuration parameter name"`
	Value     string `json:"value" jsonschema:"required,description=Value to set for the parameter"`
}

type Output struct {
	Success bool   `json:"success" jsonschema:"description=Whether the configuration was updated successfully"`
	Message string `json:"message" jsonschema:"description=Result message"`
}

type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"config_set",
			"Set Redis/Valkey server configuration parameters",
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

	if params.Value == "" {
		return nil, fmt.Errorf("value cannot be empty")
	}

	success, err := t.client.ConfigSet(ctx, params.Parameter, params.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to set config: %w", err)
	}

	message := "Configuration updated successfully"
	if !success {
		message = "Configuration update failed"
	}

	return Output{
		Success: success,
		Message: message,
	}, nil
}

func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
