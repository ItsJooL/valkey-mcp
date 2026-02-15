// Package set_is_member implements the set_is_member tool.
package set_is_member

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

// Tool implements the set_is_member functionality.
type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

// Input represents the input for set_is_member tool.
type Input struct {
	Key    string `json:"key" jsonschema:"required,description=Set key"`
	Member string `json:"member" jsonschema:"required,description=Member to check"`
}

// Output represents the output of set_is_member tool.
type Output struct {
	Key      string `json:"key"`
	Member   string `json:"member"`
	IsMember bool   `json:"is_member"`
}

// NewTool creates a new set_is_member tool.
func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"set_is_member",
			"Check if a member exists in a set",
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
	if params.Member == "" {
		return nil, fmt.Errorf("member cannot be empty")
	}

	isMember, err := t.client.CheckSetMember(ctx, params.Key, params.Member)
	if err != nil {
		return nil, fmt.Errorf("failed to check set membership for key %q: %w", params.Key, err)
	}

	return Output{
		Key:      params.Key,
		Member:   params.Member,
		IsMember: isMember,
	}, nil
}

// Init registers the tool with the registry.
func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
