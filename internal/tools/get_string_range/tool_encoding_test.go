package get_string_range

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

func TestGetStringRangeTool_BinaryValue_IsBase64InJSON(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	mockClient.SetRawBytes("binary_key", javaBinaryFixture)

	// Request full range: start=0, end=-1 (all bytes)
	inputJSON, _ := json.Marshal(map[string]interface{}{"key": "binary_key", "start": 0, "end": -1})
	result, err := tool.Execute(ctx, inputJSON)
	require.NoError(t, err)

	output, ok := result.(Output)
	require.True(t, ok)

	_, isBytes := output.Value.([]byte)
	assert.True(t, isBytes, "binary value should be []byte, not string")

	jsonBytes, err := json.Marshal(output)
	require.NoError(t, err)

	var jsonResult map[string]interface{}
	require.NoError(t, json.Unmarshal(jsonBytes, &jsonResult))

	encodedStr, ok := jsonResult["value"].(string)
	require.True(t, ok, "JSON value should be a string (base64)")

	decoded, err := base64.StdEncoding.DecodeString(encodedStr)
	require.NoError(t, err)
	assert.Equal(t, javaBinaryFixture, decoded, "decoded bytes must match original binary data exactly")
}

func TestGetStringRangeTool_TextValue_IsPlainStringInJSON(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	mockClient.SetString(ctx, "text_key", "Hello World", nil, false, false)

	inputJSON, _ := json.Marshal(map[string]interface{}{"key": "text_key", "start": 0, "end": 4})
	result, err := tool.Execute(ctx, inputJSON)
	require.NoError(t, err)

	output, ok := result.(Output)
	require.True(t, ok)

	_, isString := output.Value.(string)
	assert.True(t, isString, "text value should be string, not []byte")
}
