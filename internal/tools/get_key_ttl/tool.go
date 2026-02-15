// Package get_key_ttl implements the get_key_ttl tool.
package get_key_ttl

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

// Tool implements the get_key_ttl functionality.
type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

// Input represents the input for get_key_ttl tool.
type Input struct {
	Key string `json:"key" jsonschema:"required,description=Key to check"`
}

// Output represents the output of get_key_ttl tool.
type Output struct {
	Key        string `json:"key"`
	TTLSeconds int64  `json:"ttl_seconds" jsonschema:"description=TTL in seconds (-1=no expiry, -2=does not exist)"`
	HasExpiry  bool   `json:"has_expiry"`
	Exists     bool   `json:"exists"`
}

// NewTool creates a new get_key_ttl tool.
func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"get_key_ttl",
			"Get the time-to-live (TTL) of a key in seconds",
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

	ttl, err := t.client.GetTTL(ctx, params.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to get TTL for key %q: %w", params.Key, err)
	}

	return Output{
		Key:        params.Key,
		TTLSeconds: ttl,
		HasExpiry:  ttl > 0,
		Exists:     ttl != -2,
	}, nil
}

// Init registers the tool with the registry.
func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
