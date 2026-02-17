// Package tools contains integration tests for all tools against a live Valkey instance.
package tools

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"testing"
	"time"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/types"
	"github.com/stretchr/testify/assert"
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

// TestIntegration_BinaryHashFields_RoundTrip verifies that binary hash field values
// (e.g. Java-serialised POJOs stored by Redisson) are preserved without corruption
// through the full MCP retrieval path.
func TestIntegration_BinaryHashFields_RoundTrip(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	url, err := types.NewValkeyURL("valkey://localhost:6379")
	require.NoError(t, err)

	valkeyClient, err := client.New(ctx, client.Config{URL: url})
	if err != nil {
		t.Skipf("Skipping integration test: Valkey not available: %v", err)
	}
	defer valkeyClient.Close()

	reg := registry.NewToolRegistry()
	RegisterAll(reg, valkeyClient)

	// Known binary payloads — deliberately not valid UTF-8
	// javaSerialised simulates a Redisson-serialised Java object (POJO)
	javaSerialised := []byte{
		0xAC, 0xED, 0x00, 0x05, 0x73, 0x72, 0x00, 0x13,
		0x63, 0x6F, 0x6D, 0x2E, 0x65, 0x78, 0x61, 0x6D,
		0x70, 0x6C, 0x65, 0x2E, 0x55, 0x73, 0x65, 0x72,
		0xFF, 0xFE, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05,
	}
	// protobuf simulates a Protocol Buffer encoded message
	protobuf := []byte{
		0x08, 0x96, 0x01, 0x12, 0x07, 0x74, 0x65, 0x73,
		0x74, 0x69, 0x6E, 0x67, 0xFF, 0xFE, 0xFD,
	}
	plainText := "John Doe"

	const testKey = "integration_test:binary_hash"

	t.Cleanup(func() {
		valkeyClient.DeleteKey(ctx, testKey) //nolint:errcheck
	})
	valkeyClient.DeleteKey(ctx, testKey) //nolint:errcheck

	// Store binary data via SetMap; Go string preserves raw bytes through valkey-go
	_, err = valkeyClient.SetMap(ctx, testKey, map[string]string{
		"java_pojo": string(javaSerialised),
		"proto_msg": string(protobuf),
		"name":      plainText,
	})
	require.NoError(t, err, "failed to store binary hash data in Valkey")

	// Retrieve via the get_hash tool
	inputJSON, _ := json.Marshal(map[string]interface{}{"key": testKey})
	result, err := reg.ExecuteTool(ctx, "get_hash", inputJSON)
	require.NoError(t, err)
	require.NotNil(t, result)

	jsonBytes, err := json.Marshal(result)
	require.NoError(t, err)

	var jsonResult map[string]interface{}
	require.NoError(t, json.Unmarshal(jsonBytes, &jsonResult))

	fields, ok := jsonResult["fields"].(map[string]interface{})
	require.True(t, ok, "response should have fields map")

	// "name" must be a plain JSON string
	nameVal, ok := fields["name"].(string)
	require.True(t, ok, "name field should be a plain string")
	assert.Equal(t, "John Doe", nameVal)

	// "java_pojo" must be a base64 string that decodes back to original bytes exactly
	javaEncoded, ok := fields["java_pojo"].(string)
	require.True(t, ok, "java_pojo field should be a string (base64)")
	javaDecoded, err := base64.StdEncoding.DecodeString(javaEncoded)
	require.NoError(t, err, "java_pojo should be valid base64")
	assert.Equal(t, javaSerialised, javaDecoded,
		"java_pojo decoded bytes must be identical to original — binary data must not be corrupted")

	// "proto_msg" must be a base64 string that decodes back to original bytes exactly
	protoEncoded, ok := fields["proto_msg"].(string)
	require.True(t, ok, "proto_msg field should be a string (base64)")
	protoDecoded, err := base64.StdEncoding.DecodeString(protoEncoded)
	require.NoError(t, err, "proto_msg should be valid base64")
	assert.Equal(t, protobuf, protoDecoded,
		"proto_msg decoded bytes must be identical to original — binary data must not be corrupted")

	t.Logf("Integration test passed: binary hash fields preserved through full MCP retrieval path")
	t.Logf("  java_pojo: %d bytes stored, %d bytes retrieved", len(javaSerialised), len(javaDecoded))
	t.Logf("  proto_msg: %d bytes stored, %d bytes retrieved", len(protobuf), len(protoDecoded))
	t.Logf("  name: stored as text, retrieved as plain string")
}
