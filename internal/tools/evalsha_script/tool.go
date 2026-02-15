package evalsha_script

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

type Input struct {
	SHA  string   `json:"sha" jsonschema:"required,description=SHA1 hash of the loaded script"`
	Keys []string `json:"keys" jsonschema:"description=Keys that the script will access"`
	Args []string `json:"args" jsonschema:"description=Additional arguments for the script"`
}

type Output struct {
	Result interface{} `json:"result" jsonschema:"description=Result from script execution"`
}

type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"evalsha_script",
			"Execute a previously loaded Lua script by its SHA1 hash",
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

	if params.SHA == "" {
		return nil, fmt.Errorf("sha cannot be empty")
	}

	if params.Keys == nil {
		params.Keys = []string{}
	}
	if params.Args == nil {
		params.Args = []string{}
	}

	result, err := t.client.EvalSHA(ctx, params.SHA, params.Keys, params.Args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute script: %w", err)
	}

	return Output{
		Result: result,
	}, nil
}

func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
