package dump_key

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var javaBinaryFixture = []byte{0xAC, 0xED, 0x00, 0x05, 0x74, 0x00, 0x04, 0x54, 0x65, 0x73, 0x74}

func TestDumpKeyTool_BinaryValue_IsBase64InOutput(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	mockClient.SetRawBytes("binary_key", javaBinaryFixture)

	inputJSON, _ := json.Marshal(map[string]interface{}{"key": "binary_key"})
	result, err := tool.Execute(ctx, inputJSON)
	require.NoError(t, err)

	output, ok := result.(Output)
	require.True(t, ok)
	assert.Equal(t, len(javaBinaryFixture), output.Size)

	// dump_key always stores as base64 string — verify round-trip
	decoded, err := base64.StdEncoding.DecodeString(output.Serialized)
	require.NoError(t, err)
	assert.Equal(t, javaBinaryFixture, decoded, "decoded bytes must match original binary data exactly")

	// Verify it serializes cleanly to JSON (no encoding issues)
	jsonBytes, err := json.Marshal(output)
	require.NoError(t, err)
	assert.Contains(t, string(jsonBytes), `"serialized":"`)
}
