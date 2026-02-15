package types

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// ValkeyURL represents a validated Valkey connection URL.
type ValkeyURL string

// NewValkeyURL creates and validates a Valkey URL.
func NewValkeyURL(rawURL string) (ValkeyURL, error) {
	if rawURL == "" {
		return "", fmt.Errorf("valkey URL cannot be empty")
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid valkey URL: %w", err)
	}

	scheme := strings.ToLower(u.Scheme)
	if scheme != "valkey" && scheme != "redis" && scheme != "valkeys" && scheme != "rediss" {
		return "", fmt.Errorf("invalid scheme %q: must be valkey, redis, valkeys, or rediss", u.Scheme)
	}

	if u.Host == "" {
		return "", fmt.Errorf("valkey URL must specify a host")
	}

	return ValkeyURL(rawURL), nil
}

// String returns the string representation of the URL.
func (v ValkeyURL) String() string {
	return string(v)
}

// KeyPattern represents a validated key pattern for scanning operations.
type KeyPattern string

// NewKeyPattern creates and validates a key pattern.
func NewKeyPattern(pattern string) (KeyPattern, error) {
	if pattern == "" {
		return "", fmt.Errorf("key pattern cannot be empty")
	}
	return KeyPattern(pattern), nil
}

// String returns the string representation of the pattern.
func (k KeyPattern) String() string {
	return string(k)
}

// DBIndex represents a validated database index.
type DBIndex int

// NewDBIndex creates and validates a database index.
func NewDBIndex(index int) (DBIndex, error) {
	if index < 0 {
		return 0, fmt.Errorf("database index must be non-negative, got %d", index)
	}
	if index > 15 {
		return 0, fmt.Errorf("database index must be <= 15, got %d", index)
	}
	return DBIndex(index), nil
}

// Int returns the integer value of the database index.
func (d DBIndex) Int() int {
	return int(d)
}

// TTLSeconds represents a time-to-live in seconds.
type TTLSeconds int64

// NewTTLSeconds creates and validates a TTL value.
func NewTTLSeconds(seconds int64) (TTLSeconds, error) {
	if seconds < 0 {
		return 0, fmt.Errorf("TTL must be non-negative, got %d", seconds)
	}
	return TTLSeconds(seconds), nil
}

// Int64 returns the int64 value of the TTL.
func (t TTLSeconds) Int64() int64 {
	return int64(t)
}

// Score represents a sorted set score.
type Score float64

// NewScore creates a score value.
func NewScore(score float64) Score {
	return Score(score)
}

// Float64 returns the float64 value of the score.
func (s Score) Float64() float64 {
	return float64(s)
}

// String returns the string representation of the score.
func (s Score) String() string {
	return strconv.FormatFloat(float64(s), 'f', -1, 64)
}
