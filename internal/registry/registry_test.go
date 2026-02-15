package registry

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockTool is a simple mock tool for testing
type mockTool struct {
	name        string
	description string
	schema      interface{}
	execFunc    func(ctx context.Context, input json.RawMessage) (interface{}, error)
}

func (m *mockTool) Name() string             { return m.name }
func (m *mockTool) Description() string      { return m.description }
func (m *mockTool) InputSchema() interface{} { return m.schema }
func (m *mockTool) Execute(ctx context.Context, input json.RawMessage) (interface{}, error) {
	if m.execFunc != nil {
		return m.execFunc(ctx, input)
	}
	return map[string]string{"result": "success"}, nil
}

func TestNewToolRegistry(t *testing.T) {
	reg := NewToolRegistry()
	assert.NotNil(t, reg)
	assert.Equal(t, 0, reg.Count())
}

func TestToolRegistry_Register_Success(t *testing.T) {
	reg := NewToolRegistry()
	tool := &mockTool{name: "test_tool", description: "A test tool"}

	err := reg.Register(tool)
	assert.NoError(t, err)
	assert.Equal(t, 1, reg.Count())
}

func TestToolRegistry_Register_Duplicate(t *testing.T) {
	reg := NewToolRegistry()
	tool := &mockTool{name: "test_tool", description: "A test tool"}

	err := reg.Register(tool)
	require.NoError(t, err)

	// Try to register again
	err = reg.Register(tool)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already registered")
	assert.Equal(t, 1, reg.Count())
}

func TestToolRegistry_MustRegister_Success(t *testing.T) {
	reg := NewToolRegistry()
	tool := &mockTool{name: "test_tool", description: "A test tool"}

	assert.NotPanics(t, func() {
		reg.MustRegister(tool)
	})
	assert.Equal(t, 1, reg.Count())
}

func TestToolRegistry_MustRegister_Panic(t *testing.T) {
	reg := NewToolRegistry()
	tool := &mockTool{name: "test_tool", description: "A test tool"}

	reg.MustRegister(tool)

	assert.Panics(t, func() {
		reg.MustRegister(tool)
	})
}

func TestToolRegistry_GetTool_Success(t *testing.T) {
	reg := NewToolRegistry()
	tool := &mockTool{name: "test_tool", description: "A test tool"}
	reg.Register(tool)

	retrieved, exists := reg.GetTool("test_tool")
	assert.True(t, exists)
	assert.NotNil(t, retrieved)
	assert.Equal(t, "test_tool", retrieved.Name())
}

func TestToolRegistry_GetTool_NotFound(t *testing.T) {
	reg := NewToolRegistry()

	retrieved, exists := reg.GetTool("nonexistent")
	assert.False(t, exists)
	assert.Nil(t, retrieved)
}

func TestToolRegistry_ListTools(t *testing.T) {
	reg := NewToolRegistry()
	tool1 := &mockTool{name: "tool1"}
	tool2 := &mockTool{name: "tool2"}
	tool3 := &mockTool{name: "tool3"}

	reg.Register(tool1)
	reg.Register(tool2)
	reg.Register(tool3)

	tools := reg.ListTools()
	assert.Equal(t, 3, len(tools))
	assert.Contains(t, tools, "tool1")
	assert.Contains(t, tools, "tool2")
	assert.Contains(t, tools, "tool3")
}

func TestToolRegistry_Count(t *testing.T) {
	reg := NewToolRegistry()
	assert.Equal(t, 0, reg.Count())

	reg.Register(&mockTool{name: "tool1"})
	assert.Equal(t, 1, reg.Count())

	reg.Register(&mockTool{name: "tool2"})
	assert.Equal(t, 2, reg.Count())

	reg.Register(&mockTool{name: "tool3"})
	assert.Equal(t, 3, reg.Count())
}

func TestToolRegistry_GetAllToolInfo(t *testing.T) {
	reg := NewToolRegistry()
	tool1 := &mockTool{
		name:        "tool1",
		description: "First tool",
		schema:      map[string]string{"type": "object"},
	}
	tool2 := &mockTool{
		name:        "tool2",
		description: "Second tool",
		schema:      map[string]string{"type": "string"},
	}

	reg.Register(tool1)
	reg.Register(tool2)

	infos := reg.GetAllToolInfo()
	assert.Equal(t, 2, len(infos))

	var tool1Info *ToolInfo
	for i := range infos {
		if infos[i].Name == "tool1" {
			tool1Info = &infos[i]
			break
		}
	}
	require.NotNil(t, tool1Info)
	assert.Equal(t, "tool1", tool1Info.Name)
	assert.Equal(t, "First tool", tool1Info.Description)
	assert.NotNil(t, tool1Info.InputSchema)
}

func TestToolRegistry_ExecuteTool_Success(t *testing.T) {
	reg := NewToolRegistry()
	tool := &mockTool{
		name: "test_tool",
		execFunc: func(ctx context.Context, input json.RawMessage) (interface{}, error) {
			return map[string]string{"status": "executed"}, nil
		},
	}
	reg.Register(tool)

	ctx := context.Background()
	input := json.RawMessage(`{"key": "value"}`)
	result, err := reg.ExecuteTool(ctx, "test_tool", input)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	resultMap, ok := result.(map[string]string)
	require.True(t, ok)
	assert.Equal(t, "executed", resultMap["status"])
}

func TestToolRegistry_ExecuteTool_NotFound(t *testing.T) {
	reg := NewToolRegistry()
	ctx := context.Background()
	input := json.RawMessage(`{}`)

	result, err := reg.ExecuteTool(ctx, "nonexistent", input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

func TestToolRegistry_ExecuteTool_ExecutionError(t *testing.T) {
	reg := NewToolRegistry()
	tool := &mockTool{
		name: "failing_tool",
		execFunc: func(ctx context.Context, input json.RawMessage) (interface{}, error) {
			return nil, assert.AnError
		},
	}
	reg.Register(tool)

	ctx := context.Background()
	input := json.RawMessage(`{}`)
	result, err := reg.ExecuteTool(ctx, "failing_tool", input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to execute")
}
