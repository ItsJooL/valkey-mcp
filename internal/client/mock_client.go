package client

import (
	"context"
	"fmt"
	"sync"
)

// MockValkeyClient implements ValkeyClient interface for testing via function stubs.
type MockValkeyClient struct {
	// Server operations
	PingFunc          func(ctx context.Context) error
	GetServerInfoFunc func(ctx context.Context) (map[string]string, error)

	// String operations
	GetStringFunc       func(ctx context.Context, key string) ([]byte, bool, error)
	SetStringFunc       func(ctx context.Context, key, value string, ttlSeconds *int64, nx, xx bool) (bool, error)
	DeleteKeyFunc       func(ctx context.Context, key string) (bool, error)
	ExistsKeysFunc      func(ctx context.Context, keys []string) (map[string]bool, error)
	ExpireKeyFunc       func(ctx context.Context, key string, seconds int64) (bool, error)
	PersistKeyFunc      func(ctx context.Context, key string) (bool, error)
	RenameKeyFunc       func(ctx context.Context, oldKey, newKey string) (bool, error)
	GetTTLFunc          func(ctx context.Context, key string) (int64, error)
	IncrementNumberFunc func(ctx context.Context, key string, amount int64) (int64, error)
	DecrementNumberFunc func(ctx context.Context, key string, amount int64) (int64, error)
	AppendStringFunc    func(ctx context.Context, key, value string) (int64, error)
	GetRangeFunc        func(ctx context.Context, key string, start, end int64) ([]byte, error)

	// Hash operations
	GetMapFunc              func(ctx context.Context, key string) (map[string][]byte, error)
	SetMapFunc              func(ctx context.Context, key string, fields map[string]string) (int64, error)
	GetMapFieldFunc         func(ctx context.Context, key, field string) ([]byte, bool, error)
	GetMapFieldsFunc        func(ctx context.Context, key string, fields []string) (map[string][]byte, error)
	GetMapFieldsMultipleFunc func(ctx context.Context, key string, fields []string) (map[string][]byte, error)
	DeleteMapFieldsFunc     func(ctx context.Context, key string, fields []string) (int64, error)
	ListMapKeysFunc         func(ctx context.Context, key string) ([]string, error)
	ListMapFieldValuesFunc  func(ctx context.Context, key string) ([][]byte, error)

	// List operations
	PushListFunc      func(ctx context.Context, key string, values []string, tail bool) (int64, error)
	PopListFunc       func(ctx context.Context, key string, count int64, tail bool) ([][]byte, error)
	GetListRangeFunc  func(ctx context.Context, key string, start, stop int64) ([][]byte, error)
	GetListLengthFunc func(ctx context.Context, key string) (int64, error)
	GetListIndexFunc  func(ctx context.Context, key string, index int64) ([]byte, bool, error)
	SetListIndexFunc  func(ctx context.Context, key string, index int64, value string) (bool, error)
	TrimListFunc      func(ctx context.Context, key string, start, stop int64) (bool, error)

	// Set operations
	AddSetFunc            func(ctx context.Context, key string, members []string) (int64, error)
	RemoveSetFunc         func(ctx context.Context, key string, members []string) (int64, error)
	ListSetMembersFunc    func(ctx context.Context, key string) ([][]byte, error)
	CheckSetMemberFunc    func(ctx context.Context, key, member string) (bool, error)
	GetSetSizeFunc        func(ctx context.Context, key string) (int64, error)
	PopSetFunc            func(ctx context.Context, key string, count int64) ([][]byte, error)
	GetRandomSetMemberFunc func(ctx context.Context, key string, count int64) ([][]byte, error)
	SetIntersectionFunc   func(ctx context.Context, keys []string) ([][]byte, error)
	SetUnionFunc          func(ctx context.Context, keys []string) ([][]byte, error)
	SetDifferenceFunc     func(ctx context.Context, firstKey string, otherKeys []string) ([][]byte, error)

	// Stream operations
	AddStreamFunc       func(ctx context.Context, key string, id string, fields map[string]string) (string, error)
	GetStreamRangeFunc  func(ctx context.Context, key string, start string, end string, count int64) ([]StreamEntry, error)
	GetStreamLengthFunc func(ctx context.Context, key string) (int64, error)
	ReadStreamFunc      func(ctx context.Context, key string, id string, count int64) ([]StreamEntry, error)

	// Serialization operations
	DumpKeyFunc    func(ctx context.Context, key string) ([]byte, error)
	RestoreKeyFunc func(ctx context.Context, key string, ttl int64, serialized []byte) (bool, error)
}

// Server operations

func (m *MockValkeyClient) Ping(ctx context.Context) error {
	if m.PingFunc != nil {
		return m.PingFunc(ctx)
	}
	return nil
}

func (m *MockValkeyClient) GetServerInfo(ctx context.Context) (map[string]string, error) {
	if m.GetServerInfoFunc != nil {
		return m.GetServerInfoFunc(ctx)
	}
	return make(map[string]string), nil
}

// String operations

func (m *MockValkeyClient) GetString(ctx context.Context, key string) ([]byte, bool, error) {
	if m.GetStringFunc != nil {
		return m.GetStringFunc(ctx, key)
	}
	return nil, false, nil
}

func (m *MockValkeyClient) SetString(ctx context.Context, key, value string, ttlSeconds *int64, nx, xx bool) (bool, error) {
	if m.SetStringFunc != nil {
		return m.SetStringFunc(ctx, key, value, ttlSeconds, nx, xx)
	}
	return false, nil
}

func (m *MockValkeyClient) DeleteKey(ctx context.Context, key string) (bool, error) {
	if m.DeleteKeyFunc != nil {
		return m.DeleteKeyFunc(ctx, key)
	}
	return false, nil
}

func (m *MockValkeyClient) ExistsKeys(ctx context.Context, keys []string) (map[string]bool, error) {
	if m.ExistsKeysFunc != nil {
		return m.ExistsKeysFunc(ctx, keys)
	}
	return map[string]bool{}, nil
}

func (m *MockValkeyClient) ExpireKey(ctx context.Context, key string, seconds int64) (bool, error) {
	if m.ExpireKeyFunc != nil {
		return m.ExpireKeyFunc(ctx, key, seconds)
	}
	return false, nil
}

func (m *MockValkeyClient) PersistKey(ctx context.Context, key string) (bool, error) {
	if m.PersistKeyFunc != nil {
		return m.PersistKeyFunc(ctx, key)
	}
	return false, nil
}

func (m *MockValkeyClient) RenameKey(ctx context.Context, oldKey, newKey string) (bool, error) {
	if m.RenameKeyFunc != nil {
		return m.RenameKeyFunc(ctx, oldKey, newKey)
	}
	return false, nil
}

func (m *MockValkeyClient) GetTTL(ctx context.Context, key string) (int64, error) {
	if m.GetTTLFunc != nil {
		return m.GetTTLFunc(ctx, key)
	}
	return -1, nil
}

func (m *MockValkeyClient) IncrementNumber(ctx context.Context, key string, amount int64) (int64, error) {
	if m.IncrementNumberFunc != nil {
		return m.IncrementNumberFunc(ctx, key, amount)
	}
	return 0, nil
}

func (m *MockValkeyClient) DecrementNumber(ctx context.Context, key string, amount int64) (int64, error) {
	if m.DecrementNumberFunc != nil {
		return m.DecrementNumberFunc(ctx, key, amount)
	}
	return 0, nil
}

func (m *MockValkeyClient) AppendString(ctx context.Context, key, value string) (int64, error) {
	if m.AppendStringFunc != nil {
		return m.AppendStringFunc(ctx, key, value)
	}
	return 0, nil
}

func (m *MockValkeyClient) GetRange(ctx context.Context, key string, start, end int64) ([]byte, error) {
	if m.GetRangeFunc != nil {
		return m.GetRangeFunc(ctx, key, start, end)
	}
	return nil, nil
}

// Hash operations

func (m *MockValkeyClient) GetMap(ctx context.Context, key string) (map[string][]byte, error) {
	if m.GetMapFunc != nil {
		return m.GetMapFunc(ctx, key)
	}
	return map[string][]byte{}, nil
}

func (m *MockValkeyClient) SetMap(ctx context.Context, key string, fields map[string]string) (int64, error) {
	if m.SetMapFunc != nil {
		return m.SetMapFunc(ctx, key, fields)
	}
	return 0, nil
}

func (m *MockValkeyClient) GetMapField(ctx context.Context, key, field string) ([]byte, bool, error) {
	if m.GetMapFieldFunc != nil {
		return m.GetMapFieldFunc(ctx, key, field)
	}
	return nil, false, nil
}

func (m *MockValkeyClient) GetMapFields(ctx context.Context, key string, fields []string) (map[string][]byte, error) {
	if m.GetMapFieldsFunc != nil {
		return m.GetMapFieldsFunc(ctx, key, fields)
	}
	return make(map[string][]byte), nil
}

func (m *MockValkeyClient) GetMapFieldsMultiple(ctx context.Context, key string, fields []string) (map[string][]byte, error) {
	if m.GetMapFieldsMultipleFunc != nil {
		return m.GetMapFieldsMultipleFunc(ctx, key, fields)
	}
	return map[string][]byte{}, nil
}

func (m *MockValkeyClient) DeleteMapFields(ctx context.Context, key string, fields []string) (int64, error) {
	if m.DeleteMapFieldsFunc != nil {
		return m.DeleteMapFieldsFunc(ctx, key, fields)
	}
	return 0, nil
}

func (m *MockValkeyClient) ListMapKeys(ctx context.Context, key string) ([]string, error) {
	if m.ListMapKeysFunc != nil {
		return m.ListMapKeysFunc(ctx, key)
	}
	return []string{}, nil
}

func (m *MockValkeyClient) ListMapFieldValues(ctx context.Context, key string) ([][]byte, error) {
	if m.ListMapFieldValuesFunc != nil {
		return m.ListMapFieldValuesFunc(ctx, key)
	}
	return [][]byte{}, nil
}

// List operations

func (m *MockValkeyClient) PushList(ctx context.Context, key string, values []string, tail bool) (int64, error) {
	if m.PushListFunc != nil {
		return m.PushListFunc(ctx, key, values, tail)
	}
	return 0, nil
}

func (m *MockValkeyClient) PopList(ctx context.Context, key string, count int64, tail bool) ([][]byte, error) {
	if m.PopListFunc != nil {
		return m.PopListFunc(ctx, key, count, tail)
	}
	return [][]byte{}, nil
}

func (m *MockValkeyClient) GetListRange(ctx context.Context, key string, start, stop int64) ([][]byte, error) {
	if m.GetListRangeFunc != nil {
		return m.GetListRangeFunc(ctx, key, start, stop)
	}
	return [][]byte{}, nil
}

func (m *MockValkeyClient) GetListLength(ctx context.Context, key string) (int64, error) {
	if m.GetListLengthFunc != nil {
		return m.GetListLengthFunc(ctx, key)
	}
	return 0, nil
}

func (m *MockValkeyClient) GetListIndex(ctx context.Context, key string, index int64) ([]byte, bool, error) {
	if m.GetListIndexFunc != nil {
		return m.GetListIndexFunc(ctx, key, index)
	}
	return nil, false, nil
}

func (m *MockValkeyClient) SetListIndex(ctx context.Context, key string, index int64, value string) (bool, error) {
	if m.SetListIndexFunc != nil {
		return m.SetListIndexFunc(ctx, key, index, value)
	}
	return false, nil
}

func (m *MockValkeyClient) TrimList(ctx context.Context, key string, start, stop int64) (bool, error) {
	if m.TrimListFunc != nil {
		return m.TrimListFunc(ctx, key, start, stop)
	}
	return false, nil
}

// Set operations

func (m *MockValkeyClient) AddSet(ctx context.Context, key string, members []string) (int64, error) {
	if m.AddSetFunc != nil {
		return m.AddSetFunc(ctx, key, members)
	}
	return 0, nil
}

func (m *MockValkeyClient) RemoveSet(ctx context.Context, key string, members []string) (int64, error) {
	if m.RemoveSetFunc != nil {
		return m.RemoveSetFunc(ctx, key, members)
	}
	return 0, nil
}

func (m *MockValkeyClient) ListSetMembers(ctx context.Context, key string) ([][]byte, error) {
	if m.ListSetMembersFunc != nil {
		return m.ListSetMembersFunc(ctx, key)
	}
	return [][]byte{}, nil
}

func (m *MockValkeyClient) CheckSetMember(ctx context.Context, key, member string) (bool, error) {
	if m.CheckSetMemberFunc != nil {
		return m.CheckSetMemberFunc(ctx, key, member)
	}
	return false, nil
}

func (m *MockValkeyClient) GetSetSize(ctx context.Context, key string) (int64, error) {
	if m.GetSetSizeFunc != nil {
		return m.GetSetSizeFunc(ctx, key)
	}
	return 0, nil
}

func (m *MockValkeyClient) PopSet(ctx context.Context, key string, count int64) ([][]byte, error) {
	if m.PopSetFunc != nil {
		return m.PopSetFunc(ctx, key, count)
	}
	return [][]byte{}, nil
}

func (m *MockValkeyClient) GetRandomSetMember(ctx context.Context, key string, count int64) ([][]byte, error) {
	if m.GetRandomSetMemberFunc != nil {
		return m.GetRandomSetMemberFunc(ctx, key, count)
	}
	return [][]byte{}, nil
}

func (m *MockValkeyClient) SetIntersection(ctx context.Context, keys []string) ([][]byte, error) {
	if m.SetIntersectionFunc != nil {
		return m.SetIntersectionFunc(ctx, keys)
	}
	return [][]byte{}, nil
}

func (m *MockValkeyClient) SetUnion(ctx context.Context, keys []string) ([][]byte, error) {
	if m.SetUnionFunc != nil {
		return m.SetUnionFunc(ctx, keys)
	}
	return [][]byte{}, nil
}

func (m *MockValkeyClient) SetDifference(ctx context.Context, firstKey string, otherKeys []string) ([][]byte, error) {
	if m.SetDifferenceFunc != nil {
		return m.SetDifferenceFunc(ctx, firstKey, otherKeys)
	}
	return [][]byte{}, nil
}

// Stream operations

func (m *MockValkeyClient) AddStream(ctx context.Context, key string, id string, fields map[string]string) (string, error) {
	if m.AddStreamFunc != nil {
		return m.AddStreamFunc(ctx, key, id, fields)
	}
	return id, nil
}

func (m *MockValkeyClient) GetStreamRange(ctx context.Context, key string, start string, end string, count int64) ([]StreamEntry, error) {
	if m.GetStreamRangeFunc != nil {
		return m.GetStreamRangeFunc(ctx, key, start, end, count)
	}
	return []StreamEntry{}, nil
}

func (m *MockValkeyClient) GetStreamLength(ctx context.Context, key string) (int64, error) {
	if m.GetStreamLengthFunc != nil {
		return m.GetStreamLengthFunc(ctx, key)
	}
	return 0, nil
}

func (m *MockValkeyClient) ReadStream(ctx context.Context, key string, id string, count int64) ([]StreamEntry, error) {
	if m.ReadStreamFunc != nil {
		return m.ReadStreamFunc(ctx, key, id, count)
	}
	return []StreamEntry{}, nil
}

// Serialization operations

func (m *MockValkeyClient) DumpKey(ctx context.Context, key string) ([]byte, error) {
	if m.DumpKeyFunc != nil {
		return m.DumpKeyFunc(ctx, key)
	}
	return []byte{}, nil
}

func (m *MockValkeyClient) RestoreKey(ctx context.Context, key string, ttl int64, serialized []byte) (bool, error) {
	if m.RestoreKeyFunc != nil {
		return m.RestoreKeyFunc(ctx, key, ttl, serialized)
	}
	return true, nil
}

func (m *MockValkeyClient) ConfigGet(ctx context.Context, parameter string) (map[string]string, error) {
	result := make(map[string]string)
	result[parameter] = "mock_value"
	return result, nil
}

func (m *MockValkeyClient) ConfigSet(ctx context.Context, parameter, value string) (bool, error) {
	return true, nil
}

func (m *MockValkeyClient) GetDatabaseSize(ctx context.Context) (int64, error) {
	return 0, nil
}

func (m *MockValkeyClient) GetSlowlog(ctx context.Context, count int64) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}

func (m *MockValkeyClient) GetClusterInfo(ctx context.Context) (map[string]string, error) {
	return map[string]string{"cluster_state": "ok"}, nil
}

func (m *MockValkeyClient) GetClusterNodes(ctx context.Context) (string, error) {
	return "mock_nodes", nil
}

func (m *MockValkeyClient) GetKeySlot(ctx context.Context, key string) (int64, error) {
	return 0, nil
}

func (m *MockValkeyClient) CountKeysInSlot(ctx context.Context, slot int64) (int64, error) {
	return 0, nil
}

func (m *MockValkeyClient) EvalScript(ctx context.Context, script string, keys []string, args []string) (interface{}, error) {
	return "OK", nil
}

func (m *MockValkeyClient) LoadScript(ctx context.Context, script string) (string, error) {
	return "abc123", nil
}

func (m *MockValkeyClient) EvalSHA(ctx context.Context, sha string, keys []string, args []string) (interface{}, error) {
	return "OK", nil
}

func (m *MockValkeyClient) KeysByPattern(ctx context.Context, pattern string) ([]string, error) {
	return []string{}, nil
}

func (m *MockValkeyClient) ExistsKey(ctx context.Context, key string) (bool, error) {
	return false, nil
}

func (m *MockValkeyClient) MemoryUsage(ctx context.Context, key string) (int64, error) {
	return 100, nil
}

func (m *MockValkeyClient) TouchKeys(ctx context.Context, keys []string) (int64, error) {
	return int64(len(keys)), nil
}

func (m *MockValkeyClient) ObjectEncoding(ctx context.Context, key string) (string, error) {
	return "raw", nil
}

func (m *MockValkeyClient) ObjectIdletime(ctx context.Context, key string) (int64, error) {
	return 0, nil
}

func (m *MockValkeyClient) GetMapLength(ctx context.Context, key string) (int64, error) {
	return 0, nil
}

func (m *MockValkeyClient) ListMapFieldNames(ctx context.Context, key string) ([]string, error) {
	return []string{}, nil
}

func (m *MockValkeyClient) MapFieldExists(ctx context.Context, key, field string) (bool, error) {
	return false, nil
}

func (m *MockValkeyClient) IncrementMapField(ctx context.Context, key, field string, amount int64) (int64, error) {
	return 0, nil
}

// Ensure MockValkeyClient implements ValkeyClient at compile time
var _ ValkeyClient = (*MockValkeyClient)(nil)

// ---------------------------------------------------------------------------
// MockClient — full in-memory implementation for unit tests
// ---------------------------------------------------------------------------

// MockClient is an in-memory mock implementation of ValkeyClient for testing.
type MockClient struct {
	mu sync.RWMutex

	// Storage — all value data held as raw bytes
	strings map[string][]byte
	hashes  map[string]map[string][]byte
	lists   map[string][][]byte
	sets    map[string]map[string]bool
	ttls    map[string]int64

	// Behavior controls
	PingError          error
	GetServerInfoError error
}

// NewMockClient creates a new mock client for testing.
func NewMockClient() *MockClient {
	return &MockClient{
		strings: make(map[string][]byte),
		hashes:  make(map[string]map[string][]byte),
		lists:   make(map[string][][]byte),
		sets:    make(map[string]map[string]bool),
		ttls:    make(map[string]int64),
	}
}

// SetRawBytes stores raw byte data for a string key — for testing binary retrieval paths.
func (m *MockClient) SetRawBytes(key string, value []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.strings[key] = value
}

// SetRawHashBytes stores raw byte hash field values for testing binary retrieval paths.
func (m *MockClient) SetRawHashBytes(key string, fields map[string][]byte) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.hashes[key] == nil {
		m.hashes[key] = make(map[string][]byte)
	}
	for k, v := range fields {
		m.hashes[key][k] = v
	}
}

// SetRawListBytes stores raw byte list elements for testing binary retrieval paths.
func (m *MockClient) SetRawListBytes(key string, values [][]byte) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.lists[key] = values
}

// Server operations

func (m *MockClient) Ping(ctx context.Context) error {
	if m.PingError != nil {
		return m.PingError
	}
	return nil
}

func (m *MockClient) GetServerInfo(ctx context.Context) (map[string]string, error) {
	if m.GetServerInfoError != nil {
		return nil, m.GetServerInfoError
	}
	return map[string]string{
		"version": "7.0.0",
		"mode":    "standalone",
	}, nil
}

// String operations

func (m *MockClient) GetString(ctx context.Context, key string) ([]byte, bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	val, exists := m.strings[key]
	if !exists {
		return nil, false, nil
	}
	return val, true, nil
}

func (m *MockClient) SetString(ctx context.Context, key, value string, ttlSeconds *int64, nx, xx bool) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if nx && xx {
		return false, fmt.Errorf("cannot specify both NX and XX")
	}

	if nx {
		if _, exists := m.strings[key]; exists {
			return false, nil
		}
	}

	if xx {
		if _, exists := m.strings[key]; !exists {
			return false, nil
		}
	}

	m.strings[key] = []byte(value)
	if ttlSeconds != nil && *ttlSeconds > 0 {
		m.ttls[key] = *ttlSeconds
	}
	return true, nil
}

func (m *MockClient) DeleteKey(ctx context.Context, key string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, existsStr := m.strings[key]
	_, existsHash := m.hashes[key]
	_, existsList := m.lists[key]
	_, existsSet := m.sets[key]

	if existsStr {
		delete(m.strings, key)
		delete(m.ttls, key)
		return true, nil
	}
	if existsHash {
		delete(m.hashes, key)
		return true, nil
	}
	if existsList {
		delete(m.lists, key)
		return true, nil
	}
	if existsSet {
		delete(m.sets, key)
		return true, nil
	}
	return false, nil
}

func (m *MockClient) ExistsKeys(ctx context.Context, keys []string) (map[string]bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]bool)
	for _, key := range keys {
		_, existsStr := m.strings[key]
		_, existsHash := m.hashes[key]
		_, existsList := m.lists[key]
		_, existsSet := m.sets[key]
		result[key] = existsStr || existsHash || existsList || existsSet
	}
	return result, nil
}

func (m *MockClient) ExpireKey(ctx context.Context, key string, seconds int64) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.strings[key]; exists {
		m.ttls[key] = seconds
		return true, nil
	}
	return false, nil
}

func (m *MockClient) PersistKey(ctx context.Context, key string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.ttls[key]; exists {
		delete(m.ttls, key)
		return true, nil
	}
	return false, nil
}

func (m *MockClient) RenameKey(ctx context.Context, oldKey, newKey string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	val, exists := m.strings[oldKey]
	if !exists {
		return false, fmt.Errorf("source key does not exist")
	}

	m.strings[newKey] = val
	delete(m.strings, oldKey)

	if ttl, hasTTL := m.ttls[oldKey]; hasTTL {
		m.ttls[newKey] = ttl
		delete(m.ttls, oldKey)
	}

	return true, nil
}

func (m *MockClient) GetTTL(ctx context.Context, key string) (int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if ttl, exists := m.ttls[key]; exists {
		return ttl, nil
	}
	if _, exists := m.strings[key]; exists {
		return -1, nil
	}
	return -2, nil
}

func (m *MockClient) IncrementNumber(ctx context.Context, key string, amount int64) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var current int64
	if val, exists := m.strings[key]; exists {
		fmt.Sscanf(string(val), "%d", &current)
	}

	current += amount
	m.strings[key] = []byte(fmt.Sprintf("%d", current))
	return current, nil
}

func (m *MockClient) DecrementNumber(ctx context.Context, key string, amount int64) (int64, error) {
	return m.IncrementNumber(ctx, key, -amount)
}

func (m *MockClient) AppendString(ctx context.Context, key, value string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	current := m.strings[key]
	m.strings[key] = append(current, []byte(value)...)
	return int64(len(m.strings[key])), nil
}

func (m *MockClient) GetRange(ctx context.Context, key string, start, end int64) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	val, exists := m.strings[key]
	if !exists {
		return []byte{}, nil
	}

	length := int64(len(val))
	if start < 0 {
		start = length + start
	}
	if end < 0 {
		end = length + end
	}
	if start < 0 {
		start = 0
	}
	if end >= length {
		end = length - 1
	}
	if start > end {
		return []byte{}, nil
	}

	result := make([]byte, end-start+1)
	copy(result, val[start:end+1])
	return result, nil
}

// Hash operations

func (m *MockClient) GetMap(ctx context.Context, key string) (map[string][]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if hash, exists := m.hashes[key]; exists {
		result := make(map[string][]byte, len(hash))
		for k, v := range hash {
			result[k] = v
		}
		return result, nil
	}
	return make(map[string][]byte), nil
}

func (m *MockClient) SetMap(ctx context.Context, key string, fields map[string]string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.hashes[key] == nil {
		m.hashes[key] = make(map[string][]byte)
	}

	var newFields int64
	for field, value := range fields {
		if _, exists := m.hashes[key][field]; !exists {
			newFields++
		}
		m.hashes[key][field] = []byte(value)
	}

	return newFields, nil
}

func (m *MockClient) GetMapField(ctx context.Context, key, field string) ([]byte, bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if hash, exists := m.hashes[key]; exists {
		val, fieldExists := hash[field]
		return val, fieldExists, nil
	}
	return nil, false, nil
}

func (m *MockClient) GetMapFields(ctx context.Context, key string, fields []string) (map[string][]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string][]byte)
	if hash, exists := m.hashes[key]; exists {
		for _, field := range fields {
			if val, ok := hash[field]; ok {
				result[field] = val
			}
		}
	}
	return result, nil
}

func (m *MockClient) DeleteMapFields(ctx context.Context, key string, fields []string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var deleted int64
	if hash, exists := m.hashes[key]; exists {
		for _, field := range fields {
			if _, ok := hash[field]; ok {
				delete(hash, field)
				deleted++
			}
		}
	}
	return deleted, nil
}

func (m *MockClient) ListMapKeys(ctx context.Context, key string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if hash, exists := m.hashes[key]; exists {
		keys := make([]string, 0, len(hash))
		for k := range hash {
			keys = append(keys, k)
		}
		return keys, nil
	}
	return []string{}, nil
}

func (m *MockClient) MapFieldExists(ctx context.Context, key, field string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if hash, exists := m.hashes[key]; exists {
		_, fieldExists := hash[field]
		return fieldExists, nil
	}
	return false, nil
}

func (m *MockClient) IncrementMapField(ctx context.Context, key, field string, amount int64) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.hashes[key] == nil {
		m.hashes[key] = make(map[string][]byte)
	}

	var current int64
	if val, exists := m.hashes[key][field]; exists {
		fmt.Sscanf(string(val), "%d", &current)
	}

	current += amount
	m.hashes[key][field] = []byte(fmt.Sprintf("%d", current))
	return current, nil
}

// List operations

func (m *MockClient) PushList(ctx context.Context, key string, values []string, tail bool) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.lists[key] == nil {
		m.lists[key] = [][]byte{}
	}

	byteValues := make([][]byte, len(values))
	for i, v := range values {
		byteValues[i] = []byte(v)
	}

	if tail {
		m.lists[key] = append(m.lists[key], byteValues...)
	} else {
		m.lists[key] = append(byteValues, m.lists[key]...)
	}

	return int64(len(m.lists[key])), nil
}

func (m *MockClient) PopList(ctx context.Context, key string, count int64, tail bool) ([][]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	list, exists := m.lists[key]
	if !exists || len(list) == 0 {
		return [][]byte{}, nil
	}

	if count <= 0 {
		count = 1
	}
	if count > int64(len(list)) {
		count = int64(len(list))
	}

	var result [][]byte
	if tail {
		result = list[len(list)-int(count):]
		m.lists[key] = list[:len(list)-int(count)]
	} else {
		result = list[:count]
		m.lists[key] = list[count:]
	}

	// Return a copy to avoid aliasing
	out := make([][]byte, len(result))
	copy(out, result)
	return out, nil
}

func (m *MockClient) GetListRange(ctx context.Context, key string, start, stop int64) ([][]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	list, exists := m.lists[key]
	if !exists {
		return [][]byte{}, nil
	}

	length := int64(len(list))
	if start < 0 {
		start = length + start
	}
	if stop < 0 {
		stop = length + stop
	}
	if start < 0 {
		start = 0
	}
	if start >= length {
		return [][]byte{}, nil
	}
	if stop >= length {
		stop = length - 1
	}
	if start > stop {
		return [][]byte{}, nil
	}

	slice := list[start : stop+1]
	result := make([][]byte, len(slice))
	copy(result, slice)
	return result, nil
}

func (m *MockClient) GetListLength(ctx context.Context, key string) (int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if list, exists := m.lists[key]; exists {
		return int64(len(list)), nil
	}
	return 0, nil
}

func (m *MockClient) GetListIndex(ctx context.Context, key string, index int64) ([]byte, bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	list, exists := m.lists[key]
	if !exists {
		return nil, false, nil
	}

	length := int64(len(list))
	if index < 0 {
		index = length + index
	}
	if index < 0 || index >= length {
		return nil, false, nil
	}

	result := make([]byte, len(list[index]))
	copy(result, list[index])
	return result, true, nil
}

func (m *MockClient) SetListIndex(ctx context.Context, key string, index int64, value string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	list, exists := m.lists[key]
	if !exists {
		return false, fmt.Errorf("no such key")
	}

	if index < 0 {
		index = int64(len(list)) + index
	}
	if index < 0 || index >= int64(len(list)) {
		return false, fmt.Errorf("index out of range")
	}

	list[index] = []byte(value)
	return true, nil
}

func (m *MockClient) TrimList(ctx context.Context, key string, start, stop int64) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	list, exists := m.lists[key]
	if !exists {
		return true, nil
	}

	length := int64(len(list))
	if start < 0 {
		start = length + start
	}
	if stop < 0 {
		stop = length + stop
	}
	if start < 0 {
		start = 0
	}
	if start >= length {
		m.lists[key] = [][]byte{}
		return true, nil
	}
	if stop >= length {
		stop = length - 1
	}
	if start > stop {
		m.lists[key] = [][]byte{}
		return true, nil
	}

	m.lists[key] = list[start : stop+1]
	return true, nil
}

// Set operations

func (m *MockClient) AddSet(ctx context.Context, key string, members []string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.sets[key] == nil {
		m.sets[key] = make(map[string]bool)
	}

	var added int64
	for _, member := range members {
		if !m.sets[key][member] {
			m.sets[key][member] = true
			added++
		}
	}

	return added, nil
}

func (m *MockClient) RemoveSet(ctx context.Context, key string, members []string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var removed int64
	if set, exists := m.sets[key]; exists {
		for _, member := range members {
			if set[member] {
				delete(set, member)
				removed++
			}
		}
	}

	return removed, nil
}

func (m *MockClient) ListSetMembers(ctx context.Context, key string) ([][]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if set, exists := m.sets[key]; exists {
		members := make([][]byte, 0, len(set))
		for member := range set {
			members = append(members, []byte(member))
		}
		return members, nil
	}
	return [][]byte{}, nil
}

func (m *MockClient) CheckSetMember(ctx context.Context, key, member string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if set, exists := m.sets[key]; exists {
		return set[member], nil
	}
	return false, nil
}

func (m *MockClient) GetSetSize(ctx context.Context, key string) (int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if set, exists := m.sets[key]; exists {
		return int64(len(set)), nil
	}
	return 0, nil
}

func (m *MockClient) PopSet(ctx context.Context, key string, count int64) ([][]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	set, exists := m.sets[key]
	if !exists || len(set) == 0 {
		return [][]byte{}, nil
	}

	if count <= 0 {
		count = 1
	}

	result := make([][]byte, 0, count)
	for member := range set {
		if int64(len(result)) >= count {
			break
		}
		result = append(result, []byte(member))
		delete(set, member)
	}

	return result, nil
}

func (m *MockClient) GetRandomSetMember(ctx context.Context, key string, count int64) ([][]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	set, exists := m.sets[key]
	if !exists || len(set) == 0 {
		return [][]byte{}, nil
	}

	if count <= 0 {
		count = 1
	}

	result := make([][]byte, 0, count)
	for member := range set {
		if int64(len(result)) >= count {
			break
		}
		result = append(result, []byte(member))
	}

	return result, nil
}

// Compile-time check to ensure MockClient implements ValkeyClient
var _ ValkeyClient = (*MockClient)(nil)

// KeysByPattern mock implementation
func (m *MockClient) KeysByPattern(ctx context.Context, pattern string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	keys := make([]string, 0)
	for key := range m.strings {
		if pattern == "*" {
			keys = append(keys, key)
		}
	}
	return keys, nil
}

// ExistsKey mock implementation
func (m *MockClient) ExistsKey(ctx context.Context, key string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, existsStr := m.strings[key]
	_, existsHash := m.hashes[key]
	_, existsList := m.lists[key]
	_, existsSet := m.sets[key]
	return existsStr || existsHash || existsList || existsSet, nil
}

// MemoryUsage mock implementation
func (m *MockClient) MemoryUsage(ctx context.Context, key string) (int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if _, exists := m.strings[key]; exists {
		return 100, nil // Mock value
	}
	return 0, nil
}

// TouchKeys mock implementation
func (m *MockClient) TouchKeys(ctx context.Context, keys []string) (int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	count := int64(0)
	for _, key := range keys {
		if _, exists := m.strings[key]; exists {
			count++
		}
	}
	return count, nil
}

// ObjectEncoding mock implementation
func (m *MockClient) ObjectEncoding(ctx context.Context, key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if _, exists := m.strings[key]; exists {
		return "raw", nil
	}
	return "", nil
}

// GetMapLength mock implementation
func (m *MockClient) GetMapLength(ctx context.Context, key string) (int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if hash, exists := m.hashes[key]; exists {
		return int64(len(hash)), nil
	}
	return 0, nil
}

// ListMapFieldNames mock implementation
func (m *MockClient) ListMapFieldNames(ctx context.Context, key string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if hash, exists := m.hashes[key]; exists {
		fields := make([]string, 0, len(hash))
		for field := range hash {
			fields = append(fields, field)
		}
		return fields, nil
	}
	return []string{}, nil
}

// ListMapFieldValues mock implementation
func (m *MockClient) ListMapFieldValues(ctx context.Context, key string) ([][]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if hash, exists := m.hashes[key]; exists {
		values := make([][]byte, 0, len(hash))
		for _, value := range hash {
			values = append(values, value)
		}
		return values, nil
	}
	return [][]byte{}, nil
}

// GetMapFieldsMultiple mock implementation
func (m *MockClient) GetMapFieldsMultiple(ctx context.Context, key string, fields []string) (map[string][]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string][]byte)
	if hash, exists := m.hashes[key]; exists {
		for _, field := range fields {
			if value, found := hash[field]; found {
				result[field] = value
			}
		}
	}
	return result, nil
}

// SetIntersection mock implementation
func (m *MockClient) SetIntersection(ctx context.Context, keys []string) ([][]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(keys) == 0 {
		return [][]byte{}, nil
	}

	result := make(map[string]bool)
	if firstSet, exists := m.sets[keys[0]]; exists {
		for member := range firstSet {
			result[member] = true
		}
	}

	for _, key := range keys[1:] {
		if set, exists := m.sets[key]; exists {
			for member := range result {
				if !set[member] {
					delete(result, member)
				}
			}
		} else {
			result = make(map[string]bool)
			break
		}
	}

	members := make([][]byte, 0, len(result))
	for member := range result {
		members = append(members, []byte(member))
	}
	return members, nil
}

// SetUnion mock implementation
func (m *MockClient) SetUnion(ctx context.Context, keys []string) ([][]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]bool)
	for _, key := range keys {
		if set, exists := m.sets[key]; exists {
			for member := range set {
				result[member] = true
			}
		}
	}

	members := make([][]byte, 0, len(result))
	for member := range result {
		members = append(members, []byte(member))
	}
	return members, nil
}

// SetDifference mock implementation
func (m *MockClient) SetDifference(ctx context.Context, firstKey string, otherKeys []string) ([][]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]bool)
	if firstSet, exists := m.sets[firstKey]; exists {
		for member := range firstSet {
			result[member] = true
		}
	}

	for _, key := range otherKeys {
		if set, exists := m.sets[key]; exists {
			for member := range set {
				delete(result, member)
			}
		}
	}

	members := make([][]byte, 0, len(result))
	for member := range result {
		members = append(members, []byte(member))
	}
	return members, nil
}

// AddStream mock implementation
func (m *MockClient) AddStream(ctx context.Context, key string, id string, fields map[string]string) (string, error) {
	return id, nil
}

// GetStreamRange mock implementation
func (m *MockClient) GetStreamRange(ctx context.Context, key string, start string, end string, count int64) ([]StreamEntry, error) {
	return []StreamEntry{}, nil
}

// GetStreamLength mock implementation
func (m *MockClient) GetStreamLength(ctx context.Context, key string) (int64, error) {
	return 0, nil
}

// ReadStream mock implementation
func (m *MockClient) ReadStream(ctx context.Context, key string, id string, count int64) ([]StreamEntry, error) {
	return []StreamEntry{}, nil
}

// DumpKey mock implementation
func (m *MockClient) DumpKey(ctx context.Context, key string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if value, exists := m.strings[key]; exists {
		result := make([]byte, len(value))
		copy(result, value)
		return result, nil
	}
	return nil, nil
}

// RestoreKey mock implementation
func (m *MockClient) RestoreKey(ctx context.Context, key string, ttl int64, serialized []byte) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.strings[key] = serialized
	return true, nil
}

// ObjectIdletime mock implementation
func (m *MockClient) ObjectIdletime(ctx context.Context, key string) (int64, error) {
	return 0, nil
}

// ConfigGet mock implementation
func (m *MockClient) ConfigGet(ctx context.Context, parameter string) (map[string]string, error) {
	result := make(map[string]string)
	result[parameter] = "mock_value"
	return result, nil
}

// ConfigSet mock implementation
func (m *MockClient) ConfigSet(ctx context.Context, parameter, value string) (bool, error) {
	return true, nil
}

// GetDatabaseSize mock implementation
func (m *MockClient) GetDatabaseSize(ctx context.Context) (int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return int64(len(m.strings)), nil
}

// GetSlowlog mock implementation
func (m *MockClient) GetSlowlog(ctx context.Context, count int64) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}

// GetClusterInfo mock implementation
func (m *MockClient) GetClusterInfo(ctx context.Context) (map[string]string, error) {
	result := make(map[string]string)
	result["cluster_state"] = "ok"
	return result, nil
}

// GetClusterNodes mock implementation
func (m *MockClient) GetClusterNodes(ctx context.Context) (string, error) {
	return "mock_cluster_nodes", nil
}

// GetKeySlot mock implementation
func (m *MockClient) GetKeySlot(ctx context.Context, key string) (int64, error) {
	return 0, nil
}

// CountKeysInSlot mock implementation
func (m *MockClient) CountKeysInSlot(ctx context.Context, slot int64) (int64, error) {
	return 0, nil
}

// EvalScript mock implementation
func (m *MockClient) EvalScript(ctx context.Context, script string, keys []string, args []string) (interface{}, error) {
	return "OK", nil
}

// LoadScript mock implementation
func (m *MockClient) LoadScript(ctx context.Context, script string) (string, error) {
	return "abc123def456", nil
}

// EvalSHA mock implementation
func (m *MockClient) EvalSHA(ctx context.Context, sha string, keys []string, args []string) (interface{}, error) {
	return "OK", nil
}
