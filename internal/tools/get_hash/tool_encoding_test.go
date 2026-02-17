package get_hash

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

func TestGetHashTool_MixedFields_TextAndBinary(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	mockClient.SetRawHashBytes("user:123", map[string][]byte{
		"name": []byte("John Doe"),
		"pojo": javaBinaryFixture,
	})

	inputJSON, _ := json.Marshal(map[string]interface{}{"key": "user:123"})
	result, err := tool.Execute(ctx, inputJSON)
	require.NoError(t, err)

	output, ok := result.(Output)
	require.True(t, ok)
	assert.Equal(t, 2, output.FieldCount)
	assert.True(t, output.Exists)

	nameVal, isString := output.Fields["name"].(string)
	assert.True(t, isString, "text field should be string")
	assert.Equal(t, "John Doe", nameVal)

	_, isBytes := output.Fields["pojo"].([]byte)
	assert.True(t, isBytes, "binary field should be []byte")

	jsonBytes, err := json.Marshal(output)
	require.NoError(t, err)

	var jsonResult map[string]interface{}
	require.NoError(t, json.Unmarshal(jsonBytes, &jsonResult))

	fields, ok := jsonResult["fields"].(map[string]interface{})
	require.True(t, ok)

	encodedPojo, ok := fields["pojo"].(string)
	require.True(t, ok, "binary field should be base64 in JSON")

	decoded, err := base64.StdEncoding.DecodeString(encodedPojo)
	require.NoError(t, err)
	assert.Equal(t, javaBinaryFixture, decoded)
}
