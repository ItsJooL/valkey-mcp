// Package tools contains integration tests for all tools against a live Valkey instance.
package tools

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/types"
	"github.com/stretchr/testify/require"
)

// TestIntegrationBasicToolsWithLiveValkey tests against a real Valkey instance on localhost:6379
func TestIntegrationBasicToolsWithLiveValkey(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create client - will skip if Valkey is not running
	url, err := types.NewValkeyURL("valkey://localhost:6379")
	require.NoError(t, err)

	valkeyClient, err := client.New(ctx, client.Config{
		URL:      url,
		Password: "",
		DB:       0,
	})
	if err != nil {
		t.Skipf("Skipping integration test: Valkey not available on localhost:6379 - %v", err)
	}
	defer valkeyClient.Close()

	// Create registry and register tools
	reg := registry.NewToolRegistry()
	RegisterAll(reg, valkeyClient)

	toolCount := reg.Count()
	t.Logf("Connected to Valkey on localhost:6379")
	t.Logf("Registered %d tools", toolCount)
	require.Greater(t, toolCount, 0, "should have registered tools")

	// Test 1: PING
	t.Run("server_ping", func(t *testing.T) {
		input := json.RawMessage(`{}`)
		result, err := reg.ExecuteTool(ctx, "server_ping", input)
		require.NoError(t, err)
		require.NotNil(t, result)
		t.Logf("server_ping: %v", result)
	})

	// Test 2: SET string
	t.Run("set_string", func(t *testing.T) {
		input := map[string]interface{}{
			"key":   "integration_test_key",
			"value": "test_value",
		}
		inputJSON, _ := json.Marshal(input)
		result, err := reg.ExecuteTool(ctx, "set_string", inputJSON)
		require.NoError(t, err)
		require.NotNil(t, result)
		t.Logf("set_string: %v", result)
	})

	// Test 3: GET string
	t.Run("get_string", func(t *testing.T) {
		input := map[string]interface{}{
			"key": "integration_test_key",
		}
		inputJSON, _ := json.Marshal(input)
		result, err := reg.ExecuteTool(ctx, "get_string", inputJSON)
		require.NoError(t, err)
		require.NotNil(t, result)
		t.Logf("get_string: %v", result)
	})

	// Test 4: DELETE keys
	t.Run("delete_keys", func(t *testing.T) {
		input := map[string]interface{}{
			"keys": []string{"integration_test_key"},
		}
		inputJSON, _ := json.Marshal(input)
		result, err := reg.ExecuteTool(ctx, "delete_keys", inputJSON)
		require.NoError(t, err)
		require.NotNil(t, result)
		t.Logf("delete_keys: %v", result)
	})

	// Test 5: LPUSH list
	t.Run("lpush_list", func(t *testing.T) {
		input := map[string]interface{}{
			"key":    "integration_test_list",
			"values": []string{"item1", "item2", "item3"},
		}
		inputJSON, _ := json.Marshal(input)
		result, err := reg.ExecuteTool(ctx, "lpush_list", inputJSON)
		require.NoError(t, err)
		require.NotNil(t, result)
		t.Logf("lpush_list: %v", result)
	})

	// Test 6: LRANGE list
	t.Run("lrange_list", func(t *testing.T) {
		input := map[string]interface{}{
			"key":   "integration_test_list",
			"start": 0,
			"stop":  -1,
		}
		inputJSON, _ := json.Marshal(input)
		result, err := reg.ExecuteTool(ctx, "lrange_list", inputJSON)
		require.NoError(t, err)
		require.NotNil(t, result)
		t.Logf("lrange_list: %v", result)
	})

	// Test 7: HSET hash
	t.Run("set_hash", func(t *testing.T) {
		input := map[string]interface{}{
			"key": "integration_test_hash",
			"fields": map[string]interface{}{
				"field1": "value1",
				"field2": "value2",
			},
		}
		inputJSON, _ := json.Marshal(input)
		result, err := reg.ExecuteTool(ctx, "set_hash", inputJSON)
		require.NoError(t, err)
		require.NotNil(t, result)
		t.Logf("set_hash: %v", result)
	})

	// Test 8: HGET hash field
	t.Run("get_hash_field", func(t *testing.T) {
		input := map[string]interface{}{
			"key":   "integration_test_hash",
			"field": "field1",
		}
		inputJSON, _ := json.Marshal(input)
		result, err := reg.ExecuteTool(ctx, "get_hash_field", inputJSON)
		require.NoError(t, err)
		require.NotNil(t, result)
		t.Logf("get_hash_field: %v", result)
	})

	// Test 9: SADD set
	t.Run("add_set", func(t *testing.T) {
		input := map[string]interface{}{
			"key":     "integration_test_set",
			"members": []string{"member1", "member2", "member3"},
		}
		inputJSON, _ := json.Marshal(input)
		result, err := reg.ExecuteTool(ctx, "add_set", inputJSON)
		require.NoError(t, err)
		require.NotNil(t, result)
		t.Logf("add_set: %v", result)
	})

	// Test 10: SMEMBERS set
	t.Run("get_set_members", func(t *testing.T) {
		input := map[string]interface{}{
			"key": "integration_test_set",
		}
		inputJSON, _ := json.Marshal(input)
		result, err := reg.ExecuteTool(ctx, "get_set_members", inputJSON)
		require.NoError(t, err)
		require.NotNil(t, result)
		t.Logf("get_set_members: %v", result)
	})

	// Cleanup
	t.Run("cleanup", func(t *testing.T) {
		input := map[string]interface{}{
			"keys": []string{
				"integration_test_key",
				"integration_test_list",
				"integration_test_hash",
				"integration_test_set",
			},
		}
		inputJSON, _ := json.Marshal(input)
		result, err := reg.ExecuteTool(ctx, "delete_keys", inputJSON)
		require.NoError(t, err)
		t.Logf("cleanup: %v", result)
	})

	t.Logf("All integration tests passed")
}
