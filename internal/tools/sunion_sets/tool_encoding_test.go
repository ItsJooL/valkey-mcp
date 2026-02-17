package sunion_sets

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

func TestSunionSetsTool_BinaryMember_IsBase64InJSON(t *testing.T) {
	mockClient := client.NewMockClient()
	tool := NewTool(mockClient)
	ctx := context.Background()

	// Binary member in first set, plain text in second — union has both
	mockClient.AddSet(ctx, "set1", []string{string(javaBinaryFixture)})
	mockClient.AddSet(ctx, "set2", []string{"plain_member"})

	inputJSON, _ := json.Marshal(map[string]interface{}{"keys": []string{"set1", "set2"}})
	result, err := tool.Execute(ctx, inputJSON)
	require.NoError(t, err)

	output, ok := result.(Output)
	require.True(t, ok)
	assert.Equal(t, int64(2), output.Count)

	// Find the binary member in the union result
	var foundBinary bool
	var encodedBinary string
	jsonBytes, err := json.Marshal(output)
	require.NoError(t, err)

	var jsonResult map[string]interface{}
	require.NoError(t, json.Unmarshal(jsonBytes, &jsonResult))

	members, ok := jsonResult["members"].([]interface{})
	require.True(t, ok)

	for _, m := range members {
		s, ok := m.(string)
		if !ok {
			continue
		}
		decoded, err := base64.StdEncoding.DecodeString(s)
		if err == nil && string(decoded) != "plain_member" {
			foundBinary = true
			encodedBinary = s
		}
	}

	assert.True(t, foundBinary, "union result should contain the binary member as base64")

	decoded, err := base64.StdEncoding.DecodeString(encodedBinary)
	require.NoError(t, err)
	assert.Equal(t, javaBinaryFixture, decoded)
}
