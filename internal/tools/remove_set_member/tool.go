// Package remove_set_member implements the remove_set_member tool.
package remove_set_member

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

// Tool implements the remove_set_member functionality.
type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

// Input represents the input for remove_set_member tool.
type Input struct {
	Key     string   `json:"key" jsonschema:"required,description=Set key"`
	Members []string `json:"members" jsonschema:"required,minItems=1,description=Members to remove from the set"`
}

// Output represents the output of remove_set_member tool.
type Output struct {
	Key            string   `json:"key"`
	MembersRemoved int64    `json:"members_removed"`
	Members        []string `json:"members"`
}

// NewTool creates a new remove_set_member tool.
func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"remove_set_member",
			"Remove members from a set",
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
	if len(params.Members) == 0 {
		return nil, fmt.Errorf("at least one member must be provided")
	}

	removed, err := t.client.RemoveSet(ctx, params.Key, params.Members)
	if err != nil {
		return nil, fmt.Errorf("failed to remove members from set %q: %w", params.Key, err)
	}

	return Output{
		Key:            params.Key,
		MembersRemoved: removed,
		Members:        params.Members,
	}, nil
}

// Init registers the tool with the registry.
func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
