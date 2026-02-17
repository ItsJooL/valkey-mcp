package base

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSafeValue_ValidUTF8ReturnsString(t *testing.T) {
	result := SafeValue([]byte("hello world"))
	assert.Equal(t, "hello world", result)
	_, isString := result.(string)
	assert.True(t, isString, "valid UTF-8 should return string type")
}

func TestSafeValue_BinaryReturnsBytes(t *testing.T) {
	binary := []byte{0xAC, 0xED, 0x00, 0x05}
	result := SafeValue(binary)
	_, isBytes := result.([]byte)
	assert.True(t, isBytes, "binary data should return []byte type")
}

func TestSafeValue_EmptyBytesReturnsEmptyString(t *testing.T) {
	result := SafeValue([]byte{})
	assert.Equal(t, "", result)
}

func TestSafeValue_NilBytesReturnsEmptyString(t *testing.T) {
	result := SafeValue(nil)
	assert.Equal(t, "", result)
}

func TestSafeValue_JSONMarshalString(t *testing.T) {
	type Wrapper struct{ Value any }
	w := Wrapper{Value: SafeValue([]byte("hello"))}
	b, err := json.Marshal(w)
	require.NoError(t, err)
	assert.JSONEq(t, `{"Value":"hello"}`, string(b))
}

func TestSafeValue_JSONMarshalBinaryIsBase64(t *testing.T) {
	type Wrapper struct{ Value any }
	binary := []byte{0xAC, 0xED, 0x00, 0x05}
	w := Wrapper{Value: SafeValue(binary)}
	b, err := json.Marshal(w)
	require.NoError(t, err)

	// Unmarshal and verify round-trip
	var result map[string]string
	require.NoError(t, json.Unmarshal(b, &result))

	decoded, err := base64.StdEncoding.DecodeString(result["Value"])
	require.NoError(t, err)
	assert.Equal(t, binary, decoded)
}

func TestSafeSlice_MixedValues(t *testing.T) {
	input := [][]byte{
		[]byte("text value"),
		{0xAC, 0xED, 0x00, 0x05},
	}
	result := SafeSlice(input)
	assert.Len(t, result, 2)
	assert.Equal(t, "text value", result[0])
	_, isBytes := result[1].([]byte)
	assert.True(t, isBytes)
}

func TestSafeMap_MixedValues(t *testing.T) {
	input := map[string][]byte{
		"name":  []byte("John Doe"),
		"photo": {0xFF, 0xD8, 0xFF, 0xE0}, // JPEG magic bytes
	}
	result := SafeMap(input)
	assert.Equal(t, "John Doe", result["name"])
	_, isBytes := result["photo"].([]byte)
	assert.True(t, isBytes)
}

func TestSafeStreamEntries_MixedFieldValues(t *testing.T) {
	entries := []client.StreamEntry{
		{
			ID: "1-0",
			FieldValues: map[string][]byte{
				"name": []byte("hello"),
				"data": {0xAC, 0xED, 0x00, 0x05},
			},
		},
	}
	result := SafeStreamEntries(entries)
	require.Len(t, result, 1)
	assert.Equal(t, "1-0", result[0]["_id"])
	assert.Equal(t, "hello", result[0]["name"])
	_, isBytes := result[0]["data"].([]byte)
	assert.True(t, isBytes)
}
