// Package lset_list implements the lset_list tool.
package lset_list

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

// Tool implements the lset_list functionality.
type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

// Input represents the input for lset_list tool.
type Input struct {
	Key   string `json:"key" jsonschema:"required,description=List key"`
	Index int64  `json:"index" jsonschema:"required,description=Index (0-based, negative for from-end)"`
	Value string `json:"value" jsonschema:"required,description=Value to set"`
}

// Output represents the output of lset_list tool.
type Output struct {
	Key     string `json:"key"`
	Index   int64  `json:"index"`
	Value   string `json:"value"`
	Success bool   `json:"success"`
}

// NewTool creates a new lset_list tool.
func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"lset_list",
			"Set the value of an element in a list by index",
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

	success, err := t.client.SetListIndex(ctx, params.Key, params.Index, params.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to set list element for key %q at index %d: %w", params.Key, params.Index, err)
	}

	return Output{
		Key:     params.Key,
		Index:   params.Index,
		Value:   params.Value,
		Success: success,
	}, nil
}

// Init registers the tool with the registry.
func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
