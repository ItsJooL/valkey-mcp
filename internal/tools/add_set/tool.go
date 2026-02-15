// Package add_set implements the add_set tool.
package add_set

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

// Tool implements the add_set functionality.
type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

// Input represents the input for add_set tool.
type Input struct {
	Key     string   `json:"key" jsonschema:"required,description=Set key"`
	Members []string `json:"members" jsonschema:"required,minItems=1,description=Members to add to the set"`
}

// Output represents the output of add_set tool.
type Output struct {
	Key          string   `json:"key"`
	MembersAdded int64    `json:"members_added"`
	Members      []string `json:"members"`
	SetSize      int64    `json:"set_size,omitempty"`
}

// NewTool creates a new add_set tool.
func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"add_set",
			"Add members to a set",
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

	added, err := t.client.AddSet(ctx, params.Key, params.Members)
	if err != nil {
		return nil, fmt.Errorf("failed to add members to set %q: %w", params.Key, err)
	}

	return Output{
		Key:          params.Key,
		MembersAdded: added,
		Members:      params.Members,
	}, nil
}

// Init registers the tool with the registry.
func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
