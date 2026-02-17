package hvals_hash

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

func TestHValsHashTool_BinaryValue_IsBase64InJSON(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	mockClient.SetRawHashBytes("myhash", map[string][]byte{
		"field1": javaBinaryFixture,
	})

	inputJSON, _ := json.Marshal(map[string]interface{}{"key": "myhash"})
	result, err := tool.Execute(ctx, inputJSON)
	require.NoError(t, err)

	output, ok := result.(Output)
	require.True(t, ok)
	require.Len(t, output.Result, 1)

	_, isBytes := output.Result[0].([]byte)
	assert.True(t, isBytes, "binary value should be []byte, not string")

	jsonBytes, err := json.Marshal(output)
	require.NoError(t, err)

	var jsonResult map[string]interface{}
	require.NoError(t, json.Unmarshal(jsonBytes, &jsonResult))

	resultArr, ok := jsonResult["result"].([]interface{})
	require.True(t, ok)
	require.Len(t, resultArr, 1)

	encodedStr, ok := resultArr[0].(string)
	require.True(t, ok, "JSON value should be a string (base64)")

	decoded, err := base64.StdEncoding.DecodeString(encodedStr)
	require.NoError(t, err)
	assert.Equal(t, javaBinaryFixture, decoded)
}
