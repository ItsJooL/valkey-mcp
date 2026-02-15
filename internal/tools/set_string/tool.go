// Package set_string implements the set_string tool.
package set_string

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

// Tool implements the set_string functionality.
type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

// Input represents the input for set_string tool.
type Input struct {
	Key        string `json:"key" jsonschema:"required,description=Key to set"`
	Value      string `json:"value" jsonschema:"required,description=Value to store"`
	TTLSeconds *int64 `json:"ttl_seconds,omitempty" jsonschema:"description=Optional TTL in seconds"`
	NX         bool   `json:"nx,omitempty" jsonschema:"description=Only set if key does not exist"`
	XX         bool   `json:"xx,omitempty" jsonschema:"description=Only set if key exists"`
}

// Output represents the output of set_string tool.
type Output struct {
	Success bool   `json:"success"`
	Key     string `json:"key"`
	Message string `json:"message,omitempty"`
}

// NewTool creates a new set_string tool.
func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"set_string",
			"Set a string value in Valkey with optional TTL and conditional flags (NX/XX)",
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
	if params.Value == "" {
		return nil, fmt.Errorf("value cannot be empty")
	}
	if params.NX && params.XX {
		return nil, fmt.Errorf("cannot specify both NX and XX flags")
	}

	success, err := t.client.SetString(ctx, params.Key, params.Value, params.TTLSeconds, params.NX, params.XX)
	if err != nil {
		return nil, fmt.Errorf("failed to set string for key %q: %w", params.Key, err)
	}

	message := "Key set successfully"
	if !success {
		message = "Key condition not met (NX or XX flag)"
	}

	return Output{
		Success: success,
		Key:     params.Key,
		Message: message,
	}, nil
}

// Init registers the tool with the registry.
func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
