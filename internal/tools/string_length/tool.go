// Package string_length implements the string_length tool.
package string_length

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

// Tool implements the string_length functionality.
type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

// Input represents the input for string_length tool.
type Input struct {
	Key string `json:"key" jsonschema:"required,description=String key"`
}

// Output represents the output of string_length tool.
type Output struct {
	Key    string `json:"key"`
	Length int64  `json:"length"`
	Exists bool   `json:"exists"`
}

// NewTool creates a new string_length tool.
func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"string_length",
			"Get the length of a string value",
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

	// Get the string value to determine its length
	value, exists, err := t.client.GetString(ctx, params.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to get string for key %q: %w", params.Key, err)
	}

	return Output{
		Key:    params.Key,
		Length: int64(len(value)),
		Exists: exists,
	}, nil
}

// Init registers the tool with the registry.
func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
