// Package delete_hash_field implements the delete_hash_field tool.
package delete_hash_field

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

// Tool implements the delete_hash_field functionality.
type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

// Input represents the input for delete_hash_field tool.
type Input struct {
	Key    string   `json:"key" jsonschema:"required,description=Hash key"`
	Fields []string `json:"fields" jsonschema:"required,minItems=1,description=Fields to delete"`
}

// Output represents the output of delete_hash_field tool.
type Output struct {
	Key           string   `json:"key"`
	FieldsDeleted int64    `json:"fields_deleted"`
	Fields        []string `json:"fields"`
}

// NewTool creates a new delete_hash_field tool.
func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"delete_hash_field",
			"Delete one or more fields from a hash",
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

	deleted, err := t.client.DeleteMapFields(ctx, params.Key, params.Fields)
	if err != nil {
		return nil, fmt.Errorf("failed to delete hash fields for key %q: %w", params.Key, err)
	}

	return Output{
		Key:           params.Key,
		FieldsDeleted: deleted,
		Fields:        params.Fields,
	}, nil
}

// Init registers the tool with the registry.
func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
