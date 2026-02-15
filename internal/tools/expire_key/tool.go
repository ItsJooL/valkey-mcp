// Package expire_key implements the expire_key tool.
package expire_key

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

// Tool implements the expire_key functionality.
type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

// Input represents the input for expire_key tool.
type Input struct {
	Key     string `json:"key" jsonschema:"required,description=Key to expire"`
	Seconds int64  `json:"seconds" jsonschema:"required,minimum=1,description=Seconds until expiration"`
}

// Output represents the output of expire_key tool.
type Output struct {
	Key     string `json:"key"`
	Seconds int64  `json:"seconds"`
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// NewTool creates a new expire_key tool.
func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"expire_key",
			"Set an expiration time (TTL) on a key",
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
	if params.Seconds <= 0 {
		return nil, fmt.Errorf("seconds must be positive")
	}

	success, err := t.client.ExpireKey(ctx, params.Key, params.Seconds)
	if err != nil {
		return nil, fmt.Errorf("failed to set expiration for key %q: %w", params.Key, err)
	}

	message := "Expiration set successfully"
	if !success {
		message = "Key does not exist"
	}

	return Output{
		Key:     params.Key,
		Seconds: params.Seconds,
		Success: success,
		Message: message,
	}, nil
}

// Init registers the tool with the registry.
func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
