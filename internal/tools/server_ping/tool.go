// Package server_ping implements the server_ping tool.
package server_ping

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

// Tool implements the server_ping functionality.
type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

// Output represents the output of server_ping tool.
type Output struct {
	Alive     bool    `json:"alive"`
	LatencyMs float64 `json:"latency_ms"`
	Message   string  `json:"message,omitempty"`
}

// NewTool creates a new server_ping tool.
func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"server_ping",
			"Test connectivity to Valkey server and measure latency",
			nil,
		),
		client: client,
	}
}

func (t *Tool) Execute(ctx context.Context, input json.RawMessage) (interface{}, error) {
	start := time.Now()
	err := t.client.Ping(ctx)
	latency := time.Since(start).Seconds() * 1000

	if err != nil {
		return Output{
			Alive:     false,
			LatencyMs: latency,
			Message:   err.Error(),
		}, nil
	}

	return Output{
		Alive:     true,
		LatencyMs: latency,
		Message:   "Server is responding",
	}, nil
}

// Init registers the tool with the registry.
func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
