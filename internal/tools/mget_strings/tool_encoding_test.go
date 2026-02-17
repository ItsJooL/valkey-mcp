package mget_strings

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

func TestMGetStringsTool_BinaryValue_IsBase64InJSON(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	mockClient.SetRawBytes("binary_key", javaBinaryFixture)
	mockClient.SetString(ctx, "text_key", "hello", nil, false, false)

	inputJSON, _ := json.Marshal(map[string]interface{}{"keys": []string{"binary_key", "text_key"}})
	result, err := tool.Execute(ctx, inputJSON)
	require.NoError(t, err)

	output, ok := result.(Output)
	require.True(t, ok)
	assert.Equal(t, 2, output.Count)

	_, isBytes := output.Values["binary_key"].([]byte)
	assert.True(t, isBytes, "binary value should be []byte, not string")

	strVal, isString := output.Values["text_key"].(string)
	assert.True(t, isString, "text value should be string, not []byte")
	assert.Equal(t, "hello", strVal)

	jsonBytes, err := json.Marshal(output)
	require.NoError(t, err)

	var jsonResult map[string]interface{}
	require.NoError(t, json.Unmarshal(jsonBytes, &jsonResult))

	values, ok := jsonResult["values"].(map[string]interface{})
	require.True(t, ok)

	encodedStr, ok := values["binary_key"].(string)
	require.True(t, ok, "JSON binary value should be a string (base64)")

	decoded, err := base64.StdEncoding.DecodeString(encodedStr)
	require.NoError(t, err)
	assert.Equal(t, javaBinaryFixture, decoded, "decoded bytes must match original binary data exactly")
}
