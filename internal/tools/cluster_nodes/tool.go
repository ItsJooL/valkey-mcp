package cluster_nodes

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

type Input struct {
}

type Output struct {
	Nodes string `json:"nodes" jsonschema:"description=Cluster nodes information (raw format)"`
}

type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool(
			"cluster_nodes",
			"Get information about all nodes in the Redis/Valkey cluster",
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

	nodes, err := t.client.GetClusterNodes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster nodes: %w", err)
	}

	return Output{
		Nodes: nodes,
	}, nil
}

func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
