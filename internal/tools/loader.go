package tools

import (
	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/add_set"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/append_string"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/client_list"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/cluster_count_keysinslot"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/cluster_info"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/cluster_keyslot"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/cluster_nodes"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/config_get"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/config_set"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/dbsize"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/decr_string"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/delete_hash_field"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/delete_keys"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/dump_key"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/eval_script"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/evalsha_script"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/exists_key"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/expire_key"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/get_hash"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/get_hash_field"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/get_hash_fields"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/get_key_ttl"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/get_key_type"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/get_list_index"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/get_list_length"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/get_random_set_member"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/get_set_cardinality"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/get_set_members"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/get_string"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/get_string_range"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/hash_field_exists"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/hkeys_hash"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/hlen_hash"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/hvals_hash"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/incr_hash_field"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/incr_string"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/keys_by_pattern"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/lpop_list"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/lpush_list"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/lrange_list"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/lset_list"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/ltrim_list"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/memory_usage"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/mget_strings"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/object_encoding"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/object_idletime"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/persist_key"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/pop_set_member"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/remove_set_member"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/rename_key"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/restore_key"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/rpop_list"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/rpush_list"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/scan_keys"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/script_load"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/sdiff_sets"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/server_info"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/server_ping"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/set_hash"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/set_is_member"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/set_string"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/sinter_sets"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/slowlog_get"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/string_length"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/sunion_sets"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/touch_keys"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/xadd_stream"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/xlen_stream"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/xrange_stream"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/xread_stream"
)

// RegisterAll registers all available tools with the registry.
func RegisterAll(reg *registry.ToolRegistry, client client.ValkeyClient) {
	server_info.Init(reg, client)
	server_ping.Init(reg, client)
	client_list.Init(reg, client)

	scan_keys.Init(reg, client)
	get_key_type.Init(reg, client)
	get_key_ttl.Init(reg, client)
	delete_keys.Init(reg, client)

	get_string.Init(reg, client)
	set_string.Init(reg, client)
	mget_strings.Init(reg, client)

	lrange_list.Init(reg, client)
	lpush_list.Init(reg, client)
	rpush_list.Init(reg, client)

	get_hash.Init(reg, client)
	set_hash.Init(reg, client)

	get_set_members.Init(reg, client)
	add_set.Init(reg, client)

	incr_string.Init(reg, client)
	append_string.Init(reg, client)
	decr_string.Init(reg, client)
	expire_key.Init(reg, client)
	persist_key.Init(reg, client)
	rename_key.Init(reg, client)

	get_list_length.Init(reg, client)
	get_list_index.Init(reg, client)
	lpop_list.Init(reg, client)
	rpop_list.Init(reg, client)
	lset_list.Init(reg, client)
	ltrim_list.Init(reg, client)

	remove_set_member.Init(reg, client)
	delete_hash_field.Init(reg, client)
	get_set_cardinality.Init(reg, client)
	set_is_member.Init(reg, client)
	get_hash_field.Init(reg, client)
	string_length.Init(reg, client)

	get_string_range.Init(reg, client)
	pop_set_member.Init(reg, client)
	get_random_set_member.Init(reg, client)
	get_hash_fields.Init(reg, client)
	hash_field_exists.Init(reg, client)
	incr_hash_field.Init(reg, client)

	keys_by_pattern.Init(reg, client)
	exists_key.Init(reg, client)
	memory_usage.Init(reg, client)
	touch_keys.Init(reg, client)
	object_encoding.Init(reg, client)
	object_idletime.Init(reg, client)
	hlen_hash.Init(reg, client)
	hkeys_hash.Init(reg, client)
	hvals_hash.Init(reg, client)
	sinter_sets.Init(reg, client)
	sunion_sets.Init(reg, client)
	sdiff_sets.Init(reg, client)

	xadd_stream.Init(reg, client)
	xrange_stream.Init(reg, client)
	xlen_stream.Init(reg, client)
	xread_stream.Init(reg, client)

	dump_key.Init(reg, client)
	restore_key.Init(reg, client)

	config_get.Init(reg, client)
	config_set.Init(reg, client)

	dbsize.Init(reg, client)
	slowlog_get.Init(reg, client)

	cluster_info.Init(reg, client)
	cluster_nodes.Init(reg, client)
	cluster_keyslot.Init(reg, client)
	cluster_count_keysinslot.Init(reg, client)

	eval_script.Init(reg, client)
	script_load.Init(reg, client)
	evalsha_script.Init(reg, client)
}
