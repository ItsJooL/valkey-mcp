package lrange_list

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

func TestLRangeListTool_BinaryElement_IsBase64InJSON(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	mockClient.SetRawListBytes("mylist", [][]byte{
		javaBinaryFixture,
		[]byte("plain text"),
	})

	inputJSON, _ := json.Marshal(map[string]interface{}{"key": "mylist", "start": 0, "stop": -1})
	result, err := tool.Execute(ctx, inputJSON)
	require.NoError(t, err)

	output, ok := result.(Output)
	require.True(t, ok)
	assert.Equal(t, 2, output.Count)

	_, isBytes := output.Values[0].([]byte)
	assert.True(t, isBytes, "binary element should be []byte, not string")

	_, isString := output.Values[1].(string)
	assert.True(t, isString, "text element should be string, not []byte")

	jsonBytes, err := json.Marshal(output)
	require.NoError(t, err)

	var jsonResult map[string]interface{}
	require.NoError(t, json.Unmarshal(jsonBytes, &jsonResult))

	values, ok := jsonResult["values"].([]interface{})
	require.True(t, ok)
	require.Len(t, values, 2)

	encodedStr, ok := values[0].(string)
	require.True(t, ok, "JSON binary element should be a string (base64)")

	decoded, err := base64.StdEncoding.DecodeString(encodedStr)
	require.NoError(t, err)
	assert.Equal(t, javaBinaryFixture, decoded)
}
