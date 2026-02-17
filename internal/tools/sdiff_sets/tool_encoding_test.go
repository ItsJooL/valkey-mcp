package sdiff_sets

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

func TestSdiffSetsTool_BinaryMember_IsBase64InJSON(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	// set1 has binary + "shared"; set2 has "shared" — diff = only binary member
	mockClient.AddSet(ctx, "set1", []string{string(javaBinaryFixture), "shared"})
	mockClient.AddSet(ctx, "set2", []string{"shared"})

	inputJSON, _ := json.Marshal(map[string]interface{}{"keys": []string{"set1", "set2"}})
	result, err := tool.Execute(ctx, inputJSON)
	require.NoError(t, err)

	output, ok := result.(Output)
	require.True(t, ok)
	assert.Equal(t, int64(1), output.Count)

	_, isBytes := output.Members[0].([]byte)
	assert.True(t, isBytes, "binary member should be []byte, not string")

	jsonBytes, err := json.Marshal(output)
	require.NoError(t, err)

	var jsonResult map[string]interface{}
	require.NoError(t, json.Unmarshal(jsonBytes, &jsonResult))

	members, ok := jsonResult["members"].([]interface{})
	require.True(t, ok)
	require.Len(t, members, 1)

	encodedStr, ok := members[0].(string)
	require.True(t, ok, "JSON binary member should be a string (base64)")

	decoded, err := base64.StdEncoding.DecodeString(encodedStr)
	require.NoError(t, err)
	assert.Equal(t, javaBinaryFixture, decoded)
}
