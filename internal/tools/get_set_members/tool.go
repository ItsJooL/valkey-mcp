// Package get_set_members implements the get_set_members tool.
package get_set_members

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

// Tool implements the get_set_members functionality.
type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

// Input represents the input for get_set_members tool.
type Input struct {
	Key string `json:"key" jsonschema:"required,description=Set key"`
}

// Output represents the output of get_set_members tool.
type Output struct {
	Key         string   `json:"key"`
	Members     []string `json:"members"`
	MemberCount int      `json:"member_count"`
	Exists      bool     `json:"exists"`
}

// NewTool creates a new get_set_members tool.
func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"get_set_members",
			"Get all members of a set",
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

	members, err := t.client.ListSetMembers(ctx, params.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to get set members for key %q: %w", params.Key, err)
	}

	return Output{
		Key:         params.Key,
		Members:     members,
		MemberCount: len(members),
		Exists:      len(members) > 0,
	}, nil
}

// Init registers the tool with the registry.
func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
