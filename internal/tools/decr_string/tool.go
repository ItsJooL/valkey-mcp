// Package decr_string implements the decr_string tool.
package decr_string

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

// Tool implements the decr_string functionality.
type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

// Input represents the input for decr_string tool.
type Input struct {
	Key    string `json:"key" jsonschema:"required,description=Key storing a numeric string"`
	Amount int64  `json:"amount,omitempty" jsonschema:"description=Amount to decrement (default: 1)"`
}

// Output represents the output of decr_string tool.
type Output struct {
	Key     string `json:"key"`
	Value   int64  `json:"value"`
	Message string `json:"message,omitempty"`
}

// NewTool creates a new decr_string tool.
func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"decr_string",
			"Decrement a numeric string value",
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

	if params.Amount == 0 {
		params.Amount = 1
	}

	newValue, err := t.client.DecrementNumber(ctx, params.Key, params.Amount)
	if err != nil {
		return nil, fmt.Errorf("failed to decrement key %q: %w", params.Key, err)
	}

	return Output{
		Key:     params.Key,
		Value:   newValue,
		Message: fmt.Sprintf("Decremented by %d", params.Amount),
	}, nil
}

// Init registers the tool with the registry.
func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
