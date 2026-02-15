// Package delete_keys implements the delete_keys tool.
package delete_keys

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

// Tool implements the delete_keys functionality.
type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

// Input represents the input for delete_keys tool.
type Input struct {
	Keys []string `json:"keys" jsonschema:"required,minItems=1,description=Keys to delete"`
}

// Output represents the output of delete_keys tool.
type Output struct {
	DeletedCount int      `json:"deleted_count"`
	Keys         []string `json:"keys"`
}

// NewTool creates a new delete_keys tool.
func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"delete_keys",
			"Delete one or more keys from Valkey",
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

	count := 0
	for _, key := range params.Keys {
		if key == "" {
			return nil, fmt.Errorf("key cannot be empty")
		}
		deleted, err := t.client.DeleteKey(ctx, key)
		if err != nil {
			return nil, fmt.Errorf("failed to delete key %q: %w", key, err)
		}
		if deleted {
			count++
		}
	}

	return Output{
		DeletedCount: count,
		Keys:         params.Keys,
	}, nil
}

// Init registers the tool with the registry.
func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
