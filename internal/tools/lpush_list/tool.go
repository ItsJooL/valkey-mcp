// Package lpush_list implements the lpush_list tool.
package lpush_list

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

// Tool implements the lpush_list functionality.
type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

// Input represents the input for lpush_list tool.
type Input struct {
	Key    string   `json:"key" jsonschema:"required,description=List key"`
	Values []string `json:"values" jsonschema:"required,minItems=1,description=Values to push to the left"`
}

// Output represents the output of lpush_list tool.
type Output struct {
	Key        string   `json:"key"`
	ListLength int64    `json:"list_length" jsonschema:"description=Length of the list after push"`
	Values     []string `json:"values"`
}

// NewTool creates a new lpush_list tool.
func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"lpush_list",
			"Push values to the left (head) of a list",
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
	if len(params.Values) == 0 {
		return nil, fmt.Errorf("at least one value must be provided")
	}

	length, err := t.client.PushList(ctx, params.Key, params.Values, false)
	if err != nil {
		return nil, fmt.Errorf("failed to push to list %q: %w", params.Key, err)
	}

	return Output{
		Key:        params.Key,
		ListLength: length,
		Values:     params.Values,
	}, nil
}

// Init registers the tool with the registry.
func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
