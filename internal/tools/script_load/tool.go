package script_load

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

type Input struct {
	Script string `json:"script" jsonschema:"required,description=Lua script to load"`
}

type Output struct {
	SHA string `json:"sha" jsonschema:"description=SHA1 hash of the loaded script"`
}

type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"script_load",
			"Load a Lua script into Redis/Valkey and return its SHA1 hash",
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

	if params.Script == "" {
		return nil, fmt.Errorf("script cannot be empty")
	}

	sha, err := t.client.LoadScript(ctx, params.Script)
	if err != nil {
		return nil, fmt.Errorf("failed to load script: %w", err)
	}

	return Output{
		SHA: sha,
	}, nil
}

func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
