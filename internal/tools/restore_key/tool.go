package restore_key

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools/base"
)

type Tool struct {
	base.BaseTool
	client client.ValkeyClient
}

type Input struct {
	Key             string `json:"key" jsonschema:"required,description=Key to restore to"`
	TTL             int64  `json:"ttl" jsonschema:"description=TTL in milliseconds (0 for no expiry)"`
	SerializedValue string `json:"serialized_value" jsonschema:"description=Base64-encoded serialized value (alternative to serialized)"`
	Serialized      string `json:"serialized" jsonschema:"description=Base64-encoded serialized value (alternative to serialized_value)"`
}

type Output struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func NewTool(client client.ValkeyClient) registry.Tool {
	return &Tool{
		BaseTool: base.NewBaseTool("restore_key", "Restore serialized value to key (accepts base64-encoded serialization)", Input{}),
		client:   client,
	}
}

func (t *Tool) Execute(ctx context.Context, input json.RawMessage) (interface{}, error) {
	var params Input
	if err := t.ParseInput(input, &params); err != nil {
		return nil, err
	}

	if params.Key == "" {
		return nil, fmt.Errorf("key cannot be empty")
	}

	// Support both serialized and serialized_value parameters
	serializedData := params.Serialized
	if serializedData == "" {
		serializedData = params.SerializedValue
	}
	
	if serializedData == "" {
		return nil, fmt.Errorf("either 'serialized' or 'serialized_value' must be provided")
	}

	// Decode base64-encoded data
	decoded, err := base64.StdEncoding.DecodeString(serializedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 serialized data: %w", err)
	}

	success, err := t.client.RestoreKey(ctx, params.Key, params.TTL, decoded)
	if err != nil {
		return nil, fmt.Errorf("failed to restore key: %w", err)
	}

	message := "Key restored successfully"
	if !success {
		message = "Key restoration completed with warnings"
	}

	return Output{
		Success: success,
		Message: message,
	}, nil
}

func Init(reg *registry.ToolRegistry, client client.ValkeyClient) {
	reg.MustRegister(NewTool(client))
}
