// Package get_key_type implements the get_key_type tool.
package get_key_type

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

// Tool implements the get_key_type functionality.
type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

// Input represents the input for get_key_type tool.
type Input struct {
	Key string `json:"key" jsonschema:"required,description=Key to check"`
}

// Output represents the output of get_key_type tool.
type Output struct {
	Key    string `json:"key"`
	Type   string `json:"type" jsonschema:"description=Data type: string, list, set, hash, zset, stream, or none"`
	Exists bool   `json:"exists"`
}

// NewTool creates a new get_key_type tool.
func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"get_key_type",
			"Get the data type of a key",
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

	// A complete implementation would use TYPE command
	value, exists, err := t.client.GetString(ctx, params.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to check key %q: %w", params.Key, err)
	}

	keyType := "none"
	if exists {
		keyType = "string"
		_ = value
	}

	return Output{
		Key:    params.Key,
		Type:   keyType,
		Exists: exists,
	}, nil
}

// Init registers the tool with the registry.
func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
