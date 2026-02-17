package base

import (
	"unicode/utf8"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
)

// SafeValue converts raw bytes to a JSON-safe value.
// If b is valid UTF-8 it returns a string, which json.Marshal encodes as a JSON string.
// If b is not valid UTF-8 it returns []byte, which json.Marshal encodes as a base64 string.
// This ensures binary data is never corrupted during JSON serialisation.
func SafeValue(b []byte) any {
	if len(b) == 0 {
		return ""
	}
	if utf8.Valid(b) {
		return string(b)
	}
	return b
}

// SafeSlice converts a slice of raw byte slices to a slice of JSON-safe values.
// Each element is independently checked for UTF-8 validity.
func SafeSlice(values [][]byte) []any {
	result := make([]any, len(values))
	for i, v := range values {
		result[i] = SafeValue(v)
	}
	return result
}

// SafeMap converts a map of field name to raw bytes to a map of field name to JSON-safe value.
// Each value is independently checked for UTF-8 validity.
// Field names (hash keys) are always treated as strings — Valkey key names are
// user-controlled and expected to be valid text.
func SafeMap(fields map[string][]byte) map[string]any {
	result := make(map[string]any, len(fields))
	for k, v := range fields {
		result[k] = SafeValue(v)
	}
	return result
}

// SafeStreamEntries converts a slice of StreamEntry (with []byte field values) to a
// slice of maps with JSON-safe values. The stream entry ID is always a string.
func SafeStreamEntries(entries []client.StreamEntry) []map[string]any {
	result := make([]map[string]any, len(entries))
	for i, entry := range entries {
		m := make(map[string]any, len(entry.FieldValues)+1)
		m["_id"] = entry.ID
		for k, v := range entry.FieldValues {
			m[k] = SafeValue(v)
		}
		result[i] = m
	}
	return result
}
