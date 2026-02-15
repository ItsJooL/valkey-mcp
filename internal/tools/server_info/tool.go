// Package server_info implements the server_info tool.
package server_info

import (
	"context"
	"encoding/json"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

// Tool implements the server_info functionality.
type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

// Output represents the output of server_info tool.
type Output struct {
	Info map[string]string `json:"info" jsonschema:"description=Server information"`
}

// NewTool creates a new server_info tool.
func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"server_info",
			"Get server information and statistics",
			nil,
		),
		client: client,
	}
}

func (t *Tool) Execute(ctx context.Context, input json.RawMessage) (interface{}, error) {
	info, err := t.client.GetServerInfo(ctx)
	if err != nil {
		return nil, err
	}

	return Output{
		Info: info,
	}, nil
}

// Init registers the tool with the registry.
func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
