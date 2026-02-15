// Package client_list implements the client_list tool.
package client_list

import (
	"context"
	"encoding/json"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

// Tool implements the client_list functionality.
type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

// Input represents the input for client_list tool.
type Input struct {
	// Empty for now - CLIENT LIST has no required parameters
}

// Output represents the output of client_list tool.
type Output struct {
	ClientCount int    `json:"client_count" jsonschema:"description=Number of connected clients"`
	RawInfo     string `json:"raw_info" jsonschema:"description=Raw CLIENT LIST output"`
	Message     string `json:"message,omitempty"`
}

// NewTool creates a new client_list tool.
func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"client_list",
			"List all client connections to the Valkey server",
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

	return Output{
		ClientCount: 1,
		RawInfo:     "CLIENT LIST not fully implemented - requires underlying client method",
		Message:     "This tool requires additional client interface methods to be fully implemented",
	}, nil
}

// Init registers the tool with the registry.
func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
