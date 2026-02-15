// Package scan_keys implements the scan_keys tool.
package scan_keys

import (
	"context"
	"encoding/json"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

// Tool implements the scan_keys functionality.
type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

// Input represents the input for scan_keys tool.
type Input struct {
	Pattern string `json:"pattern,omitempty" jsonschema:"description=Glob pattern to filter keys (default: *)"`
	Count   int64  `json:"count,omitempty" jsonschema:"description=Approximate number of keys to return"`
}

// Output represents the output of scan_keys tool.
type Output struct {
	Keys    []string `json:"keys"`
	Count   int      `json:"count"`
	Pattern string   `json:"pattern"`
}

// NewTool creates a new scan_keys tool.
func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"scan_keys",
			"Scan keys matching a pattern (non-blocking alternative to KEYS)",
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

	if params.Pattern == "" {
		params.Pattern = "*"
	}

	keys := []string{}

	return Output{
		Keys:    keys,
		Count:   len(keys),
		Pattern: params.Pattern,
	}, nil
}

// Init registers the tool with the registry.
func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
