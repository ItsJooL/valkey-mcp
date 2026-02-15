// Package registry manages tool registration and lifecycle.
package registry

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Tool represents a single MCP tool.
type Tool interface {
	Name() string
	Description() string
	InputSchema() interface{}
	Execute(ctx context.Context, input json.RawMessage) (interface{}, error)
}

// ToolRegistry manages tool registration and lifecycle.
type ToolRegistry struct {
	tools map[string]Tool
}

// NewToolRegistry creates a new tool registry.
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]Tool),
	}
}

// Register adds a tool to the registry.
func (r *ToolRegistry) Register(tool Tool) error {
	name := tool.Name()
	if _, exists := r.tools[name]; exists {
		return fmt.Errorf("tool %s already registered", name)
	}
	r.tools[name] = tool
	return nil
}

// MustRegister registers a tool or panics.
func (r *ToolRegistry) MustRegister(tool Tool) {
	if err := r.Register(tool); err != nil {
		panic(err)
	}
}

// GetTool retrieves a tool by name.
func (r *ToolRegistry) GetTool(name string) (Tool, bool) {
	tool, exists := r.tools[name]
	return tool, exists
}

// ListTools returns all registered tool names.
func (r *ToolRegistry) ListTools() []string {
	names := make([]string, 0, len(r.tools))
	for name := range r.tools {
		names = append(names, name)
	}
	return names
}

// Count returns the number of registered tools.
func (r *ToolRegistry) Count() int {
	return len(r.tools)
}

// ToolInfo contains metadata about a tool.
type ToolInfo struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema,omitempty"`
}

// GetAllToolInfo returns metadata for all registered tools.
func (r *ToolRegistry) GetAllToolInfo() []ToolInfo {
	infos := make([]ToolInfo, 0, len(r.tools))
	for _, tool := range r.tools {
		infos = append(infos, ToolInfo{
			Name:        tool.Name(),
			Description: tool.Description(),
			InputSchema: tool.InputSchema(),
		})
	}
	return infos
}

// ExecuteTool executes a tool by name with given input.
func (r *ToolRegistry) ExecuteTool(ctx context.Context, name string, input json.RawMessage) (interface{}, error) {
	tool, exists := r.GetTool(name)
	if !exists {
		return nil, fmt.Errorf("tool not found: %s", name)
	}

	result, err := tool.Execute(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to execute tool %s: %w", name, err)
	}

	return result, nil
}

// RegisterWithMCP registers all tools with the MCP server.
func (r *ToolRegistry) RegisterWithMCP(server *mcp.Server) error {
	for _, tool := range r.tools {
		if err := r.registerSingleTool(server, tool); err != nil {
			return fmt.Errorf("failed to register tool %s: %w", tool.Name(), err)
		}
	}
	return nil
}

// registerSingleTool handles MCP registration for a single tool.
func (r *ToolRegistry) registerSingleTool(server *mcp.Server, tool Tool) error {
	mcpTool := &mcp.Tool{
		Name:        tool.Name(),
		Description: tool.Description(),
		InputSchema: tool.InputSchema(),
	}

	mcp.AddTool(server, mcpTool, func(ctx context.Context, request *mcp.CallToolRequest, args map[string]interface{}) (*mcp.CallToolResult, map[string]interface{}, error) {
		var argsJSON json.RawMessage
		if args != nil && len(args) > 0 {
			argsBytes, err := json.Marshal(args)
			if err != nil {
				return &mcp.CallToolResult{
					Content: []mcp.Content{
						&mcp.TextContent{Text: fmt.Sprintf("failed to marshal arguments: %v", err)},
					},
					IsError: true,
				}, nil, nil
			}
			argsJSON = argsBytes
		}

		result, err := tool.Execute(ctx, argsJSON)
		if err != nil {
			return nil, nil, err
		}
		resultJSON, err := json.Marshal(result)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal result: %w", err)
		}

		var resultMap map[string]interface{}
		if err := json.Unmarshal(resultJSON, &resultMap); err != nil {
			return nil, nil, fmt.Errorf("failed to unmarshal result to map: %w", err)
		}

		return nil, resultMap, nil
	})

	return nil
}
