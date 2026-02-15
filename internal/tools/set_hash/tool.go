// Package set_hash implements the set_hash tool.
package set_hash

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

// Tool implements the set_hash functionality.
type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

// Input represents the input for set_hash tool.
type Input struct {
	Key    string            `json:"key" jsonschema:"required,description=Hash key"`
	Fields map[string]string `json:"fields" jsonschema:"required,description=Fields to set"`
}

// Output represents the output of set_hash tool.
type Output struct {
	Key           string `json:"key"`
	FieldsAdded   int64  `json:"fields_added" jsonschema:"description=Number of new fields added"`
	FieldsUpdated int64  `json:"fields_updated,omitempty" jsonschema:"description=Number of fields updated"`
	Message       string `json:"message,omitempty"`
}

// NewTool creates a new set_hash tool.
func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"set_hash",
			"Set multiple fields in a hash",
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
	if len(params.Fields) == 0 {
		return nil, fmt.Errorf("at least one field must be provided")
	}

	count, err := t.client.SetMap(ctx, params.Key, params.Fields)
	if err != nil {
		return nil, fmt.Errorf("failed to set hash fields for key %q: %w", params.Key, err)
	}

	return Output{
		Key:         params.Key,
		FieldsAdded: count,
		Message:     fmt.Sprintf("Set %d field(s)", len(params.Fields)),
	}, nil
}

// Init registers the tool with the registry.
func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
