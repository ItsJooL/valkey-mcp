# Valkey MCP Server

An MCP (Model Context Protocol) server for Valkey that provides AI agents like Claude with safe, structured access to Valkey databases. Query data, inspect keys, and debug Valkey instances directly from your AI assistant.

## Features

- **72 Tools** - Comprehensive coverage of Valkey operations (strings, lists, hashes, sets, sorted sets, streams, server admin, cluster management)
- **Type-safe operations** - Input validation and structured responses
- **Complete error handling** - Detailed error messages with context
- **Single binary** - No runtime dependencies (compiled Go binary)
- **Docker support** - Multi-stage build for minimal container image
- **Kubernetes ready** - Deploy with manifests and ConfigMaps
- **MCP-compatible** - Works with Claude Desktop, other MCP clients, and custom applications

## Quick Start

### Prerequisites
- **Valkey** server running (local or remote)
- **Go 1.22+** (for building from source) OR binary installation

### Installation

**Option 1: Install with Go (recommended)**
```bash
go install github.com/ItsJooL/valkey-mcp-server/cmd/valkey-mcp-server@latest
# Binary installed to: $GOPATH/bin/valkey-mcp-server
```

**Option 2: Build from source**
```bash
git clone https://github.com/ItsJooL/valkey-mcp-server
cd valkey-mcp-server
make build
# Binary: ./valkey-mcp-server
```

**Option 3: Download binary**
```bash
# Visit: https://github.com/ItsJooL/valkey-mcp-server/releases
# Download the binary for your OS
chmod +x valkey-mcp-server
```

### Running Locally

**Basic usage (connects to localhost:6379)**
```bash
valkey-mcp-server
```

**With custom connection**
```bash
valkey-mcp-server --url valkey://myhost:6379 --password mysecret --db 1
```

**With environment variables**
```bash
export VALKEY_URL=valkey://myhost:6379
export VALKEY_PASSWORD=mysecret
export VALKEY_DB=1
valkey-mcp-server
```

## Configuration

### Command-line Flags
```bash
-url string        Valkey connection URL (default: "valkey://localhost:6379")
-password string    Valkey authentication password
-db int            Database number 0-15 (default: 0)
```

### Environment Variables
- `VALKEY_URL` - Connection URL (e.g., `valkey://localhost:6379` or `redis://localhost:6379`)
- `VALKEY_PASSWORD` - Authentication password
- `VALKEY_DB` - Database number (0-15)

### URL Format Examples
```
valkey://localhost:6379           # Local Valkey
redis://localhost:6379            # Redis protocol (compatible)
valkey://user:password@host:6379  # With password in URL
rediss://localhost:6380           # TLS connection
```

## Integration with AI

Add the following to your `mcp.json`:

```json
{
  "mcpServers": {
    "valkey": {
      "command": "/path/to/valkey-mcp-server",
      "args": ["-url", "valkey://valkey-address:6379", "-db", "0"],
      "env": {
        "VALKEY_PASSWORD": "valkey-password-here"
      }
    }
  }
}
```


## Docker

**Build and run**
```bash
docker build -t valkey-mcp-server .
docker run --rm -e VALKEY_URL=valkey://host.docker.internal:6379 valkey-mcp-server
```

**Or use Mise**
```bash
mise docker-build
mise docker-run
```

## Kubernetes

**Deploy to cluster**
```bash
kubectl apply -f k8s/valkey-mcp-server.yaml
```

**Customize connection (edit ConfigMap)**
```bash
kubectl edit configmap valkey-mcp-config
kubectl rollout restart deployment/valkey-mcp-server
```

See [k8s/README.md](./k8s/README.md) for detailed Kubernetes setup.

## Development

This project uses **[Mise](https://mise.jdx.dev/installing-mise.html)** for task management and dependency management.

**Run tasks**:
```bash
mise test              # Run all tests
mise test-coverage     # With coverage report
mise lint              # Check code quality
mise build             # Build binary
mise fmt               # Format code
```

**View all available tasks**:
```bash
mise task ls
```



## Available Tools

The server provides 72 tools across these categories:

| Category | Tools | Examples |
|----------|-------|----------|
| **Server** | 5 | `server_ping`, `server_info`, `dbsize`, `config_get`, `slowlog_get` |
| **Keys** | 12 | `scan_keys`, `get_key_type`, `delete_keys`, `expire_key`, `rename_key`, `memory_usage` |
| **Strings** | 9 | `get_string`, `set_string`, `append_string`, `incr_string`, `mget_strings` |
| **Lists** | 10 | `lpush_list`, `rpush_list`, `lrange_list`, `lpop_list`, `lset_list`, `ltrim_list` |
| **Hashes** | 11 | `set_hash`, `get_hash`, `hget_hash_field`, `hdel_hash`, `hincrby_hash` |
| **Sets** | 7 | `add_set`, `remove_set_member`, `get_set_members`, `sinter_sets`, `sunion_sets` |
| **Streams** | 4 | `xadd_stream`, `xrange_stream`, `xread_stream`, `xlen_stream` |
| **Other** | 14 | Scripts, cluster commands, bit operations, etc. |

Run `valkey-mcp-server --help` or query the tool list when connected to see all available tools.

## License

See [LICENSE](./LICENSE) file.
