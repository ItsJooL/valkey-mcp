package client

import (
	"context"
)

// ValkeyClient defines the interface for Valkey operations.
// Tools depend on this interface, not concrete implementations.
// This enables easy testing with mocks and clean separation of concerns.
type ValkeyClient interface {
	// Server operations
	Ping(ctx context.Context) error
	GetServerInfo(ctx context.Context) (map[string]string, error)

	// String operations
	GetString(ctx context.Context, key string) ([]byte, bool, error)
	SetString(ctx context.Context, key, value string, ttlSeconds *int64, nx, xx bool) (bool, error)
	DeleteKey(ctx context.Context, key string) (bool, error)
	ExistsKeys(ctx context.Context, keys []string) (map[string]bool, error)
	ExpireKey(ctx context.Context, key string, seconds int64) (bool, error)
	PersistKey(ctx context.Context, key string) (bool, error)
	RenameKey(ctx context.Context, oldKey, newKey string) (bool, error)
	GetTTL(ctx context.Context, key string) (int64, error)
	IncrementNumber(ctx context.Context, key string, amount int64) (int64, error)
	DecrementNumber(ctx context.Context, key string, amount int64) (int64, error)
	AppendString(ctx context.Context, key, value string) (int64, error)
	GetRange(ctx context.Context, key string, start, end int64) ([]byte, error)

	// Hash (map) operations
	GetMap(ctx context.Context, key string) (map[string][]byte, error)
	SetMap(ctx context.Context, key string, fields map[string]string) (int64, error)
	GetMapField(ctx context.Context, key, field string) ([]byte, bool, error)
	GetMapFields(ctx context.Context, key string, fields []string) (map[string][]byte, error)
	DeleteMapFields(ctx context.Context, key string, fields []string) (int64, error)
	ListMapKeys(ctx context.Context, key string) ([]string, error)
	MapFieldExists(ctx context.Context, key, field string) (bool, error)
	IncrementMapField(ctx context.Context, key, field string, amount int64) (int64, error)

	// List operations
	PushList(ctx context.Context, key string, values []string, tail bool) (int64, error)
	PopList(ctx context.Context, key string, count int64, tail bool) ([][]byte, error)
	GetListRange(ctx context.Context, key string, start, stop int64) ([][]byte, error)
	GetListLength(ctx context.Context, key string) (int64, error)
	GetListIndex(ctx context.Context, key string, index int64) ([]byte, bool, error)
	SetListIndex(ctx context.Context, key string, index int64, value string) (bool, error)
	TrimList(ctx context.Context, key string, start, stop int64) (bool, error)

	// Set operations
	AddSet(ctx context.Context, key string, members []string) (int64, error)
	RemoveSet(ctx context.Context, key string, members []string) (int64, error)
	ListSetMembers(ctx context.Context, key string) ([][]byte, error)
	CheckSetMember(ctx context.Context, key, member string) (bool, error)
	GetSetSize(ctx context.Context, key string) (int64, error)
	PopSet(ctx context.Context, key string, count int64) ([][]byte, error)
	GetRandomSetMember(ctx context.Context, key string, count int64) ([][]byte, error)

	// Additional Key operations
	KeysByPattern(ctx context.Context, pattern string) ([]string, error)
	ExistsKey(ctx context.Context, key string) (bool, error)
	MemoryUsage(ctx context.Context, key string) (int64, error)
	TouchKeys(ctx context.Context, keys []string) (int64, error)
	ObjectEncoding(ctx context.Context, key string) (string, error)

	// Additional Hash operations
	GetMapLength(ctx context.Context, key string) (int64, error)
	ListMapFieldNames(ctx context.Context, key string) ([]string, error)
	ListMapFieldValues(ctx context.Context, key string) ([][]byte, error)
	GetMapFieldsMultiple(ctx context.Context, key string, fields []string) (map[string][]byte, error)

	// Additional Set operations
	SetIntersection(ctx context.Context, keys []string) ([][]byte, error)
	SetUnion(ctx context.Context, keys []string) ([][]byte, error)
	SetDifference(ctx context.Context, firstKey string, otherKeys []string) ([][]byte, error)

	// Stream operations
	AddStream(ctx context.Context, key string, id string, fields map[string]string) (string, error)
	GetStreamRange(ctx context.Context, key string, start string, end string, count int64) ([]StreamEntry, error)
	GetStreamLength(ctx context.Context, key string) (int64, error)
	ReadStream(ctx context.Context, key string, id string, count int64) ([]StreamEntry, error)

	// Serialization operations
	DumpKey(ctx context.Context, key string) ([]byte, error)
	RestoreKey(ctx context.Context, key string, ttl int64, serialized []byte) (bool, error)

	// Key object info
	ObjectIdletime(ctx context.Context, key string) (int64, error)

	// Configuration operations
	ConfigGet(ctx context.Context, parameter string) (map[string]string, error)
	ConfigSet(ctx context.Context, parameter, value string) (bool, error)

	// Database operations
	GetDatabaseSize(ctx context.Context) (int64, error)
	GetSlowlog(ctx context.Context, count int64) ([]map[string]interface{}, error)

	// Cluster operations
	GetClusterInfo(ctx context.Context) (map[string]string, error)
	GetClusterNodes(ctx context.Context) (string, error)
	GetKeySlot(ctx context.Context, key string) (int64, error)
	CountKeysInSlot(ctx context.Context, slot int64) (int64, error)

	// Scripting operations
	EvalScript(ctx context.Context, script string, keys []string, args []string) (interface{}, error)
	LoadScript(ctx context.Context, script string) (string, error)
	EvalSHA(ctx context.Context, sha string, keys []string, args []string) (interface{}, error)
}
