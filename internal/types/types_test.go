package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewValkeyURL_Valid(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{"simple", "valkey://localhost:6379"},
		{"with auth", "valkey://user:pass@localhost:6379"},
		{"redis scheme", "redis://localhost:6379"},
		{"with db", "valkey://localhost:6379/0"},
		{"with host only", "valkey://localhost"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := NewValkeyURL(tt.url)
			assert.NoError(t, err)
			assert.Equal(t, ValkeyURL(tt.url), url)
		})
	}
}

func TestNewValkeyURL_Invalid(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{"empty", ""},
		{"no scheme", "localhost:6379"},
		{"wrong scheme", "http://localhost:6379"},
		{"invalid format", "not a url"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewValkeyURL(tt.url)
			assert.Error(t, err)
		})
	}
}

func TestValkeyURL_String(t *testing.T) {
	url, _ := NewValkeyURL("valkey://localhost:6379")
	assert.Equal(t, "valkey://localhost:6379", url.String())
}

func TestNewDBIndex_Valid(t *testing.T) {
	tests := []struct {
		name  string
		index int
	}{
		{"zero", 0},
		{"one", 1},
		{"mid", 7},
		{"max", 15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := NewDBIndex(tt.index)
			assert.NoError(t, err)
			assert.Equal(t, tt.index, db.Int())
		})
	}
}

func TestNewDBIndex_Invalid(t *testing.T) {
	tests := []struct {
		name  string
		index int
	}{
		{"negative", -1},
		{"too large", 16},
		{"way too large", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewDBIndex(tt.index)
			assert.Error(t, err)
		})
	}
}

func TestDBIndex_Int(t *testing.T) {
	db, _ := NewDBIndex(5)
	assert.Equal(t, 5, db.Int())
}

func TestNewTTLSeconds_Valid(t *testing.T) {
	tests := []struct {
		name    string
		seconds int64
	}{
		{"zero", 0},
		{"one second", 1},
		{"one minute", 60},
		{"one hour", 3600},
		{"one day", 86400},
		{"large", 999999999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ttl, err := NewTTLSeconds(tt.seconds)
			assert.NoError(t, err)
			assert.Equal(t, tt.seconds, ttl.Int64())
		})
	}
}

func TestNewTTLSeconds_Invalid(t *testing.T) {
	tests := []struct {
		name    string
		seconds int64
	}{
		{"negative", -1},
		{"very negative", -999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTTLSeconds(tt.seconds)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "non-negative")
		})
	}
}

func TestTTLSeconds_Int64(t *testing.T) {
	ttl, _ := NewTTLSeconds(3600)
	assert.Equal(t, int64(3600), ttl.Int64())
}

func TestNewKeyPattern_Valid(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
	}{
		{"wildcard", "user:*"},
		{"specific", "user:123"},
		{"question mark", "user:?"},
		{"complex", "cache:*:data"},
		{"all", "*"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern, err := NewKeyPattern(tt.pattern)
			assert.NoError(t, err)
			assert.Equal(t, KeyPattern(tt.pattern), pattern)
		})
	}
}

func TestNewKeyPattern_Invalid(t *testing.T) {
	_, err := NewKeyPattern("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")
}

func TestKeyPattern_String(t *testing.T) {
	pattern, _ := NewKeyPattern("user:*")
	assert.Equal(t, "user:*", pattern.String())
}

func TestNewScore(t *testing.T) {
	score := NewScore(3.14)
	assert.Equal(t, 3.14, score.Float64())
	assert.Equal(t, "3.14", score.String())
}

func TestNewScore_Edge(t *testing.T) {
	tests := []struct {
		name  string
		value float64
	}{
		{"zero", 0.0},
		{"negative", -123.45},
		{"large", 999999.99},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := NewScore(tt.value)
			assert.Equal(t, tt.value, score.Float64())
		})
	}
}

func TestValkeyURL_Schemes(t *testing.T) {
	schemes := []string{"valkey", "redis", "valkeys", "rediss"}
	for _, scheme := range schemes {
		t.Run(scheme, func(t *testing.T) {
			url, err := NewValkeyURL(scheme + "://localhost:6379")
			assert.NoError(t, err)
			assert.Contains(t, url.String(), scheme)
		})
	}
}

func TestDBIndex_Range(t *testing.T) {
	// Test boundary values
	db0, err := NewDBIndex(0)
	assert.NoError(t, err)
	assert.Equal(t, 0, db0.Int())

	db15, err := NewDBIndex(15)
	assert.NoError(t, err)
	assert.Equal(t, 15, db15.Int())
}
