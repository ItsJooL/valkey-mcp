package client

// StreamEntry represents a single Valkey stream entry with binary-safe field values.
// The ID is always a plain ASCII string (Valkey stream IDs are timestamp-sequence pairs).
// FieldValues contains raw bytes to preserve binary data without UTF-8 corruption.
type StreamEntry struct {
	ID          string
	FieldValues map[string][]byte
}
