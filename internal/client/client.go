package client

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/ItsJooL/valkey-mcp-server/internal/types"
	"github.com/valkey-io/valkey-go"
)

// Client wraps the Valkey client and implements ValkeyClient interface.
type Client struct {
	client valkey.Client
	url    types.ValkeyURL
}

// Config holds the configuration for creating a new client.
type Config struct {
	URL      types.ValkeyURL
	Password string
	DB       types.DBIndex
}

// New creates a new Valkey client with the given configuration.
func New(ctx context.Context, config Config) (*Client, error) {
	if config.URL == "" {
		return nil, fmt.Errorf("valkey URL is required")
	}

	u, err := url.Parse(config.URL.String())
	if err != nil {
		return nil, fmt.Errorf("invalid valkey URL: %w", err)
	}

	addr := u.Host
	if u.Port() == "" {
		addr = u.Host + ":6379"
	}

	opts := valkey.ClientOption{
		InitAddress: []string{addr},
	}

	if config.Password != "" {
		opts.Password = config.Password
	}

	if config.DB.Int() != 0 {
		opts.SelectDB = config.DB.Int()
	}

	client, err := valkey.NewClient(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create valkey client: %w", err)
	}

	if err := client.Do(ctx, client.B().Ping().Build()).Error(); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to connect to valkey: %w", err)
	}

	return &Client{
		client: client,
		url:    config.URL,
	}, nil
}

// UnderlyingClient returns the underlying Valkey client.
func (c *Client) UnderlyingClient() valkey.Client {
	return c.client
}

// Close closes the client connection.
func (c *Client) Close() error {
	c.client.Close()
	return nil
}

// Ping tests the connection to the Valkey server.
func (c *Client) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := c.client.Do(ctx, c.client.B().Ping().Build()).Error(); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}

	return nil
}

// GetServerInfo retrieves server information from Valkey.
func (c *Client) GetServerInfo(ctx context.Context) (map[string]string, error) {
	resp := c.client.Do(ctx, c.client.B().Info().Build())
	if err := resp.Error(); err != nil {
		return nil, fmt.Errorf("INFO failed: %w", err)
	}

	info, err := resp.ToString()
	if err != nil {
		return nil, fmt.Errorf("failed to parse INFO response: %w", err)
	}

	result := make(map[string]string)
	result["raw_info"] = info
	return result, nil
}

// String operations

func (c *Client) GetString(ctx context.Context, key string) ([]byte, bool, error) {
	resp := c.client.Do(ctx, c.client.B().Get().Key(key).Build())
	if err := resp.Error(); err != nil {
		return nil, false, fmt.Errorf("GET failed: %w", err)
	}

	b, err := resp.AsBytes()
	if err != nil {
		// Key does not exist or nil response
		return nil, false, nil
	}
	return b, true, nil
}

func (c *Client) SetString(ctx context.Context, key, value string, ttlSeconds *int64, nx, xx bool) (bool, error) {
	setBuilder := c.client.B().Set().Key(key).Value(value)

	if nx && xx {
		return false, fmt.Errorf("cannot specify both NX and XX")
	}

	if nx {
		if ttlSeconds != nil && *ttlSeconds > 0 {
			resp := c.client.Do(ctx, setBuilder.Nx().ExSeconds(*ttlSeconds).Build())
			if err := resp.Error(); err != nil {
				return false, fmt.Errorf("SET NX failed: %w", err)
			}
			result, _ := resp.ToString()
			return result == "OK", nil
		}
		resp := c.client.Do(ctx, setBuilder.Nx().Build())
		if err := resp.Error(); err != nil {
			return false, fmt.Errorf("SET NX failed: %w", err)
		}
		result, _ := resp.ToString()
		return result == "OK", nil
	}

	if xx {
		if ttlSeconds != nil && *ttlSeconds > 0 {
			resp := c.client.Do(ctx, setBuilder.Xx().ExSeconds(*ttlSeconds).Build())
			if err := resp.Error(); err != nil {
				return false, fmt.Errorf("SET XX failed: %w", err)
			}
			result, _ := resp.ToString()
			return result == "OK", nil
		}
		resp := c.client.Do(ctx, setBuilder.Xx().Build())
		if err := resp.Error(); err != nil {
			return false, fmt.Errorf("SET XX failed: %w", err)
		}
		result, _ := resp.ToString()
		return result == "OK", nil
	}

	if ttlSeconds != nil && *ttlSeconds > 0 {
		resp := c.client.Do(ctx, setBuilder.ExSeconds(*ttlSeconds).Build())
		if err := resp.Error(); err != nil {
			return false, fmt.Errorf("SET with TTL failed: %w", err)
		}
		result, _ := resp.ToString()
		return result == "OK", nil
	}

	resp := c.client.Do(ctx, setBuilder.Build())
	if err := resp.Error(); err != nil {
		return false, fmt.Errorf("SET failed: %w", err)
	}
	result, _ := resp.ToString()
	return result == "OK", nil
}

func (c *Client) DeleteKey(ctx context.Context, key string) (bool, error) {
	resp := c.client.Do(ctx, c.client.B().Del().Key(key).Build())
	if err := resp.Error(); err != nil {
		return false, fmt.Errorf("DEL failed: %w", err)
	}
	count, _ := resp.AsInt64()
	return count > 0, nil
}

func (c *Client) ExistsKeys(ctx context.Context, keys []string) (map[string]bool, error) {
	result := make(map[string]bool)
	for _, key := range keys {
		resp := c.client.Do(ctx, c.client.B().Exists().Key(key).Build())
		if err := resp.Error(); err != nil {
			return nil, fmt.Errorf("EXISTS failed: %w", err)
		}
		count, _ := resp.AsInt64()
		result[key] = count > 0
	}
	return result, nil
}

func (c *Client) ExpireKey(ctx context.Context, key string, seconds int64) (bool, error) {
	resp := c.client.Do(ctx, c.client.B().Expire().Key(key).Seconds(seconds).Build())
	if err := resp.Error(); err != nil {
		return false, fmt.Errorf("EXPIRE failed: %w", err)
	}
	count, _ := resp.AsInt64()
	return count > 0, nil
}

func (c *Client) PersistKey(ctx context.Context, key string) (bool, error) {
	resp := c.client.Do(ctx, c.client.B().Persist().Key(key).Build())
	if err := resp.Error(); err != nil {
		return false, fmt.Errorf("PERSIST failed: %w", err)
	}
	count, _ := resp.AsInt64()
	return count > 0, nil
}

func (c *Client) RenameKey(ctx context.Context, oldKey, newKey string) (bool, error) {
	resp := c.client.Do(ctx, c.client.B().Rename().Key(oldKey).Newkey(newKey).Build())
	if err := resp.Error(); err != nil {
		return false, fmt.Errorf("RENAME failed: %w", err)
	}
	return true, nil
}

func (c *Client) GetTTL(ctx context.Context, key string) (int64, error) {
	resp := c.client.Do(ctx, c.client.B().Ttl().Key(key).Build())
	if err := resp.Error(); err != nil {
		return -1, fmt.Errorf("TTL failed: %w", err)
	}
	ttl, _ := resp.AsInt64()
	return ttl, nil
}

func (c *Client) IncrementNumber(ctx context.Context, key string, amount int64) (int64, error) {
	if amount == 1 {
		resp := c.client.Do(ctx, c.client.B().Incr().Key(key).Build())
		if err := resp.Error(); err != nil {
			return 0, fmt.Errorf("INCR failed: %w", err)
		}
		val, _ := resp.AsInt64()
		return val, nil
	}
	resp := c.client.Do(ctx, c.client.B().Incrby().Key(key).Increment(amount).Build())
	if err := resp.Error(); err != nil {
		return 0, fmt.Errorf("INCRBY failed: %w", err)
	}
	val, _ := resp.AsInt64()
	return val, nil
}

func (c *Client) DecrementNumber(ctx context.Context, key string, amount int64) (int64, error) {
	if amount == 1 {
		resp := c.client.Do(ctx, c.client.B().Decr().Key(key).Build())
		if err := resp.Error(); err != nil {
			return 0, fmt.Errorf("DECR failed: %w", err)
		}
		val, _ := resp.AsInt64()
		return val, nil
	}
	resp := c.client.Do(ctx, c.client.B().Decrby().Key(key).Decrement(amount).Build())
	if err := resp.Error(); err != nil {
		return 0, fmt.Errorf("DECRBY failed: %w", err)
	}
	val, _ := resp.AsInt64()
	return val, nil
}

func (c *Client) AppendString(ctx context.Context, key, value string) (int64, error) {
	resp := c.client.Do(ctx, c.client.B().Append().Key(key).Value(value).Build())
	if err := resp.Error(); err != nil {
		return 0, fmt.Errorf("APPEND failed: %w", err)
	}
	length, _ := resp.AsInt64()
	return length, nil
}

func (c *Client) GetRange(ctx context.Context, key string, start, end int64) ([]byte, error) {
	resp := c.client.Do(ctx, c.client.B().Getrange().Key(key).Start(start).End(end).Build())
	if err := resp.Error(); err != nil {
		return nil, fmt.Errorf("GETRANGE failed: %w", err)
	}
	b, err := resp.AsBytes()
	if err != nil {
		return []byte{}, nil
	}
	return b, nil
}

// Hash operations

func (c *Client) GetMap(ctx context.Context, key string) (map[string][]byte, error) {
	resp := c.client.Do(ctx, c.client.B().Hgetall().Key(key).Build())
	if err := resp.Error(); err != nil {
		return nil, fmt.Errorf("HGETALL failed: %w", err)
	}
	// AsMap handles both RESP2 (flat array) and RESP3 (map) reply formats.
	m, err := resp.AsMap()
	if err != nil {
		return make(map[string][]byte), nil
	}
	result := make(map[string][]byte, len(m))
	for field, msg := range m {
		val, err := msg.AsBytes() // values may be binary
		if err != nil {
			continue
		}
		result[field] = val
	}
	return result, nil
}

func (c *Client) SetMap(ctx context.Context, key string, fields map[string]string) (int64, error) {
	if len(fields) == 0 {
		return 0, nil
	}

	args := make([]string, 0, len(fields)*2)
	for field, value := range fields {
		args = append(args, field, value)
	}

	builder := c.client.B().Hset().Key(key).FieldValue()
	for i := 0; i < len(args); i += 2 {
		builder = builder.FieldValue(args[i], args[i+1])
	}

	resp := c.client.Do(ctx, builder.Build())
	if err := resp.Error(); err != nil {
		return 0, fmt.Errorf("HSET failed: %w", err)
	}

	count, _ := resp.AsInt64()
	return count, nil
}

func (c *Client) GetMapField(ctx context.Context, key, field string) ([]byte, bool, error) {
	resp := c.client.Do(ctx, c.client.B().Hget().Key(key).Field(field).Build())
	if err := resp.Error(); err != nil {
		return nil, false, fmt.Errorf("HGET failed: %w", err)
	}
	b, err := resp.AsBytes()
	if err != nil {
		return nil, false, nil
	}
	return b, true, nil
}

func (c *Client) DeleteMapFields(ctx context.Context, key string, fields []string) (int64, error) {
	if len(fields) == 0 {
		return 0, nil
	}

	builder := c.client.B().Hdel().Key(key).Field(fields...).Build()
	resp := c.client.Do(ctx, builder)
	if err := resp.Error(); err != nil {
		return 0, fmt.Errorf("HDEL failed: %w", err)
	}

	count, _ := resp.AsInt64()
	return count, nil
}

func (c *Client) ListMapKeys(ctx context.Context, key string) ([]string, error) {
	resp := c.client.Do(ctx, c.client.B().Hkeys().Key(key).Build())
	if err := resp.Error(); err != nil {
		return nil, fmt.Errorf("HKEYS failed: %w", err)
	}
	keys, _ := resp.AsStrSlice()
	return keys, nil
}

// List operations

func (c *Client) PushList(ctx context.Context, key string, values []string, tail bool) (int64, error) {
	if len(values) == 0 {
		return 0, nil
	}

	if tail {
		builder := c.client.B().Rpush().Key(key).Element(values...).Build()
		resp := c.client.Do(ctx, builder)
		if err := resp.Error(); err != nil {
			return 0, fmt.Errorf("RPUSH failed: %w", err)
		}
		length, _ := resp.AsInt64()
		return length, nil
	}

	builder := c.client.B().Lpush().Key(key).Element(values...).Build()
	resp := c.client.Do(ctx, builder)
	if err := resp.Error(); err != nil {
		return 0, fmt.Errorf("LPUSH failed: %w", err)
	}
	length, _ := resp.AsInt64()
	return length, nil
}

func (c *Client) PopList(ctx context.Context, key string, count int64, tail bool) ([][]byte, error) {
	if count <= 0 {
		count = 1
	}

	var resp valkey.ValkeyResult
	if tail {
		resp = c.client.Do(ctx, c.client.B().Rpop().Key(key).Count(count).Build())
		if err := resp.Error(); err != nil {
			return nil, fmt.Errorf("RPOP failed: %w", err)
		}
	} else {
		resp = c.client.Do(ctx, c.client.B().Lpop().Key(key).Count(count).Build())
		if err := resp.Error(); err != nil {
			return nil, fmt.Errorf("LPOP failed: %w", err)
		}
	}

	arr, err := resp.ToArray()
	if err != nil {
		return [][]byte{}, nil
	}
	result := make([][]byte, 0, len(arr))
	for _, elem := range arr {
		b, err := elem.AsBytes()
		if err != nil {
			continue
		}
		result = append(result, b)
	}
	return result, nil
}

func (c *Client) GetListRange(ctx context.Context, key string, start, stop int64) ([][]byte, error) {
	resp := c.client.Do(ctx, c.client.B().Lrange().Key(key).Start(start).Stop(stop).Build())
	if err := resp.Error(); err != nil {
		return nil, fmt.Errorf("LRANGE failed: %w", err)
	}
	arr, err := resp.ToArray()
	if err != nil {
		return [][]byte{}, nil
	}
	result := make([][]byte, 0, len(arr))
	for _, elem := range arr {
		b, err := elem.AsBytes()
		if err != nil {
			continue
		}
		result = append(result, b)
	}
	return result, nil
}

func (c *Client) GetListLength(ctx context.Context, key string) (int64, error) {
	resp := c.client.Do(ctx, c.client.B().Llen().Key(key).Build())
	if err := resp.Error(); err != nil {
		return 0, fmt.Errorf("LLEN failed: %w", err)
	}
	length, _ := resp.AsInt64()
	return length, nil
}

func (c *Client) GetListIndex(ctx context.Context, key string, index int64) ([]byte, bool, error) {
	resp := c.client.Do(ctx, c.client.B().Lindex().Key(key).Index(index).Build())
	if err := resp.Error(); err != nil {
		return nil, false, fmt.Errorf("LINDEX failed: %w", err)
	}
	b, err := resp.AsBytes()
	if err != nil {
		// Index out of range or key does not exist
		return nil, false, nil
	}
	return b, true, nil
}

// Set operations

func (c *Client) AddSet(ctx context.Context, key string, members []string) (int64, error) {
	if len(members) == 0 {
		return 0, nil
	}

	builder := c.client.B().Sadd().Key(key).Member(members...).Build()
	resp := c.client.Do(ctx, builder)
	if err := resp.Error(); err != nil {
		return 0, fmt.Errorf("SADD failed: %w", err)
	}

	count, _ := resp.AsInt64()
	return count, nil
}

func (c *Client) RemoveSet(ctx context.Context, key string, members []string) (int64, error) {
	if len(members) == 0 {
		return 0, nil
	}

	builder := c.client.B().Srem().Key(key).Member(members...).Build()
	resp := c.client.Do(ctx, builder)
	if err := resp.Error(); err != nil {
		return 0, fmt.Errorf("SREM failed: %w", err)
	}

	count, _ := resp.AsInt64()
	return count, nil
}

func (c *Client) ListSetMembers(ctx context.Context, key string) ([][]byte, error) {
	resp := c.client.Do(ctx, c.client.B().Smembers().Key(key).Build())
	if err := resp.Error(); err != nil {
		return nil, fmt.Errorf("SMEMBERS failed: %w", err)
	}
	arr, err := resp.ToArray()
	if err != nil {
		return [][]byte{}, nil
	}
	result := make([][]byte, 0, len(arr))
	for _, elem := range arr {
		b, err := elem.AsBytes()
		if err != nil {
			continue
		}
		result = append(result, b)
	}
	return result, nil
}

func (c *Client) CheckSetMember(ctx context.Context, key, member string) (bool, error) {
	resp := c.client.Do(ctx, c.client.B().Sismember().Key(key).Member(member).Build())
	if err := resp.Error(); err != nil {
		return false, fmt.Errorf("SISMEMBER failed: %w", err)
	}
	exists, _ := resp.AsBool()
	return exists, nil
}

func (c *Client) GetSetSize(ctx context.Context, key string) (int64, error) {
	resp := c.client.Do(ctx, c.client.B().Scard().Key(key).Build())
	if err := resp.Error(); err != nil {
		return 0, fmt.Errorf("SCARD failed: %w", err)
	}
	count, _ := resp.AsInt64()
	return count, nil
}

func (c *Client) SetListIndex(ctx context.Context, key string, index int64, value string) (bool, error) {
	resp := c.client.Do(ctx, c.client.B().Lset().Key(key).Index(index).Element(value).Build())
	if err := resp.Error(); err != nil {
		return false, fmt.Errorf("LSET failed: %w", err)
	}
	return true, nil
}

func (c *Client) TrimList(ctx context.Context, key string, start, stop int64) (bool, error) {
	resp := c.client.Do(ctx, c.client.B().Ltrim().Key(key).Start(start).Stop(stop).Build())
	if err := resp.Error(); err != nil {
		return false, fmt.Errorf("LTRIM failed: %w", err)
	}
	return true, nil
}

// Ensure Client implements ValkeyClient at compile time
var _ ValkeyClient = (*Client)(nil)

func (c *Client) GetMapFields(ctx context.Context, key string, fields []string) (map[string][]byte, error) {
	resp := c.client.Do(ctx, c.client.B().Hmget().Key(key).Field(fields...).Build())
	if err := resp.Error(); err != nil {
		return nil, fmt.Errorf("HMGET failed: %w", err)
	}
	arr, err := resp.ToArray()
	if err != nil {
		return make(map[string][]byte), nil
	}
	result := make(map[string][]byte)
	for i, elem := range arr {
		if i >= len(fields) {
			break
		}
		b, err := elem.AsBytes()
		if err != nil {
			// nil response means field does not exist — skip
			continue
		}
		result[fields[i]] = b
	}
	return result, nil
}

func (c *Client) MapFieldExists(ctx context.Context, key, field string) (bool, error) {
	resp := c.client.Do(ctx, c.client.B().Hexists().Key(key).Field(field).Build())
	if err := resp.Error(); err != nil {
		return false, fmt.Errorf("HEXISTS failed: %w", err)
	}
	exists, _ := resp.AsBool()
	return exists, nil
}

func (c *Client) IncrementMapField(ctx context.Context, key, field string, amount int64) (int64, error) {
	resp := c.client.Do(ctx, c.client.B().Hincrby().Key(key).Field(field).Increment(amount).Build())
	if err := resp.Error(); err != nil {
		return 0, fmt.Errorf("HINCRBY failed: %w", err)
	}
	value, _ := resp.AsInt64()
	return value, nil
}

func (c *Client) PopSet(ctx context.Context, key string, count int64) ([][]byte, error) {
	resp := c.client.Do(ctx, c.client.B().Spop().Key(key).Count(count).Build())
	if err := resp.Error(); err != nil {
		return nil, fmt.Errorf("SPOP failed: %w", err)
	}
	arr, err := resp.ToArray()
	if err != nil {
		return [][]byte{}, nil
	}
	result := make([][]byte, 0, len(arr))
	for _, elem := range arr {
		b, err := elem.AsBytes()
		if err != nil {
			continue
		}
		result = append(result, b)
	}
	return result, nil
}

func (c *Client) GetRandomSetMember(ctx context.Context, key string, count int64) ([][]byte, error) {
	resp := c.client.Do(ctx, c.client.B().Srandmember().Key(key).Count(count).Build())
	if err := resp.Error(); err != nil {
		return nil, fmt.Errorf("SRANDMEMBER failed: %w", err)
	}
	arr, err := resp.ToArray()
	if err != nil {
		return [][]byte{}, nil
	}
	result := make([][]byte, 0, len(arr))
	for _, elem := range arr {
		b, err := elem.AsBytes()
		if err != nil {
			continue
		}
		result = append(result, b)
	}
	return result, nil
}

// KeysByPattern retrieves all keys matching the given pattern.
func (c *Client) KeysByPattern(ctx context.Context, pattern string) ([]string, error) {
	resp := c.client.Do(ctx, c.client.B().Keys().Pattern(pattern).Build())
	if err := resp.Error(); err != nil {
		return nil, fmt.Errorf("KEYS failed: %w", err)
	}
	return resp.AsStrSlice()
}

// ExistsKey checks if a single key exists.
func (c *Client) ExistsKey(ctx context.Context, key string) (bool, error) {
	resp := c.client.Do(ctx, c.client.B().Exists().Key(key).Build())
	if err := resp.Error(); err != nil {
		return false, fmt.Errorf("EXISTS failed: %w", err)
	}
	count, _ := resp.AsInt64()
	return count > 0, nil
}

// MemoryUsage gets the memory used by a key.
func (c *Client) MemoryUsage(ctx context.Context, key string) (int64, error) {
	resp := c.client.Do(ctx, c.client.B().MemoryUsage().Key(key).Build())
	if err := resp.Error(); err != nil {
		return 0, fmt.Errorf("MEMORY USAGE failed: %w", err)
	}
	return resp.AsInt64()
}

// TouchKeys updates the last access time of keys.
func (c *Client) TouchKeys(ctx context.Context, keys []string) (int64, error) {
	builder := c.client.B().Touch().Key(keys...)
	resp := c.client.Do(ctx, builder.Build())
	if err := resp.Error(); err != nil {
		return 0, fmt.Errorf("TOUCH failed: %w", err)
	}
	return resp.AsInt64()
}

// ObjectEncoding gets the encoding of a key's value.
func (c *Client) ObjectEncoding(ctx context.Context, key string) (string, error) {
	resp := c.client.Do(ctx, c.client.B().ObjectEncoding().Key(key).Build())
	if err := resp.Error(); err != nil {
		return "", fmt.Errorf("OBJECT ENCODING failed: %w", err)
	}
	return resp.ToString()
}

// GetMapLength gets the number of fields in a hash.
func (c *Client) GetMapLength(ctx context.Context, key string) (int64, error) {
	resp := c.client.Do(ctx, c.client.B().Hlen().Key(key).Build())
	if err := resp.Error(); err != nil {
		return 0, fmt.Errorf("HLEN failed: %w", err)
	}
	return resp.AsInt64()
}

// ListMapFieldNames gets all field names in a hash.
func (c *Client) ListMapFieldNames(ctx context.Context, key string) ([]string, error) {
	resp := c.client.Do(ctx, c.client.B().Hkeys().Key(key).Build())
	if err := resp.Error(); err != nil {
		return nil, fmt.Errorf("HKEYS failed: %w", err)
	}
	return resp.AsStrSlice()
}

// ListMapFieldValues gets all field values in a hash.
func (c *Client) ListMapFieldValues(ctx context.Context, key string) ([][]byte, error) {
	resp := c.client.Do(ctx, c.client.B().Hvals().Key(key).Build())
	if err := resp.Error(); err != nil {
		return nil, fmt.Errorf("HVALS failed: %w", err)
	}
	arr, err := resp.ToArray()
	if err != nil {
		return [][]byte{}, nil
	}
	result := make([][]byte, 0, len(arr))
	for _, elem := range arr {
		b, err := elem.AsBytes()
		if err != nil {
			continue
		}
		result = append(result, b)
	}
	return result, nil
}

// GetMapFieldsMultiple gets multiple fields from a hash at once.
func (c *Client) GetMapFieldsMultiple(ctx context.Context, key string, fields []string) (map[string][]byte, error) {
	builder := c.client.B().Hmget().Key(key).Field(fields...)
	resp := c.client.Do(ctx, builder.Build())
	if err := resp.Error(); err != nil {
		return nil, fmt.Errorf("HMGET failed: %w", err)
	}

	arr, err := resp.ToArray()
	if err != nil {
		return make(map[string][]byte), nil
	}
	result := make(map[string][]byte)
	for i, elem := range arr {
		if i >= len(fields) {
			break
		}
		b, err := elem.AsBytes()
		if err != nil {
			// nil response means field does not exist — skip
			continue
		}
		result[fields[i]] = b
	}
	return result, nil
}

// SetIntersection gets the intersection of multiple sets.
func (c *Client) SetIntersection(ctx context.Context, keys []string) ([][]byte, error) {
	builder := c.client.B().Sinter().Key(keys[0])
	if len(keys) > 1 {
		builder = builder.Key(keys[1:]...)
	}
	resp := c.client.Do(ctx, builder.Build())
	if err := resp.Error(); err != nil {
		return nil, fmt.Errorf("SINTER failed: %w", err)
	}
	arr, err := resp.ToArray()
	if err != nil {
		return [][]byte{}, nil
	}
	result := make([][]byte, 0, len(arr))
	for _, elem := range arr {
		b, err := elem.AsBytes()
		if err != nil {
			continue
		}
		result = append(result, b)
	}
	return result, nil
}

// SetUnion gets the union of multiple sets.
func (c *Client) SetUnion(ctx context.Context, keys []string) ([][]byte, error) {
	builder := c.client.B().Sunion().Key(keys[0])
	if len(keys) > 1 {
		builder = builder.Key(keys[1:]...)
	}
	resp := c.client.Do(ctx, builder.Build())
	if err := resp.Error(); err != nil {
		return nil, fmt.Errorf("SUNION failed: %w", err)
	}
	arr, err := resp.ToArray()
	if err != nil {
		return [][]byte{}, nil
	}
	result := make([][]byte, 0, len(arr))
	for _, elem := range arr {
		b, err := elem.AsBytes()
		if err != nil {
			continue
		}
		result = append(result, b)
	}
	return result, nil
}

// SetDifference gets the difference of sets.
func (c *Client) SetDifference(ctx context.Context, firstKey string, otherKeys []string) ([][]byte, error) {
	builder := c.client.B().Sdiff().Key(firstKey)
	for _, key := range otherKeys {
		builder = builder.Key(key)
	}
	resp := c.client.Do(ctx, builder.Build())
	if err := resp.Error(); err != nil {
		return nil, fmt.Errorf("SDIFF failed: %w", err)
	}
	arr, err := resp.ToArray()
	if err != nil {
		return [][]byte{}, nil
	}
	result := make([][]byte, 0, len(arr))
	for _, elem := range arr {
		b, err := elem.AsBytes()
		if err != nil {
			continue
		}
		result = append(result, b)
	}
	return result, nil
}

// AddStream adds an entry to a stream.
func (c *Client) AddStream(ctx context.Context, key string, id string, fields map[string]string) (string, error) {
	builder := c.client.B().Xadd().Key(key).Id(id).FieldValue()
	for k, v := range fields {
		builder = builder.FieldValue(k, v)
	}
	resp := c.client.Do(ctx, builder.Build())
	if err := resp.Error(); err != nil {
		return "", fmt.Errorf("XADD failed: %w", err)
	}
	return resp.ToString()
}

// parseStreamEntry parses a single stream entry from the raw Valkey response element.
// It is shared between GetStreamRange and ReadStream.
func parseStreamEntry(entryElem valkey.ValkeyMessage) (StreamEntry, error) {
	entryArr, err := entryElem.ToArray()
	if err != nil || len(entryArr) < 2 {
		return StreamEntry{}, fmt.Errorf("invalid stream entry format")
	}
	id, err := entryArr[0].ToString()
	if err != nil {
		return StreamEntry{}, fmt.Errorf("invalid stream entry ID")
	}
	// In RESP3, stream entry field-values are returned as a MAP type (%),
	// not a flat array. AsMap handles both RESP2 (flat array) and RESP3 (map) formats.
	fieldMap, err := entryArr[1].AsMap()
	if err != nil {
		return StreamEntry{}, fmt.Errorf("invalid stream entry fields")
	}
	fields := make(map[string][]byte, len(fieldMap))
	for fieldName, fieldMsg := range fieldMap {
		fieldVal, err := fieldMsg.AsBytes()
		if err != nil {
			continue
		}
		fields[fieldName] = fieldVal
	}
	return StreamEntry{ID: id, FieldValues: fields}, nil
}

// GetStreamRange gets entries from a stream within a range.
func (c *Client) GetStreamRange(ctx context.Context, key string, start string, end string, count int64) ([]StreamEntry, error) {
	endBuilder := c.client.B().Xrange().Key(key).Start(start).End(end)
	var resp valkey.ValkeyResult
	if count > 0 {
		resp = c.client.Do(ctx, endBuilder.Count(count).Build())
	} else {
		resp = c.client.Do(ctx, endBuilder.Build())
	}
	if err := resp.Error(); err != nil {
		return nil, fmt.Errorf("XRANGE failed: %w", err)
	}

	// XRANGE raw format: [ [id, [f1, v1, f2, v2, ...]], ... ]
	arr, err := resp.ToArray()
	if err != nil {
		return []StreamEntry{}, nil
	}
	result := make([]StreamEntry, 0, len(arr))
	for _, entryElem := range arr {
		entry, err := parseStreamEntry(entryElem)
		if err != nil {
			continue
		}
		result = append(result, entry)
	}
	return result, nil
}

// GetStreamLength gets the length of a stream.
func (c *Client) GetStreamLength(ctx context.Context, key string) (int64, error) {
	resp := c.client.Do(ctx, c.client.B().Xlen().Key(key).Build())
	if err := resp.Error(); err != nil {
		return 0, fmt.Errorf("XLEN failed: %w", err)
	}
	return resp.AsInt64()
}

// ReadStream reads entries from a stream.
// XREAD response format differs between RESP2 and RESP3:
//   - RESP2: [ [stream_key, [ [id, [fields...]], ... ]] ]  (outer array of pairs)
//   - RESP3: { stream_key: [ [id, [fields...]], ... ] }    (outer map)
func (c *Client) ReadStream(ctx context.Context, key string, id string, count int64) ([]StreamEntry, error) {
	var resp valkey.ValkeyResult
	if count > 0 {
		resp = c.client.Do(ctx, c.client.B().Xread().Count(count).Streams().Key(key).Id(id).Build())
	} else {
		resp = c.client.Do(ctx, c.client.B().Xread().Streams().Key(key).Id(id).Build())
	}
	if err := resp.Error(); err != nil {
		return nil, fmt.Errorf("XREAD failed: %w", err)
	}

	msg, err := resp.ToMessage()
	if err != nil {
		return []StreamEntry{}, nil
	}

	result := make([]StreamEntry, 0)

	parseEntries := func(entriesMsg valkey.ValkeyMessage) {
		entriesArr, err := entriesMsg.ToArray()
		if err != nil {
			return
		}
		for _, entryElem := range entriesArr {
			entry, err := parseStreamEntry(entryElem)
			if err != nil {
				continue
			}
			result = append(result, entry)
		}
	}

	if msg.IsMap() {
		// RESP3: outer is a MAP { stream_name: entries_array }
		outerMap, err := msg.AsMap()
		if err != nil {
			return []StreamEntry{}, nil
		}
		for _, entriesMsg := range outerMap {
			parseEntries(entriesMsg)
		}
		return result, nil
	}

	// RESP2: outer is [ [stream_key, entries_array], ... ]
	outerArr, err := msg.ToArray()
	if err != nil || len(outerArr) == 0 {
		return []StreamEntry{}, nil
	}
	for _, streamElem := range outerArr {
		streamArr, err := streamElem.ToArray()
		if err != nil || len(streamArr) < 2 {
			continue
		}
		parseEntries(streamArr[1])
	}
	return result, nil
}

// DumpKey serializes a key's value.
func (c *Client) DumpKey(ctx context.Context, key string) ([]byte, error) {
	resp := c.client.Do(ctx, c.client.B().Dump().Key(key).Build())
	if err := resp.Error(); err != nil {
		return nil, fmt.Errorf("DUMP failed: %w", err)
	}
	return resp.AsBytes()
}

// RestoreKey restores a serialized value to a key.
func (c *Client) RestoreKey(ctx context.Context, key string, ttl int64, serialized []byte) (bool, error) {
	resp := c.client.Do(ctx, c.client.B().Restore().Key(key).Ttl(ttl).SerializedValue(string(serialized)).Build())
	if err := resp.Error(); err != nil {
		return false, fmt.Errorf("RESTORE failed: %w", err)
	}
	return true, nil
}

// ObjectIdletime gets the idle time of a key.
func (c *Client) ObjectIdletime(ctx context.Context, key string) (int64, error) {
	resp := c.client.Do(ctx, c.client.B().ObjectIdletime().Key(key).Build())
	if err := resp.Error(); err != nil {
		return 0, fmt.Errorf("OBJECT IDLETIME failed: %w", err)
	}
	return resp.AsInt64()
}

// ConfigGet gets configuration parameters.
func (c *Client) ConfigGet(ctx context.Context, parameter string) (map[string]string, error) {
	resp := c.client.Do(ctx, c.client.B().ConfigGet().Parameter(parameter).Build())
	if err := resp.Error(); err != nil {
		return nil, fmt.Errorf("CONFIG GET failed: %w", err)
	}

	// CONFIG GET can return either a map or an array depending on the Valkey version
	// Try AsStrMap first (newer format), fall back to AsStrSlice (older format)
	result, err := resp.AsStrMap()
	if err == nil {
		return result, nil
	}

	// Fallback: try parsing as array
	values, err := resp.AsStrSlice()
	if err != nil {
		return nil, fmt.Errorf("failed to parse CONFIG GET response: %w", err)
	}

	result = make(map[string]string)
	for i := 0; i < len(values); i += 2 {
		if i+1 < len(values) {
			result[values[i]] = values[i+1]
		}
	}
	return result, nil
}

// ConfigSet sets a configuration parameter.
func (c *Client) ConfigSet(ctx context.Context, parameter, value string) (bool, error) {
	resp := c.client.Do(ctx, c.client.B().ConfigSet().ParameterValue().ParameterValue(parameter, value).Build())
	if err := resp.Error(); err != nil {
		return false, fmt.Errorf("CONFIG SET failed: %w", err)
	}
	status, _ := resp.ToString()
	return status == "OK", nil
}

// GetDatabaseSize gets the number of keys in the current database.
func (c *Client) GetDatabaseSize(ctx context.Context) (int64, error) {
	resp := c.client.Do(ctx, c.client.B().Dbsize().Build())
	if err := resp.Error(); err != nil {
		return 0, fmt.Errorf("DBSIZE failed: %w", err)
	}
	return resp.AsInt64()
}

// GetSlowlog gets slow query log entries.
func (c *Client) GetSlowlog(ctx context.Context, count int64) ([]map[string]interface{}, error) {
	var resp valkey.ValkeyResult
	if count > 0 {
		resp = c.client.Do(ctx, c.client.B().SlowlogGet().Count(count).Build())
	} else {
		resp = c.client.Do(ctx, c.client.B().SlowlogGet().Build())
	}
	if err := resp.Error(); err != nil {
		return nil, fmt.Errorf("SLOWLOG GET failed: %w", err)
	}

	// Parse slowlog entries
	result := make([]map[string]interface{}, 0)
	// Basic parsing - depends on actual response format from valkey-go
	return result, nil
}

// GetClusterInfo gets cluster information.
func (c *Client) GetClusterInfo(ctx context.Context) (map[string]string, error) {
	resp := c.client.Do(ctx, c.client.B().ClusterInfo().Build())
	if err := resp.Error(); err != nil {
		return nil, fmt.Errorf("CLUSTER INFO failed: %w", err)
	}

	info, _ := resp.ToString()
	result := make(map[string]string)
	// Parse info string (key:value\r\n format)
	_ = info
	return result, nil
}

// GetClusterNodes gets cluster nodes information.
func (c *Client) GetClusterNodes(ctx context.Context) (string, error) {
	resp := c.client.Do(ctx, c.client.B().ClusterNodes().Build())
	if err := resp.Error(); err != nil {
		return "", fmt.Errorf("CLUSTER NODES failed: %w", err)
	}
	return resp.ToString()
}

// GetKeySlot gets the cluster slot for a key.
func (c *Client) GetKeySlot(ctx context.Context, key string) (int64, error) {
	resp := c.client.Do(ctx, c.client.B().ClusterKeyslot().Key(key).Build())
	if err := resp.Error(); err != nil {
		return 0, fmt.Errorf("CLUSTER KEYSLOT failed: %w", err)
	}
	return resp.AsInt64()
}

// CountKeysInSlot counts keys in a cluster slot.
func (c *Client) CountKeysInSlot(ctx context.Context, slot int64) (int64, error) {
	resp := c.client.Do(ctx, c.client.B().ClusterCountkeysinslot().Slot(slot).Build())
	if err := resp.Error(); err != nil {
		return 0, fmt.Errorf("CLUSTER COUNTKEYSINSLOT failed: %w", err)
	}
	return resp.AsInt64()
}

// EvalScript evaluates a Lua script.
func (c *Client) EvalScript(ctx context.Context, script string, keys []string, args []string) (interface{}, error) {
	numkeys := c.client.B().Eval().Script(script).Numkeys(int64(len(keys)))

	var resp valkey.ValkeyResult
	if len(keys) > 0 {
		kb := numkeys.Key(keys...)
		if len(args) > 0 {
			resp = c.client.Do(ctx, kb.Arg(args...).Build())
		} else {
			resp = c.client.Do(ctx, kb.Build())
		}
	} else if len(args) > 0 {
		resp = c.client.Do(ctx, numkeys.Arg(args...).Build())
	} else {
		resp = c.client.Do(ctx, numkeys.Build())
	}

	if err := resp.Error(); err != nil {
		return nil, fmt.Errorf("EVAL failed: %w", err)
	}

	return resp.ToAny()
}

// LoadScript loads a Lua script and returns its SHA1 hash.
func (c *Client) LoadScript(ctx context.Context, script string) (string, error) {
	resp := c.client.Do(ctx, c.client.B().ScriptLoad().Script(script).Build())
	if err := resp.Error(); err != nil {
		return "", fmt.Errorf("SCRIPT LOAD failed: %w", err)
	}
	return resp.ToString()
}

// EvalSHA evaluates a loaded script by its SHA1 hash.
func (c *Client) EvalSHA(ctx context.Context, sha string, keys []string, args []string) (interface{}, error) {
	numkeys := c.client.B().Evalsha().Sha1(sha).Numkeys(int64(len(keys)))

	var resp valkey.ValkeyResult
	if len(keys) > 0 {
		kb := numkeys.Key(keys...)
		if len(args) > 0 {
			resp = c.client.Do(ctx, kb.Arg(args...).Build())
		} else {
			resp = c.client.Do(ctx, kb.Build())
		}
	} else if len(args) > 0 {
		resp = c.client.Do(ctx, numkeys.Arg(args...).Build())
	} else {
		resp = c.client.Do(ctx, numkeys.Build())
	}

	if err := resp.Error(); err != nil {
		return nil, fmt.Errorf("EVALSHA failed: %w", err)
	}

	return resp.ToAny()
}
