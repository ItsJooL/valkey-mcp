// Package mget_strings implements the mget_strings tool.
package mget_strings

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

// Tool implements the mget_strings functionality.
type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

// Input represents the input for mget_strings tool.
type Input struct {
	Keys []string `json:"keys" jsonschema:"required,minItems=1,description=Keys to retrieve"`
}

// Output represents the output of mget_strings tool.
type Output struct {
	Values map[string]any `json:"values" jsonschema:"description=Key-value pairs (missing keys excluded)"`
	Count  int            `json:"count"`
}

// NewTool creates a new mget_strings tool.
func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"mget_strings",
			"Get multiple string values from Valkey by keys",
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

	if len(params.Keys) == 0 {
		return nil, fmt.Errorf("at least one key must be provided")
	}

	values := make(map[string]any)
	for _, key := range params.Keys {
		if key == "" {
			return nil, fmt.Errorf("key cannot be empty")
		}
		raw, exists, err := t.client.GetString(ctx, key)
		if err != nil {
			return nil, fmt.Errorf("failed to get string for key %q: %w", key, err)
		}
		if exists {
			values[key] = base.SafeValue(raw)
		}
	}

	return Output{
		Values: values,
		Count:  len(values),
	}, nil
}

// Init registers the tool with the registry.
func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
