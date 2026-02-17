package xrange_stream

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

func TestXRangeStreamTool_BinaryField_IsBase64InJSON(t *testing.T) {
	mockClient := &client.MockValkeyClient{}
	mockClient.GetStreamRangeFunc = func(ctx context.Context, key, start, end string, count int64) ([]client.StreamEntry, error) {
		return []client.StreamEntry{
			{
				ID: "1-0",
				FieldValues: map[string][]byte{
					"name": []byte("Alice"),
					"pojo": javaBinaryFixture,
				},
			},
		}, nil
	}

	tool := NewTool(mockClient)
	ctx := context.Background()

	inputJSON, _ := json.Marshal(map[string]interface{}{"key": "mystream", "start": "-", "end": "+"})
	result, err := tool.Execute(ctx, inputJSON)
	require.NoError(t, err)

	output, ok := result.(Output)
	require.True(t, ok)
	assert.Equal(t, int64(1), output.Count)

	entry := output.Entries[0]

	_, isString := entry["name"].(string)
	assert.True(t, isString, "text field should be string")

	_, isBytes := entry["pojo"].([]byte)
	assert.True(t, isBytes, "binary field should be []byte")

	jsonBytes, err := json.Marshal(output)
	require.NoError(t, err)

	var jsonResult map[string]interface{}
	require.NoError(t, json.Unmarshal(jsonBytes, &jsonResult))

	entries, ok := jsonResult["entries"].([]interface{})
	require.True(t, ok)
	require.Len(t, entries, 1)

	entryMap, ok := entries[0].(map[string]interface{})
	require.True(t, ok)

	encodedPojo, ok := entryMap["pojo"].(string)
	require.True(t, ok, "binary field should be base64 in JSON")

	decoded, err := base64.StdEncoding.DecodeString(encodedPojo)
	require.NoError(t, err)
	assert.Equal(t, javaBinaryFixture, decoded)
}
